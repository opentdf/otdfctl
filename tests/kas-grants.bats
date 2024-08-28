#!/usr/bin/env bats

# Tests for KAS grants

setup_file() {
      echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
      export WITH_CREDS='--with-client-creds-file ./creds.json'
      export HOST='--host http://localhost:8080'

      export KAS_URI="https://e2etestkas.com"
      export KAS_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry create --uri "$KAS_URI" --public-key-remote 'https://e2etestkas.com/pub_key' --json | jq -r '.id')
      export KAS_ID_FLAG="--kas-id $KAS_ID"
}

@test "namespace: assign grant then unassign it" {
    # assign the namespace a grant
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export NS_ID_FLAG="--namespace-id $NS_ID"
    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $KAS_ID_FLAG)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Namespace ID"* ]]
    [[ "$result" == *$NS_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should find the namespace in the grants
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants | any(.[]?; .id == $id))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants | any(.[]?; .id == $id))')"
      # unfiltered
        # json
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants? | type == "array" and any(.[]?; .id == $id))')"
        # table
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
        echo $result | grep -E "Namespace.*$NS_ID"

    # unassign the namespace grant
    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $KAS_ID_FLAG --force)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Namespace ID"* ]]
    [[ "$result" == *$NS_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should not find the namespace within any grants to namespaces
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq 'map(select(has("namespace_grants") | not))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq 'map(select(has("namespace_grants") | not))')"
      # unfiltered
        # json
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants? | type == "array" and all(.[]?; .id != $id))')"
        # table
        # result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
        # echo $result | grep -E "Namespace.*$NS_ID"
}

@test "attribute: assign grant then unassign it" {
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export ATTR_ID_FLAG="--attribute-id $ATTR_ID"
    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $KAS_ID_FLAG)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Attribute ID"* ]]
    [[ "$result" == *$ATTR_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should find the attribute in the grants
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants | any(.[]?; .id == $id))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants | any(.[]?; .id == $id))')"
      # unfiltered
        # json
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants? | type == "array" and any(.[]?; .id == $id))')"
        # table
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
        echo $result | grep -E "Definition.*$ATTR_ID"


    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $KAS_ID_FLAG --force)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Attribute ID"* ]]
    [[ "$result" == *$ATTR_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should not find the attribute within any grants to attributes
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq 'map(select(has("attribute_grants") | not))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq 'map(select(has("attribute_grants") | not))')"
      # unfiltered
        # json
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants? | type == "array" and all(.[]?; .id != $id))')"
        # table
        # result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
        # echo $result | grep -qE "Definition.*$ATTR_ID"
}

@test "value: assign grant then unassign it" {
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values[0].id')
    export VAL_ID_FLAG="--value-id $VAL_ID"
    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $VAL_ID_FLAG $KAS_ID_FLAG)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Value ID"* ]]
    [[ "$result" == *$VAL_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should find the value in the grants
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants | any(.[]?; .id == $id))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants | any(.[]?; .id == $id))')"
      # unfiltered
        # json
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants? | type == "array" and any(.[]?; .id == $id))')"
        # table
        result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
        echo $result | grep -E "Value.*$VAL_ID"

    result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $VAL_ID_FLAG $KAS_ID_FLAG --force)"
    [[ "$result" == *"SUCCESS"* ]]
    [[ "$result" == *"Value ID"* ]]
    [[ "$result" == *$VAL_ID* ]]
    [[ "$result" == *"KAS ID"* ]]
    [[ "$result" == *$KAS_ID* ]]

    # LIST should not find the value within any grants to values
      # filtered by KAS
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_ID --json | jq 'map(select(has("value_grants") | not))')"
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --kas $KAS_URI --json | jq 'map(select(has("value_grants") | not))')"
      # unfiltered
      # json
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants? | type == "array" and all(.[]?; .id != $id))')"
      # table
      # result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants list)"
      # echo $result | grep -qE "Value.*$VAL_ID"
}

@test "assign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'
    
    # simulates try/catch to avoid failed tests on expected errors
    result=''
    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"* ]]

    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"* ]]

    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants assign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"* ]]
}

@test "unassign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'

    # simulates try/catch to avoid failed tests on expected errors
    result=''
    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"* ]]

    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"* ]]

    {
      result="$(./otdfctl $HOST $WITH_CREDS policy kas-grants unassign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG)"
    } || {
      true
    }
      [[ "$result" == *"Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"* ]]
}