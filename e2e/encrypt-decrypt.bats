#!/usr/bin/env bats

# Tests for encrypt decrypt

setup() {
    load "${BATS_LIB_PATH}/bats-support/load.bash"
    load "${BATS_LIB_PATH}/bats-assert/load.bash"
}

setup_file() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export DEBUG_LEVEL="--log-level debug"
  export HOST=http://localhost:8080

  export INFILE_GO_MOD=go.mod
  export OUTFILE_GO_MOD=go.mod.tdf
  export RESULTFILE_GO_MOD=result.mod

  export SECRET_TEXT="my special secret"
  export OUT_TXT=secret.txt
  export OUTFILE_TXT=secret.txt.tdf

  NS_ID=$(./otdfctl --host $HOST $WITH_CREDS $DEBUG_LEVEL policy attributes namespaces create -n "testing-enc-dec.io" --json | jq -r '.id')
  ATTR_ID=$(./otdfctl --host $HOST $WITH_CREDS $DEBUG_LEVEL policy attributes create --namespace "$NS_ID" -n attr1 -r ALL_OF --json | jq -r '.id')
  VAL_ID=$(./otdfctl --host $HOST $WITH_CREDS $DEBUG_LEVEL policy attributes values create --attribute-id "$ATTR_ID" -v value1 --json | jq -r '.id')
  # entitles opentdf client id for client credentials CLI user
  SCS='[{"condition_groups":[{"conditions":[{"operator":1,"subject_external_values":["opentdf"],"subject_external_selector_value":".clientId"}],"boolean_operator":2}]}]'
  HS256_KEY=$(openssl rand -base64 32)
  openssl genpkey -algorithm RSA -out rs_private_key.pem -pkeyopt rsa_keygen_bits:2048
  openssl rsa -pubout -in rs_private_key.pem -out rs_public_key.pem
  RS256_PRIVATE_KEY=$(awk '{printf "%s\\n", $0}' rs_private_key.pem)
  RS256_PUBLIC_KEY=$(awk '{printf "%s\\n", $0}' rs_public_key.pem)
  ASSERTIONS='[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"}}]'
  SIGNED_ASSERTIONS_HS256='[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"},"signingKey":{"alg":"HS256","key":"'$HS256_KEY'"}}]'
  SIGNED_ASSERTION_VERIFICATON_HS256='{"keys":{"assertion1":{"alg":"HS256","key":"'$HS256_KEY'"}}}'
  SIGNED_ASSERTIONS_RS256='[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"},"signingKey":{"alg":"RS256","key":"'$RS256_PRIVATE_KEY'"}}]'
  SIGNED_ASSERTION_VERIFICATON_RS256='{"keys":{"assertion1":{"alg":"RS256","key":"'$RS256_PUBLIC_KEY'"}}}'
  SM=$(./otdfctl --host $HOST $WITH_CREDS $DEBUG_LEVEL policy subject-mappings create --action-standard DECRYPT -a "$VAL_ID" --subject-condition-set-new "$SCS")
  export FQN="https://testing-enc-dec.io/attr/attr1/value/value1"
  export MIXED_CASE_FQN="https://Testing-Enc-Dec.io/attr/Attr1/value/VALUE1"
}

teardown() {
    rm -f $OUTFILE_GO_MOD $RESULTFILE_GO_MOD $OUTFILE_TXT
}

teardown_file(){
    ./otdfctl --host "$HOST" $WITH_CREDS policy attributes namespaces unsafe delete --id "$NS_ID" --force
}

@test "roundtrip TDF3, no attributes, file" {
  ./otdfctl encrypt -o $OUTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --tdf-type tdf3 $INFILE_GO_MOD
  ./otdfctl decrypt -o $RESULTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --tdf-type tdf3 $OUTFILE_GO_MOD
  diff $INFILE_GO_MOD $RESULTFILE_GO_MOD
}

@test "roundtrip TDF3, one attribute, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN
  ./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS $OUTFILE_TXT | grep "$SECRET_TEXT"
}

@test "roundtrip TDF3, one attribute, mixed case FQN, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $MIXED_CASE_FQN
  ./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS $OUTFILE_TXT | grep "$SECRET_TEXT"
}

@test "roundtrip TDF3, assertions, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN --with-assertions "$ASSERTIONS"
  ./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS $OUTFILE_TXT | grep "$SECRET_TEXT"
}

@test "roundtrip TDF3, assertions with HS265 keys and verificaion, stdin" {
  run bash -c echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN --with-assertions "$SIGNED_ASSERTIONS_HS256"
  assert_success
  run ./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --with-assertion-verification-keys "$SIGNED_ASSERTION_VERIFICATON_HS256" $OUTFILE_TXT
  assert_success
  assert_output --partial "$SECRET_TEXT"
}

@test "roundtrip TDF3, assertions with RS256 keys and verificaion, stdin" {
  result=$(echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN --with-assertions "$SIGNED_ASSERTIONS_RS256")
  assert_success
  result=$(./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --with-assertion-verification-keys "$SIGNED_ASSERTION_VERIFICATON_RS256" $OUTFILE_TXT)
  assert_success
  assert_contains "$result" "$SECRET_TEXT"
}

@test "roundtrip NANO, no attributes, file" {
  ./otdfctl encrypt -o $OUTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --tdf-type nano $INFILE_GO_MOD
  ./otdfctl decrypt -o $RESULTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --tdf-type nano $OUTFILE_GO_MOD
  diff $INFILE_GO_MOD $RESULTFILE_GO_MOD
}

@test "roundtrip NANO, no attributes, file, ecdsa binding" {
  ./otdfctl encrypt -o $OUTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --ecdsa-binding --tdf-type nano $INFILE_GO_MOD
  ./otdfctl decrypt -o $RESULTFILE_GO_MOD --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --tdf-type nano $OUTFILE_GO_MOD
  diff $INFILE_GO_MOD $RESULTFILE_GO_MOD
  ./otdfctl --host $HOST --tls-no-verify $WITH_CREDS inspect $OUTFILE_GO_MOD
  ecdsa_enabled="$(./otdfctl --host $HOST --tls-no-verify $WITH_CREDS inspect $OUTFILE_GO_MOD | jq .ecdsaEnabled)"
  [[ "$ecdsa_enabled" == true ]]
}

@test "roundtrip NANO, one attribute, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt --tdf-type nano -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN
  ./otdfctl decrypt --tdf-type nano --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS $OUTFILE_TXT | grep "$SECRET_TEXT"
}

@test "roundtrip NANO, one attribute, mixed case FQN, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt --tdf-type nano -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $MIXED_CASE_FQN
  ./otdfctl decrypt --tdf-type nano --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS $OUTFILE_TXT | grep "$SECRET_TEXT"
}
