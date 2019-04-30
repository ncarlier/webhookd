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
bin_target=/usr/local/bin

echo "Downloading $download_url to $download_file ..."
sudo curl -o $download_file --fail -L $download_url
[ $? != 0 ] && die "Unable to download binary for your architecture."

echo "Extracting $download_file ..."
sudo tar xvzf ${download_file} -C /tmp/
[ $? != 0 ] && die "Unable to extract archive."

echo "Moving binary to $bin_target ..."
sudo mv /tmp/webhookd* $bin_target
[ $? != 0 ] && die "Unable to move binary."

echo "Making $bin_target as executable ..."
sudo chmod +x $bin_target/webhookd
[ $? != 0 ] && die "Unable to make the binary as executable."

echo "Installation done. Type 'webhookd' to start the server."
