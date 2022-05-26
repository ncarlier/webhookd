#!/bin/bash


echo "Starting background job..."

nohup ./scripts/long.sh >/tmp/long.log 2>&1  &

echo "Background job started."


