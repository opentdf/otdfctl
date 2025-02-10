#!/usr/bin/env bats
load "./helpers.bash"

# Tests for kas registry

setup_file() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export HOST="${HOST:=--host=http://localhost:8080}"
  export DEBUG_LEVEL="--log-level debug"

  export REMOTE_KEY='https://hello.world/pubkey'
  PEM='-----BEGIN CERTIFICATE-----\nMIIC/TCCAeWgAwIBAgIUMu8o8Wh2HTA6TAeLCjC2f\n9pIeIwDQYJKoZIhvcNAQEL\nBQAwDjEMMAoGA1UEAwwDa2FzMB4XDTI0MDYxODE4M\nYyN1oXDTI1MDYxODE4MzYy\nN1owDjEMMAoGA1UEAwwDa2FzMIIBIjANBgkqhkiG9\n0BAQEFAAOCAQ8AMIIBCgKC\nAQEAr1pQjo7piOvPCTtdIENfG8yVi+WV1FUN/6xTD\nrLxZTtAkZ143uHTfP9a1uq\nhW1IoayJOUjnYsnQHzuEBdkZ4Huwzdy6wRneOTRcj\nN+DwnZKmDq1uafzlGsto/B\nhftmilUF4YnnFcDN+vqj2ep3abUkjhkmIQT8pr25b\nxLaiwwOnlyM5VQc8nahgln\n0M0gNWKIWFEJwhj0Zojh1L4djmzqUiOmNHBP4QzSp\n+0+tWoxIoP2OajkJy0IcZH\nq/N9iSzVbg1K/kKg+du/PmdjP+j56lkJOSRzezh+d\n7+GhrBT3UsmPncV3cWVMi8\nEsYCKcT5EMHhaNaG0XDjJmG28wIDAQABo1MwUTAdB\nNVHQ4EFgQUgPTNFczd9j0E\nX37p6HhwPRicBj8wHwYDVR0jBBgwFoAUgPTNFczd9\n0EX37p6HhwPRicBj8wDwYD\nVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCA\nEACKeqFK0JW2a5sKbOBywZ\nik0y2jrDrZPnf0odN5Hm8meenBxmyoByVVFonPeCh\nnYFStDm2QIQ6gYPmtAaCuJ\ntUyNs6LOBmpGbJhTg5yceqWZxXcsfVFwdqqUt66tW\ncOxVTBgk7xzDQOnLgFLjd6\nJVHxMzFLWTQ0kM2UrN8gtOdLk4aeBaK7bTwZPFtFt\naFebQTm4KcfR5zsfLS+8iF\nu1fF9ZJZH6g6blCTxNtwvvyS1U3/KP0VT9YPw95fp\nV2SKOd3z3Y0dJ9A9Ld9MI3\nL/Y/+5m94FB17SIkDEzY3gvNLCIVq88vXyg+ghTHs\nscc3VqE0+Lzrfdzimo31Ed\nNA==\n-----END CERTIFICATE-----'
  export KID='my_key_123'
  export CACHED_KEY=$(printf '{"cached":{"keys":[{"kid":"%s","alg":1,"pem":"%s"}]}}' "$KID" "$PEM" )

  export FIXTURE_PUBLIC_KEY_ID="f478f1cd-df6e-4a55-9603-d961b36ea392"
}

setup() {
    setup_helper
}

