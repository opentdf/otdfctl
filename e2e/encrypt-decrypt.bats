#!/usr/bin/env bats

# Tests for encrypt decrypt

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
  ASSERTIONS='[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"}}]'
  SIGNED_ASSERTIONS='[{"id":"assertion1","type":"handling","scope":"tdo","appliesToState":"encrypted","statement":{"format":"json+stanag5636","schema":"urn:nato:stanag:5636:A:1:elements:json","value":"{\"ocl\":\"2024-10-21T20:47:36Z\"}"},"signingKey":{"alg":"HS256","key":"'$HS256_KEY'"}}]'
  SIGNED_ASSERTION_VERIFICATON='{"keys":{"assertion1":{"alg":"HS256","key":"'$HS256_KEY'"}}}'
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

@test "roundtrip TDF3, assertions with keys and verificaion, stdin" {
  echo $SECRET_TEXT | ./otdfctl encrypt -o $OUT_TXT --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS -a $FQN --with-assertions "$SIGNED_ASSERTIONS"
  ./otdfctl decrypt --host $HOST --tls-no-verify $DEBUG_LEVEL $WITH_CREDS --with-assertion-verification-keys $SIGNED_ASSERTION_VERIFICATON $OUTFILE_TXT | grep "$SECRET_TEXT"
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
