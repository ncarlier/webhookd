#!/bin/bash

echo "Running long script..."

for i in {1..20}; do
  sleep 1
  echo "running ${i} ..."
done

echo "Long script end"

exit 0
