#!/usr/bin/env bats

setup_file() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export HOST='--host http://localhost:8080'
  export DEBUG_LEVEL="--log-level debug"
  export VALID_CONFIG=$(printf '{"cached":"key"}')
  export BASE64_CONFIG=eyJjYWNoZWQiOiAia2V5In0=
}

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"

    # invoke binary with credentials
    run_otdfctl_key_pc () {
      run sh -c "./otdfctl keymanagement provider $HOST $WITH_CREDS $*"
    }
}

teardown() {
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_key_pc delete --id "$ID"
}

#########
# Create Provider Configuration
#########
# Test: Fail to create provider configuration without config
@test "fail to create provider configuration without config" {
    run_otdfctl_key_pc create --name test-value
    assert_failure
    assert_output --partial "Flag '--config' is required"
}

# Test: Fail to create provider configuration without name
@test "fail to create provider configuration without name" {
    run_otdfctl_key_pc create --config '{}'
    assert_failure
    assert_output --partial "Flag '--name' is required"
}

# Test: Fail to create provider configuration with invalid config
@test "fail to create provider configuration with invalid config" {
    run_otdfctl_key_pc create --name test-config --config test-value
    assert_failure
    assert_output --partial "Invalid JSON format for config"
}

# Test: Create provider configuration
@test "create provider configuration" {
    export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config --config "$VALID_CONFIG" --json)
    run echo $CREATED
        assert_output --partial "name"
        assert_output --partial "test-config"
        assert_output --partial "config"
        assert_output --partial "$BASE64_CONFIG"
}

 # Test: Get provider configuration by id
 @test "get provider configuration by id" {
     export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-2 --config "$VALID_CONFIG" --json)
     ID=$(echo "$CREATED" | jq -r '.id')
     run_otdfctl_key_pc get --id "$ID" --json
     assert_success
     assert_output --partial "test-config-2"
 }

 # Test: Get provider configuration by name
 @test "get provider configuration by name" {
     export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-3 --config "$VALID_CONFIG" --json)
     NAME=$(echo "$CREATED" | jq -r '.name')
     run_otdfctl_key_pc get --name "$NAME" --json
     assert_success
     assert_output --partial "test-config-3"
 }

 
 # Test: Fail to get provider configuration without flags
 @test "fail to get provider configuration - no required flags" {
     run_otdfctl_key_pc get
     assert_failure
 }
 
 # Test: Fail to get provider configuration with non-existent name
 @test "fail to get provider configuration with non-existent name" {
     run_otdfctl_key_pc get --name non-existent-config
     assert_failure
     assert_output --partial "Failed to get provider config: not_found"
 }
 
 # Test: List provider configurations
 @test "list provider configurations" {
    NAME=tst-config-4
    export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name "$NAME" --config "$VALID_CONFIG" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_key_pc list --json
        assert_output --partial "$ID"
        assert_output --partial "name"
        assert_output --partial "$NAME"
        assert_output --partial "config_json"
        assert_output --partial "$BASE64_CONFIG"

    run_otdfctl_key_pc list
        assert_output --partial "Total"
        assert_line --regexp "Current Offset.*0"
 }
 # Test: Update provider configuration - success
 @test "update provider configuration - success" {
     NAME="test-config-5"
     UPDATED_NAME="test-config-5-updated"
     UPDATED_CONFIG=$(printf '{"cached":"key-updated"}')
     BASE64_UPDATED_CONFIG=eyJjYWNoZWQiOiAia2V5LXVwZGF0ZWQifQ==
     export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-5 --config "$VALID_CONFIG" --json)
     ID=$(echo "$CREATED" | jq -r '.id')
     run_otdfctl_key_pc update --id "$ID" --name $UPDATED_NAME --config "$UPDATED_CONFIG"
     export UPDATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider update --id "$ID" --name $UPDATED_NAME --config "$UPDATED_CONFIG" --json)
     run echo $UPDATED
        assert_output --partial "id"
        assert_output --partial "$ID"
        assert_output --partial "name"
        assert_output --partial "$UPDATED_NAME"
        assert_output --partial "config"
        assert_output --partial "$BASE64_UPDATED_CONFIG"
 }

 @test "fail to update provider configuration - missing id" {
     run_otdfctl_key_pc update --name test-config
     assert_failure
     assert_output --partial "Flag '--id' is required"
 }

 @test "fail to update provider configuration - no optional flags" {
    export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-6 --config "$VALID_CONFIG" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_key_pc update --id "$ID"
    assert_failure
    assert_output --partial "At least one field (name, config, or metadata labels) must be updated"
 }

 @test "fail to update provider configuration - invalid config format" {
    export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-7 --config "$VALID_CONFIG" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_key_pc update --id "$ID" --config "{invalid: json}"
    assert_failure
    assert_output --partial "Cannot update provider config with invalid json"
 }
 
  @test "delete provider configuration -- success" {
    export CREATED=$(./otdfctl $HOST $WITH_CREDS keymanagement provider create --name test-config-7 --config "$VALID_CONFIG" --json)
    ID=$(echo "$CREATED" | jq -r '.id')
    run_otdfctl_key_pc delete --id "$ID"
    assert_success
  }

  @test "delete provider configuration fail -- no id" {
    run_otdfctl_key_pc delete
    assert_failure
    assert_output --partial "Flag '--id' is required"
  }