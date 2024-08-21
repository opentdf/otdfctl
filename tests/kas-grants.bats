#!/usr/bin/env bats

# Tests for KAS grants

setup() {
      echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
      export WITH_CREDS='--with-client-creds-file ./creds.json'
      export HOST='--host http://localhost:8080'

      if [[ "$BATS_TEST_NUMBER" -eq 1 ]]; then
        export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry create --uri 'https://e2etestkas.com' --public-key-remote 'https://e2etestkas.com/pub_key' --json | jq -r '.id')
      else 
        export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry list --json | jq -r '.[-1].id')
      fi

      export KAS_ID_FLAG="--kas-id $KAS_ID"
}

@test "assign grant to namespace then unassign it" {
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export NS_ID_FLAG="--namespace-id $NS_ID"
    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Namespace ID" ]]
    [[ "$result" =~ $NS_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Namespace ID" ]]
    [[ "$result" =~ $NS_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]
}

@test "assign grant to attribute then unassign it" {
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export ATTR_ID_FLAG="--attribute-id $ATTR_ID"
    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Attribute ID" ]]
    [[ "$result" =~ $ATTR_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Attribute ID" ]]
    [[ "$result" =~ $ATTR_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]
}

@test "assign grant to value then unassign it" {
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values[0].id')
    export VAL_ID_FLAG="--value-id $VAL_ID"
    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Value ID" ]]
    [[ "$result" =~ $VAL_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "SUCCESS" ]]
    [[ "$result" =~ "Value ID" ]]
    [[ "$result" =~ $VAL_ID ]]
    [[ "$result" =~ "KAS ID" ]]
    [[ "$result" =~ $KAS_ID ]]
}

@test "assign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'
    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]
}

@test "unassign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'
    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]

    result=$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG)
    [[ "$result" =~ "ERROR" ]]
    [[ "$result" =~ "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign" ]]
}