# webhookd: A Simple Webhook Server

[![Build Status](https://github.com/ncarlier/webhookd/actions/workflows/build.yml/badge.svg)](https://github.com/ncarlier/webhookd/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ncarlier/webhookd)](https://goreportcard.com/report/github.com/ncarlier/webhookd)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/webhookd.svg)](https://hub.docker.com/r/ncarlier/webhookd/)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.me/nunux)

A minimalist, powerful webhook server designed to easily trigger shell scripts and run external processes via HTTP requests.

![Logo](webhookd.svg)

## 🚀 At a glance

![Demo GIF](demo.gif)

## Installation

Choose the method that best suits your deployment environment:

### 1. Go Install (Recommended)
For developers who have Go installed:
```bash
$ go install github.com/ncarlier/webhookd@latest
```

### 2. Binary Download Script
Download the binary for your architecture using one of these scripts:
**Via `curl` and `bash`:**
```bash
# Option A (General):
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/webhookd/master/install.sh | bash

# Option B (Gobinary):
$ curl -sf https://gobinaries.com/ncarlier/webhookd | sh
```

### 3. Docker Container
Run webhookd in a container for quick setup:
```bash
$ docker run -d --name=webhookd \
  -v ${PWD}/scripts:/scripts \
  -p 8080:8080 \
  ncarlier/webhookd
```

> **Note on Docker:** The official image is lightweight for simple scripts. For advanced needs (e.g., interacting with the Docker daemon), consider using `ncarlier/webhookd:edge-distrib`.

### 4. Package Manager (APT)
For Debian users, install via our custom repository:
[Custom Repository Link](https://packages.azlux.fr/)

> **Systemd Setup:** If installing system-wide, the service is pre-configured. You only need to start it with: `systemctl start webhookd`. Custom environment variables can be set in `/etc/webhookd.env`.

## Configuration & Usage

Webhookd accepts configuration via command line flags or by setting environment variables. For a complete list, run `webhookd -h`.

All available parameters and environment variables are detailed in [./etc/default/webhookd.env](./etc/default/webhookd.env).

### Directory Structure (Scripts)
Webhooks are defined as executable scripts within a specific directory structure.
* **Default Path:** Scripts are executed by default from the `./scripts` directory.
* **Changing Path:** You can override this using the `WHD_HOOK_SCRIPTS` environment variable or the `-hook-scripts` parameter.

**Example Structure:**
```
/scripts
|--> /github
  |--> /build.sh
  |--> /deploy.sh
|--> /push.js
|--> /echo.sh
|--> ...
```
> **Tip:** Webhookd supports any executable file type, provided it has execute rights. For example, a Node.js script requires `#!/usr/bin/env node` as its shebang line. Sample scripts are available in the [example folder](./scripts/examples), including Gitlab and Github integrations.

### Webhook Calling & Mapping
The directory structure dictates the webhook URL (`http://localhost:8080/<path>`). You can omit the script extension; by default, webhookd will look for `.sh`. This default extension can be changed via `WHD_HOOK_DEFAULT_EXT` or `-hook-default-ext`.

#### Response Handling (Streaming vs. Blocking)
How webhookd responds depends on your request headers:

1.  **Server-Sent Events (SSE):** Used when the `Accept` header is `text/event-stream`. Provides real-time, streamed output (see [reference][sse]).
2.  **Chunked Transfer Coding:** The default mode. Used when the `X-Hook-Mode` header is set to `chunked`. Also provides streamed output (see [reference][chunked]).
3.  **Blocking Mode:** Use this if no streaming is required by setting the `X-Hook-Mode` header to `buffered`. The request blocks until the script finishes, returning a summary payload.

#### Exit Code Mapping (Only in Blocking Mode)
Webhookd maps the script's exit code (0-255) to an HTTP status code:
*   **0:** `200 OK`
*   **1 - 99:** `500 Internal Server Error`
*   **100 - 255:** Adds 300 to the exit code (resulting in a 4xx or 5xx range).

> **Example Exit Code Calculation:** An exit status of `118` results in HTTP status `418 I'm a teapot`.

**Streaming Examples:**

*   **Server-sent events (SSE):**
    ```bash
    $ curl -v --header "Accept: text/event-stream" -XGET http://localhost:8080/foo/bar
    # ... output showing data: lines in real time ...
    error: exit status 118
    ```

*   **Chunked Transfer Coding (Default):**
    ```bash
    $ curl -v -XPOST --header "X-Hook-Mode: chunked" http://localhost:8080/foo/bar
    # ... output showing data lines immediately ...
    error: exit status 118
    ```

*   **Blocking Request:**
    ```bash
    $ curl -v -XPOST --header "X-Hook-Mode: buffered" http://localhost:8080/foo/bar
    # HTTP/1.1 418 I m a teapot (The status code)
    # ... script output followed by error details...
    ```

### Webhook Parameters (Input Handling)
Webhookd automatically converts various incoming request data into script variables:

*   **Query Parameters:** Converted directly to script variables.
*   **HTTP Headers:** Converted, following the snake\_case convention. (*e.g., `CONTENT-TYPE` becomes `content_type`*).
*   **Request Body:**
    *   `application/x-www-form-urlencoded`: Keys/values are mapped to variables.
    *   `text/*` or `application/json`: The entire payload is passed as the **first script parameter (`$1`)**.

**Built-in Parameters Added by Webhookd:**
| Variable | Description |
| :--- | :--- |
| `hook_id` | Unique hook ID (auto-increment) |
| `hook_name` | Name associated with the webhook call |
| `hook_method` | HTTP request method used |
| `x_forwarded_for` | Client IP address |
| `x_webauth_user` | Username if authentication is enabled |

**Example Usage:**
```bash
$ curl --data @test.json -H 'Content-Type: application/json' http://localhost:8080/echo?foo=bar
# Script output shows variables mapped correctly
Hook information: name=echo, id=1, method=POST
Query parameter: foo=bar
Header parameter: user-agent=curl/...
Script parameters: {"message": "this is a test"}
```

### Advanced Configuration

#### Timeout Control
*   **Global Timeout:** Set the default timeout for all hooks using `WHD_HOOK_TIMEOUT` (seconds).
*   **Per-Request Override:** Use the HTTP header `X-Hook-Timeout` (seconds) to override the global setting.

```bash
$ curl -H "X-Hook-Timeout: 5" http://localhost:8080/echo?foo=bar
```

#### Log Retrieval
While logs stream in real time, you can retrieve historical output using the hook ID: `http://localhost:8080/<NAME>/<ID>`.

The current execution's unique ID is returned via the `X-Hook-Id` header. (Logs can also be redirected to the server output using `WHD_LOG_MODULES=hook`).

#### Post-Hook Notifications
The script output can be collected and sent to external notification services.
*   **Configuration:** Set `WHD_NOTIFICATION_URI` or use `--notification-uri`.
*   **Filtering:** Only lines prefixed with "notify:" are sent. You can override this prefix via a query parameter (e.g., `?prefix="foo:"`).

**Example Script Snippet:**
```bash
#!/bin/bash

echo "notify: Success message for deployment." # Will be notified
echo "This is debug output, will be ignored."  # Will not trigger notification
```

##### Supported Notification Channels

###### Email Notification
*   **Configuration URI:** `mailto:foo@bar.com`
*   **Options (Query Params):**
    *   `prefix`: Filter output log lines by this prefix.
    *   `smtp`, `username`, `password`: Credentials for the SMTP relay.
    *   `conn`: Connection type (`tls`, `plain`, etc.).

###### HTTP Notification
*   **Configuration URI:** `http://example.org/endpoint`
*   **Payload:** A JSON object is POSTed to the target URL, suitable for Mattermost, Slack, or Discord webhooks.

```json
{
  "id": "42",
  "name": "echo",
  "text": "Script output content...",
  "error": "Error details..."
}
```

#### Security Features

##### Basic Authentication (Auth)
Restrict access using standard HTTP basic authentication.
1. Create the password file: `htpasswd -B -c .htpasswd api`
2. Set/Use the path: `export WHD_PASSWD_FILE=/etc/webhookd/users.htpasswd`
3. Usage requires credentials:
```bash
$ curl -u api:test -XPOST "http://localhost:8080/echo?msg=hello"
```

##### Upstream Authentication

In scenarios where authentication is handled by an upstream reverse proxy or API gateway (e.g., Authelia, Pomerium, Traefik), you can configure webhookd to rely on upstream headers.

Set the allowed upstream headers using `WHD_ALLOWED_UPSTREAM_HEADERS` (or `--allowed-upstream-headers`). By default this is `Accept,Content-Type,Content-Length,User-Agent,X-Forwarded-For`.
You can set `*` which forwards all HTTP headers to the webhook scripts.

For instance, to accept authentication from an upstream proxy using the `x-webauthn-user` header:

```bash
export WHD_ALLOWED_UPSTREAM_HEADERS="x-webauthn-user,content-type,user-agent"
```

##### Signature Verification (Integrity)
Ensure message authenticity using cryptographic signatures. Webhookd supports two methods:

1.  **HTTP Signatures:** Uses standards defined by IETF draft-cavage.
2.  **Ed25519 Signature:** Used by services like Discord.

To activate, set the truststore file location:
```bash
$ export WHD_TRUSTSTORE_FILE=/etc/webhookd/pubkey.pem
# Or use command flag: webhookd --truststore-file /path/to/pubkey.pem
```
*Calls must include appropriate headers (e.g., `Signature:` or `X-Signature-Ed25519`).*

#### TLS Support
Secure communications by enabling SSL/TLS.

*   **Simple Enable:**
    ```bash
    export WHD_TLS_ENABLED=true
    # Or: webhookd --tls-enabled
    ```
*   **Custom Certificates:** Provide specific files using flags (or environment variables for path): `-tls-cert-file` and `-tls-key-file`.
*   **ACME Support (Automatic SSL):** Enable by specifying a fully qualified domain name:
    ```bash
    export WHD_TLS_ENABLED=true
    export WHD_TLS_DOMAIN=hook.example.com
    # Or: webhookd --tls-enabled --tls-domain=hook.example.com
    ```

**⚠️ Networking Note:** To listen on privileged ports (80/443) on Linux, remember to use `setcap` for the binary: `sudo setcap CAP_NET_BIND_SERVICE+ep webhookd`

## License & Credits

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.


[sse]: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events
[chunked]: https://datatracker.ietf.org/doc/html/rfc2616#section-3.6.1