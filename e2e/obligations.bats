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
    run_otdfctl_obl_triggers () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy obligations triggers $*"
    }

    run_otdfctl_action () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy actions $*"
    }

    run_otdfctl_attr() {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy attributes $*"
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

@test "Create an obligation value - Bad" {
  # bad obligation value names
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value ends_underscored_
    assert_failure
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value -first-char-hyphen
    assert_failure
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value inval!d.chars
    assert_failure

  # missing flag
  run_otdfctl_obl_values create
    assert_failure
    assert_output --partial "Flag '--obligation' is required"
  run_otdfctl_obl_values create --obligation "$OBL_ID"
    assert_failure
    assert_output --partial "Flag '--value' is required"

  # non-existent obligation fqn
  run_otdfctl_obl_values create --obligation invalid_fqn --value test_create_obl_val
    assert_failure
    assert_output --partial "obligation_fqn: value must be a valid URI [string.uri]"
  
  # conflict
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_create_obl_val_conflict
    assert_output --partial "SUCCESS"
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_create_obl_val_conflict
      assert_failure
      assert_output --partial "already_exists"

  # cleanup
  run_otdfctl_obl_values delete --id $created_id --force
}

@test "Get an obligation value - Good" {
  # setup an obligation value to get
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_get_obl_val
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # get by id
  run_otdfctl_obl_values get --id "$created_id" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.value')" = "test_get_obl_val" ]

  # get by fqn
  run_otdfctl_obl_values get --fqn "https://$NS_NAME/obl/$OBL_NAME/value/test_get_obl_val" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$created_id" ]
    [ "$(echo "$output" | jq -r '.value')" = "test_get_obl_val" ]

  # cleanup
  run_otdfctl_obl_values delete --id $created_id --force
}

@test "Get an obligation value - Bad" {
  run_otdfctl_obl_values get
    assert_failure
    assert_output --partial "one of id, fqn must be set [message.oneof]"

  # invalid id
  run_otdfctl_obl_values get --id 'not_a_uuid'
    assert_failure
    assert_output --partial "must be a valid UUID"

  # invalid fqn
  run_otdfctl_obl_values get --fqn 'not_a_fqn'
    assert_failure
    assert_output --partial "must be a valid URI"
}

@test "Update obligation values" {
  # setup an obligation value to update
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value test_update_obl_val
    assert_success
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)

  # force replace labels
  run_otdfctl_obl_values update --id "$created_id" -l key=other --force-replace-labels
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Value.*test_update_obl_val"
    assert_line --regexp "Labels.*key: other"
    refute_output --regexp "Labels.*key: value"
    refute_output --regexp "Labels.*test: true"
    refute_output --regexp "Labels.*test: true"

  # renamed
  run_otdfctl_obl_values update --id "$created_id" --value test_renamed_obl_val
    assert_success
    assert_line --regexp "Id.*$created_id"
    assert_line --regexp "Value.*test_renamed_obl_val"
    refute_output --regexp "Value.*test_update_obl_val"

  # cleanup
  run_otdfctl_obl_values delete --id $created_id --force
}

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

# Tests for obligation triggers

