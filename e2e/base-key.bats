#!/usr/bin/env bats

setup_file() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
  export HOST='--host http://localhost:8080'

  # Create a KAS registry entry for testing base keys
  export KAS_NAME_BASE_KEY_TEST="kas-registry-for-base-key-tests"
  export KAS_URI_BASE_KEY_TEST="https://test-kas-for-base-keys.com"
  export KAS_REGISTRY_ID_BASE_KEY_TEST=$(./otdfctl $HOST $WITH_CREDS policy kas-registry create --name "${KAS_NAME_BASE_KEY_TEST}" --uri "${KAS_URI_BASE_KEY_TEST}" --public-key-remote "${KAS_URI_BASE_KEY_TEST}" --json | jq -r '.id')

  # Create a regular KAS key to be set as a base key
  # This key will be used by the 'set' command tests
  export REGULAR_KEY_ID_FOR_BASE_TEST="regular-key-for-base-$(date +%s)"
  export WRAPPING_KEY="gp6TcYb/ZrgkQOYPdiYFRj11jZwbevy+r2KFbAYM0GE="
  export KAS_KEY_SYSTEM_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry key create --kasId "${KAS_REGISTRY_ID_BASE_KEY_TEST}" --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}" --alg rsa:2048 --mode local --wrappingKey "${WRAPPING_KEY}" --wrappingKeyId "wrapping-key-id" --json | jq -r '.key.id')
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials for base key commands
    run_otdfctl_base_key () {
      run sh -c "./otdfctl policy kas-registry key base $HOST $WITH_CREDS $*"
    }
}

teardown_file() {
  unset HOST WITH_CREDS KAS_REGISTRY_ID_BASE_KEY_TEST KAS_NAME_BASE_KEY_TEST KAS_URI_BASE_KEY_TEST REGULAR_KEY_ID_FOR_BASE_TEST WRAPPING_KEY KAS_KEY_SYSTEM_ID
  rm -f creds.json
}

# --- get base key tests ---

@test "base-key: get (initially no base key should be set for a new KAS)" {
  run_otdfctl_base_key get
  assert_failure # Expecting failure or specific message indicating no base key
  assert_output --partial "No base key found" # Or similar error message
}

# --- set base key tests ---

@test "base-key: set by --id" {
  run_otdfctl_base_key set --id "${KAS_KEY_SYSTEM_ID}" --json
  assert_success
  # Verify the new base key part of the response
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "rsa:2048"
  # Verify previous base key is null or not present if this is the first set
   assert_equal "$(echo "$output" | jq -r .previous_base_key)" "null"
}



@test "base-key: set by --keyId and --kasId" {
  run_otdfctl_base_key set --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}" --kasId "${KAS_REGISTRY_ID_BASE_KEY_TEST}" --json
  assert_success
  # Verify the new base key part of the response
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "rsa:2048"
}

@test "base-key: get (after setting a base key)" {
  run_otdfctl_base_key set --id "${KAS_KEY_SYSTEM_ID}" --json
  assert_success

  run_otdfctl_base_key get --json
  assert_success
  assert_equal "$(echo "$output" | jq -r .public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_equal "$(echo "$output" | jq -r .kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
  assert_not_equal "$(echo "$output" | jq -r .public_key.pem)" ""
  assert_not_equal "$(echo "$output" | jq -r .public_key.pem)" "null"
  assert_equal "$(echo "$output" | jq -r .public_key.algorithm)" "rsa:2048"
}

@test "base-key: set by --keyId and --kasName" {
  run_otdfctl_base_key set --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}" --kasName "${KAS_NAME_BASE_KEY_TEST}" --json
  assert_success
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}" # KAS URI should remain the same for the KAS Name
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "rsa:2048"
}

