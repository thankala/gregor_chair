#!/bin/bash

# Iterate over all child directories of the current directory
for dir in */ ; do
    # Check if a go.mod file exists in the directory
    if [[ -f "${dir}go.mod" ]]; then
        echo "Found go.mod in ${dir}, running 'go mod tidy'"
        # Navigate into the directory
        cd "$dir" || exit
        # Run 'go mod tidy'
        go mod tidy
        # Navigate back to the parent directory
        cd ..
    else
        echo "No go.mod found in ${dir}, skipping"
    fi
done
