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

**Or** use APT:

Finally, it is possible to install Webhookd using the Debian packaging system through this [custom repository](https://packages.azlux.fr/).

> Note that custom configuration variables can be set into `/etc/webhoockd.env` file.
> Sytemd service is already set and enable, you just have to start it with `systemctl start webhoockd`.

## Configuration

Webhookd can be configured by using command line parameters or by setting environment variables.

Type `webhookd -h` to display all parameters and related environment variables.

All configuration variables are described in [etc/default/webhookd.env](./etc/default/webhookd.env) file.

## Usage

### Directory structure

Webhooks are simple scripts within a directory structure.

By default inside the `./scripts` directory.
You can override the default using the `WHD_SCRIPTS` environment variable or `-script` parameter.

*Example:*

```
/scripts
|--> /github
  |--> /build.sh
  |--> /deploy.sh
|--> /push.js
|--> /echo.sh
|--> ...
```

Note that Webhookd is able to run any type of file in this directory as long as the file is executable.
For example, you can execute a Node.js file if you give execution rights to the file and add the appropriate `#!` header (in this case: `#!/usr/bin/env node`).

### Webhook URL

The directory structure define the webhook URL.

You can omit the script extension. If you do, webhookd will search for a `.sh` file.
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
```

### Webhook timeout configuration

By default a webhook has a timeout of 10 seconds.
This timeout is globally configurable by setting the environment variable:
`WHD_HOOK_TIMEOUT` (in seconds).

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
(configured by the `WHD_LOG_DIR` environment variable).

Once the script is executed, you can send the result and this log file to a notification channel.
Currently, only two channels are supported: `Email` and `HTTP`.

Notifications configuration can be done as follow:

```bash
$ export WHD_NOTIFICATION_URI=http://requestb.in/v9b229v9
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

You can override the notification prefix by adding `prefix` as a query parameter to the configuration URL.

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
$ export WHD_PASSWD_FILE=/etc/webhookd/users.htpasswd
$ # or
$ webhookd --passwd-file /etc/webhookd/users.htpasswd
```

Once configured, you must call webhooks using basic authentication:

```bash
$ curl -u api:test -XPOST "http://localhost:8080/echo?msg=hello"
```

### Signature

You can ensure message integrity (and authenticity) with [HTTP Signatures](https://www.ietf.org/archive/id/draft-cavage-http-signatures-12.txt).

To activate HTTP signature verification, you have to configure the key store:

```bash
$ export WHD_KEY_STORE_URI=file:///etc/webhookd/keys
$ # or
$ webhookd --key-store-uri file:///etc/webhookd/keys
```

Note that only `file://` URI s currently supported.
All public keys stored in PEM format in the targeted directory will be loaded.

Once configured, you must call webhooks using a valid HTTP signature:

```bash
$ curl -X POST \
  -H 'Date: <req-date>' \
  -H 'Signature: keyId=<key-id>,algorithm="rsa-sha256",headers="(request-target) date",signature=<signature-string>' \
  -H 'Accept: application/json' \
  "http://loclahost:8080/echo?msg=hello"
```

### TLS

You can activate TLS to secure communications:

```bash
$ export WHD_TLS_LISTEN_ADDR=:8443
$ # or
$ webhookd --tls-listen-addr=:8443
```

This will disable HTTP port.

By default webhookd is expecting a certificate and key file (`./server.pem` and `./server.key`).
You can provide your own certificate and key with `-tls-cert-file` and `-tls-key-file`.

Webhookd also support [ACME](https://ietf-wg-acme.github.io/acme/) protocol.
You can activate ACME by setting a fully qualified domain name:

```bash
$ export WHD_TLS_LISTEN_ADDR=:8443
$ export WHD_TLS_DOMAIN=hook.example.com
$ # or
$ webhookd --tls-listen-addr=:8443 --tls-domain=hook.example.com
```

**Note:**
On *nix, if you want to listen on ports 80 and 443, don't forget to use `setcap` to privilege the binary:

```bash
sudo setcap CAP_NET_BIND_SERVICE+ep webhookd
```

---
