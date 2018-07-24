#!/bin/sh

URL=http://localhost:8081

echo "Test URL: $URL"

echo "Testing bad request..."
curl -H "Content-Type: application/json" \
    --data @test.json \
    $URL/bad/action

echo "Testing nominal case..."
curl -H "Content-Type: application/json" \
    -H "X-API-Key: test" \
    --data @test.json \
    $URL/test?firstname=obi-wan\&lastname=kenobi

echo "Testing parallel request..."
curl -XPOST $URL/test &
curl -XPOST $URL/test &
curl -XPOST $URL/test &

wait

echo "Done"
