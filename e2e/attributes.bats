#!/usr/bin/env bats
load "./helpers.bash"

# Tests for attributes

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
     export HOST="${HOST:---host http://localhost:8080}"

    # Create the namespace to be used by other tests

    export NS_NAME="testing-attr.co"
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
}

# always create a randomly named attribute
setup() {
    setup_helper

    # invoke binary with credentials
    run_otdfctl_attr () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy attributes $*"
    }

    export ATTR_NAME_RANDOM=$(LC_ALL=C tr -dc 'a-zA-Z' < /dev/urandom | head -c 16)
    export ATTR_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes create --namespace "$NS_ID" --name "$ATTR_NAME_RANDOM" --rule ANY_OF -l key=value --json | jq -r '.id')
}

# always unsafely delete the created attribute
teardown() {
    ./otdfctl $HOST $WITH_CREDS policy attributes unsafe delete --force --id "$ATTR_ID"

    cleanup_helper
}

teardown_file() {
  # remove the namespace
  ./otdfctl $HOST $WITH_CREDS policy attributes namespaces unsafe delete --id "$NS_ID" --force

  # clear out all test env vars
  unset HOST WITH_CREDS NS_NAME NS_ID ATTR_NAME_RANDOM
}

@test "Create an attribute - With Values" {
    run_otdfctl_attr create --name attrWithValues --namespace "$NS_ID" --rule HIERARCHY -v val1 -v val2 --json
      assert_success 
      [ "$( echo "$output" | jq -r '.values[0].value' )" = "val1" ]
      [ "$( echo "$output" | jq -r '.values[1].value' )" = "val2" ]
}

@test "Create an attribute - Bad" {
  # bad rule
    run_otdfctl_attr create --name attr1 --namespace "$NS_ID" --rule NONEXISTENT
      assert_failure 
      assert_output --partial "invalid attribute rule: NONEXISTENT, must be one of [ALL_OF, ANY_OF, HIERARCHY]"

  # missing flags
    run_otdfctl_attr create --name attr1 --rule ALL_OF
      assert_failure
      run_otdfctl_attr create --name attr1 --namespace "$NS_ID"
      assert_failure
      run_otdfctl_attr create --rule HIERARCHY --namespace "$NS_ID"
      assert_failure
}

@test "Get an attribute definition - Good" {
  LOWERED=$(echo "$ATTR_NAME_RANDOM" | awk '{print tolower($0)}')

   run_otdfctl_attr get --id "$ATTR_ID"
     assert_success
     assert_line --regexp "Id.*$ATTR_ID"
     assert_line --regexp "Name.*$LOWERED"
     assert_output --partial "ANY_OF"
     assert_line --regexp "Namespace.*$NS_NAME"

  run_otdfctl_attr get --id "$ATTR_ID" --json
    assert_success
    [ "$(echo "$output" | jq -r '.id')" = "$ATTR_ID" ]
    [ "$(echo "$output" | jq -r '.name')" = "$LOWERED" ]
    [ "$(echo "$output" | jq -r '.rule')" = 2 ]
    [ "$(echo "$output" | jq -r '.namespace.id')" = "$NS_ID" ]
    [ "$(echo "$output" | jq -r '.namespace.name')" = "$NS_NAME" ]
    [ "$(echo "$output" | jq -r '.metadata.labels.key')" = "value" ]
}

@test "Get an attribute definition - Bad" {
  # no id flag
   run_otdfctl_attr get
    assert_failure  
}

@test "Update an attribute definition (Safe) - Good" {
  # replace labels
    run_otdfctl_attr update --force-replace-labels -l key=somethingElse --id "$ATTR_ID" --json
    assert_success
    [ "$(echo $output | jq -r '.metadata.labels.key')" = "somethingElse" ]

  # extend labels
    run_otdfctl_attr update -l other=testing  --id "$ATTR_ID" --json
    assert_success
    [ "$(echo $output | jq -r '.metadata.labels.other')" = "testing" ]
    [ "$(echo $output | jq -r '.metadata.labels.key')" = "somethingElse" ]
}

@test "Update an attribute definition (Safe) - Bad" {
  # no id
  run_otdfctl_attr update
  assert_failure
}

@test "List attribute definitions" {
  run_otdfctl_attr list
  assert_success
  assert_output --partial "$ATTR_ID"
  assert_output --partial "Total"
  assert_line --regexp "Current Offset.*0"

  run_otdfctl_attr list --state active
  assert_success
  assert_output --partial "$ATTR_ID"
  assert_output --partial "Total"
  assert_line --regexp "Current Offset.*0"

  run_otdfctl_attr list --state inactive
  assert_success
  refute_output --partial "$ATTR_ID"
  assert_output --partial "Total"
  assert_line --regexp "Current Offset.*0"
}

@test "List - comprehensive pagination tests" {
  # create 10 random attributes so we have confidence there are >= 10 attribute definitions
  for i in {1..10}; do
    random_name=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 12)
    run_otdfctl_attr create --name "$random_name" --namespace "$NS_ID" --rule ANY_OF
    assert_success
  done

  run_otdfctl_attr list --limit 2
    assert_success
    assert_line --regexp "Current Offset.*0"
    assert_line --regexp "Next Offset.*2"
  
  run_otdfctl_attr list --limit 5 --offset 2
    assert_success
    assert_line --regexp "Current Offset.*2"
    assert_line --regexp "Next Offset.*7"

  run_otdfctl_attr list --offset 2
    assert_success
    assert_line --regexp "Current Offset.*2"

  run_otdfctl_attr list --limit 500
    assert_success
    refute_output --partial "Next Offset"
}

