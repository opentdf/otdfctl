#!/usr/bin/env bats

# Tests for registered resources

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # create registered resource used in registered resource values tests
    export RR_ID=$(./otdfctl $HOST $WITH_CREDS policy registered-resources create --name test_rr_for_values --json | jq -r '.id')

    # create custom action to be used in registered resource values tests
    export CUSTOM_ACTION_NAME="test_action_for_values"
    export CUSTOM_ACTION_ID=$(./otdfctl $HOST $WITH_CREDS policy actions create --name "$CUSTOM_ACTION_NAME" --json | jq -r '.id')

    # create attribute value to be used in registered resource values tests
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create --name "test-reg-res.org" --json | jq -r '.id')
    attr_id=$(./otdfctl $HOST $WITH_CREDS policy attributes create --namespace "$NS_ID" --name test_reg_res_attr --rule ANY_OF -l key=value --json | jq -r '.id')
    export ATTR_VAL_1_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$attr_id" --value test_reg_res_attr__val_1 --json | jq -r '.id')
    export ATTR_VAL_2_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes values create --attribute-id "$attr_id" --value test_reg_res_attr__val_2 --json | jq -r '.id')
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_reg_res () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy registered-resources $*"
    }
    run_otdfctl_reg_res_values () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy registered-resources values $*"
    }
}

teardown_file() {
  # remove the registered resource used in registered resource values tests
  ./otdfctl $HOST $WITH_CREDS policy registered-resources delete --id "$RR_ID" --force

  # remove the custom action used in registered resource values tests
  ./otdfctl $HOST $WITH_CREDS policy actions delete --id "$CUSTOM_ACTION_ID" --force

  # remove the namespace and cascade delete attributes and values used in registered resource values tests
  ./otdfctl $HOST $WITH_CREDS policy attributes namespaces unsafe delete --id "$NS_ID" --force

  # clear out all test env vars
  unset HOST WITH_CREDS RR_ID CUSTOM_ACTION_NAME CUSTOM_ACTION_ID NS_ID ATTR_VAL_1_ID ATTR_VAL_2_ID
}

@test "Create a registered resource - Good" {
  run_otdfctl_reg_res create --name test_create_rr
  assert_output --partial "SUCCESS"
  assert_line --regexp "Name.*test_create_rr"
  assert_output --partial "Id"
  assert_output --partial "Created At"
  assert_line --partial "Updated At"

  # cleanup
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_reg_res delete --id $created_id --force
}

@test "Create a registered resource - Bad" {
  # bad resource names
  run_otdfctl_reg_res create --name ends_underscored_
    assert_failure
  run_otdfctl_reg_res create --name -first-char-hyphen
    assert_failure
  run_otdfctl_reg_res create --name inval!d.chars
    assert_failure

  # missing flag
  run_otdfctl_reg_res create
    assert_failure
    assert_output --partial "Flag '--name' is required"
  
  # conflict
  run_otdfctl_reg_res create --name test_create_rr_conflict
    assert_output --partial "SUCCESS"
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_reg_res create --name test_create_rr_conflict
      assert_failure
      assert_output --partial "AlreadyExists"

  # cleanup
  run_otdfctl_reg_res delete --id $created_id --force
}

@test "Get a registered resource - Good" {
  # setup a resource to get
  run_otdfctl_reg_res create --name test_get_rr
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # get by id
  run_otdfctl_reg_res get --id "$created_id" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.name')" = "test_get_rr" ]

  # get by name
  run_otdfctl_reg_res get --name test_get_rr --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.name')" = "test_get_rr" ]

  # cleanup
  run_otdfctl_reg_res delete --id $created_id --force
}

@test "Get a registered resource - Bad" {
  run_otdfctl_reg_res get
    assert_failure
    assert_output --partial "Either 'id' or 'name' must be provided"

  run_otdfctl_reg_res get --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"
}

@test "List registered resources" {
  # setup registered resources to list
  run_otdfctl_reg_res create --name test_list_rr_1
  reg_res1_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_reg_res create --name test_list_rr_2
  reg_res2_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  run_otdfctl_reg_res list
    assert_success
    assert_output --partial "$reg_res1_id"
    assert_output --partial "test_list_rr_1"
    assert_output --partial "$reg_res2_id"
    assert_output --partial "test_list_rr_2"
    assert_output --partial "Total"
    assert_line --regexp "Current Offset.*0"

  # cleanup
  run_otdfctl_reg_res delete --id $reg_res1_id --force
  run_otdfctl_reg_res delete --id $reg_res2_id --force
}

@test "Update registered resource" {
  # setup a resource to update
  run_otdfctl_reg_res create --name test_update_rr
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # force replace labels
  run_otdfctl_reg_res update --id "$created_id" -l key=other --force-replace-labels
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Name.*test_update_rr"
    assert_line --regexp "Labels.*key: other"
    refute_output --regexp "Labels.*key: value"
    refute_output --regexp "Labels.*test: true"
    refute_output --regexp "Labels.*test: true"

  # renamed
  run_otdfctl_reg_res update --id "$created_id" --name test_renamed_rr
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Name.*test_renamed_rr"
    refute_output --regexp "Name.*test_update_rr"

  # cleanup
  run_otdfctl_reg_res delete --id $created_id --force
}

@test "Delete registered resource - Good" {
  # setup a resource to delete
  run_otdfctl_reg_res create --name test_delete_rr
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  run_otdfctl_reg_res delete --id "$created_id" --force
    assert_success
}

@test "Delete registered resource - Bad" {
  # no id
  run_otdfctl_reg_res delete
    assert_failure
    assert_output --partial "Flag '--id' is required"

  # invalid id
  run_otdfctl_reg_res delete --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"
}
