# webhookd

[![Image size](https://images.microbadger.com/badges/image/ncarlier/webhookd.svg)](https://microbadger.com/images/ncarlier/webhookd)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/webhookd.svg)](https://hub.docker.com/r/ncarlier/webhookd/)

A very simple webhook server to launch shell scripts.

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/webhookd/webhookd
```

**Or** download the binary regarding your architecture:

```bash
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/webhookd/master/install.sh | bash
```

**Or** use Docker:

```bash
$ docker run -d --name=webhookd \
  --env-file .env \
  -v ${PWD}/scripts:/var/opt/webhookd/scripts \
  -p 8080:8080 \
  ncarlier/webhookd
```

Check the provided environment file [.env](.env) for configuration details.

> Note that this image extends `docker:dind` Docker image. Therefore you are
> able to interact with a Docker daemon with yours shell scripts.

## Configuration

You can configure the daemon by:

### Setting environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_LISTEN_ADDR` | `:8080` | HTTP service address |
| `APP_NB_WORKERS` | `2` | The number of workers to start |
| `APP_HOOK_TIMEOUT` | `10` | Hook maximum delay before timeout (in second) |
| `APP_SCRIPTS_DIR` | `./scripts` | Scripts directory |
| `APP_SCRIPTS_GIT_URL` | none | GIT repository that contains scripts (Note: this is only used by the Docker image or by using the Docker entrypoint script) |
| `APP_SCRIPTS_GIT_KEY` | none | GIT SSH private key used to clone the repository (Note: this is only used by the Docker image or by using the Docker entrypoint script) |
| `APP_WORKING_DIR` | `/tmp` (OS temp dir) | Working directory (to store execution logs) |
| `APP_NOTIFIER` | none | Post script notification (`http` or `smtp`) |
| `APP_NOTIFIER_FROM` | none | Sender of the notification |
| `APP_NOTIFIER_TO` | none | Recipient of the notification |
| `APP_HTTP_NOTIFIER_URL` | none | URL of the HTTP notifier |
| `APP_SMTP_NOTIFIER_HOST` | none | Hostname of the SMTP relay |
| `APP_DEBUG` | `false` | Output debug logs |

### Using command parameters:

| Parameter | Default | Description |
|----------|---------|-------------|
| `-l <address> or --listen <address>` | `:8080` | HTTP service address |
| `-d or --debug` | false | Output debug logs |
| `--nb-workers <workers>` | `2` | The number of workers to start |
| `--scripts <dir>` | `./scripts` | Scripts directory |
| `--timeout <timeout>` | `10` | Hook maximum delay before timeout (in second) |

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

You have several way to provide parameters to your webhook script:

- URL query parameters and HTTP headers are converted into environment
  variables.
  Variable names follows "snakecase" naming convention.
  Therefore the name can be altered.

  *ex: `CONTENT-TYPE` will become `content_type`.*

- Body content (text/plain or application/json) is transmit to the script as
  parameter.

*Example:*

The script:

```bash
#!/bin/bash

echo "Query parameter: foo=$foo"
echo "Header parameter: user-agent=$user_agent"
echo "Script parameters: $1"
```

```bash
$ curl --data @test.json http://localhost/echo?foo=bar
data: Hook work request "echo" queued...

data: Running echo script...

data: Query parameter: foo=bar

data: Header parameter: user-agent=curl/7.52.1

data: Script parameter: {"foo": "bar"}

data: done
```

### Webhook timeout configuration

By default a webhook as a timeout of 10 seconds.
This timeout is globally configurable by setting the environment variable:
`APP_HOOK_TIMEOUT` (in seconds).

You can override this global behavior per request by setting the HTTP header:
`X-Hook-Timeout` (in seconds).

*Example:*

```bash
$ curl -XPOST -H "X-Hook-Timeout: 5" http://localhost/echo?foo=bar
```

### Post hook notifications

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


