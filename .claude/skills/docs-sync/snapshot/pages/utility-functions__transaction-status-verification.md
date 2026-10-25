
# Transaction Status Verification

In some situations, the merchant account might to confirm the status of the transaction. This document resource describes the process of retrieving the transaction details (status inclusive)

Whenever we get the final transaction status from the processor (telecom/bank/card processor), we'll send the notifications to your designated callback URLs and these requests will be retried using the exponential back-off strategy. This guarantees that even when your URL is unreachable OR the http request fails, there will be automated retries at different intervals. Our belief is that; if we don't have the final status yet, checking status from the merchant side doesn't help. However, if status check is crucial for your workflow, please note the following.

1. Status check requests will be denied if they're done less than **2 minutes** from the time the transaction was initiated.
2. We recommend that an interval of **5 minutes** is implemented between status check requests for the same transaction.

### Get Transaction Status

<mark style="color:blue;">`GET`</mark> `https://sandboxapi.dusupay.com/data/transaction/verify/{reference}`

Returns the transaction details and the parameter `data.transaction_status` holds the current status of the transaction. Replace `{reference}` with the merchant reference that was used prior when initiating the transaction.

#### Headers

<table><thead><tr><th width="150">Name</th><th width="107">Type</th><th>Description</th></tr></thead><tbody><tr><td>public-key<mark style="color:red;">*</mark></td><td>String</td><td>The merchant account Public Key</td></tr></tbody></table>

### Sample Transaction Status Request

```powershell
curl https://sandboxapi.dusupay.com/data/transaction/verify/MCTREFT2WMNWZ23SBN6Y \
   -H "Accept: application/json" \
   -H "x-api-version: 1" \
   -H "public-key: your-public-key"
```

{% tabs %}
{% tab title="200: Request successful" %}

```json
{
    "code": 200,
    "status": "success",
    "message": "Request completed successfully.",
    "data": {
        "id": 20760,
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
        "transaction_type": "COLLECTION",
        "request_currency": "UGX",
        "transaction_amount": 2000000,
        "transaction_currency": "UGX",
        "transaction_charge": 60000,
        "transaction_account": "256787008803",
        "charge_customer": false,
        "total_credit": 1940000,
        "provider_code": "mtn_ug",
        "request_amount": 2000000,
        "customer_name": "JOHN DOE",
        "transaction_status": "COMPLETED",
        "status_message": "Transaction Completed Successfully"
    }
}
```

{% endtab %}

{% tab title="404: Transaction Not Found" %}

```json
{
    "code": 404,
    "status": "error",
    "message": "Transaction with reference: MCTREF2BUBWGEVDZ47DQFH1 Not Found. Contact Support",
    "data": {}
}
```

{% endtab %}

{% tab title="403: Request before 2mins grace period" %}

```json
{
    "code": 403,
    "status": "error",
    "message": "Status check denied. Try again after 2 minutes.",
    "data": {}
}
```

{% endtab %}

{% tab title="429: Violation of recommended interval" %}

```json
{
    "code": 429,
    "status": "error",
    "message": "Request Failed. Too many requests in a short time.",
    "data": {}
}
```

{% endtab %}
{% endtabs %}

### Transaction Status Descriptions

<table><thead><tr><th width="190">Status</th><th>Description</th></tr></thead><tbody><tr><td>PENDING</td><td>The transaction has been logged on the DusuPay platform, awaiting escalation to the last mile processor (telecom, bank, card processor etc)</td></tr><tr><td>PROCESSING</td><td>The transaction is currently being processed and the DusuPay platform is waiting for the final status from the last mile processor</td></tr><tr><td>FAILED</td><td>The transaction is failed. Reason for failure will always be in the <code>data.status_message</code> parameter.</td></tr><tr><td>COMPLETED</td><td>The transaction was processed successfully.</td></tr><tr><td>CANCELLED</td><td>The transaction was cancelled. This is usually done by the customer.</td></tr></tbody></table>
