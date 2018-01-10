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

download_url=$base_download_url/webhookd-$os-$arch
bin_target=/usr/local/bin/webhookd

echo "Downloading $download_url to $bin_target ..."
sudo curl -o $bin_target --fail -L $download_url
[ $? != 0 ] && die "Unable download binary for your architecture."

echo "Making $bin_target as executable ..."
sudo chmod +x $bin_target
[ $? != 0 ] && die "Unable to make the binary as executable."

echo "Installation done. Type 'webhookd' to start the server."