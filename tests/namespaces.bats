#!/usr/bin/env bats

# Tests for namespaces

setup_file() {
      echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
      export WITH_CREDS='--with-client-creds-file ./creds.json'
      export HOST='--host http://localhost:8080'

      bats_require_minimum_version 1.5.0

      if [[ $(which bats) == *"homebrew"* ]]; then
          BATS_LIB_PATH=$(brew --prefix)/lib
        fi

      # Check if BATS_LIB_PATH environment variable exists
      if [ -z "${BATS_LIB_PATH}" ]; then
        # Check if bats bin has homebrew in path name
        if [[ $(which bats) == *"homebrew"* ]]; then
          BATS_LIB_PATH=$(dirname $(which bats))/../lib
        elif [ -d "/usr/lib/bats-support" ]; then
          BATS_LIB_PATH="/usr/lib"
        elif [ -d "/usr/local/lib/bats-support" ]; then
          # Check if bats-support exists in /usr/local/lib
          BATS_LIB_PATH="/usr/local/lib"
        fi
      fi
      echo "BATS_LIB_PATH: $BATS_LIB_PATH"
      export BATS_LIB_PATH=$BATS_LIB_PATH
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl () {
      run sh -c "./otdfctl $HOST $WITH_CREDS $*"
    }
}

teardown_file() {
  # clear out all test env vars
  unset HOST WITH_CREDS NS_NAME NS_FQN NS_ID NS_ID_FLAG
}

@test "Create a namespace - Good" {
  # Create the namespace to be used by other tests
  export NS_NAME="creating-test-ns.net"
  export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')

  run_otdfctl policy attributes namespaces create --name throwaway.test
  assert_output --partial "SUCCESS"
  assert_output --regexp "Name.*throwaway.test"
  assert_output --partial "Id"
  assert_output --partial "Created At"
  assert_output --regexp "Updated At"
}

@test "Create a namespace - Bad" {
  # bad namespace names
    run_otdfctl policy attributes namespaces create --name no_domain_extension
    assert_failure
    run_otdfctl policy attributes namespaces create --name -first-char-hyphen.co
    assert_failure
    run_otdfctl policy attributes namespaces create --name last-char-hyphen-.co
    assert_failure

  # missing flag
    run_otdfctl policy attributes namespaces create
    assert_failure
    assert_output --partial "Flag '--name' is required"
}

# Get namesapce

# Update namespace

# List namespaces

# Deactivate namespace

# Unsafe namespace

# Unsafe namespace

# Cleanup - delete everything