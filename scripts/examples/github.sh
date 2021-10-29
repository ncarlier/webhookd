#!/bin/sh

# Functions
die() { echo "error: $@" 1>&2 ; exit 1; }
confDie() { echo "error: $@ Check the server configuration!" 1>&2 ; exit 2; }
debug() {
  [ "$debug" = "true" ] && echo "debug: $@"
}

# Validate global configuration
[ -z "$GITHUB_SECRET" ] && confDie "GITHUB_SECRET not set."

# Validate Github hook
signature=$(echo -n "$1" | openssl sha1 -hmac "$GITHUB_SECRET" | sed -e 's/^.* //')
[ "sha1=$signature" != "$x_hub_signature" ] && die "bad hook signature: expecting $x_hub_signature and got $signature"

# Validate parameters
payload=$1
[ -z "$payload" ] && die "missing request payload"
payload_type=$(echo $payload | jq type -r)
[ $? != 0 ] && die "bad body format: expecting JSON"
[ ! $payload_type = "object" ] && die "bad body format: expecting JSON object but having $payload_type"

debug "received payload: $payload"

# Extract values
action=$(echo $payload | jq .action -r)
[ $? != 0 -o "$action" = "null" ] && die "unable to extract 'action' from JSON payload"

# Do something with the payload:
# Here create a simple notification when an issue has been published
if [ "$action" = "opened" ]
then
  issue_url=$(echo $payload | jq .issue.url -r)
  sender=$(echo $payload | jq .sender.login -r)
  echo "notify: New issue from $sender: $issue_url"
fi
