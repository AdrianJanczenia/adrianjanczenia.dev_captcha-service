# CAPTCHA service

This Go service issues signed proof-of-work (PoW) challenges, generates image CAPTCHAs, and verifies their answers. Redis stores short-lived challenge state. The gateway calls this service over HTTP; the content service reads the solved CAPTCHA state from the same Redis instance when issuing a CV download token.

## HTTP API

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/pow` | Return a signed `seed` and `signature`. |
| `POST` | `/captcha` | Accept `seed`, `signature`, and `nonce`; return `captchaId` and `captchaImg`. |
| `POST` | `/verify` | Accept `captchaId` and `captchaValue`; mark a correct answer as solved. |

The client finds a nonce for which the SHA-256 hash of `seed + nonce` begins with the configured number of zero hexadecimal characters. The service checks the HMAC signature and timestamp before accepting the work. Redis `SETNX` prevents the same seed from creating multiple CAPTCHAs, including under concurrent requests.

The `captcha:<id>` key holds the answer, remaining attempts, and solved flag, with expiry enforced by its Redis TTL. The short script in `internal/service/redis/scripts/verify_captcha.lua` verifies an answer and updates this state atomically without extending the TTL. It is embedded in the Go binary. Incorrect answers consume an attempt; a correct answer marks the challenge as solved. Verification itself does not issue a CV token or remove the solved state: the content service handles that next step.

The code separates HTTP handlers, use-case processes and tasks, and the Redis adapter. Dependencies are assembled in `internal/app` and passed to constructors; there is no process-wide configuration variable. The image library's required store is discarded because Redis is the source of challenge state.

Error responses use a JSON `error` slug with the appropriate HTTP status. Invalid input and unsupported methods have separate slugs; a `405` response includes the allowed method in the `Allow` header.

## Configuration

The service uses Go 1.25.1 (see `go.mod`). `APP_ENV=production` selects `config/production/config.yml`; otherwise the local file is used. Supply `REDIS_URL` and `HMAC_SECRET` through environment variables; their YAML fields are intentionally empty. The configuration also defines the HTTP port, PoW difficulty and lifetime, CAPTCHA lifetime and attempt limit, and startup retry policy. Startup validates these settings and checks Redis with a bounded ping.

Keep the PoW difficulty consistent with the front-end implementation: its current JavaScript solver expects four leading zero hexadecimal characters.

## Run locally

1. Start a local Redis instance and set `REDIS_URL` to its connection URL. Set `HMAC_SECRET` to a private value in the environment. Neither value is supplied by `config/local/config.yml`.
2. From this repository root, run `go run .`. With the local configuration, HTTP listens on port 8083.
3. In another terminal, check that the service issues a PoW challenge:

   ```sh
   curl -i http://127.0.0.1:8083/pow
   ```

The `/captcha` and `/verify` endpoints form a sequence: first solve the work from `/pow`, then submit the returned CAPTCHA answer. The front end normally calls them through the gateway. For container use, the `Dockerfile` builds the service; its Redis address and HMAC secret still need to be supplied at runtime.

## Verification

```sh
go test -race ./...
```

Opt-in integration tests exercise Lua state transitions, TTL preservation, and concurrent attempts against a real Redis instance:

```sh
go test -race -tags=integration ./internal/service/redis
```

Set `REDIS_TEST_URL` to an isolated test Redis before running the integration tests. They skip when it is unset and create temporary keys in the selected database.
