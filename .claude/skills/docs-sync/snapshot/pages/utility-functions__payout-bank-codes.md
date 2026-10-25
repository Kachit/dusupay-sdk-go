
# Payout Bank Codes

For all countries/currencies where DusuPay supports bank payouts, there's need to specify the unique code for the bank to which the funds are to be sent. These codes are available via the API

The bank codes are attached to the respective payment providers designated for paying out funds to the supported banks. Therefore, a provider code, obtained from [here](/utility-functions/payment-options.md) can be used to query the list of supported banks. Any provider that supports BANK payouts will have the list of banks and respective codes.

### Get Payout Banks List

<mark style="color:blue;">`GET`</mark> `https://sandboxapi.dusupay.com/data/payout-bank-codes`

Returns the list of payout banks based on the provider code.

#### Query Parameters

<table><thead><tr><th width="162">Name</th><th width="80">Type</th><th>Description</th></tr></thead><tbody><tr><td>provider_code<mark style="color:red;">*</mark></td><td>String</td><td>The payment provider code for the provider whose bank codes are required</td></tr></tbody></table>

#### Headers

<table><thead><tr><th width="150">Name</th><th width="107">Type</th><th>Description</th></tr></thead><tbody><tr><td>public-key<mark style="color:red;">*</mark></td><td>String</td><td>The merchant account Public Key</td></tr></tbody></table>

### Sample Payment Options Request

```powershell
curl https://sandboxapi.dusupay.com/data/payout-bank-codes?provider_code=bank_ng \
   -H "Accept: application/json" \
   -H "x-api-version: 1" \
   -H "public-key: your-public-key"
```

{% tabs %}
{% tab title="200: OK Returns the list of payout banks" %}

```json
{
    "code": 200,
    "status": "success",
    "message": "Request completed successfully.",
    "data": {
        "payout_banks": [
            {
                "bank_name": "Access Bank Nigeria",
                "bank_code": "access_bank_ng",
                "is_active": true
            },
            {
                "bank_name": "Guaranty Trust Bank",
                "bank_code": "gtbank_ng",
                "is_active": true
            }
        ]
    }
}
```

{% tabs %}
{% tab title="Response Description" %}

<table><thead><tr><th width="137">Parameter</th><th width="96">Type</th><th>Description</th></tr></thead><tbody><tr><td>bank_code</td><td>String</td><td>The unique code for the payout bank. It is required for the bank payout transaction requests</td></tr><tr><td>bank_name</td><td>String</td><td>The name of the bank</td></tr><tr><td>is_active</td><td>Boolean</td><td>Whether or not the bank is active for usage</td></tr></tbody></table>
{% endtab %}
{% endtabs %}
{% endtab %}

{% tab title="400: Bad Request When the request is not formed as expected" %}

```json
{
  "code": 400,
  "status": "error",
  "message": "provider_code is required",
  "data": {}
}
```

{% endtab %}
{% endtabs %}
