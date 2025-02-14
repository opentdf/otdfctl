# helpers.bash
#!/usr/bin/env bash

# OTDFCTL Helper Functions
run_otdfctl_kasr() {
    run sh -c "./otdfctl policy kas-registry $HOST $WITH_CREDS $*"
}

create_kas() {
    log_debug "Creating KAS... $1 $2"

    run_otdfctl_kasr create --uri "$1" -n "$2" --json

    log_debug "Created KAS: $output" # Debug log: the output of the create command

    KAS_ID=$(echo "$output" | jq -r '.id')
}

create_public_key() {
    local kas="$1"
    local key_id="$2"
    local algorithm="$3"
    local key_content
    local label_args="$4"

    log_debug "Creating public key..."

    # Select the appropriate key generation function based on the algorithm
    case "$algorithm" in
    "$RSA_2048_ALG")
        eval "$(gen_rsa_2048)"
        key_content="$RSA_2048_PUBLIC_KEY"
        ;;
    "$RSA_4096_ALG")
        eval "$(gen_rsa_4096)"
        key_content="$RSA_4096_PUBLIC_KEY"
        ;;
    "$EC_256_ALG")
        eval "$(gen_ec256)"
        key_content="$EC_256_PUBLIC_KEY"
        ;;
    "$EC_384_ALG")
        eval "$(gen_ec384)"
        key_content="$EC_384_PUBLIC_KEY"
        ;;
    "$EC_521_ALG")
        eval "$(gen_ec521)"
        key_content="$EC_521_PUBLIC_KEY"
        ;;
    *)
        log_error "Unsupported algorithm: $algorithm"
        return 1
        ;;
    esac

    # Verify key content is not empty
    if [ -z "$key_content" ]; then
        log_info "Empty key content for algorithm: $algorithm"
        return 1
    fi

    # Base64 encode the key content
    key_content=$(echo "$key_content" | base64 -w 0)

    log_debug "Running ${run_otdfctl_kasr} public-key create --kas $kas --key "$key_content" --key-id $key_id --algorithm $algorithm $label_args --json"
    run_otdfctl_kasr public-key create \
        --kas "$kas" \
        --key "$key_content" \
        --key-id "$key_id" \
        --algorithm "$algorithm" \
        $label_args \
        --json

    if [ -z "$output" ]; then
        log_info "Failed to create public key"
        return 1
    fi

    log_debug "Created public key: $output"

    PUBLIC_KEY_ID=$(echo "$output" | jq -r '.id')
    PUBLIC_KEY_IDS+=("$PUBLIC_KEY_ID")
}

# Setup Helper
setup_helper() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # Initialize IDs to empty strings in case creation fails
    KAS_ID=""
    PUBLIC_KEY_ID=""
    PUBLIC_KEY_IDS=() # Initialize an empty array

    KAS_URI="https://testing-public-key.io"
    KAS_NAME="public-key-kas"

    RSA_2048_ALG="rsa:2048"
    RSA_4096_ALG="rsa:4096"
    EC_256_ALG="ec:secp256r1"
    EC_384_ALG="ec:secp384r1"
    EC_521_ALG="ec:secp521r1"
}

# Cleanup Helper
cleanup_helper() {
    # Iterate over the array of public key IDs and delete them
    for PUBLIC_KEY_ID in "${PUBLIC_KEY_IDS[@]}"; do
        if [ -n "$PUBLIC_KEY_ID" ]; then
            log_debug "Running ${run_otdfctl_kasr} public-key unsafe delete --id $PUBLIC_KEY_ID --force --json"
            run_otdfctl_kasr public-key unsafe delete --id "$PUBLIC_KEY_ID" --force --json
            log_debug "$output"
            if [ $? -ne 0 ]; then
                log_info "Error: Failed to delete public key with ID: $PUBLIC_KEY_ID"
            fi
            log_debug "Deleted public key with ID: $PUBLIC_KEY_ID"
        fi
    done
    if [ -n "$KAS_ID" ]; then
        log_debug "Running ${run_otdfctl_kasr} delete --id $KAS_ID --force --json"
        run_otdfctl_kasr delete --id "$KAS_ID" --force --json
        log_debug "$output"
        if [ $? -ne 0 ]; then
            log_info "Error: Failed to delete KAS registry with ID: $KAS_ID"
        fi
        log_debug "Deleted KAS registry with ID: $KAS_ID"
    fi
}

# Helper function for debug logging
log_debug() {
    if [[ "${BATS_DEBUG:-0}" == "1" ]]; then
        echo "DEBUG($BATS_TEST_NAME): $1" >&3
    fi
}

# Helper function for info logging
log_info() {
    echo "INFO($BATS_TEST_NAME): $1" >&3
}

# Helper function to generate a rsa 2048 key pair
gen_rsa_2048() {
    log_debug "Generating RSA 2048 key pair"
    local private_key public_key

    # Generate private key
    private_key=$(openssl genrsa 2048)

    # Extract public key
    public_key=$(echo "$private_key" | openssl rsa -pubout)

    # Output using proper escaping
    printf 'export RSA_2048_PUBLIC_KEY=%q\n' "$public_key"
}

# Helper function to generate a rsa 4096 key pair
gen_rsa_4096() {
    log_debug "Generating RSA 4096 key pair"
    local private_key public_key

    # Generate private key
    private_key=$(openssl genrsa 4096)

    # Extract public key
    public_key=$(echo "$private_key" | openssl rsa -pubout)

    printf 'export RSA_4096_PUBLIC_KEY=%q\n' "$public_key"
}

# Helper function to generate an EC 256 key pair
gen_ec256() {
    log_debug "Generating EC 256 key pair"
    local private_key public_key

    # Generate private key
    private_key=$(openssl ecparam -name prime256v1 -genkey)

    # Extract public key
    public_key=$(echo "$private_key" | openssl ec -pubout)

    printf 'export EC_256_PUBLIC_KEY=%q\n' "$public_key"
}

# Helper function to generate an EC 384 key pair
gen_ec384() {
    log_debug "Generating EC 384 key pair"
    local private_key public_key

    # Generate private key
    private_key=$(openssl ecparam -name secp384r1 -genkey)

    # Extract public key
    public_key=$(echo "$private_key" | openssl ec -pubout)

    printf 'export EC_384_PUBLIC_KEY=%q\n' "$public_key"
}

# Helper function to generate an EC 521 key pair
gen_ec521() {
    log_debug "Generating EC 521 key pair"
    local private_key public_key

    # Generate private key
    private_key=$(openssl ecparam -name secp521r1 -genkey)

    # Extract public key
    public_key=$(echo "$private_key" | openssl ec -pubout)

    printf 'export EC_521_PUBLIC_KEY=%q\n' "$public_key"
}
