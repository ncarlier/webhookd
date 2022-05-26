#!/bin/bash

# Usage: http POST :8080/echo msg==hello foo=bar

echo "This is a simple echo hook."

echo "Hook information: name=$hook_name, id=$hook_id, method=$hook_method"

echo "Command result: hostname=`hostname`"

echo "Header variable: User-Agent=$user_agent"

echo "Query parameter: msg=$msg"

echo "Body payload: $1"