@test "Deactivate then unsafe reactivate an attribute definition" {
  run_otdfctl_attr deactivate
  assert_failure

  run_otdfctl_attr get --id "$ATTR_ID" --json
  assert_success
  [ "$(echo "$output" | jq -r '.active.value')" = true ]

  run_otdfctl_attr deactivate --id "$ATTR_ID" --force
  assert_success

  run_otdfctl_attr get --id "$ATTR_ID" --json
  assert_success
  [ "$(echo "$output" | jq -r '.active')" = {} ]

  run_otdfctl_attr unsafe reactivate
  assert_failure

  run_otdfctl_attr unsafe reactivate --id "$ATTR_ID" --force
  assert_success

  run_otdfctl_attr get --id "$ATTR_ID" --json
  assert_success
  [ "$(echo "$output" | jq -r '.active.value')" = true ]
}

@test "Unsafe Update an attribute definition" {
  # create with two values
  run_otdfctl_attr create --name created --namespace "$NS_ID" --rule HIERARCHY -v val1 -v val2 --json
    CREATED_ID=$(echo "$output" | jq -r '.id')
    VAL1_ID=$(echo "$output" | jq -r '.values[0].id')
    VAL2_ID=$(echo "$output" | jq -r '.values[1].id')

  run_otdfctl_attr unsafe update --name updated --id "$CREATED_ID" --json --force
    assert_success
  run_otdfctl_attr get --id "$CREATED_ID" --json
    assert_success
    [ "$(echo "$output" | jq -r '.name')" = "updated" ]

  run_otdfctl_attr unsafe update --rule ALL_OF --id "$CREATED_ID" --json --force
    assert_success
  run_otdfctl_attr get --id "$CREATED_ID" --json
    assert_success
    [ "$(echo "$output" | jq -r '.rule')" = 1 ]

  run_otdfctl_attr unsafe update --id "$CREATED_ID" --json --values-order "$VAL2_ID" --values-order "$VAL1_ID" --force
    assert_success
  run_otdfctl_attr get --id "$CREATED_ID" --json
    assert_success
    [ "$(echo "$output" | jq -r '.values[0].value')" = "val2" ]
    [ "$(echo "$output" | jq -r '.values[1].value')" = "val1" ]
}

@test "add_remove_key_to_definition" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    ALG="rsa:2048"
    KID="test"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Add the key to the attribute definition
    log_info "Running ${run_otdfctl_attr} keys add --definition $ATTR_ID --public-key-id $KID"
    run_otdfctl_attr keys add --definition "$ATTR_ID" --public-key-id "$PUBLIC_KEY_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    # Check that the key was added to the attribute definition
    log_info "Running ${run_otdfctl_attr} get --id $ATTR_ID"
    run_otdfctl_attr get --id "$ATTR_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    echo "$output" | jq -r '.keys[].id' | while read -r id; do
        log_debug "Checking PK ID: $id against $PUBLIC_KEY_ID"
        [ "$id" = "$PUBLIC_KEY_ID" ] || fail "KAS ID does not match"
    done

    # Remove the key from the attribute definition
    log_info "Running ${run_otdfctl_attr} keys remove --definition $ATTR_ID --public-key-id $PUBLIC_KEY_ID"
    run_otdfctl_attr keys remove --definition "$ATTR_ID" --public-key-id "$PUBLIC_KEY_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    # Check that the key was removed from the attribute definition
    log_info "Running ${run_otdfctl_attr} get --id $ATTR_ID"
    run_otdfctl_attr get --id "$ATTR_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    echo "$output" | jq -e 'has("keys") | not' || fail "KAS ID still present"
}

@test "add_remove_key_to_value" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    ALG="rsa:2048"
    KID="test"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Add value to the attribute definition
    log_info "Running ${run_otdfctl_attr} values create --attribute-id $ATTR_ID --value val1"
    run_otdfctl_attr values create --attribute-id "$ATTR_ID" --value val1 --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    VALUE_ID=$(echo "$output" | jq -r '.id')

    # Add the key to the attribute value
    log_info "Running ${run_otdfctl_attr} values keys add --value $VALUE_ID --public-key-id $KID"
    run_otdfctl_attr values keys add --value "$VALUE_ID" --public-key-id "$PUBLIC_KEY_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    # Check that the key was added to the attribute value
    log_info "Running ${run_otdfctl_attr} values get --id $VALUE_ID"
    run_otdfctl_attr values get --id "$VALUE_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    echo "$output" | jq -r '.keys[].id' | while read -r id; do
        log_debug "Checking PK ID: $id against $PUBLIC_KEY_ID"
        [ "$id" = "$PUBLIC_KEY_ID" ] || fail "KAS ID does not match"
    done

    # Remove the key from the attribute value
    log_info "Running ${run_otdfctl_attr} keys remove --value $VALUE_ID --public-key-id $PUBLIC_KEY_ID"
    run_otdfctl_attr values keys remove --value "$VALUE_ID" --public-key-id "$PUBLIC_KEY_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    # Check that the key was removed from the attribute value
    log_info "Running ${run_otdfctl_attr} values get --id $VALUE_ID"
    run_otdfctl_attr values get --id "$VALUE_ID" --json

    log_debug "Raw output:"
    log_debug "$output"

    assert_success

    echo "$output" | jq -e 'has("keys") | not' || fail "KAS ID still present"
}