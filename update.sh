for dir in */ ; do
    # Check if a go.mod file exists in the directory
    if [[ -f "${dir}go.mod" ]]; then
        echo "Found go.mod in ${dir}, running 'go get -u all'"
        # Navigate into the directory
        cd "$dir" || exit

        go get -u ./...
        # Navigate back to the parent directory
        cd ..
    else
        echo "No go.mod found in ${dir}, skipping"
    fi
done
