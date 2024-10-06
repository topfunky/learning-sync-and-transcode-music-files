#!/bin/bash

# Source the timeout_func.sh file
source ./timeout-func.sh

EXIT_STATUS_PROCESS_TIMEOUT=124
EXIT_STATUS_SUCCESS=0

# A function that exits with an error code if the test fails
function fail() {
    echo "💀 FAIL Test failed: $1"
    exit 1
}

function run_test() {
    local description=$1
    local command=$2
    local timeout=$3
    local expected_duration=$4

    echo "🏁 Running test: $description"
    start_time=$(date +%s)
    timeout_func "$command" "$timeout"
    local result=$?
    end_time=$(date +%s)
    duration=$((end_time - start_time))

    if [ "$duration" -eq "$expected_duration" ]; then
        echo "😎 Test passed: '$description'"
    else
        fail "'$description' expected duration $expected_duration, got $duration"
    fi
}


# Test cases
# run_test "Invalid timeout value" "sleep 2" "invalid" $EXIT_STATUS_PROCESS_TIMEOUT 2
run_test "Command should complete" "sleep 2" 3 2
run_test "Long running command with sufficient timeout" "sleep 5" 10 5
run_test "Long command should timeout" "sleep 5" 3 3
run_test "Zero timeout should timeout" "sleep 2" 0 0