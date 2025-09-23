#!/usr/bin/env bats

# Tests for obligations

setup_file() {
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # create attribute value to be used in obligation values tests
    export NS_NAME="test-obl.org"
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create --name "$NS_NAME" --json | jq -r '.id')
   
    # create obligation used in obligation values tests
    export OBL_NAME="test_obl_for_values"
    export OBL_ID=$(./otdfctl $HOST $WITH_CREDS policy obligations create --name "$OBL_NAME" --namespace "$NS_ID" --json | jq -r '.id')
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_obl () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy obligations $*"
    }
    run_otdfctl_obl_values () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy obligations values $*"
    }
}

teardown_file() {
  # remove the obligation used in obligation values tests
  ./otdfctl $HOST $WITH_CREDS policy obligations delete --id "$OBL_ID" --force

  # # remove the custom action used in obligation values tests
  # ./otdfctl $HOST $WITH_CREDS policy actions delete --id "$CUSTOM_ACTION_ID" --force

  # remove the namespace used in obligation values tests
  ./otdfctl $HOST $WITH_CREDS policy attributes namespaces unsafe delete --id "$NS_ID" --force

  # clear out all test env vars
  unset HOST WITH_CREDS OBL_NAME OBL_ID
}

@test "Create a obligation - Good" {
  run_otdfctl_obl create --name test_create_obl --namespace "$NS_ID"
    assert_output --partial "SUCCESS"
    assert_line --regexp "Name.*test_create_obl"
    assert_output --partial "Id"
    assert_output --partial "Created At"
    assert_line --partial "Updated At"

  # cleanup
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_obl delete --id $created_id --force
}

@test "Create a obligation - Bad" {
  # bad obligation names
  run_otdfctl_obl create --name ends_underscored_ --namespace "$NS_ID"
    assert_failure
  run_otdfctl_obl create --name -first-char-hyphen --namespace "$NS_ID"
    assert_failure
  run_otdfctl_obl create --name inval!d.chars --namespace "$NS_ID"
    assert_failure

  # missing flag
  run_otdfctl_obl create
    assert_failure
    assert_output --partial "Flag '--name' is required"
  
  # conflict
  run_otdfctl_obl create --name test_create_obl_conflict --namespace "$NS_ID"
    assert_output --partial "SUCCESS"
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_obl create --name test_create_obl_conflict --namespace "$NS_ID"
      assert_failure
      assert_output --partial "already_exists"

  # cleanup
  run_otdfctl_obl delete --id $created_id --force
}

@test "Get an obligation - Good" {
  # setup an obligation to get
  run_otdfctl_obl create --name test_get_obl --namespace "$NS_ID"
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # get by id
  run_otdfctl_obl get --id "$created_id" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.name')" = "test_get_obl" ]

  # get by fqn
  run_otdfctl_obl get --fqn "https://${NS_NAME}/obl/test_get_obl" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.name')" = "test_get_obl" ]

  # cleanup
  run_otdfctl_obl delete --id $created_id --force
}

@test "Get an obligation - Bad" {
  run_otdfctl_obl get
    assert_failure
    assert_output --partial "one of id, fqn must be set [message.oneof]"

  run_otdfctl_obl get --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"
}

@test "List obligations" {
  # setup obligations to list
  run_otdfctl_obl create --name test_list_obl_1 --namespace "$NS_ID"
  obl1_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_obl create --name test_list_obl_2 --namespace "$NS_ID"
  obl2_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  run_otdfctl_obl list
    assert_success
    assert_output --partial "$obl1_id"
    assert_output --partial "test_list_obl_1"
    assert_output --partial "$obl2_id"
    assert_output --partial "test_list_obl_2"
    assert_output --partial "Total"
    assert_line --regexp "Current Offset.*0"

  # cleanup
  run_otdfctl_obl delete --id $obl1_id --force
  run_otdfctl_obl delete --id $obl2_id --force
}

@test "Update obligation" {
  # setup an obligation to update
  run_otdfctl_obl create --name test_update_obl --namespace "$NS_ID"
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # force replace labels
  run_otdfctl_obl update --id "$created_id" -l key=other --force-replace-labels
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Name.*test_update_obl"
    assert_line --regexp "Labels.*key: other"
    refute_output --regexp "Labels.*key: value"
    refute_output --regexp "Labels.*test: true"
    refute_output --regexp "Labels.*test: true"

  # renamed
  run_otdfctl_obl update --id "$created_id" --name test_renamed_obl
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Name.*test_renamed_obl"
    refute_output --regexp "Name.*test_update_obl"

  # cleanup
  run_otdfctl_obl delete --id $created_id --force
}

@test "Delete obligation - Good" {
  # setup an obligation to delete
  run_otdfctl_obl create --name test_delete_obl --namespace "$NS_ID"
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  run_otdfctl_obl delete --id "$created_id" --force
    assert_success
}

