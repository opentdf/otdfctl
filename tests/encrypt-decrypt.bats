#!/usr/bin/env bats

# Tests for encrypt decrypt

setup_creds_json() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
}

@test "roundtrip TDF3" {
  setup_creds_json
  ./otdfctl encrypt -o sensitive.yaml.tdf --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 otdfctl.yaml
  [ "$status" -eq 0 ]
  ./otdfctl decrypt -o result.yaml --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 sensitive.yaml.tdf
  [ "$status" -eq 0 ]
  diff otdfctl.yaml result.yaml
  [ "$status" -eq 0 ]
}

@test "roundtrip NANO" {
  setup_creds_json
  ./otdfctl encrypt -o sensitive.yaml.tdf --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano otdfctl.yaml
  [ "$status" -eq 0 ]
  ./otdfctl decrypt -o result.yaml --host http://localhost:8080/kas --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano sensitive.yaml.tdf
  [ "$status" -eq 0 ]
  diff otdfctl.yaml result.yaml
  [ "$status" -eq 0 ]
}

# Future Tests

# Encrypt and decrypt with attributes:
# Create an attribute and a subject mapping for the specific clientId then roundtrip trip w it