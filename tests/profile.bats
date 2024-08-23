#!/usr/bin/env bats

setup() {
  bats_require_minimum_version 1.5.0
  
  OTDFCTL_BIN=./otdfctl_testbuild

  # Check if BATS_SUPPORT_PATH environment variable exists
  if [ -z "${BATS_SUPPORT_PATH}" ]; then
    # Check if bats is installed via homebrew
    if [[ $(which bats) == *"/homebrew/"* ]]; then
      FINAL_BATS_SUPPORT_PATH=$(brew --prefix)/lib
    else
      FINAL_BATS_SUPPORT_PATH=$(dirname $(which bats))/../lib
    fi
  else
    FINAL_BATS_SUPPORT_PATH="${BATS_SUPPORT_PATH}"
  fi
  echo "FINAL_BATS_SUPPORT_PATH: $FINAL_BATS_SUPPORT_PATH"
  load "${FINAL_BATS_SUPPORT_PATH}/bats-support/load.bash"
  load "${FINAL_BATS_SUPPORT_PATH}/bats-assert/load.bash"

  set_test_profile() {
    auth=""
    # if 3rd argument is empty, then don't include it
    if [ -n "$3" ]; then
      auth=",\"auth\":$3"
    fi
    echo "{\"profile\":\"$1\",\"endpoint\":\"$2\"$auth}"
  }

  set_test_profile_auth() {
    authType=$1
    clientId=$2
    clientSecret=$3
    accessToken=$4
    echo "{\"authType\":\"$authType\",\"clientId\":\"$clientId\",\"clientSecret\":\"$clientSecret\",\"accessToken\":\"$accessToken\"}"
  }

  set_test_profile_auth_access_token() {
    publicClientId=$1
    accessToken=$2
    refreshToken=$3
    expiration=$4
    echo "{\"publicClientId\":\"$publicClientId\",\"accessToken\":\"$accessToken\",\"refreshToken\":\"$refreshToken\",\"expiration\":\"$expiration\"}"
  }

  set_test_config() {
    defaultProfile=$1
    shift 1
    profiles=""
    for i in "$@"; do
      # if first profile just set it
      if [ -z "$profiles" ]; then
        profiles="$i"
      else
        profiles="$profiles,$i"
      fi
    done
    export OTDFCTL_TEST_PROFILE="{\"defaultProfile\":\"$defaultProfile\",\"profiles\":[$profiles]}"
  }

  run_otdfctl() {
    run sh -c "./$OTDFCTL_BIN $*"
  }

  assert_no_profile_set() {
    assert_output --partial "No default profile set"
  }

  # Set the keyring provider to in-memory
  export OTDFCTL_KEYRING_PROVIDER="in-memory"
}

teardown() {
  unset OTDFCTL_KEYRING_PROVIDER
}

@test "profile create" {
  run_otdfctl profile create test http://localhost:8080
  assert_output --regexp "Creating profile .* ok"

  run_otdfctl profile create test localhost:8080
  assert_output --regexp "Failed .* invalid scheme"

  # TODO figure out how to test the case where the profile already exists
}

@test "profile list" {
  run_otdfctl profile list
  assert_no_profile_set

  # export OTDFCTL_TEST_CONFIG='{"defaultProfile":"test","profiles":[{"profile": "test","endpoint":"http://localhost:8080"}]}'
  set_test_config "test2" $(set_test_profile "test" "http://localhost:8080") $(set_test_profile "test2" "http://localhost:8081")
  run_otdfctl profile list
  assert_line --index 5 --regexp "test$"
  assert_line --index 6 --regexp "\* test2$"
}

@test "profile get" {
  run_otdfctl profile get test
  assert_no_profile_set

  set_test_config "test2" $(set_test_profile "test" "http://localhost:8080") $(set_test_profile "test2" "http://localhost:8081")
  run_otdfctl profile get test
  assert_line --index 8 --regexp "Profile\s+|\s*test\s*"
  assert_output --regexp "Endpoint\s+|\s*http://localhost:8080"
  assert_output --regexp "default\s+|\s*false"
  # TODO check auth
}

@test "profile delete" {
  run_otdfctl profile delete test
  assert_no_profile_set

  # TODO test deleting the default profile

  set_test_config "test2" $(set_test_profile "test" "http://localhost:8080") $(set_test_profile "test2" "http://localhost:8081")
  run_otdfctl profile delete test
  assert_output --partial "Deleting profile test... ok"
}

@test "profile set-default" {
  run_otdfctl profile set-default test
  assert_no_profile_set

  set_test_config "test2" $(set_test_profile "test" "http://localhost:8080") $(set_test_profile "test2" "http://localhost:8081")
  run_otdfctl profile set-default test
  assert_output --partial "Setting profile test as default... ok"
}

@test "profile set-endpoint" {
  run_otdfctl profile set-endpoint test http://localhost:8081
  assert_no_profile_set

  set_test_config "test2" $(set_test_profile "test" "http://localhost:8080") $(set_test_profile "test2" "http://localhost:8081")
  run_otdfctl profile set-endpoint test http://localhost:8081
  assert_output --partial "Setting endpoint for profile test... ok"
}