@test "Delete obligation - Bad" {
  # no id
  run_otdfctl_obl delete
    assert_failure
    assert_output --partial "one of id, fqn must be set [message.oneof]"

  # invalid id
  run_otdfctl_obl delete --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"
}

# Tests for obligation values

@test "Create an obligation value - Good" {
  # simple by obligation ID
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_create_obl_val
    assert_output --partial "SUCCESS"
    assert_line --regexp "Value.*test_create_obl_val"
    assert_output --partial "Id"
    assert_output --partial "Created At"
    assert_line --partial "Updated At"
  created_id_simple=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # simple by obligation FQN
  run_otdfctl_obl_values create --obligation "https://$NS_NAME/obl/$OBL_NAME" --value test_create_obl_val_by_obl_fqn
    assert_output --partial "SUCCESS"
    assert_line --regexp "Value.*test_create_obl_val"
    assert_output --partial "Id"
    assert_output --partial "Created At"
    assert_line --partial "Updated At"
  created_id_simple_by_fqn=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # cleanup
  run_otdfctl_obl_values delete --id $created_id_simple --force
  run_otdfctl_obl_values delete --id $created_id_simple_by_fqn --force
}

# @test "Create a registered resource value - Bad" {
#   # bad resource value names
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value ends_underscored_
#     assert_failure
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value -first-char-hyphen
#     assert_failure
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value inval!d.chars
#     assert_failure

#   # missing flag
#   run_otdfctl_reg_res_values create
#     assert_failure
#     assert_output --partial "Flag '--resource' is required"
#   run_otdfctl_reg_res_values create --resource "$RR_ID"
#     assert_failure
#     assert_output --partial "Flag '--value' is required"

#   # bad action attribute value arg separator (not a semicolon)
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_create_rr_val_bad_aav --action-attribute-value "\"$READ_ACTION_ID:$ATTR_VAL_1_ID\""
#     assert_failure
#     assert_output --partial "Invalid action attribute value arg format"

#   # non-existent resource name
#   run_otdfctl_reg_res_values create --resource invalid_rr --value test_create_rr_val_bad_aav_action_name
#     assert_failure
#     assert_output --partial "Failed to find registered resource (name: invalid_rr)"
  
#   # conflict
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_create_rr_val_conflict
#     assert_output --partial "SUCCESS"
#   created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_create_rr_val_conflict
#       assert_failure
#       assert_output --partial "already_exists"

#   # cleanup
#   run_otdfctl_reg_res_values delete --id $created_id --force
# }

# @test "Get a registered resource value - Good" {
#   # setup a resource value to get
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_get_rr_val --action-attribute-value "\"$READ_ACTION_ID;$ATTR_VAL_1_ID\""
#     assert_success
#   created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

#   # get by id
#   run_otdfctl_reg_res_values get --id "$created_id" --json
#     assert_success
#     [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
#     [ "$(echo "$output" | jq -r '.value')" = "test_get_rr_val" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].action.id')" = "$READ_ACTION_ID" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].action.name')" = "$READ_ACTION_NAME" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].attribute_value.id')" = "$ATTR_VAL_1_ID" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].attribute_value.fqn')" = "$ATTR_VAL_1_FQN" ]

#   # get by fqn
#   run_otdfctl_reg_res_values get --fqn "https://reg_res/$RR_NAME/value/test_get_rr_val" --json
#     assert_success
#     [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
#     [ "$(echo "$output" | jq -r '.value')" = "test_get_rr_val" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].action.id')" = "$READ_ACTION_ID" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].action.name')" = "$READ_ACTION_NAME" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].attribute_value.id')" = "$ATTR_VAL_1_ID" ]
#     [ "$(echo "$output" | jq -r '.action_attribute_values[0].attribute_value.fqn')" = "$ATTR_VAL_1_FQN" ]

#   # cleanup
#   run_otdfctl_reg_res_values delete --id $created_id --force
# }

# @test "Get a registered resource value - Bad" {
#   run_otdfctl_reg_res_values get
#     assert_failure
#     assert_output --partial "Either 'id' or 'fqn' must be provided"

#   # invalud id
#   run_otdfctl_reg_res_values get --id 'not_a_uuid'
#     assert_failure
#     assert_output --partial "must be a valid UUID"

#   # invalid fqn
#   run_otdfctl_reg_res_values get --fqn 'not_a_fqn'
#     assert_failure
#     assert_output --partial "must be a valid URI"
# }

# @test "List registered resource values - Good" {
#   # setup values to list
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_list_rr_val_1 --action-attribute-value "\"$READ_ACTION_ID;$ATTR_VAL_1_ID\""
#   reg_res_val1_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_list_rr_val_2
#   reg_res_val2_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

