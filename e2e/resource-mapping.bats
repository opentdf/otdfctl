#!/usr/bin/env bats

# Tests for resource mappings

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create two namespaced values to be used in other tests
        NS_NAME="resource-mappings.io"
        export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
        ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes create --namespace "$NS_ID" --name attr1 --rule ANY_OF --json | jq -r '.id')
        export VAL1_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$ATTR_ID" --value val1 --json | jq -r '.id')
        export VAL2_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$ATTR_ID" --value val2 --json | jq -r '.id')

    # Create a single resource mapping to val1 - comma separated
        export RM1_TERMS="valueone,valuefirst,first,one"
        export RM1_ID=$(./otdfctl $HOST $WITH_CREDS policy resource-mappings create --attribute-value-id "$VAL1_ID" --terms "$RM1_TERMS" --json | jq -r '.id')
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_rm () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy resource-mappings $*"
    }

}

teardown_file() {
    # remove the created namespace with all underneath upon test suite completion
    ./otdfctl $HOST $WITH_CREDS policy attributes namespaces unsafe delete --force --id "$NS_ID"

    unset HOST WITH_CREDS VAL1_ID VAL2_ID NS_ID RM1_TERMS RM1_ID
}

@test "Create resource mapping" {
    # create with multiple terms flags instead of comma-separated
    run_otdfctl_rm create --attribute-value-id "$VAL2_ID" --terms "second" --terms "TWO"
    assert_success
    assert_output --partial "second"
    assert_output --partial "TWO"
    assert_output --regexp "Attribute Value Id.*$VAL2_ID"

    # value id flag must be uuid
    run_otdfctl_rm create --attribute-value-id "val2" --terms "testing"
    assert_failure
    assert_output --partial "must be a valid UUID"

    # terms are required
    run_otdfctl_rm create --attribute-value-id $VAL2_ID
    assert_failure
    assert_output --partial "must have at least 1 non-empty values"
}

@test "Get resource mapping" {
    spaced_terms=$(echo $RM1_TERMS | sed 's/,/, /g')
    # table
    run_otdfctl_rm get --id "$RM1_ID"
        assert_success
        assert_output --regexp "Id.*$RM1_ID"
        assert_output --regexp "Attribute Value Id.*$VAL1_ID"
        assert_output --regexp "Terms.*$spaced_terms"
    
    # json
    run_otdfctl_rm get --id "$RM1_ID" --json
        assert_success
        [ $(echo $output | jq -r '.id') = "$RM1_ID" ]
        [ $(echo $output | jq -r '.attribute_value.id') = "$VAL1_ID" ]
        [ $(echo $output | jq -r '.terms | join (",")') = "$RM1_TERMS" ]
    
    # id required
    run_otdfctl_rm get
        assert_failure
        assert_output --partial "is required"
    run_otdfctl_rm get --id "test"
        assert_failure
        assert_output --partial "must be a valid UUID"
}

@test "Update a resource mapping" {
    NEW_RM_ID=$(./otdfctl $HOST $WITH_CREDS policy resource-mappings create --attribute-value-id "$VAL2_ID" --terms test --terms found --json | jq -r '.id')
    
    # replace the terms
    run_otdfctl_rm update --id "$NEW_RM_ID" --terms replaced,new 
        assert_success
        refute_output --partial "test"
        refute_output --partial "found"
        assert_output --partial "replaced"
        assert_output --partial "new"
        assert_output --partial "$VAL2_ID"

    # reassign the attribute value being mapped
    run_otdfctl_rm update --id "$NEW_RM_ID" --attribute-value-id "$VAL1_ID"
        assert_success
        refute_output --partial "test"
        refute_output --partial "found"
        assert_output --partial "replaced"
        assert_output --partial "new"
        refute_output --partial "$VAL2_ID"
        assert_output --partial "$VAL1_ID"
}

@test "List resource mappings" {
    run_otdfctl_rm list
        assert_success
        assert_output --partial "$RM1_ID"
        assert_output --partial "$VAL1_ID"
        assert_output --partial "valueone, valuefirst, first"
}

@test "Delete resource mapping" {
    spaced_terms=$(echo $RM1_TERMS | sed 's/,/, /g')
    # --force to avoid indefinite hang waiting for confirmation
    run_otdfctl_rm delete --id "$RM1_ID" --force
        assert_success
        assert_output --regexp "Id.*$RM1_ID"
        assert_output --regexp "Attribute Value Id.*$VAL1_ID"
        assert_output --regexp "Terms.*$spaced_terms"
}