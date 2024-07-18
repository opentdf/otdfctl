#!/usr/bin/env bats

# Tests for encrypt decrypt

setup() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
}

@test "roundtrip TDF3" {
  ./otdfctl encrypt -o sensitive.yaml.tdf --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 otdfctl.yaml
  ./otdfctl decrypt -o result.yaml --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 sensitive.yaml.tdf
  diff otdfctl.yaml result.yaml
}

@test "roundtrip NANO" {
  ./otdfctl encrypt -o sensitive.yaml.tdf --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano otdfctl.yaml
  ./otdfctl decrypt -o result.yaml --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano sensitive.yaml.tdf
  diff otdfctl.yaml result.yaml
}

# Future Tests

# Encrypt and decrypt with attributes:
# Create an attribute and a subject mapping for the specific clientId then roundtrip trip w it