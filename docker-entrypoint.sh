#!/bin/sh

# Error function
die() { echo "error: $@" 1>&2 ; exit 1; }

# When running as root, install any extra packages (this needs root and happens
# once, at startup, before any webhook is served) then permanently drop to the
# unprivileged webhookd user so the script clone and the server never run as root.
if [ "$(id -u)" = "0" ]
then
  if [ ! -z "$WHD_EXTRA_PACKAGES" ]
  then
    # Accept a space or comma separated list of packages
    packages=$(echo "$WHD_EXTRA_PACKAGES" | tr ',' ' ')

    echo "Installing extra packages: $packages ..."
    apk add --no-cache $packages
    [ $? != 0 ] && die "Unable to install extra packages"
  fi

  # Re-exec as webhookd; unset so the second pass doesn't warn about the now-done install
  unset WHD_EXTRA_PACKAGES
  exec su-exec webhookd "$0" "$@"
elif [ ! -z "$WHD_EXTRA_PACKAGES" ]
then
  echo "warning: WHD_EXTRA_PACKAGES is set but not running as root; skipping package installation" 1>&2
fi

if [ ! -z "$WHD_SCRIPTS_GIT_URL" ]
then
  [ ! -f "$WHD_SCRIPTS_GIT_KEY" ] && die "Git clone key not found."

  export WHD_HOOK_SCRIPTS=${WHD_HOOK_SCRIPTS:-/opt/scripts-git}
  export GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

  mkdir -p $WHD_HOOK_SCRIPTS

  echo "Cloning $WHD_SCRIPTS_GIT_URL into $WHD_HOOK_SCRIPTS ..."
  ssh-agent sh -c 'ssh-add ${WHD_SCRIPTS_GIT_KEY}; git clone --depth 1 --single-branch ${WHD_SCRIPTS_GIT_URL} ${WHD_HOOK_SCRIPTS}'
  [ $? != 0 ] && die "Unable to clone repository"
fi

exec "$@"

