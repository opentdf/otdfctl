#!/usr/bin/env bats

# Tests for kas registry

setup_file() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export DEBUG_LEVEL="--log-level debug"
  export HOST='--host http://localhost:8080'

  export REMOTE_KEY='https://hello.world/pubkey'
  export PEM='LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvVENDQWVXZ0F3SUJBZ0lVTXU4bzhXaDJIVEE2VEFlTENqQzJmVjlwSWVJd0RRWUpLb1pJaHZjTkFRRUwKQlFBd0RqRU1NQW9HQTFVRUF3d0RhMkZ6TUI0WERUSTBNRFl4T0RFNE16WXlOMW9YRFRJMU1EWXhPREU0TXpZeQpOMW93RGpFTU1Bb0dBMVVFQXd3RGEyRnpNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDCkFRRUFyMXBRam83cGlPdlBDVHRkSUVOZkc4eVZpK1dWMUZVTi82eFREcXJMeFpUdEFrWjE0M3VIVGZQOWExdXEKaFcxSW9heUpPVWpuWXNuUUh6dUVCZGtaNEh1d3pkeTZ3Um5lT1RSY2o3TitEd25aS21EcTF1YWZ6bEdzdG8vQgpoZnRtaWxVRjRZbm5GY0ROK3ZxajJlcDNhYlVramhrbUlRVDhwcjI1YkZ4TGFpd3dPbmx5TTVWUWM4bmFoZ2xuCjBNMGdOV0tJV0ZFSndoajBab2poMUw0ZGptenFVaU9tTkhCUDRRelNwNyswK3RXb3hJb1AyT2Fqa0p5MEljWkgKcS9OOWlTelZiZzFLL2tLZytkdS9QbWRqUCtqNTZsa0pPU1J6ZXpoK2R5NytHaHJCVDNVc21QbmNWM2NXVk1pOApFc1lDS2NUNUVNSGhhTmFHMFhEakptRzI4d0lEQVFBQm8xTXdVVEFkQmdOVkhRNEVGZ1FVZ1BUTkZjemQ5ajBFClgzN3A2SGh3UFJpY0JqOHdId1lEVlIwakJCZ3dGb0FVZ1BUTkZjemQ5ajBFWDM3cDZIaHdQUmljQmo4d0R3WUQKVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFDS2VxRkswSlcyYTVzS2JPQnl3WgppazB5MmpyRHJaUG5mMG9kTjVIbThtZWVuQnhteW9CeVZWRm9uUGVDaG1uWUZTdERtMlFJUTZnWVBtdEFhQ3VKCnRVeU5zNkxPQm1wR2JKaFRnNXljZXFXWnhYY3NmVkZ3ZHFxVXQ2NnRXdmNPeFZUQmdrN3h6RFFPbkxnRkxqZDYKSlZIeE16RkxXVFEwa00yVXJOOGd0T2RMazRhZUJhSzdiVHdaUEZ0RnR1YUZlYlFUbTRLY2ZSNXpzZkxTKzhpRgp1MWZGOVpKWkg2ZzZibENUeE50d3Z2eVMxVTMvS1AwVlQ5WVB3OTVmcElWMlNLT2QzejNZMGRKOUE5TGQ5TUkzCkwvWS8rNW05NEZCMTdTSWtERXpZM2d2TkxDSVZxODh2WHlnK2doVEhzSHNjYzNWcUUwK0x6cmZkemltbzMxRWQKTkE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t'
  export KID='my_key_123'
  export CACHED_KEY="{\"cached\":{\"keys\":[{\"pem\":\"$PEM\",\"kid\":\"$KID\",\"alg\":1}]}}"
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_kasr () {
      run sh -c "./otdfctl policy kas-registry $HOST $WITH_CREDS $*"
    }
}

teardown() {
    echo "created $CREATED"
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr delete --id "$ID"
    # run_otdfctl_kasr delete --id "$ID" --force
}

@test "create registration of a KAS with remote key" {
    URI="https://testing-create-remote.co"
    run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY" --json
        assert_output --partial "$REMOTE_KEY"
        assert_output --partial "$URI"
    export CREATED="$output"
}

@test "create registration of a KAS with cached key" {
    echo "cached: $CACHED_KEY"
    URI="https://testing-create-cached.co"
    run_otdfctl_kasr create --uri "$URI" -c "$CACHED_KEY" --json
        assert_output --partial "$KID"
        assert_output --partial "$PEM"
        assert_output --partial "$URI"
    export CREATED="$output"
}

@test "get registered KAS" {
    URI="https://testing-get.gov"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr get --id "$ID"
        assert_output --partial "$ID"
        assert_output --partial "$URI"
        assert_output --partial "URI"
        assert_output --partial "pem"
}

@test "update registered KAS" {
    URI="https://testing-update.net"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr update --id "$ID" -u "https://newuri.com" --public-key-remote "$REMOTE_KEY" --json
        assert_output --partial "$ID"
        assert_output --partial "https://newuri.com"
        assert_output --partial "$REMOTE_KEY"
        assert_output --partial "uri"
        refute_output --partial "pem"
        refute_output --partial "cached"
}

@test "list registered KASes" {
    URI="https://testing-list.io"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr list --json
    assert_output --partial "$ID"
    assert_output --partial "uri"
    assert_output --partial "$URI"
}
