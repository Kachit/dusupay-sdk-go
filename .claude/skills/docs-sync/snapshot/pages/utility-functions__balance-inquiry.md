
# Balance Inquiry

The API allows the merchant to check the balance(s) for the different currencies they might hold.

<mark style="color:blue;">`GET`</mark> `https://sandboxapi.dusupay.com/data/wallet-balances`

The endpoint retrieves the balance(s) by currency. It returns currency wallets belonging to the merchant account

#### Query Parameters

<table><thead><tr><th width="121">Name</th><th width="88">Type</th><th>Description</th></tr></thead><tbody><tr><td>currency</td><td>String</td><td>The currency whose balances should be retrieved (Optional)</td></tr></tbody></table>

#### Headers

<table><thead><tr><th width="149">Name</th><th width="112">Type</th><th>Description</th></tr></thead><tbody><tr><td>public-key<mark style="color:red;">*</mark></td><td>String</td><td>The merchant account Public Key</td></tr><tr><td>secret-key<mark style="color:red;">*</mark></td><td>String</td><td>The merchant account Secret Key</td></tr></tbody></table>

### Sample Balance Inquiry Request

```powershell
curl https://sandboxapi.dusupay.com/data/wallet-balances \
   -H "Accept: application/json" \
   -H "x-api-version: 1" \
   -H "public-key: your-public-key" \
   -H "secret-key: your-secret-key"
```

{% tabs %}
{% tab title="200: OK Returns array of wallet balances based on the currency" %}

```json
{
    "code": 200,
    "status": "success",
    "message": "Request completed successfully.",
    "data": {
        "merchant_wallets": [
            {
                "available_balance": 40000,
                "actual_balance": 40000,
                "currency": "UGX",
                "uncleared_balance": 40000,
                "reserve_balance": 0,
                "wallet_type": "MERCHANT_COLLECTION"
            },
            {
                "available_balance": 0,
                "actual_balance": 0,
                "currency": "UGX",
                "uncleared_balance": 0,
                "reserve_balance": 0,
                "wallet_type": "MERCHANT_PAYOUT"
            }
        ]
    }
}

```

{% endtab %}

{% tab title="400: Bad Request When the request isn" %}

```javascript
{
  "code": 400,
  "status": "error",
  "message": "secret-key is required",
  "data": {}
}
```

{% endtab %}

{% tab title="401: Unauthorized When one of the keys is invalid" %}

```javascript
{
  "code": 401,
  "status": "error",
  "message": "Unauthorized API access. Invalid Merchant Secret Key",
  "data": {}
}
```

{% endtab %}
{% endtabs %}

### Wallet Types

<table><thead><tr><th width="253">Wallet Type</th><th>Description</th></tr></thead><tbody><tr><td>MERCHANT_COLLECTION</td><td>This wallet holds funds accumulated from funds collection of a specific currency. Funds directly deposited into the merchant account for a given currency would also live in a wallet of this type.</td></tr><tr><td>MERCHANT_PAYOUT</td><td>All the funds the merchant account wishes to designate for disbursements will be stored in a special wallet of this type. The merchant will be allowed to provision funds here by transferring from the collection account to this one.</td></tr></tbody></table>

### Balance Descriptions

<table><thead><tr><th width="190">Balance Name</th><th>Description</th></tr></thead><tbody><tr><td>available_balance</td><td>The amount/balance available for use. This the amount the merchant account has access to for any other platform operations</td></tr><tr><td>actual_balance</td><td>This is the total balance of a given currency for a certain merchant account wallet. It's the summation of all the balances in the given wallet.</td></tr><tr><td>uncleared_balance</td><td>This is the amount yet to be made available to the merchant account for use.</td></tr><tr><td>reserve_balance</td><td>In the event transactions are done using payment options for which a rolling reserve is implemented (e.g. card options), the accumulated rolling reserve amounts would be deposited under this balance</td></tr></tbody></table>
