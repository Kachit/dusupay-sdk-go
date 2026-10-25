
# Getting Started

Collections are payments that the merchant receives into their account. The sections here will highlight how the DusuPay platform supports these.

The DusuPay API supports funds collection via a number of methods e.g. Mobile Money, Bank Transfer, Card, Crypto and many other methods will be supported along the way. This is to allow for the flexibility of customer payments, such that customers can make payments in the ways that are most convenient. The funds collection section will describe the process flows for every available payment method. However, every collection will follow the steps described below, in the order of occurrence.

### The collection API request

The merchant will form the collection request payload (based on the selected payment method) and send it to the API for processing. If there are any errors with the parameters, these would be communicated immediately in the response otherwise the API would respond with an acknowledgement. Take the following the mobile money request and a sample acknowledgement.

```json
{
    "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
    "transaction_method": "MOBILE_MONEY",
    "currency": "UGX",
    "amount": 1000,
    "provider_code": "mtn_ug",
    "msisdn": "256777000001",
    "customer_name": "JOHN DOE",
    "description": "Test Collection",
    "allow_final_status_change": true
}

/*sample acknowledgement would be something like the following*/
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

However, should this process fail in any way even after the acknowledgement e.g. if the merchant is making too many requests to the same phone number, the following failure callback will be sent to the collection callback URL.

```json
{
    "event": "request.failed",
    "payload": {
        "id": 20760,
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
        "transaction_type": "COLLECTION",
        "request_currency": "UGX",
        "request_amount": 1000,
        "transaction_status": "FAILED",
        "status_message": "Failed to initialize the transaction. System exception!"
    }
} 
```

### Redirect to Payment URL - Optional

In special circumstances e.g. Mobile Money via the hosted page, bank, card and crypto collections, the initial response would have the payment URL to which the customer needs to be redirected to complete the transaction. A sample of this response payload is shared below.

```json
{
    "code": 202,
    "status": "accepted",
    "message": "Request Accepted",
    "data": {
        "internal_reference": "DUSUPAYZAERK2SBE6WJAG",
        "merchant_reference": "MCTREFQSSBHMHU3RCTLA",
        "payment_url": "https://devpay.dusupay.com/pay/collection/DUSUPAYZAERK2SBE6WJAG"
    }
}
```

{% hint style="info" %}
When the `payment_url` is part of the response, be sure to redirect the customer to that URL as that will be the only way for the customer to complete the transaction.
{% endhint %}

### Payment instructions notification - Optional

In special circumstances e.g. Mobile Money (direct charge) and Card (S2S), this notification is sent when the transaction is still in progress BUT the merchant needs to act on the instructions shared. Sometimes, the merchant will be required to submit some extra parameters from the customer or simply display the instructions to the user. These scenarios will be described in the specific document sections. Below is a sample payload

```json
{
    "event": "transaction.processing",
    "payload": {
        "id": 27594,
        "merchant_reference": "MCTREF9WVCPGFRHYCWSK",
        "internal_reference": "DUSUPAYSTPDKZ6MLY5WPQ",
        "transaction_type": "COLLECTION",
        "request_currency": "UGX",
        "transaction_amount": 5000,
        "transaction_currency": "UGX",
        "transaction_charge": 150,
        "transaction_account": "256777000005",
        "charge_customer": false,
        "total_credit": 4850,
        "provider_code": "mtn_ug",
        "request_amount": 5000,
        "customer_name": "JOHN DOE",
        "transaction_status": "PROCESSING",
        "status_message": "Transaction Initiated successfully. Approve Payment by entering the PIN",
        "authorization": {
            "mode": "NONE",
            "required_parameters": []
        },
        "payment_instructions": "<p>You will receive a prompt on the mobile number <b>256777000005</b>.<br/> Enter your PIN to authorize your payment of <b>UGX 5,000</b></p>"
    }
}
```

### Final status notification

This is the notification that communicates the very final status of the transaction, and this would mark the close of the transaction process. Sample completed and failed transaction callbacks/notifications shared below.

```json
/*sample completed transaction callback*/
{
    "event": "transaction.completed",
    "payload": {
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

/*sample failed transaction callback*/
{
      "event": "transaction.failed",
      "payload": {
        "id": 26609,
        "merchant_reference": "MCTREFNRFRTQA6SCWT5X",
        "internal_reference": "DUSUPAYWYUF8CR3ZRCGYU",
        "transaction_type": "COLLECTION",
        "request_currency": "UGX",
        "transaction_amount": 2000000,
        "transaction_currency": "UGX",
        "transaction_charge": 0,
        "transaction_account": "256777000002",
        "charge_customer": false,
        "total_credit": 0,
        "provider_code": "mtn_ug",
        "request_amount": 2000000,
        "customer_name": "JOHN DOE",
        "transaction_status": "FAILED",
        "status_message": "Balance Insufficient for the transaction"
      }
    }
```

### Transaction final status change

In rare conditions, it's possible for the final status of the transaction to change especially from **FAILED** to **COMPLETED**. In such situations, it's usually because the DusuPay automation could fail some transactions that take over 1hr without a final status from the last mile processor (telecom/bank/card processor etc) and then later, it is discovered that the customer actually made the payment. The default action from DusuPay is that the status is updated and the merchant updated accordingly with a new callback. However, some merchants do not like this behaviour and prefer that the final status never changes. Such merchants are advised to make use of the collection Boolean parameter `allow_final_status_change` and set the value to `false`. By doing this, DusuPay is instructed not to alter the final transaction status.

{% hint style="info" %}
**Cross-currency Collections**

If the request currency is different e.g. USD from the payment provider currency (transaction currency) e.g. KES, the API will convert the USD to KES and at the end of the transaction, your KES balance will be credited with the converted amount. More details [here](/appendix/cross-currency-transactions.md).
{% endhint %}
