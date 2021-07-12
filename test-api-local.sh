#!/bin/bash

PORT=9999 go run ./server.go >/dev/null 2>&1 &
server_pid=$!

URL=http://localhost:9999 strest
testCode=$?

pkill -P $server_pid
exit $testCode
