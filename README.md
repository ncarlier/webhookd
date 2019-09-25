# webhookd

[![Build Status](https://travis-ci.org/ncarlier/webhookd.svg?branch=master)](https://travis-ci.org/ncarlier/webhookd)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncarlier/webhookd)](https://goreportcard.com/report/github.com/ncarlier/webhookd)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/webhookd.svg)](https://hub.docker.com/r/ncarlier/webhookd/)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/nunux)

A very simple webhook server to launch shell scripts.

![Logo](webhookd.svg)

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/webhookd
```

**Or** download the binary regarding your architecture:

```bash
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/webhookd/master/install.sh | bash
```

**Or** use Docker:

```bash
$ docker run -d --name=webhookd \
  -v ${PWD}/scripts:/var/opt/webhookd/scripts \
  -p 8080:8080 \
  ncarlier/webhookd \
  webhookd --scripts=/var/opt/webhookd/scripts
```

> Note that this image extends `docker:dind` Docker image.
> Therefore you are able to interact with a Docker daemon with yours shell scripts.

## Configuration

You can configure the daemon by:

### Setting environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_LISTEN_ADDR` | `:8080` | HTTP service address |
| `APP_PASSWD_FILE` | `.htpasswd` | Password file for HTTP basic authentication |
| `APP_NB_WORKERS` | `2` | The number of workers to start |
| `APP_HOOK_TIMEOUT` | `10` | Hook maximum delay before timeout (in second) |
| `APP_SCRIPTS_DIR` | `./scripts` | Scripts directory |
| `APP_SCRIPTS_GIT_URL` | none | GIT repository that contains scripts (Note: this is only used by the Docker image or by using the Docker entrypoint script) |
| `APP_SCRIPTS_GIT_KEY` | none | GIT SSH private key used to clone the repository (Note: this is only used by the Docker image or by using the Docker entrypoint script) |
| `APP_LOG_DIR` | `/tmp` (OS temp dir) | Directory to store execution logs |
| `APP_NOTIFICATION_URI` | none | Notification configuration URI |
| `APP_DEBUG` | `false` | Output debug logs |

### Using command parameters:

| Parameter | Default | Description |
|----------|---------|-------------|
| `-l <address> or --listen <address>` | `:8080` | HTTP service address |
| `-p or --passwd <htpasswd file>` | `.htpasswd` | Password file for HTTP basic authentication
| `-d or --debug` | false | Output debug logs |
| `--nb-workers <workers>` | `2` | The number of workers to start |
| `--scripts <dir>` | `./scripts` | Scripts directory |
| `--timeout <timeout>` | `10` | Hook maximum delay before timeout (in second) |
| `--notification-uri <uri>` |  | Notification configuration URI |
| `--log-dir <dir>` | `/tmp` | Directory to store execution logs |

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

If the script exists, the output the will be streamed to the HTTP response.

The streaming technology depends on the HTTP method used.
With `POST` the response will be chunked.
With `GET` the response will use [Server-sent events][sse].

[sse]: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events

*Example:*

The script: `./scripts/foo/bar.sh`

```bash
#!/bin/bash

echo "foo foo foo"
echo "bar bar bar"
```

Output using `POST` (`Chunked transfer encoding`):

```bash
$ curl -v -XPOST http://localhost:8080/foo/bar
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Transfer-Encoding: chunked
< X-Hook-Id: 7
foo foo foo
bar bar bar
done
```

Output using  `GET` (`Server-sent events`):

```bash
$ curl -v -XGET http://localhost:8080/foo/bar
< HTTP/1.1 200 OK
< Content-Type: text/event-stream
< Transfer-Encoding: chunked
< X-Hook-Id: 8
data: foo foo foo

data: bar bar bar

data: done
```

### Webhook parameters

You have several way to provide parameters to your webhook script:

- URL query parameters and HTTP headers are converted into environment variables.
  Variable names follows "snakecase" naming convention.
  Therefore the name can be altered.

  *ex: `CONTENT-TYPE` will become `content_type`.*

- When using `POST`, body content (text/plain or application/json) is transmit to the script as parameter.

*Example:*

The script:

```bash
#!/bin/bash

echo "Query parameter: foo=$foo"
echo "Header parameter: user-agent=$user_agent"
echo "Script parameters: $1"
```

The result:

```bash
$ curl --data @test.json http://localhost:8080/echo?foo=bar
Query parameter: foo=bar
Header parameter: user-agent=curl/7.52.1
Script parameter: {"foo": "bar"}
done
```

### Webhook timeout configuration

By default a webhook has a timeout of 10 seconds.
This timeout is globally configurable by setting the environment variable:
`APP_HOOK_TIMEOUT` (in seconds).

You can override this global behavior per request by setting the HTTP header:
`X-Hook-Timeout` (in seconds).

*Example:*

```bash
$ curl -H "X-Hook-Timeout: 5" http://localhost:8080/echo?foo=bar
```

### Webhook logs

As mentioned above, web hook logs are stream in real time during the call.
However, you can retrieve the logs of a previous call by using the hook ID: `http://localhost:8080/<NAME>/<ID>`

The hook ID is returned as an HTTP header with the Webhook response: `X-Hook-ID`

*Example:*

```bash
$ # Call webhook
$ curl -v http://localhost:8080/echo?foo=bar
...
< HTTP/1.1 200 OK
< Content-Type: text/event-stream
< X-Hook-Id: 2
...
$ # Retrieve logs afterwards
$ curl http://localhost:8080/echo/2
```

### Post hook notifications

The output of the script is collected and stored into a log file
(configured by the `APP_LOG_DIR` environment variable).

Once the script is executed, you can send the result and this log file to a notification channel.
Currently, only two channels are supported: `Email` and `HTTP`.

Notifications configuration can be done as follow:

```bash
$ export APP_NOTIFICATION_URI=http://requestb.in/v9b229v9
$ # or
$ webhookd --notification-uri=http://requestb.in/v9b229v9
```

Note that only the output of the script prefixed by "notify:" is sent to the notification channel.
If the output does not contain a prefixed line, no notification will be sent.

**Example:**

```bash
#!/bin/bash

echo "notify: Hello World" # Will be notified
echo "Goodbye"             # Will not be notified
```

You can overide the notification prefix by adding `prefix` as a query parameter to the configuration URL.

**Example:** http://requestb.in/v9b229v9?prefix="foo:"

#### HTTP notification

Configuration URI: `http://example.org`

Options (using query parameters):

- `prefix`: Prefix to filter output log

The following JSON payload is POST to the target URL:

```json
{
  "id": "42",
  "name": "echo",
  "text": "foo\nbar...\n",
  "error": "Error cause... if present",
}
```

Note that because the payload have a `text` attribute, you can use a [Mattermost][mattermost] webhook endpoint.

[mattermost]: https://docs.mattermost.com/developer/webhooks-incoming.html

#### Email notification

Configuration URI: `mailto:foo@bar.com`

Options (using query parameters):

- `prefix`: Prefix to filter output log
- `smtp`: SMTP host to use (by default: `localhost:25`)
- `from`: Sender email (by default: `webhookd <noreply@nunux.org>`)

### Authentication

You can restrict access to webhooks using HTTP basic authentication.

To activate basic authentication, you have to create a `htpasswd` file:

```bash
$ # create passwd file the user 'api'
$ htpasswd -B -c .htpasswd api
```
This command will ask for a password and store it in the htpawsswd file.

Please note that by default, the daemon will try to load the `.htpasswd` file.

But you can override this behavior by specifying the location of the file:

```bash
$ APP_PASSWD_FILE=/etc/webhookd/users.htpasswd
$ # or
$ webhookd -p /etc/webhookd/users.htpasswd
```

Once configured, you must call webhooks using basic authentication:

```bash
$ curl -u api:test -XPOST "http://localhost:8080/echo?msg=hello"
```

---


