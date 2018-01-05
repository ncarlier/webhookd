# webhookd

[![Image size](https://img.shields.io/imagelayers/image-size/ncarlier/webhookd/latest.svg)](https://hub.docker.com/r/ncarlier/webhookd/)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/webhookd.svg)](https://hub.docker.com/r/ncarlier/webhookd/)

A very simple webhook server to launch shell scripts.

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/webhookd/webhookd
```

**Or** download the binary regarding your architecture:

```bash
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/webhookd/master/install.sh | sh
```

**Or** use Docker:

```bash
$ docker run -d --name=webhookd \
  --env-file .env \
  -v ${PWD}/scripts:/var/opt/webhookd/scripts \
  -p 8080:8080 \
  ncarlier/webhookd
```

Check the provided environment file [.env](.env) for details.

> Note that this image extends `docker:dind` Docker image. Therefore you are
> able to interact with a Docker daemon with yours shell scripts.

## Usage

### Directory structure

Webhooks are simple scripts dispatched into a directory structure.

By default inside the `./scripts` directory.
You can override the default using the `APP_SCRIPTS_DIR` environment variable.

*Example:*

```
/scripts
|--> /github
  |--> /build.sh
  |--> /deploy.sh
|--> /ping.sh
|--> ...
```

### Webhook URL

The directory structure define the webhook URL.
The Webhook can only be call with HTTP POST verb.
If the script exists, the HTTP response will be a `text/event-stream` content
type (Server-sent events).

*Example:*

The script: `./scripts/foo/bar.sh`

```bash
#!/bin/bash

echo "foo foo foo"
echo "bar bar bar"
```

```bash
$ curl -XPOST http://localhost/foo/bar
data: Hook work request "foo/bar" queued...

data: Running foo/bar script...

data: foo foo foo

data: bar bar bar

data: done
```

### Webhook parameters

You can add query parameters to the webhook URL.
Those parameters will be available as environment variables into the shell
script.
You can also send a payload (text/plain or application/json) as request body.
This payload will be transmit to the shell script as first parameter.

*Example:*

The script:

```bash
#!/bin/bash

echo "Environment parameters: foo=$foo"
echo "Script parameters: $1"
```

```bash
$ curl --data @test.json http://localhost/echo?foo=bar
data: Hook work request "echo" queued...

data: Running echo script...

data: Environment parameters: foo=bar

data: Script parameters: {"foo": "bar"}

data: done
```

### Notifications

The script's output is collected and stored into a log file (configured by the
`APP_WORKING_DIR` environment variable).

Once the script executed, you can send the result and this log file to a
notification channel. Currently only two channels are supported: Email and HTTP.

#### HTTP notification

HTTP notification configuration:

- **APP_NOTIFIER**=http
- **APP_NOTIFIER_FROM**=webhookd <noreply@nunux.org>
- **APP_NOTIFIER_TO**=hostmaster@nunux.org
- **APP_HTTP_NOTIFIER_URL**=http://requestb.in/v9b229v9

> Note that the HTTP notification is compatible with
[Mailgun](https://mailgun.com) API.

#### Email notification

SMTP notification configuration:

- **APP_NOTIFIER**=smtp
- **APP_SMTP_NOTIFIER_HOST**=localhost:25

The log file will be sent as an GZIP attachment.

---


