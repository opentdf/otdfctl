package migrations

import (
	"context"
	"fmt"
	"testing"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/registeredresources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockMigrationHandler implements MigrationHandler for testing.
type MockMigrationHandler struct {
	Resources      []*policy.RegisteredResource
	ResourceValues map[string][]*policy.RegisteredResourceValue // keyed by resource ID
	Namespaces     []*policy.Namespace

	// Track calls
	CreatedResources      []createdResourceCall
	CreatedResourceValues []createdResourceValueCall
	DeletedResourceIDs    []string

	// Control behavior
	CreateResourceErr      error
	CreateResourceValueErr error
	DeleteResourceErr      error
}

type createdResourceCall struct {
	Namespace string
	Name      string
	Values    []string
	Metadata  *common.MetadataMutable
}

type createdResourceValueCall struct {
	ResourceID          string
	Value               string
	ActionAttributeVals []*registeredresources.ActionAttributeValue
	Metadata            *common.MetadataMutable
}

func (m *MockMigrationHandler) ListRegisteredResources(_ context.Context, limit, offset int32, _ string) (*registeredresources.ListRegisteredResourcesResponse, error) {
	start := int(offset)
	if start >= len(m.Resources) {
		return &registeredresources.ListRegisteredResourcesResponse{}, nil
	}
	end := start + int(limit)
	if end > len(m.Resources) {
		end = len(m.Resources)
	}
	return &registeredresources.ListRegisteredResourcesResponse{
		Resources: m.Resources[start:end],
	}, nil
}

func (m *MockMigrationHandler) ListRegisteredResourceValues(_ context.Context, resourceID string, limit, offset int32) (*registeredresources.ListRegisteredResourceValuesResponse, error) {
	values := m.ResourceValues[resourceID]
	start := int(offset)
	if start >= len(values) {
		return &registeredresources.ListRegisteredResourceValuesResponse{}, nil
	}
	end := start + int(limit)
	if end > len(values) {
		end = len(values)
	}
	return &registeredresources.ListRegisteredResourceValuesResponse{
		Values: values[start:end],
	}, nil
}

func (m *MockMigrationHandler) CreateRegisteredResource(_ context.Context, namespace, name string, values []string, metadata *common.MetadataMutable) (*policy.RegisteredResource, error) {
	m.CreatedResources = append(m.CreatedResources, createdResourceCall{
		Namespace: namespace,
		Name:      name,
		Values:    values,
		Metadata:  metadata,
	})
	if m.CreateResourceErr != nil {
		return nil, m.CreateResourceErr
	}

	// Build response with values
	rrValues := make([]*policy.RegisteredResourceValue, 0, len(values))
	for i, v := range values {
		rrValues = append(rrValues, &policy.RegisteredResourceValue{
			Id:    fmt.Sprintf("new-value-%d", i),
			Value: v,
		})
	}

	return &policy.RegisteredResource{
		Id:     "new-resource-id",
		Name:   name,
		Values: rrValues,
		Namespace: &policy.Namespace{
			Id:  "ns-id",
			Fqn: namespace,
		},
	}, nil
}

func (m *MockMigrationHandler) CreateRegisteredResourceValue(_ context.Context, resourceID string, value string, actionAttributeValues []*registeredresources.ActionAttributeValue, metadata *common.MetadataMutable) (*policy.RegisteredResourceValue, error) {
	m.CreatedResourceValues = append(m.CreatedResourceValues, createdResourceValueCall{
		ResourceID:          resourceID,
		Value:               value,
		ActionAttributeVals: actionAttributeValues,
		Metadata:            metadata,
	})
	if m.CreateResourceValueErr != nil {
		return nil, m.CreateResourceValueErr
	}
	return &policy.RegisteredResourceValue{
		Id:    "new-recreated-value-id",
		Value: value,
	}, nil
}

func (m *MockMigrationHandler) DeleteRegisteredResource(_ context.Context, id string) error {
	m.DeletedResourceIDs = append(m.DeletedResourceIDs, id)
	if m.DeleteResourceErr != nil {
		return m.DeleteResourceErr
	}
	return nil
}

func (m *MockMigrationHandler) ListNamespaces(_ context.Context, _ common.ActiveStateEnum, limit, offset int32) (*namespaces.ListNamespacesResponse, error) {
	start := int(offset)
	if start >= len(m.Namespaces) {
		return &namespaces.ListNamespacesResponse{}, nil
	}
	end := start + int(limit)
	if end > len(m.Namespaces) {
		end = len(m.Namespaces)
	}
	return &namespaces.ListNamespacesResponse{
		Namespaces: m.Namespaces[start:end],
	}, nil
}

func TestBuildRegisteredResourcePlan(t *testing.T) {
	t.Run("builds plan with resources lacking namespaces", func(t *testing.T) {
		mock := &MockMigrationHandler{
			Resources: []*policy.RegisteredResource{
				{Id: "res-1", Name: "resource-one"},
				{Id: "res-2", Name: "resource-two", Namespace: &policy.Namespace{Id: "ns-1", Fqn: "https://example.com"}},
				{Id: "res-3", Name: "resource-three"},
			},
			ResourceValues: map[string][]*policy.RegisteredResourceValue{
				"res-1": {
					{Id: "val-1", Value: "value-one"},
					{Id: "val-2", Value: "value-two"},
				},
				"res-3": {
					{Id: "val-3", Value: "value-three"},
				},
			},
		}

		plan, err := buildRegisteredResourcePlan(context.Background(), mock)
		require.NoError(t, err)

		// Should only include resources without namespaces (res-1 and res-3)
		assert.Len(t, plan, 2)
		assert.Equal(t, "res-1", plan[0].Resource.GetId())
		assert.Equal(t, "resource-one", plan[0].Resource.GetName())
		assert.Len(t, plan[0].Values, 2)
		assert.Equal(t, "res-3", plan[1].Resource.GetId())
		assert.Len(t, plan[1].Values, 1)
	})

	t.Run("returns empty plan when no resources exist", func(t *testing.T) {
		mock := &MockMigrationHandler{}

		plan, err := buildRegisteredResourcePlan(context.Background(), mock)
		require.NoError(t, err)
		assert.Empty(t, plan)
	})

	t.Run("returns empty plan when all resources have namespaces", func(t *testing.T) {
		mock := &MockMigrationHandler{
			Resources: []*policy.RegisteredResource{
				{Id: "res-1", Name: "resource-one", Namespace: &policy.Namespace{Id: "ns-1"}},
				{Id: "res-2", Name: "resource-two", Namespace: &policy.Namespace{Id: "ns-2"}},
			},
		}

		plan, err := buildRegisteredResourcePlan(context.Background(), mock)
		require.NoError(t, err)
		assert.Empty(t, plan)
	})
}

func TestCommitRegisteredResourceMigration(t *testing.T) {
	t.Run("creates resource with correct namespace and name", func(t *testing.T) {
		mock := &MockMigrationHandler{}

		plan := RegisteredResourceMigrationPlan{
			Resource: &policy.RegisteredResource{
				Id:   "old-id",
				Name: "my-resource",
				Metadata: &common.Metadata{
					Labels: map[string]string{"env": "prod"},
				},
			},
			Values: []*policy.RegisteredResourceValue{
				{Id: "old-val-1", Value: "val-a"},
				{Id: "old-val-2", Value: "val-b"},
			},
			TargetNamespace: "https://example.com",
			Commit:          true,
		}

		err := commitRegisteredResourceMigration(context.Background(), mock, plan)
		require.NoError(t, err)

		// Verify resource was created without values (values are created individually)
		require.Len(t, mock.CreatedResources, 1)
		assert.Equal(t, "https://example.com", mock.CreatedResources[0].Namespace)
		assert.Equal(t, "my-resource", mock.CreatedResources[0].Name)
		assert.Nil(t, mock.CreatedResources[0].Values)
		assert.Equal(t, map[string]string{"env": "prod"}, mock.CreatedResources[0].Metadata.GetLabels())

		// Verify values were created individually
		require.Len(t, mock.CreatedResourceValues, 2)
		assert.Equal(t, "new-resource-id", mock.CreatedResourceValues[0].ResourceID)
		assert.Equal(t, "val-a", mock.CreatedResourceValues[0].Value)
		assert.Equal(t, "val-b", mock.CreatedResourceValues[1].Value)

		// Verify old resource was deleted
		assert.Contains(t, mock.DeletedResourceIDs, "old-id")
	})

	t.Run("re-creates values with action-attribute mappings", func(t *testing.T) {
		mock := &MockMigrationHandler{}

		plan := RegisteredResourceMigrationPlan{
			Resource: &policy.RegisteredResource{
				Id:   "old-id",
				Name: "my-resource",
			},
			Values: []*policy.RegisteredResourceValue{
				{
					Id:    "old-val-1",
					Value: "val-a",
					ActionAttributeValues: []*policy.RegisteredResourceValue_ActionAttributeValue{
						{
							Id:             "aav-1",
							Action:         &policy.Action{Id: "action-1"},
							AttributeValue: &policy.Value{Id: "attr-val-1"},
						},
					},
				},
			},
			TargetNamespace: "https://example.com",
			Commit:          true,
		}

		err := commitRegisteredResourceMigration(context.Background(), mock, plan)
		require.NoError(t, err)

		// Should have re-created the value with AAVs
		require.Len(t, mock.CreatedResourceValues, 1)
		assert.Equal(t, "new-resource-id", mock.CreatedResourceValues[0].ResourceID)
		assert.Equal(t, "val-a", mock.CreatedResourceValues[0].Value)
		require.Len(t, mock.CreatedResourceValues[0].ActionAttributeVals, 1)
	})

	t.Run("returns error when create fails", func(t *testing.T) {
		mock := &MockMigrationHandler{
			CreateResourceErr: fmt.Errorf("create failed"),
		}

		plan := RegisteredResourceMigrationPlan{
			Resource: &policy.RegisteredResource{
				Id:   "old-id",
				Name: "my-resource",
			},
			TargetNamespace: "https://example.com",
			Commit:          true,
		}

		err := commitRegisteredResourceMigration(context.Background(), mock, plan)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create failed")

		// Old resource should NOT have been deleted
		assert.Empty(t, mock.DeletedResourceIDs)
	})

	t.Run("returns error when delete fails", func(t *testing.T) {
		mock := &MockMigrationHandler{
			DeleteResourceErr: fmt.Errorf("delete failed"),
		}

		plan := RegisteredResourceMigrationPlan{
			Resource: &policy.RegisteredResource{
				Id:   "old-id",
				Name: "my-resource",
			},
			Values: []*policy.RegisteredResourceValue{
				{Id: "old-val-1", Value: "val-a"},
			},
			TargetNamespace: "https://example.com",
			Commit:          true,
		}

		err := commitRegisteredResourceMigration(context.Background(), mock, plan)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})

	t.Run("returns error when plan not ready", func(t *testing.T) {
		mock := &MockMigrationHandler{}

		plan := RegisteredResourceMigrationPlan{
			Resource: &policy.RegisteredResource{Id: "old-id"},
			Commit:   false,
		}

		err := commitRegisteredResourceMigration(context.Background(), mock, plan)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not ready for commit")
	})
}

func TestBuildNamespaceOptions(t *testing.T) {
	t.Run("builds options from namespaces with FQN", func(t *testing.T) {
		nsList := []*policy.Namespace{
			{Id: "ns-1", Name: "example", Fqn: "https://example.com"},
			{Id: "ns-2", Name: "other", Fqn: "https://other.org"},
		}

		opts := buildNamespaceOptions(nsList)
		assert.Len(t, opts, 2)
	})

	t.Run("returns empty options for empty namespace list", func(t *testing.T) {
		opts := buildNamespaceOptions(nil)
		assert.Empty(t, opts)
	})
}

func TestConvertActionAttributeValues(t *testing.T) {
	t.Run("converts action-attribute values correctly", func(t *testing.T) {
		aavs := []*policy.RegisteredResourceValue_ActionAttributeValue{
			{
				Id:             "aav-1",
				Action:         &policy.Action{Id: "action-1"},
				AttributeValue: &policy.Value{Id: "attr-val-1"},
			},
			{
				Id:             "aav-2",
				Action:         &policy.Action{Id: "action-2"},
				AttributeValue: &policy.Value{Id: "attr-val-2"},
			},
		}

		result := convertActionAttributeValues(aavs)
		require.Len(t, result, 2)
		assert.Equal(t, "action-1", result[0].GetActionId())
		assert.Equal(t, "attr-val-1", result[0].GetAttributeValueId())
		assert.Equal(t, "action-2", result[1].GetActionId())
		assert.Equal(t, "attr-val-2", result[1].GetAttributeValueId())
	})

	t.Run("handles empty input", func(t *testing.T) {
		result := convertActionAttributeValues(nil)
		assert.Empty(t, result)
	})
}
