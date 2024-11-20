#!/usr/bin/env bats

# Tests for kas registry

setup_file() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export HOST='--host http://localhost:8080'
  export DEBUG_LEVEL="--log-level debug"

  export REMOTE_KEY='https://hello.world/pubkey'
  export PEM='-----BEGIN CERTIFICATE-----
MIIC/TCCAeWgAwIBAgIUMu8o8Wh2HTA6TAeLCjC2fV9pIeIwDQYJKoZIhvcNAQEL
BQAwDjEMMAoGA1UEAwwDa2FzMB4XDTI0MDYxODE4MzYyN1oXDTI1MDYxODE4MzYy
N1owDjEMMAoGA1UEAwwDa2FzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAr1pQjo7piOvPCTtdIENfG8yVi+WV1FUN/6xTDqrLxZTtAkZ143uHTfP9a1uq
hW1IoayJOUjnYsnQHzuEBdkZ4Huwzdy6wRneOTRcj7N+DwnZKmDq1uafzlGsto/B
hftmilUF4YnnFcDN+vqj2ep3abUkjhkmIQT8pr25bFxLaiwwOnlyM5VQc8nahgln
0M0gNWKIWFEJwhj0Zojh1L4djmzqUiOmNHBP4QzSp7+0+tWoxIoP2OajkJy0IcZH
q/N9iSzVbg1K/kKg+du/PmdjP+j56lkJOSRzezh+dy7+GhrBT3UsmPncV3cWVMi8
EsYCKcT5EMHhaNaG0XDjJmG28wIDAQABo1MwUTAdBgNVHQ4EFgQUgPTNFczd9j0E
X37p6HhwPRicBj8wHwYDVR0jBBgwFoAUgPTNFczd9j0EX37p6HhwPRicBj8wDwYD
VR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEACKeqFK0JW2a5sKbOBywZ
ik0y2jrDrZPnf0odN5Hm8meenBxmyoByVVFonPeChmnYFStDm2QIQ6gYPmtAaCuJ
tUyNs6LOBmpGbJhTg5yceqWZxXcsfVFwdqqUt66tWvcOxVTBgk7xzDQOnLgFLjd6
JVHxMzFLWTQ0kM2UrN8gtOdLk4aeBaK7bTwZPFtFtuaFebQTm4KcfR5zsfLS+8iF
u1fF9ZJZH6g6blCTxNtwvvyS1U3/KP0VT9YPw95fpIV2SKOd3z3Y0dJ9A9Ld9MI3
L/Y/+5m94FB17SIkDEzY3gvNLCIVq88vXyg+ghTHsHscc3VqE0+Lzrfdzimo31Ed
NA==
-----END CERTIFICATE-----'
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
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr delete --id "$ID" --force
}

@test "create registration of a KAS with remote key" {
    URI="https://testing-create-remote.co"
    NAME="my_kas_name"
    run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY" -n "$NAME" --json
        assert_output --partial "$REMOTE_KEY"
        assert_output --partial "$URI"
        assert_output --partial "$NAME"
    export CREATED="$output"
}

@test "create KAS registration with invalid URI - fails" {
    BAD_URIS=(
        "no-scheme.co"
        "localhost"
        "http://example.com:abc"
        "https ://example.com"
    )

    for URI in "${BAD_URIS[@]}"; do
        run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY"
            assert_failure
            assert_output --partial "Failed to create Registered KAS"
            assert_output --partial "uri: "
    done
}

@test "create KAS registration with duplicate URI - fails" {
    URI="https://testing-duplication.io"
    run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY"
        assert_success
    run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY"
        assert_failure
        assert_output --partial "Failed to create Registered KAS entry"
        assert_output --partial "AlreadyExists"
}

@test "create KAS registration with duplicate name - fails" {
    NAME="duplicate_name_kas"
    run_otdfctl_kasr create --uri "https://testing-duplication.name.io" -r "$REMOTE_KEY" -n "$NAME"
        assert_success
    run_otdfctl_kasr create --uri "https://testing-duplication.name.net" -r "$REMOTE_KEY" -n "$NAME"
        assert_failure
        assert_output --partial "Failed to create Registered KAS entry"
        assert_output --partial "AlreadyExists"
}

@test "create KAS registration with invalid name - fails" {
    URI="http://creating.kas.invalid.name/kas"
    BAD_NAMES=(
        "-bad-name"
        "bad-name-"
        "_bad_name"
        "bad_name_"
        "name@with!special#chars"
        "$(printf 'a%.0s' {1..254})" # Generates a string of 254 'a' characters
    )

    for NAME in "${BAD_NAMES[@]}"; do
        echo "testing $NAME"
        run_otdfctl_kasr create --uri "$URI" -r "$REMOTE_KEY" -n "$NAME"
            assert_failure
            assert_output --partial "Failed to create Registered KAS"
            assert_output --partial "name: "
    done
}

@test "create KAS with cached key then get it" {
    URI="https://testing-get.gov"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run echo $CREATED
        assert_output --partial "$URI"
        assert_output --partial "uri"
        assert_output --partial "pem"
        assert_output --partial "$PEM"
        assert_output --partial "$KID"

    run_otdfctl_kasr get -i "$ID" --json
        assert_output --partial "$ID"
        assert_output --partial "$URI"
        assert_output --partial "uri"
        assert_output --partial "$PEM"
        assert_output --partial "pem"
        assert_output --partial "$KID"
}

@test "update registered KAS" {
    URI="https://testing-update.net"
    NAME="new-kas-testing-update"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" -n "$NAME" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr update --id "$ID" -u "https://newuri.com" -n "newer-name" --public-key-remote "$REMOTE_KEY" --json
        assert_output --partial "$ID"
        assert_output --partial "https://newuri.com"
        assert_output --partial "$REMOTE_KEY"
        assert_output --partial "newer-name"
        assert_output --partial "uri"
        refute_output --partial "pem"
        refute_output --partial "$NAME"
        refute_output --partial "cached"
}

@test "update registered KAS with invalid URI - fails" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "https://bad-update.uri.kas" -c "$CACHED_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    BAD_URIS=(
        "no-scheme.co"
        "localhost"
        "http://example.com:abc"
        "https ://example.com"
    )

    for URI in "${BAD_URIS[@]}"; do
        run_otdfctl_kasr update -i "$ID" --uri "$URI"
            assert_failure
            assert_output --partial "$ID"
            assert_output --partial "Failed to update Registered KAS entry"
            assert_output --partial "uri: "
    done
}

@test "update registered KAS with invalid name - fails" {
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "https://bad-update.name.kas" -r "$REMOTE_KEY" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    BAD_NAMES=(
        "-bad-name"
        "bad-name-"
        "_bad_name"
        "bad_name_"
        "name@with!special#chars"
        "$(printf 'a%.0s' {1..254})" # Generates a string of 254 'a' characters
    )

    for NAME in "${BAD_NAMES[@]}"; do
        run_otdfctl_kasr update --name "$NAME" -c "$CACHED_KEY" --id "$ID"
            assert_failure
            assert_output --partial "Failed to update Registered KAS"
            assert_output --partial "name: "
    done
}

@test "list registered KASes" {
    URI="https://testing-list.io"
    NAME="listed-kas"
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -c "$CACHED_KEY" -n "$NAME" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr list --json
    assert_output --partial "$ID"
    assert_output --partial "uri"
    assert_output --partial "$URI"
    assert_output --partial "name"
    assert_output --partial "$NAME"
}
