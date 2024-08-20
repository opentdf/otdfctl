#!/usr/bin/env bats

# Tests for KAS grants

setup() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
  export HOST='--host http://localhost:8080'
}

@test "assign grant to namespace then unassign it" {
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[0].id')
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --namespace-id '$NS_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Namespace ID'
    assert_output --partial '$NS_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --namespace-id '$NS_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Namespace ID'
    assert_output --partial '$NS_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign grant to attribute then unassign it" {
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[0].id')
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --attribute-id '$ATTR_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Attribute ID'
    assert_output --partial '$ATTR_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --attribute-id '$ATTR_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Attribute ID'
    assert_output --partial '$ATTR_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign grant to value then unassign it" {
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[1].id')
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Value ID'
    assert_output --partial '$VAL_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Value ID'
    assert_output --partial '$VAL_ID'
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign rejects more than one type of grant at once" {
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values.[0].id')
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[1].id')

    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --attribute-id '$ATTR_ID' --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --namespace-id '$NS_ID' --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign --attribute-id '$ATTR_ID' --namespace-id '$NS_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'
}

@test "unassign rejects more than one type of grant at once" {
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values.[0].id')
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[1].id')

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --attribute-id '$ATTR_ID' --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --namespace-id '$NS_ID' --value-id '$VAL_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign --attribute-id '$ATTR_ID' --namespace-id '$NS_ID' --kas-id $KAS_ID
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'
}