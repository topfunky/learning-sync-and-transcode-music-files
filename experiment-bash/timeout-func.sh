#!/bin/bash

# Function to run another function with a timeout
function timeout_func() {
    local timeout=$2
    local function_name=$1

    echo "Running function: $function_name with timeout: $timeout seconds"
    timeout --preserve-status $timeout bash -c "$function_name"
    local status=$?

    if [ $status -eq 124 ]; then
        echo "Function timed out after $timeout seconds"
    elif [ $status -eq 137 ]; then
        echo "Function terminated by SIGKILL"
    else
        echo "Function exited with status $status"
    fi

    return $status
}

# Function to run a command with a timeout and terminate if it takes too long
function old_timeout_func() { 
    cmd="$1"; timeout="$2";
    grep -qP '^\d+$' <<< "$timeout" || timeout=10

    ( 
        eval "$cmd" &
        child=$!
        echo "Running command: '$cmd' (pid $child)"
        trap -- "" SIGTERM 
        (       
                sleep "$timeout"
                # If the child process is still running, kill it
                if ps -p $child > /dev/null; then
                    echo "Terminating child process (pid $child)"
                    kill -SIGTERM $child
                fi                
        ) &
        sleeper=$!
        wait $child
        child_exit_status=$?

        # If the sleep command is still running, kill it
        if ps -p $sleeper > /dev/null; then
            echo "Terminating sleeper (pid $sleeper)"
            # Terminate without triggering `set -e`
            kill -9 $sleeper
        fi

        if [ $child_exit_status -eq 143 ]; then
            echo "Command terminated due to timeout"
        elif [ $child_exit_status -ne 0 ]; then
            echo "Command failed with exit status $child_exit_status"
        else
            echo "Command completed successfully (pid $child)"
        fi

        return 0
    )
}