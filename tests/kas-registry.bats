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
  export URI="https://end-to-end-kas.com"
}

teardown() {
    ID=$(echo "$CREATED" | jq -r '.id')
    printf "y" | ./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry delete --id "$ID" --force
}

@test "create registration of a KAS with remote key" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -r "$REMOTE_KEY" --json)
    echo $CREATED | grep "$REMOTE_KEY"
    echo $CREATED | grep "$URI"
}

@test "create registration of a KAS with cached key" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    echo $CREATED | grep "$KID"
    echo $CREATED | grep "$PEM"
    echo $CREATED | grep "$URI"
}

@test "get registered KAS" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    RESULT=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry get --id "$ID")
    echo $RESULT | grep "$ID"
    echo $RESULT | grep "$URI"
    echo $RESULT | grep -i "uri"
    echo $RESULT | grep "pem"
}

@test "update registered KAS" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    RESULT=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry update --id "$ID" -u "https://newuri.com" --public-key-remote "$REMOTE_KEY" --json)
    echo $RESULT | grep "$ID"
    echo $RESULT | grep "https://newuri.com"
    echo $RESULT | grep "$REMOTE_KEY"
    echo $RESULT | grep -i "uri"
    [ "$(echo "$RESULT" | grep -c "pem")" -eq 0 ]
    [ "$(echo "$RESULT" | grep -c "cached")" -eq 0 ]
}

@test "list registered KASes" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    RESULT=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry list --json)
    echo $RESULT | grep "$ID"
    echo $RESULT | grep -i "uri"
    echo $RESULT | grep -i "$URI"
}
