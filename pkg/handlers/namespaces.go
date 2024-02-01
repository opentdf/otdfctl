package handlers

import (
	"github.com/opentdf/opentdf-v2-poc/sdk/namespaces"
)

func (h Handler) GetNamespace(id string) (*namespaces.Namespace, error) {
	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, &namespaces.GetNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (h Handler) ListNamespaces() ([]*namespaces.Namespace, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(h.ctx, &namespaces.ListNamespacesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Namespaces, nil
}

func (h Handler) CreateNamespace(name string) (*namespaces.Namespace, error) {
	resp, err := h.sdk.Namespaces.CreateNamespace(h.ctx, &namespaces.CreateNamespaceRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return resp.Namespace, nil
}

func (h Handler) UpdateNamespace(id string, name string) (*namespaces.Namespace, error) {
	resp, err := h.sdk.Namespaces.UpdateNamespace(h.ctx, &namespaces.UpdateNamespaceRequest{
		Id:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return resp.Namespace, nil
}

func (h Handler) DeleteNamespace(id string) error {
	_, err := h.sdk.Namespaces.DeleteNamespace(h.ctx, &namespaces.DeleteNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return err
	}

	return nil
}
