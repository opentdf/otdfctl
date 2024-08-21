#!/usr/bin/env bats

# Tests for KAS grants

setup() {
      if [[ "$BATS_TEST_NUMBER" -eq 1 ]]; then
        export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry create --uri 'https://e2etestkas.com' --public-key-remote 'https://e2etestkas.com/pub_key' --json | jq -r '.id')
      else 
        export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[-1].id')
      fi
      
      echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
      export WITH_CREDS='--with-client-creds-file ./creds.json'
      export HOST='--host http://localhost:8080'

      export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
      export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
      export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values[0].id')

      export KAS_ID_FLAG="--kas-id $KAS_ID"
      export NS_ID_FLAG="--namespace-id $NS_ID"
      export ATTR_ID_FLAG="--attribute-id $ATTR_ID"
      export VAL_ID_FLAG="--value-id $VAL_ID"
}

@test "assign grant to namespace then unassign it" {
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Namespace ID'
    assert_output --partial $NS_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Namespace ID'
    assert_output --partial $NS_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign grant to attribute then unassign it" {
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Attribute ID'
    assert_output --partial $ATTR_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Attribute ID'
    assert_output --partial $ATTR_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign grant to value then unassign it" {
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Value ID'
    assert_output --partial $VAL_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'SUCCESS'
    assert_output --partial 'Value ID'
    assert_output --partial $VAL_ID
    assert_output --partial 'KAS ID'
    assert_output --partial $KAS_ID
}

@test "assign rejects more than one type of grant at once" {
    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'
}

@test "unassign rejects more than one type of grant at once" {
    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'

    ./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG
    assert_output --partial 'ERROR'
    assert_output --partial 'Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign'
}