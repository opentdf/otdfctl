#!/usr/bin/env bats

# Tests for registered resources

setup_file() {
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
    export WITH_CREDS='--with-client-creds-file ./creds.json'
    export HOST='--host http://localhost:8080'

    # Create the namespace to be used by other tests

    export NS_NAME="testing-reg-res.co"
    export NS_ID=$(./otdfctl $HOST $WITH_CREDS policy attributes namespaces create -n "$NS_NAME" --json | jq -r '.id')
}

# always create a randomly named attribute
setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_reg_res () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy registered-resources $*"
    }
    run_otdfctl_reg_res_values () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy registered-resources values $*"
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

@test "Create a registered resource" {
  run_otdfctl_reg_res create --name testRegRes --json
    assert_success
    [ "$( echo "$output" | jq -r '.name' )" = "testRegRes" ]
}
