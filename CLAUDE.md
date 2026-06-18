# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
# Run all tests (CGO_ENABLED not required; omit -race on Windows without CGO)
go test ./...

# Run a single test or suite
go test -run TestFunctionName ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Lint
golangci-lint run ./...

# Build
go build ./...

# Regenerate RSA stub signature (after changing test payload or keys)
go run ./main
```

## Architecture

This is a single-package Go SDK (`package dusupay`) that wraps the Dusupay payment gateway REST API v2. The module is `github.com/kachit/dusupay-sdk-go/v2`.

**Request flow:** `Client` → resource method → `Transport.Get/Post` → `RequestBuilder.BuildRequest` → `http.Client.Do` → `unmarshalResponse`

**Key types:**

- `Config` (`config.go`) — holds credentials (`PublicKey`, `SecretKey`, `SigningKey`) and the API base `URI`. Use `NewConfig` for production or `NewConfigSandbox` for sandbox. The two base URLs are declared in `env.go`.
- `Transport` / `RequestBuilder` (`http.go`) — handles all HTTP mechanics. `RequestBuilder` sets auth headers (`public-key`, `secret-key`, `x-api-version: 1`). `ResponseBody` is the shared envelope (`code`, `status`, `message`). Construct with `NewHTTPTransport`.
- `ResourceAbstract` (`resource.go`) — embedded by every resource struct, carries `*Transport` and `*Config`. Access fields as `r.tr` and `r.cfg` (not `r.ResourceAbstract.tr`).
- Resource files (`collections.go`, `payouts.go`, `merchants.go`, `providers.go`, `banks.go`) — each exposes typed request/response structs and methods that call `r.tr.Get` or `r.tr.Post`, then `unmarshalResponse`. API responses are async (202 Accepted); data contains only `internal_reference` and `merchant_reference`.
- `webhooks.go` — `CallbackPayloadInterface` with `BuildPayloadString()`. Concrete types: `CollectionCallbackPayload`, `PayoutCallbackPayload`. Payload format: `event:merchant_reference:internal_reference:transaction_type:transaction_status`.
- `signature.go` — `SignatureValidator` (RSA-SHA256) and `HMACSignatureValidator` (HMAC-SHA256) for webhook verification. HMAC signature format: `t=timestamp,s=hash`.
- `common.go` — shared typed constants (`TransactionTypeCode`, `TransactionMethodCode`, `CountryCode`, `CurrencyCode`) and `transformStructToMap` helper.

**Testing conventions:**

- Tests use `httpmock` and `testify/suite`; no real API calls are made.
- Stub helpers live in `stubs_test.go`: `BuildStubConfig`, `BuildStubHTTPTransport`, `BuildStubResponseFromFile`.
- JSON response fixtures are under `stubs/<resource>/<action>/` (e.g. `stubs/collections/initialize/`). Error fixtures are in `stubs/errors/`. The RSA public key for signature tests is `stubs/rsa/public-key.pem`.
- Each test file follows: mock the HTTP call → call the resource method → assert on the typed response struct.
- `defer func() { _ = resp.Body.Close() }()` is used instead of bare `defer resp.Body.Close()` to satisfy errcheck.