teardown() {
    log_debug "Running teardown for: $BATS_TEST_NAME"
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr delete --id "$ID" --force
        
    cleanup_helper
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

@test "create KAS registration with invalid key - fails" {
    BAD_CACHED=(
        '{"cached":{"keys":[{"pem":"bad"}]}}'
        '{"cached":[]}'
        '{"cached":{"keys":[{]}}'
    )

    for BAD_KEY in "${BAD_CACHED[@]}"; do
        URI='https://bad.pem/kas'
        run_otdfctl_kasr create --uri "$URI" --public-keys "$BAD_KEY"
            assert_failure
            assert_output --partial "KAS registry key is invalid"
    done
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
    export CREATED="$output"
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
    export CREATED=$(./otdfctl $HOST $DEBUG_LEVEL $WITH_CREDS policy kas-registry create --uri "$URI" -r "$REMOTE_KEY" -n "$NAME" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_kasr update --id "$ID" -u "https://newuri.com" -n "newer-name" -c '"$CACHED_KEY"' --json
        assert_output --partial "$ID"
        assert_output --partial "https://newuri.com"
        assert_output --partial "kid"
        assert_output --partial "pem"
        assert_output --partial "alg"
        assert_output --partial "newer-name"
        refute_output --partial "$NAME"
        refute_output --partial "$URI"
        refute_output --partial "remote"
        refute_output --partial "$REMOTE_KEY"
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
        run_otdfctl_kasr update -i "$ID" -r "$REMOTE_KEY" --uri "$URI"
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
        run_otdfctl_kasr update -c '"$CACHED_KEY"' --id "$ID" --name "$NAME"
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

    run_otdfctl_kasr list
        assert_output --partial "Total"
        assert_line --regexp "Current Offset.*0"
}

@test "create_public_key_success" {
    log_info "Starting test: $BATS_TEST_NAME"
    
    create_kas "$KAS_URI" "$KAS_NAME"

    # # Generate the key pair and set variables
    # create_
    # log_debug "RSA_2048_PUBLIC_KEY=\"$RSA_2048_PUBLIC_KEY\""
    # eval "$(gen_rsa_4096)"
    # log_debug "RSA_4096_PUBLIC_KEY=\"$RSA_4096_PUBLIC_KEY\""
    # eval "$(gen_ec256)"
    # log_debug "EC_256_PUBLIC_KEY=\"$EC_256_PUBLIC_KEY\""
    # eval "$(gen_ec384)"
    # log_debug "EC_384_PUBLIC_KEY=\"$EC_384_PUBLIC_KEY\""
    # eval "$(gen_ec521)"
    # log_debug "EC_521_PUBLIC_KEY=\"$EC_521_PUBLIC_KEY\""
    
    KID="test_key_123" 

    ########## Creating RSA 2048 public keys ##########
    create_public_key "$KAS_ID" "$KID" "$RSA_2048_ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas $KAS_ID --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_2048_ALG" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_2048_ALG" --json

    # log_debug "Raw Output:" # Debug log: Raw output
    # log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get the ID of the public key
    # PUBLIC_KEY_ID=$(echo "$output" | jq -r '.id')
    # PUBLIC_KEY_IDS+=("$PUBLIC_KEY_ID")

    ########## Creating RSA 4096 public keys ##########
    create_public_key "$KAS_ID" "$KID" "$RSA_4096_ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas $KAS_ID --key \"$(echo "$RSA_4096_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_4096_ALG" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_4096_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_4096_ALG" --json

    # log_debug "Raw Output:" # Debug log: Raw output
    # log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get the ID of the public key
    # PUBLIC_KEY_ID=$(echo "$output" | jq -r '.id')
    # PUBLIC_KEY_IDS+=("$PUBLIC_KEY_ID")

    ########## Creating EC 256 public keys ##########
    create_public_key "$KAS_ID" "$KID" "$EC_256_ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas $KAS_ID --key \"$(echo "$EC_256_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$EC_256_ALG" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$EC_256_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$EC_256_ALG" --json

    # log_debug "Raw Output:" # Debug log: Raw output
    # log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get the ID of the public key
    # PUBLIC_KEY_ID=$(echo "$output" | jq -r '.id')
    # PUBLIC_KEY_IDS+=("$PUBLIC_KEY_ID")

    ########## Creating EC 384 public keys ##########
    create_public_key "$KAS_ID" "$KID" "$EC_384_ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas $KAS_ID --key \"$(echo "$EC_384_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$EC_384_ALG" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$EC_384_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$EC_384_ALG" --json

    # log_debug "Raw Output:" # Debug log: Raw output
    # log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get the ID of the public key
    # PUBLIC_KEY_ID=$(echo "$output" | jq -r '.id')
    # PUBLIC_KEY_IDS+=("$PUBLIC_KEY_ID")
}

@test "create_public_key_required_flags" {
    log_info "Starting test: $BATS_TEST_NAME"
    
    create_kas "$KAS_URI" "$KAS_NAME"

    # Generate the key pair and set variables
    eval "$(gen_rsa_2048)"
    log_debug "RSA_2048_PUBLIC_KEY=\"$RSA_2048_PUBLIC_KEY\""

    ALG="rsa:2048"
    KID="test_key_123" 

    # Missing KAS Flag
    run_otdfctl_kasr public-key create --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$ALG" --json
    assert_failure # Check if the command failed requiring the KAS flag
    assert_output --partial "Flag '--kas' is required"

    # Missing Key Flag
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key-id "$KID" --algorithm "$ALG" --json
    assert_failure # Check if the command failed requiring the Key flag
    assert_output --partial "Flag '--key' is required"

    # Missing Key ID Flag
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --algorithm "$ALG" --json
    assert_failure # Check if the command failed requiring the Key ID flag
    assert_output --partial "Flag '--key-id' is required"

    # Missing Algorithm Flag
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --json
    assert_failure # Check if the command failed requiring the Algorithm flag
    assert_output --partial "Flag '--algorithm' is required"
}

@test "create_public_key_invalid_algorithm" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    # Generate the key pair and set
    eval "$(gen_rsa_2048)"
    log_debug "RSA_2048_PUBLIC_KEY=\"$RSA_2048_PUBLIC_KEY\""

    ALG="rsa:2048"
    KID="test_key_123"

    # Invalid Algorithm
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "rsa:4096" --json
    assert_failure # Check if the command failed with an invalid algorithm
    assert_output --partial "invalid rsa key size"

    # Invalid Algorithm
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "ec:secp256r1" --json
    assert_failure # Check if the command failed with an invalid algorithm
    assert_output --partial "ey algorithm does not match the provided algorithm"

    # Unsupported Algorithm
    run_otdfctl_kasr public-key create --kas "$KAS_ID" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "rsa:1024" --json
    assert_failure # Check if the command failed with an unsupported algorithm
    assert_output --partial "unsupported algorithm"
}

@test "add_public_key_by_kas_uri" {
    log_info "Starting test: $BATS_TEST_NAME"
    
    create_kas "$KAS_URI" "$KAS_NAME"

    # Generate the key pair and set variables
    # eval "$(gen_rsa_2048)"
    # log_debug "PUBLIC_KEY=\"$RSA_2048_PUBLIC_KEY\""

    ALG="rsa:2048"
    KID="test_key_123" 

    create_public_key "$KAS_URI" "$KID" "$ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas $KAS_ID --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_2048_ALG" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_URI" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$RSA_2048_ALG" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

}

@test "add_public_key_by_kas_name" {
    log_info "Starting test: $BATS_TEST_NAME"
    
    create_kas "$KAS_URI" "$KAS_NAME"

    # Generate the key pair and set variables
    # eval "$(gen_rsa_2048)"
    # log_debug "PUBLIC_KEY=\"$RSA_2048_PUBLIC_KEY\""

    ALG="rsa:2048"
    KID="test_key_123" 

    create_public_key "$KAS_NAME" "$KID" "$ALG"
    # log_debug "Running ${run_otdfctl_kasr} public-key create --kas "$KAS_NAME" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "rsa:2048" --json"
    # run_otdfctl_kasr public-key create --kas "$KAS_NAME" --key \"$(echo "$RSA_2048_PUBLIC_KEY" | base64)\" --key-id "$KID" --algorithm "$ALG" --json

    # log_debug "Raw Output:" # Debug log: Raw output
    # log_debug "$output"

    assert_success # Check if the command ran successfully
}

@test "update_public_key_labels" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    ALG="rsa:2048"
    KID="test"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Update the public key with labels
    log_debug "Running ${run_otdfctl_kasr} public-key update --id $PUBLIC_KEY_ID --label test=test --json"
    run_otdfctl_kasr public-key update --id "$PUBLIC_KEY_ID" --label test=test --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Get public key by ID and check if the labels are set
    log_debug "Running ${run_otdfctl_kasr} public-key get --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key get --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Check json response for the labels
    echo "$output" | jq -e '.metadata.labels | has("test")' || fail "Label not found"
}

