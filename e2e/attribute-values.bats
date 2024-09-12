#!/usr/bin/env bats

# Tests for attribute values

setup_file() {
    export CREDSFILE=creds.json
    echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
    export WITH_CREDS="--with-client-creds-file $CREDSFILE"
    export HOST='--host http://localhost:8080'
    export DEBUG_LEVEL="--log-level debug"
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_av () {
      run sh -c "./otdfctl $HOST $WITH_CREDS policy attributes values $*"
    }
}