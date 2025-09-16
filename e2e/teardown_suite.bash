#!/bin/bash

####
# Remove the creds file if it exists
####

setup_suite(){
    rm -f ./bats_creds.json
}