#   # by resource ID
#   run_otdfctl_reg_res_values list --resource "$RR_ID"
#     assert_success
#     assert_output --partial "$reg_res_val1_id"
#     assert_output --partial "test_list_rr_val_1"
#     # check for partial FQN due to possible trimmed output
#     assert_output --partial "$READ_ACTION_NAME -> https://$NS_NAME/attr/$ATTR_NAME"
#     assert_output --partial "$reg_res_val2_id"
#     assert_output --partial "test_list_rr_val_2"
#     assert_output --partial "Total"
#     assert_line --regexp "Current Offset.*0"

#   # by resource name
#   run_otdfctl_reg_res_values list --resource "$RR_NAME"
#     assert_success
#     assert_output --partial "$reg_res_val1_id"
#     assert_output --partial "test_list_rr_val_1"
#     # check for partial FQN due to possible trimmed output
#     assert_output --partial "$READ_ACTION_NAME -> https://$NS_NAME/attr/$ATTR_NAME"
#     assert_output --partial "$reg_res_val2_id"
#     assert_output --partial "test_list_rr_val_2"
#     assert_output --partial "Total"
#     assert_line --regexp "Current Offset.*0"

#   # cleanup
#   run_otdfctl_reg_res_values delete --id $reg_res_val1_id --force
#   run_otdfctl_reg_res_values delete --id $reg_res_val2_id --force
# }

# @test "List registered resource values - Bad" {
#   # non-existent resource name
#   run_otdfctl_reg_res_values list --resource 'invalid_rr'
#     assert_failure
#     assert_output --partial "Failed to find registered resource (name: invalid_rr)"
# }

# @test "Update registered resource values" {
#   # setup a resource value to update
#   run_otdfctl_reg_res_values create --resource "$RR_ID" --value test_update_rr_val --action-attribute-value "\"$READ_ACTION_ID;$ATTR_VAL_1_ID\""
#     assert_success
#   created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

#   # force replace labels
#   run_otdfctl_reg_res_values update --id "$created_id" -l key=other --force-replace-labels
#     assert_success
#     assert_line --regexp "Id.*$created_id"
#     assert_line --regexp "Value.*test_update_rr_val"
#     assert_line --regexp "Labels.*key: other"
#     refute_output --regexp "Labels.*key: value"
#     refute_output --regexp "Labels.*test: true"
#     refute_output --regexp "Labels.*test: true"

#   # renamed
#   run_otdfctl_reg_res_values update --id "$created_id" --value test_renamed_rr_val
#     assert_success
#     assert_line --regexp "Id.*$created_id"
#     assert_line --regexp "Value.*test_renamed_rr_val"
#     refute_output --regexp "Value.*test_update_rr_val"

#   # ensure previous updates without action attribute value args did not clear action attribute values
#   run_otdfctl_reg_res_values get --id "$created_id" --json
#     [ "$(echo "$output" | jq -r 'any(.action_attribute_values[]; .action.id == "'"$READ_ACTION_ID"'" and .action.name == "'"$READ_ACTION_NAME"'" and .attribute_value.id == "'"$ATTR_VAL_1_ID"'" and .attribute_value.fqn == "'"$ATTR_VAL_1_FQN"'")')" = "true" ]

#   # update action attribute values
#   run_otdfctl_reg_res_values update --id "$created_id" --action-attribute-value "\"$READ_ACTION_NAME;$ATTR_VAL_1_FQN\"" --action-attribute-value "\"$CUSTOM_ACTION_ID;$ATTR_VAL_2_ID\"" --force --json
#     assert_success
#     [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
#     [ "$(echo "$output" | jq -r 'any(.action_attribute_values[]; .action.id == "'"$READ_ACTION_ID"'" and .action.name == "'"$READ_ACTION_NAME"'" and .attribute_value.id == "'"$ATTR_VAL_1_ID"'" and .attribute_value.fqn == "'"$ATTR_VAL_1_FQN"'")')" = "true" ]
#     [ "$(echo "$output" | jq -r 'any(.action_attribute_values[]; .action.id == "'"$CUSTOM_ACTION_ID"'" and .action.name == "'"$CUSTOM_ACTION_NAME"'" and .attribute_value.id == "'"$ATTR_VAL_2_ID"'" and .attribute_value.fqn == "'"$ATTR_VAL_2_FQN"'")')" = "true" ]

#   # cleanup
#   run_otdfctl_reg_res_values delete --id $created_id --force
# }

@test "Delete obligation value - Good" {
  # setup a value to delete
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_delete_obl_val
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  run_otdfctl_obl_values delete --id "$created_id" --force
    assert_success
}

@test "Delete obligation value - Bad" {
  # no id
  run_otdfctl_obl_values delete
    assert_failure
    assert_output --partial "one of id, fqn must be set [message.oneof]"

  # invalid id
  run_otdfctl_obl_values delete --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"
}