@test "Create an obligation trigger - Required Only - IDs - Success" {
  # setup an attribute value to use
  run_otdfctl_attr create --name "test_attr_for_trigger" --namespace "$NS_ID" --rule "HIERARCHY" -v "test_val_for_trigger" --json
  attr_val_id=$(echo "$output" | jq -r '.values[0].id')
  attr_id=$(echo "$output" | jq -r '.id')

  # setup an action to use
  run_otdfctl_action create --name "test_action_for_trigger" --json
  action_id=$(echo "$output" | jq -r '.id')

  # setup an obligation value to use
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value "test_obl_val_for_trigger" --json
  obl_val_id=$(echo "$output" | jq -r '.id')

  # create trigger
  run_otdfctl_obl_triggers create --attribute-value "$attr_val_id" --action "$action_id" --obligation-value "$obl_val_id" --json
  assert_success
  [ "$(echo "$output" | jq -r '.id')" != "null" ]
  trigger_id=$(echo "$output" | jq -r '.id')
  assert_equal "$(echo "$output" | jq -r '.attribute_value.id')" "$attr_val_id"
  assert_equal "$(echo "$output" | jq -r '.attribute_value.fqn')" "https://$NS_NAME/attr/test_attr_for_trigger/value/test_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.action.id')" "$action_id"
  assert_equal "$(echo "$output" | jq -r '.action.name')" "test_action_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.id')" "$obl_val_id"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.value')" "test_obl_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.id')" "$OBL_ID"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.namespace.fqn')" "https://$NS_NAME"


  # cleanup
  run_otdfctl_obl_triggers delete --id "$trigger_id" --force
  assert_success
  run_otdfctl_obl_values delete --id "$obl_val_id" --force
  assert_success
  run_otdfctl_action delete --id "$action_id" --force
  assert_success
  run_otdfctl_attr unsafe delete --id "$attr_id" --force
  assert_success
}

@test "Create an obligation trigger - Required Only - FQNs - Success" {
  # setup an attribute value to use
  run_otdfctl_attr create --name "test_attr_for_trigger" --namespace "$NS_ID" --rule "HIERARCHY" -v "test_val_for_trigger" --json
  attr_val_id=$(echo "$output" | jq -r '.values[0].id')
  attr_id=$(echo "$output" | jq -r '.id')
  attr_val_fqn=$(echo "$output" | jq -r '.values[0].fqn')

  # setup an action to use
  action_name="test_action_for_trigger"
  run_otdfctl_action create --name "$action_name" --json
  action_id=$(echo "$output" | jq -r '.id')

  # setup an obligation value to use
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value "test_obl_val_for_trigger" --json
  obl_val_id=$(echo "$output" | jq -r '.id')
  obl_val_fqn="https://$NS_NAME/obl/$OBL_NAME/value/test_obl_val_for_trigger"

  # create trigger
  run_otdfctl_obl_triggers create --attribute-value "$attr_val_fqn" --action "$action_name" --obligation-value "$obl_val_fqn" --json
  assert_success
  [ "$(echo "$output" | jq -r '.id')" != "null" ]
  trigger_id=$(echo "$output" | jq -r '.id')
  assert_equal "$(echo "$output" | jq -r '.attribute_value.id')" "$attr_val_id"
  assert_equal "$(echo "$output" | jq -r '.attribute_value.fqn')" "https://$NS_NAME/attr/test_attr_for_trigger/value/test_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.action.id')" "$action_id"
  assert_equal "$(echo "$output" | jq -r '.action.name')" "test_action_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.id')" "$obl_val_id"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.value')" "test_obl_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.id')" "$OBL_ID"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.namespace.fqn')" "https://$NS_NAME"


  # cleanup
  run_otdfctl_obl_triggers delete --id "$trigger_id" --force
  assert_success
  run_otdfctl_obl_values delete --id "$obl_val_id" --force
  assert_success
  run_otdfctl_action delete --id "$action_id" --force
  assert_success
  run_otdfctl_attr unsafe delete --id "$attr_id" --force
  assert_success
}

