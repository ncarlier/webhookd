#!/bin/sh

IP=`sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' webhookd`

echo "Test bad URL"
curl -H "Content-Type: application/json" \
    --data @assets/bitbucket.json \
    http://$IP:8080/bad/action

echo "Test Bitbucket hook"
curl -H "Content-Type: application/json" \
    --data @assets/bitbucket.json \
    http://$IP:8080/bitbucket/echo

echo "Test Bitbucket hook"
curl -H "Content-Type: application/x-www-form-urlencoded" \
    --data @assets/bitbucket.raw \
    http://$IP:8080/bitbucket/echo

echo "Test Github hook"
curl -H "Content-Type: application/json" \
    --data @assets/github.json \
    http://$IP:8080/github/echo

echo "Test Docker hook"
curl -H "Content-Type: application/json" \
    --data @assets/docker.json \
    http://$IP:8080/docker/echo

