#!/usr/bin/env bats

# Tests for namespaces

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create the namespace to be used by other tests

    export NS_NAME="creating-test-ns.net"
    export NS_NAME_UPDATE="updated-test-ns.net"
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
    export NS_ID_FLAG="--id $NS_ID"
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_ns () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy attributes namespaces $*"
    }
}

teardown_file() {
  # clear out all test env vars
  unset HOST WITH_CREDS NS_NAME NS_FQN NS_ID NS_ID_FLAG
}

@test "Create a namespace - Good" {
  run_otdfctl_ns create --name throwaway.test
  assert_output --partial "SUCCESS"
  assert_line --regexp "Name.*throwaway.test"
  assert_output --partial "Id"
  assert_output --partial "Created At"
  assert_line --partial "Updated At"

  # cleanup
  created_id=$(echo "$output" | grep Id | awk -F'│' '{print $3}' | xargs)
  run_otdfctl_ns unsafe delete --id $created_id --force
}

@test "Create a namespace - Bad" {
  # bad namespace names
    run_otdfctl_ns create --name no_domain_extension
    assert_failure
    run_otdfctl_ns create --name -first-char-hyphen.co
    assert_failure
    run_otdfctl_ns create --name last-char-hyphen-.co
    assert_failure

  # missing flag
    run_otdfctl_ns create
    assert_failure
    assert_output --partial "Flag '--name' is required"
  
  # conflict
    run_otdfctl_ns create -n "$NS_NAME"
    assert_failure
    assert_output --partial "already_exists"
}

@test "Get a namespace - Good" {
  run_otdfctl_ns get "$NS_ID_FLAG"
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  assert_line --regexp "Name.*$NS_NAME"

  run_otdfctl_ns get "$NS_ID_FLAG" --json
  assert_success
  [ "$(echo "$output" | jq -r '.id')" = "$NS_ID" ]
  [ "$(echo "$output" | jq -r '.name')" = "$NS_NAME" ]
}

@test "Get a namespace - Bad" {
  run_otdfctl_ns get
  assert_failure
  assert_output --partial "Flag '--id' is required"

  run_otdfctl_ns get --id 'example.com'
  assert_failure
  assert_output --partial "Flag '--id' received value 'example.com' must be a valid UUID"

  run_otdfctl_ns get --id 'demo.com' --json
  assert_failure
  assert_output --partial "Flag '--id' received value 'demo.com' must be a valid UUID"
}

@test "List namespaces - when active" {
  run_otdfctl_ns list --json 
  echo $output | jq --arg id "$NS_ID" '.[] | select(.[]? | type == "object" and .id == $id)'

  run_otdfctl_ns list --state inactive --json
  refute_output --partial "$NS_ID"

  run_otdfctl_ns list --state active
  assert_output --partial "$NS_ID"
  assert_output --partial "Total"
  assert_line --regexp "Current Offset.*0"
  
}

@test "Update namespace - Safe" {
  # extend labels
  run_otdfctl_ns update "$NS_ID_FLAG" -l key=value --label test=true
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  assert_line --regexp "Name.*$NS_NAME"
  assert_line --regexp "Labels.*key: value"
  assert_line --regexp "Labels.*test: true"

  # force replace labels
  run_otdfctl_ns update "$NS_ID_FLAG" -l key=other --force-replace-labels
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  assert_line --regexp "Name.*$NS_NAME"
  assert_line --regexp "Labels.*key: other"
  refute_output --regexp "Labels.*key: value"
  refute_output --regexp "Labels.*test: true"
}

@test "Update namespace - Unsafe" {
  run_otdfctl_ns unsafe update "$NS_ID_FLAG" -n "$NS_NAME_UPDATE" --force
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  run_otdfctl_ns get "$NS_ID_FLAG"
  assert_line --regexp "Name.*$NS_NAME_UPDATE"
  refute_output --regexp "Name.*$NS_NAME"
}

@test "Deactivate namespace" {
  run_otdfctl_ns deactivate "$NS_ID_FLAG" --force
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  assert_line --regexp "Name.*$NS_NAME_UPDATE"
}

@test "List namespaces - when inactive" {
  run_otdfctl_ns list --json 
  echo $output | jq --arg id "$NS_ID" '.[] | select(.[]? | type == "object" and .id == $id)'

  # json
    run_otdfctl_ns list --state inactive --json
    echo $output | assert_output --partial "$NS_ID"

    run_otdfctl_ns list --state active --json
    echo $output | refute_output --partial "$NS_ID"
  # table
    run_otdfctl_ns list --state inactive
    echo $output | assert_output --partial "$NS_ID"

    run_otdfctl_ns list --state active
    echo $output | refute_output --partial "$NS_ID"
}

@test "Unsafe reactivate namespace" {
  run_otdfctl_ns unsafe reactivate "$NS_ID_FLAG" --force
  assert_success
  assert_line --regexp "Id.*$NS_ID"
}

@test "List namespaces - when reactivated" {
  run_otdfctl_ns list --json 
  echo $output | jq --arg id "$NS_ID" '.[] | select(.[]? | type == "object" and .id == $id)'

  run_otdfctl_ns list --state inactive --json
  echo $output | refute_output --partial "$NS_ID"

  run_otdfctl_ns list --state active
  echo $output | assert_output --partial "$NS_ID"
}

@test "Unsafe delete namespace" {
  run_otdfctl_ns unsafe delete "$NS_ID_FLAG" --force
  assert_success
  assert_line --regexp "Id.*$NS_ID"
  assert_line --regexp "Name.*$NS_NAME_UPDATE"
}

@test "List namespaces - when deleted" {
  run_otdfctl_ns list --json 
  echo $output | refute_output --partial "$NS_ID"

  run_otdfctl_ns list --state inactive --json
  echo $output | refute_output --partial "$NS_ID"

  run_otdfctl_ns list --state active
  echo $output | refute_output --partial "$NS_ID"
}