@test "Create an obligation trigger - Optional Fields - Success" {
  # setup an attribute value to use
  run_otdfctl_attr create --name "test_attr_for_trigger" --namespace "$NS_ID" --rule "HIERARCHY" -v "test_val_for_trigger" --json
  attr_val_id=$(echo "$output" | jq -r '.values[0].id')
  attr_id=$(echo "$output" | jq -r '.id')

  # setup an action to use
  run_otdfctl_action create --name "test_action_for_trigger" --json
  action_id=$(echo "$output" | jq -r '.id')

  # setup an obligation value to use
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value "test_obl_val_for_trigger" --json
  obl_val_id=$(echo "$output" | jq -r '.id')

  # create trigger
  client_id="a-pep"
  run_otdfctl_obl_triggers create --attribute-value "$attr_val_id" --action "$action_id" --obligation-value "$obl_val_id" --client-id "$client_id" --label "my=label" --json
  assert_success
  assert_not_equal "$(echo "$output" | jq -r '.id')" "null"
  trigger_id=$(echo "$output" | jq -r '.id')
  assert_equal "$(echo "$output" | jq -r '.attribute_value.id')" "$attr_val_id"
  assert_equal "$(echo "$output" | jq -r '.attribute_value.fqn')" "https://$NS_NAME/attr/test_attr_for_trigger/value/test_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.action.id')" "$action_id"
  assert_equal "$(echo "$output" | jq -r '.action.name')" "test_action_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.id')" "$obl_val_id"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.value')" "test_obl_val_for_trigger"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.id')" "$OBL_ID"
  assert_equal "$(echo "$output" | jq -r '.obligation_value.obligation.namespace.fqn')" "https://$NS_NAME"
  assert_equal "$(echo "$output" | jq -r '.metadata.labels.my')" "label"



  # cleanup
  run_otdfctl_obl_triggers delete --id "$trigger_id" --force
  assert_success
  run_otdfctl_obl_values delete --id "$obl_val_id" --force
  assert_success
  run_otdfctl_action delete --id "$action_id" --force
  assert_success
  run_otdfctl_attr unsafe delete --id "$attr_id" --force
  assert_success
}

@test "Create an obligation trigger - Bad" {
  # missing flags
  run_otdfctl_obl_triggers create --attribute-value "http://example.com/attr/attr_name/value/attr_value" --action "read" 
  assert_failure 
  assert_output --partial "Flag '--obligation-value' is required"
  
  run_otdfctl_obl_triggers create --obligation-value "http://example.com/attr/attr_name/value/attr_value" --action "read"
  assert_failure
  assert_output --partial "Flag '--attribute-value' is required"

  run_otdfctl_obl_triggers create --obligation-value "http://example.com/attr/attr_name/value/attr_value" --attribute-value "http://example.com/attr/attr_name/value/attr_value"
  assert_failure
  assert_output --partial "Flag '--action' is required"
}

@test "Delete an obligation trigger - Good" {
  # setup an attribute value to use
  run_otdfctl_attr create --name "test_attr_for_del_trigger" --namespace "$NS_ID" --rule "HIERARCHY" -v "test_val_for_del_trigger" --json
  assert_success
  attr_val_id=$(echo "$output" | jq -r '.values[0].id')
  attr_id=$(echo "$output" | jq -r '.id')

  # setup an action to use
  run_otdfctl_action create --name "test_action_for_del_trigger" --json
  assert_success
  action_id=$(echo "$output" | jq -r '.id')

  # setup an obligation value to use
  run_otdfctl_obl_values create --obligation "$OBL_ID" --value "test_obl_val_for_del_trigger" --json
  assert_success
  obl_val_id=$(echo "$output" | jq -r '.id')

  # create trigger
  run_otdfctl_obl_triggers create --attribute-value "$attr_val_id" --action "$action_id" --obligation-value "$obl_val_id" --json
  assert_success
  assert_not_equal "$(echo "$output" | jq -r '.id')" "null"
  trigger_id=$(echo "$output" | jq -r '.id')

  # delete trigger
  run_otdfctl_obl_triggers delete --id "$trigger_id" --force --json
  assert_success
  assert_equal "$(echo "$output" | jq -r '.id')" "$trigger_id"

  # cleanup
  run_otdfctl_obl_values delete --id "$obl_val_id" --force
  assert_success
  run_otdfctl_action delete --id "$action_id" --force
  assert_success
  run_otdfctl_attr unsafe delete --id "$attr_id" --force
}