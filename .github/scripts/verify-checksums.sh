#!/bin/bash

# Check if the required arguments are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <output_directory> <checksum_file>"
    exit 1
fi

# Assign arguments to variables
output_dir="$1"
checksum_file="$2"
checksum_path="${output_dir}/${checksum_file}"  # Full path to the checksum file
lock_file="${checksum_path}.lock"  # Append .lock to the full path of the checksum file

# Ensure the checksum file exists
if [ ! -f "$checksum_path" ]; then
    echo "ERROR: Checksum file $checksum_path does not exist."
    exit 1
fi

# Wait for the lock file to be available for reading
exec 200<"$lock_file"  # Open lock file descriptor for reading
flock -s 200           # Acquire shared lock (will wait if exclusive lock is held)

echo "Verifying checksums..."
echo "Looking for checksum file: $checksum_path"

# Iterate over each line in the checksum file
while read -r line; do
    # Extract checksum and filename from the line
    expected_checksum=$(echo "$line" | awk '{print $1}')
    filename=$(echo "$line" | awk '{print $2}')

    # Construct the full path to the file
    file_path="$output_dir/$filename"

    # Check if the file exists
    if [ ! -f "$file_path" ]; then
        echo "ERROR: File $filename not found in $output_dir"
        continue
    fi

    # Calculate the actual checksum of the file
    actual_checksum=$(shasum -a 256 "$file_path" | awk '{print $1}')

    # Compare the expected and actual checksums
    if [ "$expected_checksum" != "$actual_checksum" ]; then
        echo "ERROR: Checksum for $filename does not match."
    else
        echo "Checksum for $filename is correct."
    fi
done < "$checksum_path"

# Release the lock and close the lock file descriptor
flock -u 200
exec 200>&-

# Clean up the lock file
rm -f "$lock_file"

echo "Checksum verification completed."
