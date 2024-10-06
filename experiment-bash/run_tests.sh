#!/bin/bash


# for all files in `test` directory, run the tests
# and print the results

for file in test/*; do
    if [ -f "$file" ]; then
        echo "Running tests in $file"
        source "$file"
    fi
done

# Check exit value of the script
exit_value=$?
if [ $exit_value -eq 0 ]; then
    echo "All tests passed"
else
    echo "Some tests failed"
fi

