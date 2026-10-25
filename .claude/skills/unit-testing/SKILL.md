---
name: unit-testing
description: Write and fix unit tests in the dusupay-sdk-go repo the project's way, using testify suites, httpmock (no real API calls), and JSON fixtures under stubs/ copied from the DusuPay docs, with a check that coverage doesn't drop. Use this whenever the user asks to add, write, update or fix tests; to cover a new resource method, request field, validation rule, callback type or signature logic; when `go test` is red; when a fixture or stub is out of date; or after changing SDK code without touching the tests ("test this", "why is TestCreateSuccess failing", "add coverage for the new field"), even if they don't say "unit test".
---

# unit-testing

Tests in this repo protect two things: that the SDK sends exactly what the DusuPay API expects,
and that it decodes exactly what the API returns. Users can't see the wire format, so a wrong
JSON tag or a dropped field becomes a silent production bug: money doesn't move and nobody knows
why. Write tests that would catch that, not tests that just execute lines.

Before writing anything, read `references/patterns.md`. It has copy-ready templates for every
test kind in the repo. Match those suites rather than inventing a new style: one scenario per
suite method, `assert.*(suite.T(), ...)`, files named `<resource>_test.go`, and the helpers from
`stubs_test.go` (`BuildStubConfig`, `BuildStubHTTPTransport`, `LoadStubResponseData`).

## Mode A: tests for new or changed code

### 1. Find what changed and what it needs

```bash
git status --short
git diff
```

For each changed exported thing, work out which kinds of test it needs:

| change | tests to write (templates in references/patterns.md) |
|---|---|
| new resource method | the four standard cases: success, JSON error (401), non-JSON error (500.html), validation error with `//nolint:bodyclose` |
| new or changed request field | JSON shape test via `transformStructToMap`; a validation test if it's required or conditionally required |
| new query/filter parameter | `httpmock.RegisterResponderWithQuery` so the test fails if the param isn't sent; also the nil/empty-filter case |
| new response field | add it to the fixture and assert it in the success test |
| field → pointer, or `omitempty` added | shape test that zero values (`false`, `0`, `""`) are sent or omitted as intended |
| new validation rule | one test per rule; for conditional rules, both sides (condition true → error, false → ok) |
| callback/webhook type | unmarshal a fixture, assert both levels, test `BuildPayloadString` (including a nil payload) |
| crypto (signatures, card encryption) | a round trip with the repo test keypair, a tampered input that fails, a malformed key |

Also check existing tests whose expectations the change invalidates. A changed signature or a
new required field often breaks "valid request" fixtures in other tests.

### 2. Get fixture data from the DusuPay docs

Fixtures should be the API's real shapes, not guesses. If you invent a field, the test passes and
production still breaks. When a test needs a new or changed response, callback or request:

1. Find the endpoint's page. The index is `https://developer.dusupay.com/llms.txt`, and any page
   is available as markdown by appending `.md`. The OpenAPI spec linked from
   `https://developer.dusupay.com/appendix/api-reference.md` has exact field names and types.
2. Copy the documented JSON example into `stubs/<resource>/<action>/<case>.json`. Keep the
   documented field names and nesting exactly. You may change values (references, amounts) to
   keep them consistent with existing tests.
3. If the docs disagree with themselves (for example, the spec says `otp` is a number but the
   curl example sends a string), follow what the SDK code does, and tell the user about the
   discrepancy rather than resolving it silently.
4. If you can't reach the docs, say so and model the fixture on the closest existing stub. Don't
   present invented shapes as documented.

Don't edit `stubs/webhooks/callback/collection-success.json`'s event or references casually:
`stubSignature` in `stubs_test.go` is an RSA signature over them. If you must change them,
regenerate the signature with `go run ./main` and update the constant.

### 3. Write the tests, then make them earn their keep

After writing a test, ask whether it would fail if the code were wrong. Quick ways to check:

- Temporarily break the code (wrong JSON tag, removed validation line), run the test, and
  confirm it goes red. Then restore the code. Do this for at least the most important new test.
- Assert specific values (`assert.Equal(..., "DUSUPAYRMGRXNNYBWATKJ", ...)`). Don't stop at
  `NotEmpty`/`NoError` on decoded data: those pass even when every field decodes to its zero value.
- Use `assert.EqualError` with the full message for validation errors. The message is
  user-facing API.

### 4. Run and measure

```bash
go test ./...
bash .claude/skills/unit-testing/scripts/coverage.sh HEAD
golangci-lint run ./...
```

`coverage.sh HEAD` measures the current tree against the last commit in a temporary git worktree
(your working tree is never stashed or reset). It prints the delta and the functions still below
100%. The bar is that **total coverage must not drop** below the baseline. If it drops (exit
code 2), add tests for the uncovered new branches it lists before finishing. Don't chase 100%;
an uncovered `json.Marshal` error branch isn't worth a contrived test.

The lint config relaxes funlen, lll, goconst and wsl in `_test.go` files, but `bodyclose`,
`errcheck` and `gosec` still apply. Close bodies with `defer func() { _ = resp.Body.Close() }()`.

## Mode B: fixing a failing test

The question to answer first: **is the code wrong, or is the test wrong?** Both happen here, and
"fixing" the test when the code is wrong hides a real bug. Never edit an assertion just to make
it match the current output without knowing which side is right.

1. Reproduce narrowly and read the full failure:
   ```bash
   go test -run 'TestPayoutsResourceTestSuite/TestCreateSuccess' -v .
   ```
   Suite tests are addressed as `TestSuiteRunner/TestMethod`.
2. Classify the failure:
   - **Code regression:** the test encodes documented behaviour and the code changed. Fix the code.
   - **Intentional change, stale test:** the code was deliberately changed (new signature, renamed
     field, new required param) and the test still encodes the old behaviour. Update the test and
     fixtures to the new behaviour, and check the docs to confirm the new behaviour is right.
   - **API drift:** the fixture no longer matches what DusuPay documents. Re-sync the fixture from
     the docs (step 2 above), then see whether the code needs to change too.
   - **Environment:** for example, a skipped test because `main/private.pem` is missing (it's
     gitignored), or a httpmock "no responder found" from a URL or query mismatch. Fix the setup,
     not the assertion.
3. Tell the user which category it was and why, especially when you changed code rather than the
   test, or when you suspect a real bug you didn't fix.
4. Run the whole suite afterwards (`go test ./...`). Shared fixtures mean a fix in one place can
   break another test.

## Finish

Report briefly:
- which tests you added or changed, and what behaviour each one protects
- which fixtures you added or changed, and the doc page each came from
- the coverage line from `coverage.sh` (before → after)
- any doc discrepancy or suspected bug you found but didn't resolve

Don't commit. If the user wants that, the `git-commit` skill handles it.