@test "base-key: set by --keyId and --kasUri" {
  # This will set REGULAR_KEY_ID_FOR_BASE_TEST back as the base key
  run_otdfctl_base_key set --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}" --kasUri "${KAS_URI_BASE_KEY_TEST}" --json
  assert_success
  # Verify the new base key
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}" # KAS URI should remain the same for the KAS Name
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
  assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
  assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "rsa:2048"
}

@test "base-key: set, get, and verify previous base key" {
    run_otdfctl_base_key set --id "${KAS_KEY_SYSTEM_ID}" --json
    assert_success
    assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
    assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
    assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "rsa:2048"


    run_otdfctl_base_key get --json
    assert_success
    assert_equal "$(echo "$output" | jq -r .public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
    assert_equal "$(echo "$output" | jq -r .kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
    assert_not_equal "$(echo "$output" | jq -r .public_key.pem)" ""
    assert_not_equal "$(echo "$output" | jq -r .public_key.pem)" "null"
    assert_equal "$(echo "$output" | jq -r .public_key.algorithm)" "rsa:2048"


    SECOND_KEY_ID_FOR_BASE_TEST="second-key-for-base-$(date +%s)"
    SECOND_KAS_KEY_SYSTEM_ID=$(./otdfctl $HOST $WITH_CREDS policy kas-registry key create --kasId "${KAS_REGISTRY_ID_BASE_KEY_TEST}" --keyId "${SECOND_KEY_ID_FOR_BASE_TEST}" --alg ec:secp256r1 --mode local --wrappingKey "${WRAPPING_KEY}" --wrappingKeyId "test-key" --json | jq -r '.key.id')

    run_otdfctl_base_key set --id "${SECOND_KAS_KEY_SYSTEM_ID}" --json
    assert_success
    assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.kid)" "${SECOND_KEY_ID_FOR_BASE_TEST}"
    assert_equal "$(echo "$output" | jq -r .new_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
    assert_equal "$(echo "$output" | jq -r .new_base_key.public_key.algorithm)" "ec:secp256r1"
    # Verify previous base key
    assert_equal "$(echo "$output" | jq -r .previous_base_key.public_key.kid)" "${REGULAR_KEY_ID_FOR_BASE_TEST}"
    assert_equal "$(echo "$output" | jq -r .previous_base_key.kas_uri)" "${KAS_URI_BASE_KEY_TEST}"
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" ""
    assert_not_equal "$(echo "$output" | jq -r .new_base_key.public_key.pem)" "null"
    assert_equal "$(echo "$output" | jq -r .previous_base_key.public_key.algorithm)" "rsa:2048"
}


@test "base-key: set (missing kas identifier: kasId, kasName, or kasUri)" {
  run_otdfctl_base_key set --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}"
  assert_failure
  assert_output --partial "at least one of 'kasId', 'kasName', or 'kasUri' must be provided"
}

@test "base-key: set (missing key identifier: id or keyId)" {
  run_otdfctl_base_key set --kasId "${KAS_REGISTRY_ID_BASE_KEY_TEST}"
  assert_failure
  assert_output --partial "Error: at least one of the flags in the group [id keyId] is required"
}

@test "base-key: set (using non-existent keyId)" {
  NON_EXISTENT_KEY_ID="this-key-does-not-exist-12345"
  run_otdfctl_base_key set --keyId "${NON_EXISTENT_KEY_ID}" --kasId "${KAS_REGISTRY_ID_BASE_KEY_TEST}"
  assert_failure
  # The exact error message might depend on the backend implementation
  assert_output --partial "not_found" # Or a more specific "key not found" error
}

@test "base-key: set (using non-existent kasId)" {
  NON_EXISTENT_KAS_ID="a1b2c3d4-e5f6-7890-1234-567890abcdef"
  run_otdfctl_base_key set --keyId "${REGULAR_KEY_ID_FOR_BASE_TEST}" --kasId "${NON_EXISTENT_KAS_ID}"
  assert_failure
  assert_output --partial "not_found" # Or a more specific "KAS not found" error
}

