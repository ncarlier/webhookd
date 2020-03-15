# httpsig

A simple HTTP client with HTTP signature support.

# Usage

- Generate an RSA key: `make key`

- Add `key_id` header to public key:
  
```pem
-----BEGIN PUBLIC KEY-----
key_id: my-key

MIIEowIBAAKCAQEAwdCB5DZD0cFeJYUu1W3IlNN9y+NZC/Jqktdkn8/WHlXec07n
...
-----END PUBLIC KEY-----
```

- Start Webhookd with HTTP signature support:

```bash
$ webhookd --trust-store-file ./key-pub.pem
```

- Make HTTP signed request:

```bash
$ ./release/httpsig \
  --key-id my-key \
  --key-file ./key.pem \
  http://localhost:8080/echo`
```