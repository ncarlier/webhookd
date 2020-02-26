#!/bin/bash

echo "Running simple test script..."

echo "Testing parameters..."
[ -z "$name" ] && echo "Name variable undefined" && exit 1
[ -z "$user_agent" ] && echo "User-Agent variable undefined" && exit 1
[ "$user_agent" != "test" ] && echo "Invalid User-Agent variable: $user_agent" && exit 1

echo "Testing payload..."
[ -z "$1" ] && echo "Payload undefined" && exit 1
[ "$1" != "{\"foo\": \"bar\"}" ] && echo "Invalid payload: $1" && exit 1

echo "notify: OK"

exit 0
