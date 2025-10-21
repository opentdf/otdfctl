package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/pem"
    "errors"
    "fmt"

    "github.com/opentdf/platform/lib/ocrypto"
    "github.com/opentdf/platform/protocol/go/policy"
)

// ValidatePublicKeyPEM validates a PEM-encoded public key block and ensures it
// matches the expected algorithm. The input should be raw PEM bytes (not base64).
func ValidatePublicKeyPEM(pemBytes []byte, expected policy.Algorithm) error {
    var info KeyInfo
    if len(pemBytes) == 0 {
        return info, errors.New("empty pem input")
    }

    enc, err := ocrypto.FromPublicPEM(string(pemBytes))
    if err != nil {
        return fmt.Errorf("invalid public key pem: %w", err)
    }

    switch expected { //nolint:exhaustive
    case policy.Algorithm_ALGORITHM_RSA_2048:
        if enc.KeyType() != ocrypto.RSA2048Key {
            return info, errors.New("algorithm mismatch: expected RSA 2048")
        }
    case policy.Algorithm_ALGORITHM_RSA_4096:
        if enc.KeyType() != ocrypto.RSA4096Key {
            return info, errors.New("algorithm mismatch: expected RSA 4096")
        }
    case policy.Algorithm_ALGORITHM_EC_P256:
        if enc.KeyType() != ocrypto.EC256Key {
            return info, errors.New("algorithm mismatch: expected EC P-256")
        }
    case policy.Algorithm_ALGORITHM_EC_P384:
        if enc.KeyType() != ocrypto.EC384Key {
            return info, errors.New("algorithm mismatch: expected EC P-384")
        }
        info.Alg = "ec:secp384r1"
        info.Curve = "secp384r1"
    case policy.Algorithm_ALGORITHM_EC_P521:
        if enc.KeyType() != ocrypto.EC521Key {
            return info, errors.New("algorithm mismatch: expected EC P-521")
        }
    default:
        return info, errors.New("unsupported or unspecified algorithm")
    }

    return nil
}
