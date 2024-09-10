#!/usr/bin/env bats

# Tests for attributes

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create the namespace to be used by other tests

    export NS_NAME="creating-test-ns.net"
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

# Create attribute

# Get Attribute

# Update attribute

# List attributes

# Deactivate Attribute

# Unsafe Reactivate

# Unsafe Delete

# Cleanup -- delete everything created here