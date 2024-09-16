#!/usr/bin/env bats

# Tests for subject mappings

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create two namespaced values to be used in other tests
        NS_NAME="subject-mappings.net"
        export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
        ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes create --namespace "$NS_ID" --name attr1 --rule ANY_OF --json | jq -r '.id')
        export VAL1_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$ATTR_ID" --value val1 --json | jq -r '.id')
        export VAL2_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$ATTR_ID" --value value2 --json | jq -r '.id')

    export SCS_1='[{"condition_groups":[{"conditions":[{"operator":1,"subject_external_values":["ShinyThing"],"subject_external_selector_value":".team.name"},{"operator":2,"subject_external_values":["marketing"],"subject_external_selector_value":".org.name"}],"boolean_operator":1}]}]'
    export SCS_2='[{"condition_groups":[{"conditions":[{"operator":2,"subject_external_values":["CoolTool","RadService"],"subject_external_selector_value":".team.name"},{"operator":1,"subject_external_values":["sales"],"subject_external_selector_value":".org.name"}],"boolean_operator":2}]}]'
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_sm () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy subject-mappings $*"
    }

}

teardown_file() {
    # remove the created namespace with all underneath upon test suite completion
    ./otdfctl $HOST $WITH_CREDS policy attributes namespaces unsafe delete --force --id "$NS_ID"

    unset HOST WITH_CREDS VAL1_ID VAL2_ID NS_ID SCS_1 SCS_2
}

@test "Create subject mapping" {
    # create with simultaneous new SCS
    run_otdfctl_sm create --attribute-value-id "$VAL1_ID" --action-standard DECRYPT -s TRANSMIT --subject-condition-set-new "$SCS_2"
    assert_success
    # assert_output --partial "Subject Condition Set: Id"
    # assert_output --partial '"Standard":1"
    # assert_output --partial '"Standard":2"
    # assert_output --regexp 'Subject Condition Set.*\.team\.name"
    # assert_output --regexp "Attribute Value Id.*$VAL1_ID"

    # scs is required
    run_otdfctl_sm create --attribute-value-id "$VAL2_ID" -s TRANSMIT
    assert_failure
    assert_output --partial "At least one Subject Condition Set flag [--subject-condition-set-id, --subject-condition-set-new] must be provided"

    # action is required
    run_otdfctl_sm create -a "$VAL1_ID" --subject-condition-set-new "$SCS_2"
    assert_failure
    assert_output --partial "At least one Standard or Custom Action [--action-standard, --action-custom] is required"
}

@test "Get subject mapping" {
    new_scs=$(./otdfctl $HOST $WITH_CREDS policy scs create -s "$SCS_2" --json | jq -r '.id')
    created=$(./otdfctl $HOST $WITH_CREDS policy sm create -a "$VAL2_ID" -s TRANSMIT --subject-condition-set-id "$new_scs" --json | jq -r '.id')
    # table
    run_otdfctl_sm get --id "$created"
        assert_success
        assert_output --regexp "Id.*$created"
        assert_output --regexp "Attribute Value: Id.*$VAL2_ID"
        assert_output --regexp "Attribute Value: Value.*value2"
        assert_output --regexp "Subject Condition Set: Id.*$created"

    # json
    run_otdfctl_sm get --id "$created" --json
        assert_success
        [ "$(echo $output | jq -r '.id')" = "$created" ]
        [ "$(echo $output | jq -r '.attribute_value.id')" = "$VAL2_ID" ]
        [ "$(echo $output | jq -r '.subject_condition_set.id')" = "$new_scs" ]
}

@test "Update a subject mapping" {
    created=$(./otdfctl $HOST $WITH_CREDS policy sm create -a "$VAL1_ID" -s DECRYPT --subject-condition-set-new "$SCS_1" --json | jq -r '.id')
    additional_scs=$(./otdfctl $HOST $WITH_CREDS policy scs create -s "$SCS_2" --json | jq -r '.id')

    # replace the action (always destructive replacement)
    run_otdfctl_sm update --id "$created" -s TRANSMIT --json
        assert_success
        [ "$(echo $output | jq -r '.id')" = "$created" ]
        [ "$(echo $output | jq -r '.actions[0].Value.Standard')" = 2 ]

    # reassign the SCS being mapped to
    run_otdfctl_sm update --id "$created" --subject-condition-set-id "$additional_scs" --json
        assert_success
        [ "$(echo $output | jq -r '.id')" = "$created" ]
        [ "$(echo $output | jq -r '.subject_condition_set.id')" = "$additional_scs" ]
}

@test "List subject mappings" {
    created=$(./otdfctl $HOST $WITH_CREDS policy sm create -a "$VAL1_ID" -s TRANSMIT --subject-condition-set-new "$SCS_2" --json | jq -r '.id')

    run_otdfctl_sm list
        assert_success
        assert_output --partial "$VAL2_ID"
        assert_output --partial "$VAL1_ID"
        assert_output --partial "$created"
}

@test "Delete subject mapping" {
    first_listed=$(./otdfctl $HOST $WITH_CREDS policy sm list --json | jq -r '.[0].id')
    # --force to avoid indefinite hang waiting for confirmation
    run_otdfctl_sm delete --id "$first_listed" --force
        assert_success
        assert_output --regexp "Id.*$first_listed"
}