#!/bin/sh

release_url="https://api.github.com/repos/ncarlier/webhookd/releases/latest"
download_url=`curl -s $release_url | grep browser_download_url | head -n 1 | cut -d '"' -f 4`

sudo curl -o /usr/local/bin/webhookd -L $download_url
sudo chmod +x /usr/local/bin/webhookd

