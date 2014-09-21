webhookd
=========

A very simple webhook server to launch shell scripts.

Installation
------------

Linux binaries for release [0.0.1](https://github.com/ncarlier/webhookd/releases)

* [amd64](https://github.com/ncarlier/webhookd/releases/download/v0.0.1/webhookd-linux-amd64-v0.0.1.tar.gz)

Download the version you need, untar, and install to your PATH.

    $ wget https://github.com/ncarlier/webhookd/releases/download/v0.0.1/webhookd-linux-amd64-v0.0.1.tar.gz
    $ tar xvzf webhookd-linux-amd64-v0.0.1.tar.gz
    $ ./webhookd

Usage
-------

Create your own scripts template in the **scripts** directory.

Respect the following structure:

    /scripts
    |--> /bitbucket
      |--> /echo.sh
      |--> /build.sh
    |--> /github
    |--> /docker

The directory name right under the **scripts** directory defined the hookname.

For now, supported hooks are:

- GitHub
- Bitbucket
- Docker Hub

The scripts under the **hook** directory defined the actions.

The action script take parameters. These parameters are extract from the payload of the hook. For instance the GitHub hook extract the repository URL and name. Then pass them by parameter to the action script.

Check the scripts directory for samples.

Once the action script created, you can trigger the webhook :

    $ curl -H "Content-Type: application/json" \
        --data @payload.json \
        http://localhost:8080/<hookname>/<action>

The action script's output is collected and sent by email or by HTTP request.

The HTTP notification need some configuration:

- **APP_NOTIFIER**=http
- **APP_NOTIFIER_FROM**=webhookd <noreply@nunux.org>
- **APP_NOTIFIER_TO**=hostmaster@nunux.org
- **APP_HTTP_NOTIFIER_URL**=http://requestb.in/v9b229v9

> Note that the HTTP notification is compatible with [Mailgun](https://mailgun.com) API.

As the smtp notification:

- **APP_NOTIFIER**=smtp
- **APP_SMTP_NOTIFIER_HOST**=localhost:25


