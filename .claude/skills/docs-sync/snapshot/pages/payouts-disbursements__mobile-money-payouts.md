
# Mobile Money Payouts

The DusuPay API enables you to transfer funds directly from your available balance to a mobile money account. The documentation below explains further.

{% hint style="info" %}
We recommend checking out the [Getting Started](/payouts-disbursements/getting-started.md) section to understand the basics of payouts first and the general workflow. This guide assumes that you have read that
{% endhint %}

## Overview

The currently supported mobile money channels are listed [here](/getting-started/supported-countries-regions.md) (to be updated from time to time). Test mobile money phone numbers are also described in [this](/getting-started/sandbox-test-accounts.md) section. It's very important that you track the available balance on your PAYOUT wallet before initiating these transactions.

## Step 1: Obtain the required data for the payment request

The table below describes the request parameters that are used for the payout/disbursement request. Most/all will be collected from the paying customer.

<table><thead><tr><th width="197">Parameter</th><th width="93">Type</th><th width="98">Required</th><th>Description</th></tr></thead><tbody><tr><td>merchant_reference</td><td>String</td><td>true</td><td>The unique reference for this request. It must be at least 8 characters long.</td></tr><tr><td>transaction_method</td><td>String</td><td>true</td><td>The transaction method to be used. This will be <em><strong>MOBILE_MONEY</strong></em> for this request</td></tr><tr><td>currency</td><td>String</td><td>true</td><td>The 3-character ISO currency code for the request currency</td></tr><tr><td>amount</td><td>Number</td><td>true</td><td>The amount to be transferred</td></tr><tr><td>provider_code</td><td>String</td><td>true</td><td>The provider code as obtained from the payment options <a href="/utility-functions/payment-options.md#get-payment-options-list">list</a></td></tr><tr><td>account_number</td><td>String</td><td>true</td><td>The phone number of the recipient. This should be sent in international format e.g. 256777000001 for Ugandan numbers</td></tr><tr><td>customer_name</td><td>String</td><td>true</td><td>The name of the customer</td></tr><tr><td>description</td><td>String</td><td>true</td><td>The description/narration for the transaction. Between 10-30 characters</td></tr></tbody></table>

After collecting the necessary mobile money payment information from your customer, prepare your request payload as demonstrated below.

```json
{
	"merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
	"transaction_method": "MOBILE_MONEY",
	"currency": "UGX",
	"amount": 4000,
	"provider_code": "mtn_ug",
	"account_number": "256777000001",
	"customer_name": "JOHN DOE",
	"description": "Test Payout"
}
```

<mark style="color:green;">`POST`</mark> `https://sandboxapi.dusupay.com/payout/send-funds`

The request is sent as a JSON body as demonstrated by the sample request below. Sample responses (acknowledgement and failure) are also shared.

```powershell
curl -X POST "https://sandboxapi.dusupay.com/payout/send-funds" \
   -H 'Content-Type: application/json' \
   -H "x-api-version: 1" \
   -H "public-key: your-public-key" \
   -H "secret-key: your-secret-key" \
   -d '{
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "transaction_method": "MOBILE_MONEY",
        "currency": "UGX",
        "amount": 10000,
        "provider_code": "mtn_ug",
        "account_number": "256777000001",
        "customer_name": "JOHN DOE",
        "description": "Test Payout"
    }'
```

{% tabs %}
{% tab title="202: Accepted - Request acknowledged for processing" %}

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

{% endtab %}

{% tab title="400: Bad Request - Request is not formed as expected" %}

```json
{
    "code": 400,
    "status": "error",
    "message": "256752000001 is not a valid MTN Mobile Money Uganda (mtn_ug) phone number",
    "data": {}
}
```

{% endtab %}
{% endtabs %}

## Step 2: Handshake - Verify transaction

{% hint style="info" %}
The API will attempt to verify the transaction from your platform as described in [this](/payouts-disbursements/getting-started.md#step-2-payout-request-verification) section. Ensure that the correct URL is configured on your merchant account and that you're handling the verification request appropriately.
{% endhint %}

## Step 3: Handle the final status webhook

Every merchant account is expected to have configured a callback/webhook URL for payouts. For all payouts that transition to the final state (COMPLETED or FAILED), a JSON POST request will be made to the callback URL. Sample callback payloads (request bodies) are shared below. Be sure to check out [Handling Notifications](/utility-functions/handling-notifications-callbacks.md) to see how you should verify the signature(s) in the request headers and how to respond.

{% tabs %}
{% tab title="Successful Mobile Money Payout" %}

```json
{
    "event": "transaction.completed",
    "payload": {
        "id": 20760,
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
        "transaction_type": "PAYOUT",
        "request_currency": "UGX",
        "transaction_amount": 10000,
        "transaction_currency": "UGX",
        "transaction_charge": 1000,
        "transaction_account": "256777000001",
        "charge_customer": false,
        "total_debit": 11000,
        "provider_code": "mtn_ug",
        "request_amount": 10000,
        "customer_name": "JOHN DOE",
        "transaction_status": "COMPLETED",
        "status_message": "Transaction Completed Successfully"
    }
}
```

{% endtab %}

{% tab title="Failed Mobile Money Payout" %}

```json
{
    "event": "transaction.failed",
    "payload": {
        "id": 20760,
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
        "transaction_type": "PAYOUT",
        "request_currency": "UGX",
        "transaction_amount": 10000,
        "transaction_currency": "UGX",
        "transaction_charge": 0,
        "transaction_account": "256777000002",
        "charge_customer": false,
        "total_debit": 0,
        "provider_code": "mtn_ug",
        "request_amount": 10000,
        "customer_name": "JOHN DOE",
        "transaction_status": "FAILED",
        "status_message": "Balance Insufficient for the transaction"
    }
}
```

{% endtab %}
{% endtabs %}
