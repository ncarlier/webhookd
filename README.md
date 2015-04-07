webhookd
=========

A very simple webhook server to launch shell scripts.

It can be used as a cheap alternative of Docker hub in order to build private Docker images.

Installation
------------

Binaries
------

Linux binaries for release [0.0.1](https://github.com/ncarlier/webhookd/releases)

* [amd64](https://github.com/ncarlier/webhookd/releases/download/v0.0.1/webhookd-linux-amd64-v0.0.1.tar.gz)

Download the version you need, untar, and install to your PATH.

    $ wget https://github.com/ncarlier/webhookd/releases/download/v0.0.1/webhookd-linux-amd64-v0.0.1.tar.gz
    $ tar xvzf webhookd-linux-amd64-v0.0.1.tar.gz
    $ ./webhookd

Docker
----

Use the following make command to start docker containers:

- **make build** will build the webhookd image
- **make volume** (optionnal) will create volume that will contain the scripts *folder* inside the container (usefull for dev)
- **make run** will run a container from the freshly build image. Optionaly, you can use:
    - **make dev run** to use the volume container (dev)
    - **make key dev run** to mount the volume and link the ssh folder containing technical ssh key


Usage
-------

Create your own scripts template in the **scripts** directory.

Respect the following structure:

    /scripts
    |--> /bitbucket
      |--> /script_1.sh
      |--> /script_2.sh
    |--> /github
    |--> /gitlab
    |--> /docker

The hookname you will use will be related to the hook you want to use (github, bitbucket, ...) and the script name you want to call:
For instance if you are **gitlab** and want to call **build.sh** then you will need to use:

    http://webhook_ip:port/gitlab/build

It is important to use the right hook in order for your script to received parameters extract from the hook payload.


For now, supported hooks are:

- GitHub
- Gitlab
- Bitbucket
- Docker Hub


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


