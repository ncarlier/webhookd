#!/bin/bash

# Usage: http POST :8080/echo msg==hello foo=bar

echo "Echo script:"

echo "Command result: hostname=`hostname`"

echo "Header variable: User-Agent=$user_agent"

echo "Query parameter: msg=$msg"

echo "Body payload: $1"
