package handlers

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
)

// ConvertPEMToX5C converts a PEM-encoded certificate to x5c format (base64-encoded DER)
func ConvertPEMToX5C(pemData []byte) (string, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM certificate")
	}
	if block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("PEM block is not a certificate, got: %s", block.Type)
	}

	// Validate it's a valid certificate
	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Convert DER to base64 (x5c format)
	x5c := base64.StdEncoding.EncodeToString(block.Bytes)
	return x5c, nil
}

// ConvertX5CToPEM converts an x5c format certificate back to PEM
func ConvertX5CToPEM(x5c string) ([]byte, error) {
	derBytes, err := base64.StdEncoding.DecodeString(x5c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode x5c: %w", err)
	}

	// Validate it's a valid certificate
	_, err = x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

// GetCertificate retrieves a certificate by ID
func (h Handler) GetCertificate(ctx context.Context, id string) (*policy.Certificate, error) {
	// Note: This would require adding a GetCertificate RPC to the namespace service
	// For now, we'll get it via the namespace
	return nil, fmt.Errorf("GetCertificate not yet implemented in platform API")
}

// AssignCertificateToNamespace assigns a certificate to a namespace
func (h Handler) AssignCertificateToNamespace(ctx context.Context, namespaceID, x5c string, labels map[string]string) (*namespaces.AssignCertificateToNamespaceResponse, error) {
	metadata := &common.MetadataMutable{}
	if labels != nil {
		metadata.Labels = labels
	}

	req := &namespaces.AssignCertificateToNamespaceRequest{
		NamespaceId: namespaceID,
		X5C:         x5c,
		Metadata:    metadata,
	}

	resp, err := h.sdk.Namespaces.AssignCertificateToNamespace(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to assign certificate to namespace [%s]: %w", namespaceID, err)
	}

	return resp, nil
}

// RemoveCertificateFromNamespace removes a certificate from a namespace
func (h Handler) RemoveCertificateFromNamespace(ctx context.Context, namespaceID, certID string) error {
	req := &namespaces.RemoveCertificateFromNamespaceRequest{
		NamespaceCertificate: &namespaces.NamespaceCertificate{
			NamespaceId:   namespaceID,
			CertificateId: certID,
		},
	}

	_, err := h.sdk.Namespaces.RemoveCertificateFromNamespace(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to remove certificate [%s] from namespace [%s]: %w", certID, namespaceID, err)
	}

	return nil
}

// ListNamespaceCertificates lists certificates for a namespace
func (h Handler) ListNamespaceCertificates(ctx context.Context, namespaceID string) ([]*policy.Certificate, error) {
	ns, err := h.GetNamespace(ctx, namespaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	return ns.GetRootCerts(), nil
}

// GetNamespaceWithCertificates retrieves a namespace with its certificates
func (h Handler) GetNamespaceWithCertificates(ctx context.Context, identifier string) (*policy.Namespace, error) {
	req := &namespaces.GetNamespaceRequest{}

	// Try to parse as UUID first, otherwise use as FQN
	req.Identifier = &namespaces.GetNamespaceRequest_Fqn{
		Fqn: identifier,
	}

	resp, err := h.sdk.Namespaces.GetNamespace(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace [%s]: %w", identifier, err)
	}

	return resp.GetNamespace(), nil
}
