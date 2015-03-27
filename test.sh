#!/bin/sh

IP=`sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' webhookd`
PORT=${1:-8080}

echo "Test URL: http://$IP:$PORT"
echo "Test bad URL"
curl -H "Content-Type: application/json" \
    --data @assets/bitbucket.json \
    http://$IP:$PORT/bad/action

echo "Test Bitbucket hook"
curl -H "Content-Type: application/json" \
    --data @assets/bitbucket.json \
    http://$IP:$PORT/bitbucket/echo

echo "Test Bitbucket hook"
curl -H "Content-Type: application/x-www-form-urlencoded" \
    --data @assets/bitbucket.raw \
    http://$IP:$PORT/bitbucket/echo

echo "Test Github hook"
curl -H "Content-Type: application/json" \
    --data @assets/github.json \
    http://$IP:$PORT/github/echo

echo "Test Gitlab hook"
curl -H "Content-Type: application/json" \
    --data @assets/gitlab.json \
    http://$IP:$PORT/gitlab/echo

echo "Test Docker hook"
curl -H "Content-Type: application/json" \
    --data @assets/docker.json \
    http://$IP:$PORT/docker/echo

