#!/bin/sh

# Error function
die() { echo "error: $@" 1>&2 ; exit 1; }

if [ ! -z "$APP_SCRIPTS_GIT_URL" ]
then
  [ ! -f "$APP_SCRIPTS_GIT_KEY" ] && die "Git clone key not found."

  export APP_SCRIPTS_DIR=${APP_SCRIPTS_DIR:-/opt/scripts-git}
  export GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

  mkdir -p $APP_SCRIPTS_DIR

  echo "Cloning $APP_SCRIPTS_GIT_URL into $APP_SCRIPTS_DIR ..."
  ssh-agent sh -c 'ssh-add ${APP_SCRIPTS_GIT_KEY}; git clone --depth 1 --single-branch ${APP_SCRIPTS_GIT_URL} ${APP_SCRIPTS_DIR}'
  [ $? != 0 ] && die "Unable to clone repository"
fi

exec "$@"

