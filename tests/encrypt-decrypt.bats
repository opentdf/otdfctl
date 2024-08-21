#!/usr/bin/env bats

# Tests for encrypt decrypt

setup() {
  export CREDSFILE=creds.json
  echo -n '{"clientId":"opentdf","clientSecret":"secret"}' > $CREDSFILE
  export WITH_CREDS="--with-client-creds-file $CREDSFILE"
  export HOST=http://localhost:8080

  export INFILE=go.mod
  export OUTFILE=go.mod.tdf
  export RESULTFILE=result.mod
}

teardown() {
    rm -f $OUTFILE $RESULTFILE $CREDSFILE
}

@test "roundtrip TDF3" {
  ./otdfctl encrypt -o $OUTFILE --host $HOST --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 $INFILE
  ./otdfctl decrypt -o $RESULTFILE --host $HOST --tls-no-verify --log-level debug $WITH_CREDS --tdf-type tdf3 $OUTFILE
  diff $INFILE $RESULTFILE
}

@test "roundtrip NANO" {
  ./otdfctl encrypt -o $OUTFILE --host $HOST --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano $INFILE
  ./otdfctl decrypt -o $RESULTFILE --host $HOST --tls-no-verify --log-level debug $WITH_CREDS --tdf-type nano $OUTFILE
  diff $INFILE $RESULTFILE
}

# Future Tests

# Encrypt and decrypt with attributes:
# Create an attribute and a subject mapping for the specific clientId then roundtrip trip w it
