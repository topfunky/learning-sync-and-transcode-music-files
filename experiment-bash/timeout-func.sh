#!/bin/bash



# Function to run a command with a timeout and terminate if it takes too long
function timeout_func() { 
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