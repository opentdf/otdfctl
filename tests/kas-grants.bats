#!/usr/bin/env bats

# Tests for KAS grants

setup() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
  export HOST='--host http://localhost:8080'
}

@test "assign grant to namespace" {
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[0].id')
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --namespace-id $NS_ID --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial $NS_ID
    assert_output --partial $KAS_ID
}