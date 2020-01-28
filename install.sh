#!/bin/bash

die() { echo "error: $@" 1>&2 ; exit 1; }

# Getting operating system
os=`uname -s`
os=${os,,}

# Getting architecture
arch=`uname -m`
case "$arch" in
"armv7l")
    arch="arm"
    ;;
"x86_64")
    arch="amd64"
    ;;
esac

release_url="https://api.github.com/repos/ncarlier/webhookd/releases/latest"
artefact_url=`curl -s $release_url | grep browser_download_url | head -n 1 | cut -d '"' -f 4`
[ -z "$artefact_url" ] && die "Unable to extract artefact URL"
base_download_url=`dirname $artefact_url`

download_url=$base_download_url/webhookd-$os-${arch}.tgz
download_file=/tmp/webhookd-$os-${arch}.tgz
bin_target=${1:-$HOME/.local/bin}

echo "Downloading $download_url to $download_file ..."
curl -o $download_file --fail -L $download_url
[ $? != 0 ] && die "Unable to download binary for your architecture."

echo "Extracting $download_file to $bin_target ..."
[ -d $bin_target ] || mkdir -p $bin_target
tar xvzf ${download_file} -C $bin_target
[ $? != 0 ] && die "Unable to extract archive."

echo "Cleaning..."
rm $download_file \
   $bin_target/LICENSE \
   $bin_target/README.md \
   $bin_target/CHANGELOG.md
[ $? != 0 ] && die "Unable to clean installation files."

echo "Installation done. Type '$bin_target/webhookd' to start the server."
