package utils

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "testing"

    "github.com/opentdf/platform/protocol/go/policy"
    "github.com/stretchr/testify/require"
)

func pemBlockForKey(t *testing.T, pub interface{}) []byte {
    t.Helper()
    var der []byte
    var err error
    switch k := pub.(type) {
    case *rsa.PublicKey:
        der, err = x509.MarshalPKIXPublicKey(k)
        require.NoError(t, err)
    case *ecdsa.PublicKey:
        der, err = x509.MarshalPKIXPublicKey(k)
        require.NoError(t, err)
    default:
        t.Fatalf("unsupported key type")
    }
    return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func TestValidatePublicKeyPEM_RSA2048_OK(t *testing.T) {
    k, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)
    pub := &k.PublicKey
    pemBytes := pemBlockForKey(t, pub)

    info, err := ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_RSA_2048)
    require.NoError(t, err)
    require.Equal(t, "rsa:2048", info.Alg)
    require.Equal(t, 2048, info.Bits)
    require.NotEmpty(t, info.Fingerprint)
}

func TestValidatePublicKeyPEM_RSA_SizeMismatch(t *testing.T) {
    k, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)
    pemBytes := pemBlockForKey(t, &k.PublicKey)

    _, err = ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_RSA_4096)
    require.Error(t, err)
}

func TestValidatePublicKeyPEM_EC_P256_OK(t *testing.T) {
    k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    require.NoError(t, err)
    pemBytes := pemBlockForKey(t, &k.PublicKey)

    info, err := ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_EC_P256)
    require.NoError(t, err)
    require.Equal(t, "ec:secp256r1", info.Alg)
    require.Equal(t, "secp256r1", info.Curve)
    require.NotEmpty(t, info.Fingerprint)
}

func TestValidatePublicKeyPEM_EC_Mismatch(t *testing.T) {
    k, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
    require.NoError(t, err)
    pemBytes := pemBlockForKey(t, &k.PublicKey)

    _, err = ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_EC_P256)
    require.Error(t, err)
}

func TestValidatePublicKeyPEM_InvalidPEM(t *testing.T) {
    _, err := ValidatePublicKeyPEM([]byte("not a pem"), policy.Algorithm_ALGORITHM_RSA_2048)
    require.Error(t, err)
}

func TestValidatePublicKeyPEM_RSA_PKCS1_Block_OK(t *testing.T) {
    k, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)
    der := x509.MarshalPKCS1PublicKey(&k.PublicKey)
    pemBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: der})
    info, err := ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_RSA_2048)
    require.NoError(t, err)
    require.Equal(t, "rsa:2048", info.Alg)
}

func TestValidatePublicKeyPEM_UnsupportedBlockType(t *testing.T) {
    // Create a CERTIFICATE block to simulate wrong type
    k, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)
    // self-signed-like DER but malformed on purpose: use public key DER as bytes under CERTIFICATE type
    der, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
    require.NoError(t, err)
    pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
    _, err = ValidatePublicKeyPEM(pemBytes, policy.Algorithm_ALGORITHM_RSA_2048)
    require.Error(t, err)
}

func TestRoundTripBase64DecodeThenValidate(t *testing.T) {
    k, err := rsa.GenerateKey(rand.Reader, 2048)
    require.NoError(t, err)
    pemBytes := pemBlockForKey(t, &k.PublicKey)
    b64 := base64.StdEncoding.EncodeToString(pemBytes)
    decoded, err := base64.StdEncoding.DecodeString(b64)
    require.NoError(t, err)
    _, err = ValidatePublicKeyPEM(decoded, policy.Algorithm_ALGORITHM_RSA_2048)
    require.NoError(t, err)
}
