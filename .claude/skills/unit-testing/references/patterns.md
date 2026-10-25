# Test patterns for dusupay-sdk-go

Copy-ready templates taken from the existing suites. Rename the identifiers, keep the shape.

## Contents
1. Resource suite skeleton
2. Resource method: success / JSON error / non-JSON error / validation
3. GET with query parameters
4. Request validation (`isValid`) tests
5. JSON shape tests (omitempty, pointers)
6. Callback / webhook tests
7. Signature tests
8. Fixture file template

---

## 1. Resource suite skeleton

Each resource has two suites: `<Name>TestSuite` for pure logic (validation, structs) and
`<Name>ResourceTestSuite` for HTTP calls through httpmock.

```go
type PayoutsResourceTestSuite struct {
	suite.Suite
	cfg      *Config
	ctx      context.Context
	testable *PayoutsResource
}

func (suite *PayoutsResourceTestSuite) SetupTest() {
	cfg := BuildStubConfig()
	transport := BuildStubHTTPTransport()
	suite.cfg = cfg
	suite.ctx = context.Background()
	suite.testable = &PayoutsResource{NewResourceAbstract(transport, cfg)}
	httpmock.Activate()
}

func (suite *PayoutsResourceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func TestPayoutsResourceTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutsResourceTestSuite))
}
```

`httpmock.Activate()` hooks `http.DefaultTransport`. `BuildStubHTTPTransport` uses `&http.Client{}`,
which goes through the default transport, so no real network call can escape.

## 2. Resource method: the four standard cases

`validPayoutRequest()` below is shorthand for "a request with every required field set". Existing
tests inline the struct literal. Either inline it or add a small helper in the test file.

**Success.** Asserts the typed fields and that the raw body is still readable afterwards.
`unmarshalResponse` re-wraps the body, and callers rely on that.

```go
func (suite *PayoutsResourceTestSuite) TestCreateSuccess() {
	body, _ := LoadStubResponseData("stubs/payouts/send-funds/success.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds",
		httpmock.NewBytesResponder(http.StatusAccepted, body))

	result, resp, err := suite.testable.Create(suite.ctx, validPayoutRequest())
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.True(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), http.StatusAccepted, result.Code)
	assert.Equal(suite.T(), "DUSUPAYRMGRXNNYBWATKJ", result.Data.InternalReference)

	defer func() { _ = resp.Body.Close() }()
	bodyRsp, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), body, bodyRsp)
}
```

**API error with a JSON body.** The error text comes from the API `message`.

```go
func (suite *PayoutsResourceTestSuite) TestCreateJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/401.json")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds",
		httpmock.NewBytesResponder(http.StatusUnauthorized, body))

	result, resp, err := suite.testable.Create(suite.ctx, validPayoutRequest())
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.False(suite.T(), result.IsSuccess())
	assert.Equal(suite.T(), "Unauthorized API access. Unknown Merchant", err.Error())
}
```

**Non-JSON error** (gateway HTML page). The result is nil and the response is still returned.

```go
func (suite *PayoutsResourceTestSuite) TestCreateNonJsonError() {
	body, _ := LoadStubResponseData("stubs/errors/500.html")
	httpmock.RegisterResponder(http.MethodPost, suite.cfg.URI+"/payout/send-funds",
		httpmock.NewBytesResponder(http.StatusInternalServerError, body))

	result, resp, err := suite.testable.Create(suite.ctx, validPayoutRequest())
	defer func() { _ = resp.Body.Close() }()
	assert.Error(suite.T(), err)
	assert.NotEmpty(suite.T(), resp)
	assert.Empty(suite.T(), result)
}
```

**Validation error.** No HTTP call happens, so the response is nil. The `//nolint:bodyclose`
comment is required because the linter can't see that there is no body to close.

```go
func (suite *PayoutsResourceTestSuite) TestCreateInvalidRequest() {
	result, rsp, err := suite.testable.Create(suite.ctx, &PayoutRequest{}) //nolint:bodyclose // returns nil response on validation error
	assert.Nil(suite.T(), result)
	assert.Nil(suite.T(), rsp)
	assert.EqualError(suite.T(), err, `PayoutsResource.Create error: parameter "currency" is empty`)
}
```

