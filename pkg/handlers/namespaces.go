package handlers

import (
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
)

func (h Handler) GetNamespace(id string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, &namespaces.GetNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (h Handler) ListNamespaces() ([]*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(h.ctx, &namespaces.ListNamespacesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Namespaces, nil
}

func (h Handler) CreateNamespace(name string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.CreateNamespace(h.ctx, &namespaces.CreateNamespaceRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (h Handler) UpdateNamespace(id string, name string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.UpdateNamespace(h.ctx, &namespaces.UpdateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Namespace, nil
}

func (h Handler) DeactivateNamespace(id string) error {
	_, err := h.sdk.Namespaces.DeactivateNamespace(h.ctx, &namespaces.DeactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return err
	}

	return nil
}
