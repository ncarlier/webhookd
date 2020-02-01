#!/bin/sh

# Error function
die() { echo "error: $@" 1>&2 ; exit 1; }

if [ ! -z "$WHD_SCRIPTS_GIT_URL" ]
then
  [ ! -f "$WHD_SCRIPTS_GIT_KEY" ] && die "Git clone key not found."

  export WHD_SCRIPTS=${WHD_SCRIPTS:-/opt/scripts-git}
  export GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

  mkdir -p $WHD_SCRIPTS

  echo "Cloning $WHD_SCRIPTS_GIT_URL into $WHD_SCRIPTS ..."
  ssh-agent sh -c 'ssh-add ${WHD_SCRIPTS_GIT_KEY}; git clone --depth 1 --single-branch ${WHD_SCRIPTS_GIT_URL} ${WHD_SCRIPTS}'
  [ $? != 0 ] && die "Unable to clone repository"
fi

exec "$@"

