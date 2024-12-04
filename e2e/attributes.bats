#!/usr/bin/env bats

# Tests for attributes

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create the namespace to be used by other tests

    export NS_NAME="testing-attr.co"
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
}

# always create a randomly named attribute
setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

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
     assert_output --regexp "Id.*$ATTR_ID"
     assert_output --regexp "Name.*$LOWERED"
     assert_output --partial "ANY_OF"
     assert_output --regexp "Namespace.*$NS_NAME"

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
  assert_output --regexp "Current Offset       0"

  run_otdfctl_attr list --state active
  assert_success
  assert_output --partial "$ATTR_ID"
  assert_output --partial "Total"
  assert_output --regexp "Current Offset       0"

  run_otdfctl_attr list --state inactive
  assert_success
  refute_output --partial "$ATTR_ID"
  assert_output --partial "Total"
  assert_output --regexp "Current Offset       0"
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