@test "update_public_key_labels_force" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    ALG="rsa:2048"
    KID="test"

    create_public_key "$KAS_ID" "$KID" "$ALG" "--label test=test"

    # Update the public key with labels
    log_debug "Running ${run_otdfctl_kasr} public-key update --id $PUBLIC_KEY_ID --label test1=test1 --force-replace-labels --json"
    run_otdfctl_kasr public-key update --id "$PUBLIC_KEY_ID" --label test1=test1 --force-replace-labels --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Get public key by ID and check if the labels are set
    log_debug "Running ${run_otdfctl_kasr} public-key get --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key get --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Check json response for the labels
    echo "$output" | jq -e '.metadata.labels | (has("test") | not) and has("test1")' || fail "Labels check failed"
}

@test "get_public_key" {
    log_info "Starting test: $BATS_TEST_NAME"

    log_debug "Running ${run_otdfctl_kasr} public-key get --id $FIXTURE_PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key get --id $FIXTURE_PUBLIC_KEY_ID --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    if ! echo "$output" | jq -e ; then
      fail "Output is not valid JSON"
    fi

    # Parse the JSON output using jq
    output_json=$(echo "$output" | jq -c '.')

    echo "$output_json" | jq -e '
    .[0] as $root |
    [
        if $root.id then empty else "id" end,
        if $root.is_active then empty else "is_active" end,
        if $root.was_mapped then empty else "was_mapped" end,
        if $root.public_key then empty else "public_key" end,
        if $root.public_key then (
            if $root.public_key.pem then empty else "public_key.pem" end,
            if $root.public_key.kid then empty else "public_key.kid" end,
            if $root.public_key.alg then empty else "public_key.alg" end
        ) else empty end,
        if $root.kas then empty else "kas" end,
        if $root.kas then (
            if $root.kas.id then empty else "kas.id" end,
            if $root.kas.uri then empty else "kas.uri" end,
            if $root.kas.name then empty else "kas.name" end
        ) else empty end,
        if $root.metadata then empty else "metadata" end
    ] | if length > 0 then error("Missing fields: " + join(", ")) else true end
' || fail "Structure validation failed"
}

