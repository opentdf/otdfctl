#!/bin/bash

# Check if the required arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <binary_directory> <output_directory>"
    exit 1
fi

binary_dir="$1"
output_dir="$2"

# Create the output directory if it doesn't exist
mkdir -p "$output_dir"

# Iterate over each binary file
for binary_file in "$binary_dir"/*; do
    if [[ $binary_file == *.exe ]]; then
        # If the file is a Windows binary, zip it
        filename=$(basename "$binary_file")
        zip -j "$output_dir/${filename%.exe}.zip" "$binary_file"
    else
        # For other binaries, tar and gzip them
        filename=$(basename "$binary_file")
        tar -czf "$output_dir/${filename}.tar.gz" "$binary_file"
    fi
done