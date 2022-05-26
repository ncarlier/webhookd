#!/bin/bash

echo "Running simple test script..."

echo "Testing parameters..."
[ -z "$name" ] && echo "Name variable undefined" && exit 1
[ -z "$user_agent" ] && echo "User-Agent variable undefined" && exit 1
[ "$user_agent" != "test" ] && echo "Invalid User-Agent variable: $user_agent" && exit 1
[ -z "$hook_id" ] && echo "Hook ID variable undefined" && exit 1
[ "$hook_name" != "test_simple" ] && echo "Invalid hook name variable: $hook_name" && exit 1
[ "$hook_method" != "GET" ] && echo "Invalid hook method variable: $hook_method" && exit 1

echo "Testing payload..."
[ -z "$1" ] && echo "Payload undefined" && exit 1
[ "$1" != "{\"foo\": \"bar\"}" ] && echo "Invalid payload: $1" && exit 1

echo "notify: OK"

exit 0