@test "get_public_key_required_flags" {
    log_info "Starting test: $BATS_TEST_NAME"
    # Missing ID Flag
    run_otdfctl_kasr public-key get --json
    assert_failure # Check if the command failed requiring the ID flag
    assert_output --partial "Flag '--id' is required"
}


@test "list_public_keys" {
    log_info "Starting test: $BATS_TEST_NAME"

    log_debug "Running ${run_otdfctl_kasr} public-key list --json"
    run_otdfctl_kasr public-key list --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Parse the JSON output using jq
    output_json=$(echo "$output" | jq -c '.')

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output_json" | jq 'length')" -gt 0 ]

    echo "$output_json" | jq -e '
    .[0] as $root |
    [
        if $root.id then empty else "id" end,
        if $root.is_active then empty else "is_active" end,
        if $root.was_mapped then empty else "was_mapped" end,
        if $root.public_key then empty else "public_key" end,
        if $root.public_key then (
            if $root.public_key.pem then empty else "public_key.pem" end,
            if $root.public_key.kid then empty else "public_key.kid" end,
            if $root.public_key.alg then empty else "public_key.alg" end
        ) else empty end,
        if $root.kas then empty else "kas" end,
        if $root.kas then (
            if $root.kas.id then empty else "kas.id" end,
            if $root.kas.uri then empty else "kas.uri" end,
            if $root.kas.name then empty else "kas.name" end
        ) else empty end,
        if $root.metadata then empty else "metadata" end
    ] | if length > 0 then error("Missing fields: " + join(", ")) else true end
' || fail "Structure validation failed"
}

@test "list_public_keys_by_kas" {
    # Create a KAS to Filter By
    ALG="rsa:2048"
    
    create_kas "$KAS_URI" "$KAS_NAME"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Filter By ID
    log_debug "Running ${run_otdfctl_kasr} public-key list --kas "$KAS_ID" --json"
    run_otdfctl_kasr public-key list --kas "$KAS_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas.id' | while read -r id; do
        log_debug "Checking KAS ID: $id against $KAS_ID"
        [ "$id" = "$KAS_ID" ] || fail "KAS ID does not match"
    done


     # Filter By URI
    log_debug "Running ${run_otdfctl_kasr} public-key list --kas "$KAS_URI" --json"
    run_otdfctl_kasr public-key list --kas "$KAS_URI" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas.uri' | while read -r uri; do
        log_debug "Checking KAS ID: $uri against $KAS_URI"
        [ "$uri" = "$KAS_URI" ] || fail "KAS ID does not match"
    done


     # Filter By Name
    log_debug "Running ${run_otdfctl_kasr} public-key list --kas "$KAS_NAME" --json"
    run_otdfctl_kasr public-key list --kas "$KAS_NAME" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas.name' | while read -r name; do
        log_debug "Checking KAS ID: $name against $KAS_NAME"
        [ "$name" = "$KAS_NAME" ] || fail "KAS ID does not match"
    done
}