## 3. GET with query parameters

`RegisterResponder` ignores the query string, so a test using it can't tell whether a filter was
sent. To prove a parameter is sent, use `RegisterResponderWithQuery`. An unmatched request then
fails with "no responder found".

```go
httpmock.RegisterResponderWithQuery(http.MethodGet, suite.cfg.URI+"/data/wallet-balances",
	"currency=UGX", httpmock.NewBytesResponder(http.StatusOK, body))
```

To prove a parameter is *not* sent, register with the exact expected query, or inspect
`httpmock.GetCallCountInfo()`.

## 4. Request validation tests

Write one method per rule, named `Test<Struct>IsValid<Condition>`. Start from a fully valid
request and remove or alter exactly one field, so each test isolates the rule under test.

```go
func (suite *PayoutsTestSuite) TestPayoutRequestIsValidBankEmptyBankCode() {
	request := PayoutRequest{ /* all required fields */ TransactionMethod: TransactionMethodBank}
	assert.EqualError(suite.T(), request.isValid(), `parameter "extra_params.bank_code" is empty`)

	request.ExtraParams = &PayoutExtraParams{BankCode: "access_bank_ng"}
	assert.NoError(suite.T(), request.isValid())
}
```

Conditional rules ("required only when X") need both sides tested: X true with the field missing
gives an error, and X false with the field missing gives no error.

## 5. JSON shape tests

What actually goes over the wire is what matters. Test it through `transformStructToMap`, the
same helper the resources use:

```go
allow := false
result, err := transformStructToMap(&CollectionRequest{AllowFinalStatusChange: &allow})
assert.NoError(suite.T(), err)
assert.Equal(suite.T(), false, result["allow_final_status_change"]) // pointer: false is sent
assert.NotContains(suite.T(), result, "card_cipher")                 // omitempty: unset is omitted
```

Add a test like this whenever a field gains `omitempty`, becomes a pointer, or nests (for
example `extra_params`). Nested objects come back as `map[string]interface{}`.

## 6. Callback / webhook tests

Unmarshal a fixture into the wrapper type and check both levels:

```go
var callback CollectionCallback
body, _ := LoadStubResponseData("stubs/webhooks/callback/collection-success.json")
assert.NoError(suite.T(), json.Unmarshal(body, &callback))
assert.Equal(suite.T(), CallbackEventTransactionCompleted, callback.Event)
assert.Equal(suite.T(), "DUSUPAY405GZM1G5JXGA71IK", callback.Payload.InternalReference)
assert.Equal(suite.T(),
	"transaction.completed:76859aae-f148-48c5-9901-2e474cf19b71:DUSUPAY405GZM1G5JXGA71IK:COLLECTION:COMPLETED",
	callback.BuildPayloadString())
```

## 7. Signature tests

- **RSA with the repo test key:** `stubs/rsa/public-key.pem` pairs with `main/private.pem`.
  `main/` is gitignored, so guard with `suite.T().Skip("main/private.pem not available")`.
- `stubSignature` (in `stubs_test.go`) signs the collection-success fixture's payload string.
  If you change that fixture's references, regenerate it with `go run ./main` and update the constant.
- **Real DusuPay samples:** the docs publish a sandbox RSA signature and an HMAC example (key
  `SGNKYUEMYFDEHRWGPEUG`). `signature_test.go` verifies both against `SandboxPublicKey`. Keep these.
  They are the only tests that check the SDK against DusuPay's own data.

## 8. Fixture file template

Path: `stubs/<resource>/<action>/<case>.json`, for example `stubs/payouts/send-funds/success.json`.
The envelope is always `code`, `status`, `message`, `data`:

```json
{
  "code": 202,
  "status": "accepted",
  "message": "Request Accepted",
  "data": {
    "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
    "merchant_reference": "MCTREFT2WMNWZ23SBN6Y"
  }
}
```

Errors are `"status": "error"` with `"data": {}`. Callbacks have no envelope, just
`{"event": ..., "payload": {...}}`.
