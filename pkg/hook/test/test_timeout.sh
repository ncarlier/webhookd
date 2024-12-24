#!/bin/sh

echo "Running timeout test script..."

for i in `seq 5`; do
  sleep .5
  echo "running..."
done

echo "This line should not be executed!"

exit 0
