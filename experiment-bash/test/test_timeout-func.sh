#!/bin/bash

# Source the timeout_func.sh file
source ./timeout-func.sh

# Function to run a test case
run_test() {
    local description=$1
    local command=$2
    local timeout=$3
    local expected=$4

    echo "Running test: $description"
    timeout_func "$command" "$timeout"
    local result=$?

    if [ "$result" -eq "$expected" ]; then
        echo "Test passed"
    else
        echo "Test failed: expected $expected, got $result"
    fi
    echo
}

# Test cases
run_test "Command should timeout" "sleep 5" 3 124
run_test "Command should complete" "sleep 2" 3 0
run_test "Invalid timeout value" "sleep 2" "invalid" 124

# Add more test cases as needed