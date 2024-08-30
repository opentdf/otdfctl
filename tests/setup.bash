#!/bin/bash

####
# Make sure we can load BATS dependencies
####
bats_require_minimum_version 1.5.0

if [[ $(which bats) == *"homebrew"* ]]; then
    BATS_LIB_PATH=$(brew --prefix)/lib
fi

# Check if BATS_LIB_PATH environment variable exists
if [ -z "${BATS_LIB_PATH}" ]; then
# Check if bats bin has homebrew in path name
if [[ $(which bats) == *"homebrew"* ]]; then
    BATS_LIB_PATH=$(dirname $(which bats))/../lib
elif [ -d "/usr/lib/bats-support" ]; then
    BATS_LIB_PATH="/usr/lib"
elif [ -d "/usr/local/lib/bats-support" ]; then
    # Check if bats-support exists in /usr/local/lib
    BATS_LIB_PATH="/usr/local/lib"
fi
fi
echo "BATS_LIB_PATH: $BATS_LIB_PATH"
export BATS_LIB_PATH=$BATS_LIB_PATH

####
# Make sure we have a terminal size large enough to test table output
####
set_terminal_size_linux() {
    # Using resize command if available
    if command -v resize &> /dev/null; then
        resize -s 40 200
    else
        export COLUMNS=200
        export LINES=40
    fi
}

set_terminal_size_mac() {
    printf '\e[8;40;200t'
}

set_terminal_size_windows() {
    # Check if running in Git Bash or similar environment
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        # Assuming Git Bash
        printf '\e[8;40;200t'
    else
        # Assuming Windows Command Prompt or PowerShell
        cmd.exe /c "mode con: cols=200 lines=40"
    fi
}

# Detect the OS and set the terminal size appropriately
case "$OSTYPE" in
    linux*)
        set_terminal_size_linux
        ;;
    darwin*)
        set_terminal_size_mac
        ;;
    msys* | cygwin* | win*)
        set_terminal_size_windows
        ;;
    *)
        echo "Unsupported OS: $OSTYPE"
        ;;
esac