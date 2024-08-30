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

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_kasg () {
      run sh -c "./otdfctl $HOST $WITH_CREDS $*"
    }
}

teardown_file() {
  # clear out all test env vars
  unset HOST WITH_CREDS KAS_ID KAS_ID_FLAG KAS_URI NS_ID NS_ID_FLAG ATTR_ID ATTR_ID_FLAG VAL_ID VAL_ID_FLAG
}

@test "namespace: assign grant then unassign it" {
    # assign the namespace a grant
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces list --json | jq -r '.[0].id')
    export NS_ID_FLAG="--namespace-id $NS_ID"

    run_otdfctl_kasg assign "$NS_ID_FLAG" "$KAS_ID_FLAG"
      assert_output --partial "SUCCESS"
      assert_output --partial "Namespace ID"
      assert_output --partial $NS_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should find the namespace in the grants
      # filtered by KAS
        # json
          run_otdfctl_kasg list --kas $KAS_ID --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants | any(.[]?; .id == $id))'
          assert_success
          run_otdfctl_kasg list --kas $KAS_URI --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants | any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list --kas $KAS_ID
          assert_output --regexp "$KAS_URI.*Namespace.*$NS_ID"
          run_otdfctl_kasg list --kas $KAS_URI
          assert_output --regexp "$KAS_URI.*Namespace.*$NS_ID"


      # unfiltered (all KASes)
        # json
          run_otdfctl_kasg list --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants? | type == "array" and any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          assert_output --regexp "$KAS_URI.*Namespace.*$NS_ID"

    # unassign the namespace grant
    run_otdfctl_kasg unassign $NS_ID_FLAG $KAS_ID_FLAG --force
      assert_output --partial "SUCCESS"
      assert_output --partial "Namespace ID"
      assert_output --partial $NS_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should not find the namespace within any grants to namespaces
      # filtered by KAS
        # json
          run_otdfctl_kasg list --kas $KAS_ID --json | jq 'map(select(has("namespace_grants") | not))'
          assert_success
          run_otdfctl_kasg list --kas $KAS_URI --json | jq 'map(select(has("namespace_grants") | not))'
          assert_success
        # table
          run_otdfctl_kasg list
          refute_output --regexp "$KAS_URI.*Namespace.*$NS_ID"
      # unfiltered
        # json
          run_otdfctl_kasg list --json | jq --arg id "$NS_ID" '.[] | select(.namespace_grants? | type == "array" and all(.[]?; .id != $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          refute_output --regexp "$KAS_URI.*Namespace.*$NS_ID"
}

@test "attribute: assign grant then unassign it" {
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].id')
    export ATTR_ID_FLAG="--attribute-id $ATTR_ID"
    run_otdfctl_kasg assign "$ATTR_ID_FLAG" "$KAS_ID_FLAG"
      assert_output --partial "SUCCESS"
      assert_output --partial "Attribute ID"
      assert_output --partial $ATTR_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should find the attribute in the grants
      # filtered by KAS
        # json
          run_otdfctl_kasg list --kas $KAS_ID --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants | any(.[]?; .id == $id))'
          assert_success
          run_otdfctl_kasg list --kas $KAS_URI --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants | any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list --kas $KAS_URI
          assert_output --regexp "$KAS_URI.*Definition.*$ATTR_ID"
          run_otdfctl_kasg list --kas $KAS_ID
          assert_output --regexp "$KAS_URI.*Definition.*$ATTR_ID"
      # unfiltered
        # json
          run_otdfctl_kasg list --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants? | type == "array" and any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          assert_output --regexp "$KAS_URI.*Definition.*$ATTR_ID"

    run_otdfctl_kasg unassign $ATTR_ID_FLAG $KAS_ID_FLAG --force
      assert_output --partial "SUCCESS"
      assert_output --partial "Attribute ID"
      assert_output --partial $ATTR_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should not find the attribute within any grants to attributes
      # filtered by KAS
        run_otdfctl_kasg list --kas $KAS_ID --json | jq 'map(select(has("attribute_grants") | not))'
        assert_success
        run_otdfctl_kasg list --kas $KAS_URI --json | jq 'map(select(has("attribute_grants") | not))'
        assert_success
      # unfiltered
        # json
          run_otdfctl_kasg list --json | jq --arg id "$ATTR_ID" '.[] | select(.attribute_grants? | type == "array" and all(.[]?; .id != $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          refute_output --regexp "$KAS_URI.*Definition.*$ATTR_ID"
}

@test "value: assign grant then unassign it" {
    export VAL_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes list --json | jq -r '.[0].values[0].id')
    export VAL_ID_FLAG="--value-id $VAL_ID"
    run_otdfctl_kasg assign "$VAL_ID_FLAG" "$KAS_ID_FLAG"
      assert_output --partial "SUCCESS"
      assert_output --partial "Value ID"
      assert_output --partial $VAL_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should find the value in the grants
      # filtered by KAS
        # json
          run_otdfctl_kasg list --kas $KAS_ID --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants | any(.[]?; .id == $id))'
          assert_success
          run_otdfctl_kasg list --kas $KAS_URI --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants | any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list --kas $KAS_ID
          assert_output --regexp "$KAS_URI.*Value.*$VAL_ID"
          run_otdfctl_kasg list --kas $KAS_URI
          assert_output --regexp "$KAS_URI.*Value.*$VAL_ID"

      # unfiltered
        # json
          run_otdfctl_kasg list --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants? | type == "array" and any(.[]?; .id == $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          assert_output --regexp "$KAS_URI.*Value.*$VAL_ID"

    run_otdfctl_kasg unassign $VAL_ID_FLAG $KAS_ID_FLAG --force
      assert_output --partial "SUCCESS"
      assert_output --partial "Value ID"
      assert_output --partial $VAL_ID
      assert_output --partial "KAS ID"
      assert_output --partial $KAS_ID

    # LIST should not find the value within any grants to values
      # filtered by KAS
        # json
          run_otdfctl_kasg list --kas $KAS_ID --json | jq 'map(select(has("value_grants") | not))'
          assert_success
          run_otdfctl_kasg list --kas $KAS_URI --json | jq 'map(select(has("value_grants") | not))'
          assert_success
        # table
          run_otdfctl_kasg list --kas $KAS_ID
          refute_output --regexp "$KAS_URI.*Value.*$VAL_ID"
          run_otdfctl_kasg list --kas $KAS_URI
          refute_output --regexp "$KAS_URI.*Value.*$VAL_ID"
      # unfiltered
        # json
          run_otdfctl_kasg list --json | jq --arg id "$VAL_ID" '.[] | select(.value_grants? | type == "array" and all(.[]?; .id != $id))'
          assert_success
        # table
          run_otdfctl_kasg list
          refute_output --regexp "$KAS_URI.*Value.*$VAL_ID"
    }

@test "assign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'
    
    run_otdfctl_kasg assign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"

    run_otdfctl_kasg assign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"

    run_otdfctl_kasg assign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to assign"
}

@test "unassign rejects more than one type of grant at once" {
    export NS_ID_FLAG='--namespace-id hello'
    export ATTR_ID_FLAG='--attribute-id world'
    export VAL_ID_FLAG='--value-id goodnight'

    run_otdfctl_kasg unassign $ATTR_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"
    
    run_otdfctl_kasg unassign $NS_ID_FLAG $VAL_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"
    
    run_otdfctl_kasg unassign $ATTR_ID_FLAG $NS_ID_FLAG $KAS_ID_FLAG
      assert_failure
      assert_output --partial "Must specify exactly one Attribute Namespace ID, Definition ID, or Value ID to unassign"
}