@test "list_public_key_mappings" {
    log_info "Starting test: $BATS_TEST_NAME"

    log_debug "Running ${run_otdfctl_kasr} public-key list-mappings --json"
    run_otdfctl_kasr public-key list-mappings --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -e '
    .[0].public_keys[] |  (
        (.key | has("id")) and
        (.key | has("is_active")) and
        (.key | has("was_mapped")) and
        (.key | has("public_key")) and
        (.key.public_key | has("pem") and has("kid") and has("alg"))
    ) or error("Missing required public key fields in response structure")
    ' || fail "Structure validation failed"
}

@test "list_public_key_mappings_by_kas" {
    log_info "Starting test: $BATS_TEST_NAME"

    # Create a KAS to Filter By
    ALG="rsa:2048"

    create_kas "$KAS_URI" "$KAS_NAME"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Filter By ID
    log_debug "Running ${run_otdfctl_kasr} public-key list-mappings --kas "$KAS_ID" --json"
    run_otdfctl_kasr public-key list-mappings --kas "$KAS_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas_id' | while read -r id; do
        log_debug "Checking KAS ID: $id against $KAS_ID"
        [ "$id" = "$KAS_ID" ] || fail "KAS ID does not match"
    done

    # Filter By URI
    log_debug "Running ${run_otdfctl_kasr} public-key list-mappings --kas "$KAS_URI" --json"
    run_otdfctl_kasr public-key list-mappings --kas "$KAS_URI" --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas_uri' | while read -r uri; do
        log_debug "Checking KAS ID: $uri against $KAS_URI"
        [ "$uri" = "$KAS_URI" ] || fail "KAS ID does not match"
    done


    # Filter By Name
    log_debug "Running ${run_otdfctl_kasr} public-key list-mappings --kas "$KAS_NAME" --json"
    run_otdfctl_kasr public-key list-mappings --kas "$KAS_NAME" --json

    log_debug "Raw Output:" # Debug log: Raw output

    assert_success # Check if the command ran successfully

    # Check if the output is valid JSON and is an array (without using jq -t)
    if ! echo "$output" | jq -e '.[0]'; then
      fail "Output is not a JSON array"
    fi

    # Check if the output is not empty (contains at least one key)
    [ "$(echo "$output" | jq 'length')" -gt 0 ]

    echo "$output" | jq -r '.[].kas_name' | while read -r name; do
        log_debug "Checking KAS ID: $name against $KAS_NAME"
        [ "$name" = "$KAS_NAME" ] || fail "KAS ID does not match"
    done
}

@test "activate_deactivate_public_key" {
    log_info "Starting test: $BATS_TEST_NAME"

    create_kas "$KAS_URI" "$KAS_NAME"

    ALG="rsa:2048"
    KID="test"

    create_public_key "$KAS_ID" "$KID" "$ALG"

    # Deactivate the public key
    log_debug "Running ${run_otdfctl_kasr} public-key deactivate --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key deactivate --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get public key by ID and check if the key is deactivated
    log_debug "Running ${run_otdfctl_kasr} public-key get --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key get --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check json response for the is_active flag
    echo "$output" | jq -e '.is_active == {}' || fail "Public key is still active"

    # Activate the public key
    log_debug "Running ${run_otdfctl_kasr} public-key activate --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key activate --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Get public key by ID and check if the key is activated
    log_debug "Running ${run_otdfctl_kasr} public-key get --id $PUBLIC_KEY_ID --json"
    run_otdfctl_kasr public-key get --id "$PUBLIC_KEY_ID" --json

    log_debug "Raw Output:" # Debug log: Raw output
    log_debug "$output"

    assert_success # Check if the command ran successfully

    # Check json response for the is_active flag
    echo "$output" | jq -e '.is_active.value' || fail "Public key is not active"
}