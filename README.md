# webhookd

A very simple webhook server to launch shell scripts.

It can be used as a cheap alternative of Docker hub in order to build private Docker images.

## Installation

### Binaries

Linux binaries for release [0.0.3](https://github.com/ncarlier/webhookd/releases)

* [amd64](https://github.com/ncarlier/webhookd/releases/download/v0.0.3/webhookd-linux-amd64-v0.0.3.tar.gz)

Download the version you need, untar, and install to your PATH.

```
$ wget https://github.com/ncarlier/webhookd/releases/download/v0.0.3/webhookd-linux-amd64-v0.0.3.tar.gz
$ tar xvzf webhookd-linux-amd64-v0.0.3.tar.gz
$ ./webhookd
```

### Docker

Start the container mounting your scripts directory:

```
$ docker run -d --name=webhookd \
  --env-file etc/env.conf \
  -v ${PWD}/scripts:/var/opt/webhookd/scripts \
  -p 8080:8080 \
  ncarlier/webhookd
```

The provided environment file (`etc/env.conf`) is used to configure the app.
Check [sample configuration](etc/env_sample.com) for details.

## Usage

Create your own scripts template in the **scripts** directory.

Respect the following structure:

```
/scripts
|--> /bitbucket
  |--> /script_1.sh
  |--> /script_2.sh
|--> /github
|--> /gitlab
|--> /docker
```

The hookname you will use will be related to the hook you want to use (github, bitbucket, ...) and the script name you want to call:
For instance if you are **gitlab** and want to call **build.sh** then you will need to use:

```
http://webhook_ip:port/gitlab/build
```

It is important to use the right hook in order for your script to received parameters extract from the hook payload.


For now, supported hooks are:

- GitHub
- Gitlab
- Bitbucket
- Docker Hub


Check the scripts directory for samples.

Once the action script created, you can trigger the webhook :

```
$ curl -H "Content-Type: application/json" \
  --data @payload.json \
  http://localhost:8080/<hookname>/<action>
```

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


