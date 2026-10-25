
# Cross Currency Transactions

The DusuPay API allows for cross currency transactions whereby the request currency does not match the currency of the selected payment option. The platform handles the conversion appropriately.

The payment options provided by the DusuPay API have specific transaction currencies as seen from the response [here](/utility-functions/payment-options.md). These currencies therefore correspond to the currencies in which DusuPay is able to hold balances. The platform however allows the merchant to send any currency in the collection/payout requests for example the request could be for USD 100 but the selected payment option has transaction currency as GHS or XAF etc.

The DusuPay API has a robust workflow to ensure the conversion happens to the transaction currency/amount since that's the currency the payment provider supports. When the transaction completes successfully, the merchant balance would be updated in the transaction currency. The API uses the prevailing market exchange rates to do the conversion.

Consider a sample Mobile Money collection request where the requested currency is USD and the payment option is mtn\_ug (MTN Mobile Money Uganda). In the example, we assume the conversion rate to be 1 USD = 3800 UGX and that the transaction charge is 1%

```json
{
    "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
    "transaction_method": "MOBILE_MONEY",
    "currency": "USD",
    "amount": 100,
    "provider_code": "mtn_ug",
    "msisdn": "256777000001",
    "customer_name": "JOHN DOE",
    "description": "Test Collection"
}
```

On completion, the callback would look like the following. Take critical note on the callback parameters; `request_currency`, `request_amount`, `transaction_currency` and `transaction_amount`

```json
{
    "event": "transaction.completed",
    "payload": {
        "id": 20760,
        "merchant_reference": "MCTREFT2WMNWZ23SBN6Y",
        "internal_reference": "DUSUPAYRMGRXNNYBWATKJ",
        "transaction_type": "COLLECTION",
        "request_currency": "USD",
        "request_amount": 100,
        "transaction_currency": "UGX",
        "transaction_amount": 380000,
        "transaction_charge": 3800,
        "transaction_account": "256777000001",
        "charge_customer": false,
        "total_credit": 376200,
        "provider_code": "mtn_ug",
        "customer_name": "JOHN DOE",
        "transaction_status": "COMPLETED",
        "status_message": "Transaction Completed Successfully"
    }
}
```

{% hint style="info" %}
This workflow gives merchants the flexibility to bill their customers in other currencies, while allowing the customers to pay/get paid in the currencies at their disposal, and using the payment options that are most convenient. This workflow applies to both collections and payouts
{% endhint %}
