
# Payment Options

Before any transaction request is performed, it's recommended that the payment options available to the merchant account are known. This endpoint returns that list.

The API description below shows how the list of payment options can be obtained. The request requires that the merchant filters by currency and transaction type to get the appropriate list of the payment options.

## Get Payment Options List

<mark style="color:blue;">`GET`</mark> `https://sandboxapi.dusupay.com/data/payment-providers`

Returns the list of payment options available to the merchant account based on currency and transaction type

#### Query Parameters

<table><thead><tr><th width="183">Name</th><th width="89">Type</th><th>Description</th></tr></thead><tbody><tr><td>currency<mark style="color:red;">*</mark></td><td>String</td><td>The 3-character ISO currency code e.g. UGX</td></tr><tr><td>transaction_type<mark style="color:red;">*</mark></td><td>String</td><td>The transaction type. Allowed types at the moment are <strong>COLLECTION</strong> and <strong>PAYOUT</strong></td></tr><tr><td>country_code</td><td>String</td><td>The 2-character country code e.g. SN (Optional). Useful for situations where multiple countries have the same currency</td></tr></tbody></table>

#### Headers

<table><thead><tr><th width="136">Name</th><th width="88">Type</th><th>Description</th></tr></thead><tbody><tr><td>public-key<mark style="color:red;">*</mark></td><td>String</td><td>The merchant account Public Key</td></tr></tbody></table>

### Sample Payment Options Request

```powershell
curl https://sandboxapi.dusupay.com/data/payment-providers?currency=UGX&transaction_type=COLLECTION \
   -H "Accept: application/json" \
   -H "x-api-version: 1" \
   -H "public-key: your-public-key"
```

{% tabs %}
{% tab title="200: OK Returns the list of payment options matching the query" %}

```json
{
    "code": 200,
    "status": "success",
    "message": "Request completed successfully.",
    "data": {
        "payment_providers": [
            {
                "provider_name": "MTN Mobile Money Uganda",
                "provider_code": "mtn_ug",
                "transaction_method": "MOBILE_MONEY",
                "transaction_method_name": "Mobile Money",
                "transaction_type": "Collection",
                "transaction_currency": "UGX",
                "min_amount": 500,
                "max_amount": 3000000,
                "max_daily_amount": 7000000,
                "is_active": true,
                "option_prefixes": [
                    "77",
                    "78",
                    "76",
                    "39"
                ]
            },
            {
                "provider_name": "Airtel Money Uganda",
                "provider_code": "airtel_ug",
                "transaction_method": "MOBILE_MONEY",
                "transaction_method_name": "Mobile Money",
                "transaction_type": "Collection",
                "transaction_currency": "UGX",
                "min_amount": 500,
                "max_amount": 3000000,
                "max_daily_amount": 7000000,
                "is_active": true,
                "option_prefixes": [
                    "74",
                    "75",
                    "70"
                ]
            }
        ]
    }
}

```

{% tabs %}
{% tab title="Response Description" %}

<table><thead><tr><th width="209">Parameter</th><th width="83">Type</th><th>Description</th></tr></thead><tbody><tr><td>provider_code</td><td>string</td><td>The unique code for the payment option. It is required for the transaction requests</td></tr><tr><td>provider_name</td><td>string</td><td>The name of the option</td></tr><tr><td>transaction_type</td><td>string</td><td>The transaction type</td></tr><tr><td>transaction_method</td><td>string</td><td>The transaction method supported by this option</td></tr><tr><td>transaction_currency</td><td>string</td><td>The currency in which the option transacts</td></tr><tr><td>min_amount</td><td>float</td><td>The minimum amount allowed per transaction when using this option</td></tr><tr><td>max_amount</td><td>float</td><td>The maximum amount allowed per transaction when using this option</td></tr><tr><td>max_daily_amount</td><td>float</td><td>The maximum amount a single account number is allowed to transact in 24hrs</td></tr><tr><td>option_prefixes</td><td>array</td><td>String prefixes that validate the phone numbers. Applies to mobile money options</td></tr></tbody></table>
{% endtab %}
{% endtabs %}
{% endtab %}

{% tab title="400: Bad Request When the request is not formed as expected" %}

```javascript
{
  "code": 400,
  "status": "error",
  "message": "transaction_type is required",
  "data": {}
}
```

{% endtab %}
{% endtabs %}
