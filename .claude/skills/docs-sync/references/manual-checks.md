# Manual checks the drift tool can't do

The drift tool compares field *names and kinds*. These things live in prose, tables or
special pages, so read them yourself. Limit yourself to pages that `snapshot.sh compare`
reports as changed. On a first sync (no snapshot), check all of them.

| What | Where in the docs | Compare with |
|---|---|---|
| Base URLs | `introduction.md` | `ProdAPIUrl`, `SandboxAPIUrl` in `env.go` |
| Public keys | `getting-started__dusupay-public-keys.md` | `SandboxPublicKey`, `ProdPublicKey` in `env.go` (command below) |
| Auth headers, API version | `getting-started__authentication.md` | `RequestBuilder.buildHeaders` in `http.go` (`x-api-version: 1`) |
| Error codes and envelope | `getting-started__error-handling.md`, error tabs on endpoint pages | `ResponseBody`, `IsSuccess` in `http.go` |
| Conditional / prose validation ("required for bank transfers", "10-30 characters", "at least 8 characters", "auto") | parameter-table descriptions on `funds-collection__*`, `payouts-disbursements__*` | `isValid()` methods |
| Callback events and statuses | `utility-functions__handling-notifications-callbacks__callback-events.md`, `utility-functions__transaction-status-verification.md` | `CallbackEvent*`, `TransactionStatus*` in `common.go` |
| Callback envelope / redirect query params | `funds-collection__getting-started.md`, `*hosted-page*.md`, `*direct-card-payment-s2s.md` | `webhooks.go` (`CollectionCallback`, `RedirectPayload`) |
| Signature payload format and headers | `callbacks__hmac-signature-verification.md`, `callbacks__rsa-signature-verification.md` | `signature.go`, `BuildPayloadString` in `webhooks.go` |
| Published signature samples | same pages (sample RSA signature, HMAC key `SGNKYUEMYFDEHRWGPEUG`) | `TestValidateSignatureDocsSandboxSample`, `TestValidateSignatureDocsSample` in `signature_test.go` |
| Card encryption | `funds-collection__card-payments__direct-card-payment-s2s.md` | `cards.go` (PKCS#1 v1.5; the doc samples disagree on padding, so keep it unless the docs settle it) |
| Countries, currencies, methods | `getting-started__supported-countries-regions.md` | `CountryCode*`, `CurrencyCode*`, `TransactionMethod*` in `common.go` |
| Wallet types | `utility-functions__balance-inquiry.md` | `WalletType*` in `common.go` |
| Sandbox test numbers/OTPs | `getting-started__sandbox-test-accounts.md`, direct-charge page | test values in `*_test.go` (informational) |

## Public key check

```bash
D=<fetched dir>
diff <(awk '/BEGIN PUBLIC KEY/,/END PUBLIC KEY/' "$D/pages/getting-started__dusupay-public-keys.md" | tr -d '\r`' | sed 's/^const [A-Za-z]* = //') \
     <(awk '/BEGIN PUBLIC KEY/,/END PUBLIC KEY/' env.go | tr -d '\r`' | sed 's/^const [A-Za-z]* = //') \
  && echo "public keys match env.go"
```

A key change matters a lot. Every merchant's signature verification and card encryption would
break. Treat it as top of the report, and also update `main/dusupay.public.key.pem` if it exists.

## The Postman collection

`postman.json` (from https://postman.dusupay.com/) is the closest thing to real traffic: complete
request bodies per flow, plus which headers each request sends (for example `secret-key` only on
balances and payouts). It has no response examples. When it disagrees with the pages (for
example, payouts use `"merchant_reference": "auto"` there, but the payout pages say "at least 8
characters" and never mention `auto`), that's a conflict for the report, not a fact to act on.
