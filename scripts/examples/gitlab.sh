#!/bin/sh

# Functions
die() { echo "error: $@" 1>&2 ; exit 1; }
confDie() { echo "error: $@ Check the server configuration!" 1>&2 ; exit 2; }
debug() {
  [ "$debug" = "true" ] && echo "debug: $@"
}

# Validate global configuration
[ -z "$GITLAB_TOKEN" ] && confDie "GITLAB_TOKEN not set."

# Validate Gitlab hook
[ "$x_gitlab_token" != "$GITLAB_TOKEN" ] && die "bad hook token"

# Validate parameters
payload=$1
payload="$(echo "$payload"|tr -d '\n')"
[ -z "$payload" ] && die "missing request payload"
payload_type=$(echo $payload | jq type -r)
[ $? != 0 ] && die "bad body format: expecting JSON"
[ ! $payload_type = "object" ] && die "bad body format: expecting JSON object but having $payload_type"

debug "received payload: $payload"

# Extract values
object_kind=$(echo $payload | jq .object_kind -r)
[ $? != 0 -o "$object_kind" = "null" ] && die "unable to extract 'object_kind' from JSON payload"

# Do something with the payload:
# Here create a simple notification when a push has occured
if [ "$object_kind" = "push" ]
then
  username=$(echo $payload | jq .user_name -r)
  git_ssh_url=$(echo $payload | jq .project.git_ssh_url -r)
  total_commits_count=$(echo $payload | jq .total_commits_count -r)
  echo "notify: $username push $total_commits_count commit(s) on $git_ssh_url"
fi
