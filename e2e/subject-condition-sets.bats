#!/usr/bin/env bats

# Tests for subject condition sets

setup_file() {
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > creds.json
  export WITH_CREDS='--with-client-creds-file ./creds.json'
  export HOST='--host http://localhost:8080'

  export SCS_1='[{"condition_groups":[{"conditions":[{"operator":1,"subject_external_values":["marketing"],"subject_external_selector_value":".org.name"},{"operator":1,"subject_external_values":["ShinyThing"],"subject_external_selector_value":".team.name"}],"boolean_operator":1}]}]'
  export SCS_2='[{"condition_groups":[{"conditions":[{"operator":3,"subject_external_values":["piedpiper.com","hooli.com"],"subject_external_selector_value":".emailAddress"},{"operator":1,"subject_external_values":["sales"],"subject_external_selector_value":".department"}],"boolean_operator":2}]}]'
  export SCS_3='[{"condition_groups":[{"conditions":[{"operator":2,"subject_external_values":["CoolTool","RadService"],"subject_external_selector_value":".team.name"}],"boolean_operator":2}]}]'
}

setup() {
  load "${BATS_LIB_PATH}/bats-support/load.bash"
  load "${BATS_LIB_PATH}/bats-assert/load.bash"

  # invoke binary with credentials
  run_otdfctl_scs () {
    run sh -c "./otdfctl $HOST $WITH_CREDS policy subject-condition-sets $*"
  }

  run_delete_scs () {
     # Capture the first argument as the ID
    local id="$1" 

    run sh -c "./otdfctl $HOST $WITH_CREDS policy scs delete --id $id --force" 
  }
}

teardown_file() {
  # clear out all test env vars
  unset HOST WITH_CREDS NS_NAME NS_ID ATTR_NAME_RANDOM
}

@test "Create a Subject Condition Set (SCS) - from file" {
  echo -n "$SCS_1" > scs.json

  run_otdfctl_scs create --subject-sets-file-json scs.json -l fromfile=true
  assert_success
  assert_output --partial "Id"
  assert_output --partial "SubjectSets"
  assert_output --partial ".org.name"
  assert_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_IN"
  assert_output --regexp "fromfile: true"
}

@test "Create a Subject Condition Set (SCS) - from flag value JSON" {
  run ./otdfctl $HOST $WITH_CREDS policy scs create --subject-sets "$SCS_2"
  assert_success
  assert_output --partial "Id"
  assert_output --partial "SubjectSets"
  assert_output --partial ".emailAddress"
  assert_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_IN"
}

@test "Get a SCS" {
  CREATED_ID=$(./otdfctl $HOST $WITH_CREDS policy scs add -s "$SCS_3" -l hello=world --json | jq -r '.id')
  run_otdfctl_scs get --id "$CREATED_ID"
  assert_success
  assert_output --regexp "Id.*$CREATED_ID"
  assert_output --partial "Labels"
  assert_output --partial "hello: world"
  assert_output --partial "Created At"
  assert_output --partial "Updated At"
  assert_output --partial ".team.name"
  assert_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN"

  run_delete_scs "$CREATED_ID"
}

@test "Update a SCS - from flag value JSON" {
  echo -n "$SCS_1" > scs.json
  CREATED_ID=$(./otdfctl $HOST $WITH_CREDS policy scs create --subject-sets-file-json scs.json -l fromfile=true --json | jq -r '.id')

  run ./otdfctl $HOST $WITH_CREDS policy scs update --subject-sets "$SCS_2" --id "$CREATED_ID"
  assert_success
  assert_output --partial ".emailAddress"
  assert_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_IN"
  assert_output --partial "fromfile: true"
  refute_output --partial ".org.name"

  run_delete_scs "$CREATED_ID"
}

@test "Update a SCS - from file" {
  CREATED_ID=$(./otdfctl $HOST $WITH_CREDS policy scs create --subject-sets "$SCS_2" -l fromfile=false --json | jq -r '.id')

  echo -n "$SCS_3" > scs.json

  run ./otdfctl $HOST $WITH_CREDS policy scs update --subject-sets-file-json scs.json --id "$CREATED_ID" -l fromfile=true
  assert_success
  refute_output --partial ".emailAddress"
  refute_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_IN"
  assert_output --partial ".team.name"
  assert_output --partial "fromfile: true"
  assert_output --partial "SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN"

  run_delete_scs "$CREATED_ID"
}

@test "List SCS" {
  CREATED_ID=$(./otdfctl $HOST $WITH_CREDS policy scs create --subject-sets "$SCS_2" -l fromfile=false --json | jq -r '.id')

  run_otdfctl_scs list
  assert_success
  assert_output --partial "$CREATED_ID"

  run_otdfctl_scs list --json
  assert_success
  assert_output --partial ".department"
  assert_output --partial ".emailAddress"
  assert_output --partial ".team.name"
  assert_output --partial ".org.name"
  [ $(echo "$output" | jq -r '.[-1].subject_sets[0].condition_groups[0].conditions[0].subject_external_values[0]') = "piedpiper.com" ]
  [ $(echo "$output" | jq -r '.[-1].id') = "$CREATED_ID" ]
  [ $(echo "$output" | jq -r '.[-1].metadata.labels.fromfile') = "false" ]

  # validate deletion
  run_delete_scs "$CREATED_ID"
  assert_success
  assert_output --partial "$CREATED_ID"
}
