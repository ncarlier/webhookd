#!/bin/sh

URL=http://localhost:8080

echo "Test URL: $URL"
echo "Test bad URL"
curl -H "Content-Type: application/json" \
    --data @test.json \
    $URL/bad/action

echo "Test hook"
curl -H "Content-Type: application/json" \
    -H "X-API-Key: test" \
    --data @test.json \
    $URL/test?firstname=obi-wan\&lastname=kenobi
