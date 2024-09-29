#!/bin/bash

# Read the PID from the daemon.pid file
PID=$(cat daemon.pid)

# Send the termination signal
kill -SIGTERM $PID

echo "Daemon process with PID $PID has been terminated."