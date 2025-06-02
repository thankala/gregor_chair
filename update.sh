for dir in */ ; do
    # Check if a go.mod file exists in the directory
    if [[ -f "${dir}go.mod" ]]; then
        echo "Found go.mod in ${dir}, running 'go get -u all'"
        # Navigate into the directory
        cd "$dir" || exit

        # Get the current branch name
        branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)

        # If git fails (not a git repo), default to master or skip
        if [[ -z "$branch" || "$branch" == "HEAD" ]]; then
            echo "Not a git repository or detached HEAD in ${dir}, skipping 'go get' for branch."
        else
            # Run 'go get' with the current branch
            echo "Using branch: $branch"
            go get -v -u "github.com/thankala/gregor_chair_common"
        fi

        go get -u all
        # Navigate back to the parent directory
        cd ..
    else
        echo "No go.mod found in ${dir}, skipping"
    fi
done
