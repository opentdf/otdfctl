#!/usr/bin/env bats

setup_file() {
  # Use the test build of the CLI
  OTDFCTL_BIN=./otdfctl_testbuild # Let's use the testbuild for the keyring.
  export OTDFCTL_BIN

  # Prefix for all profiles created in this file to avoid clashing
  PROFILE_TEST_PREFIX="bats-profile-$(date +%s)"
  export PROFILE_TEST_PREFIX
}

setup() {
  load "${BATS_LIB_PATH}/bats-support/load.bash"
  load "${BATS_LIB_PATH}/bats-assert/load.bash"

  run_otdfctl() {
    run sh -c "./otdfctl $*"
  }
}

teardown() {
  run_otdfctl profile delete-all --force
}

@test "profile create" {
  profile="${PROFILE_TEST_PREFIX}-create"
  run_otdfctl profile create "$profile" http://localhost:8080
  assert_success
  assert_output --partial "Creating profile ${profile}..."
  assert_output --partial "ok"

  # Invalid endpoint should fail with a helpful message
  run_otdfctl profile create "$profile" localhost:8080
  assert_failure
  assert_output --partial "Failed to create profile"
  assert_output --partial "invalid scheme"
}

@test "profile list shows profiles and default" {
  profile1="${PROFILE_TEST_PREFIX}-list-1"
  profile2="${PROFILE_TEST_PREFIX}-list-2"

  run_otdfctl profile create "$profile1" http://localhost:8080
  assert_success

  run_otdfctl profile create "$profile2" http://localhost:8080 --set-default
  assert_success

  run_otdfctl profile list
  assert_success
  assert_output --partial "Listing profiles from filesystem"
  assert_output --partial "  ${profile1}"
  assert_output --partial "* ${profile2}"
}

@test "profile get shows profile details" {
  profile="${PROFILE_TEST_PREFIX}-get"

  run_otdfctl profile create "$profile" http://localhost:8080
  assert_success

  run_otdfctl profile get "$profile"
  assert_success
  assert_output --partial "Profile"
  assert_output --partial "$profile"
  assert_output --partial "Endpoint"
  assert_output --partial "http://localhost:8080"
  assert_output --partial "Is default"
  assert_output --partial "true"
}

@test "profile delete removes profile" {
  base="${PROFILE_TEST_PREFIX}-delete"
  default_profile="${base}-default"
  target_profile="${base}-target"

  run_otdfctl profile create "$default_profile" http://localhost:8080 --set-default
  assert_success

  run_otdfctl profile create "$target_profile" http://localhost:8080
  assert_success

  run_otdfctl profile delete "$target_profile"
  assert_success
  assert_output --partial "Deleting profile ${target_profile}, from filesystem..."
  assert_output --partial "ok"

  run_otdfctl profile list
  assert_success
  refute_output --partial "$target_profile"
}

@test "profile set-default updates default profile" {
  base="${PROFILE_TEST_PREFIX}-set-default"
  profile1="${base}-1"
  profile2="${base}-2"

  run_otdfctl profile create "$profile1" http://localhost:8080 --set-default
  assert_success

  run_otdfctl profile create "$profile2" http://localhost:8081
  assert_success

  run_otdfctl profile set-default "$profile2"
  assert_success
  assert_output --partial "Setting profile ${profile2} as default..."
  assert_output --partial "ok"

  run_otdfctl profile list
  assert_success
  assert_output --partial "* ${profile2}"
}

@test "profile set-endpoint updates endpoint" {
  profile="${PROFILE_TEST_PREFIX}-set-endpoint"

  run_otdfctl profile create "$profile" http://localhost:8080
  assert_success

  run_otdfctl profile set-endpoint "$profile" http://localhost:8081
  assert_success
  assert_output --partial "Setting endpoint for profile ${profile}... "
  assert_output --partial "ok"

  run_otdfctl profile get "$profile"
  assert_success
  assert_output --partial "http://localhost:8081"
}

@test "profile delete-all deletes all profiles" {
  base="${PROFILE_TEST_PREFIX}-delete-all"
  profile1="${base}-1"
  profile2="${base}-2"

  run_otdfctl profile create "$profile1" http://localhost:8080 --set-default
  assert_success

  run_otdfctl profile create "$profile2" http://localhost:8081
  assert_success

  run_otdfctl profile delete-all --force
  assert_success
  assert_output --partial "profiles from filesystem..."
  assert_output --partial "ok"

  run_otdfctl profile list
  assert_success
  refute_output --partial "$profile1"
  refute_output --partial "$profile2"
}
