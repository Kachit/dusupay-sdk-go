
# Callback Events

Every callback payload has a parameter named event and its purpose is to communicate the reason for the callback. The section below describes the various values the parameter can take.

Depending on the workflow (collection, invoice generation/payment, disbursement etc) being used, the gateway will send appropriate events to the appropriate callback URLs. For every workflow type, the applicable events are described below;

### Collections

<table><thead><tr><th width="218">Event</th><th width="176">Transaction Status</th><th>Description</th></tr></thead><tbody><tr><td>request.failed</td><td>FAILED</td><td>Request failed and could not be processed by the gateway. Failure reason is part of the callback payload</td></tr><tr><td>transaction.charges</td><td>PENDING</td><td>Communicates success of the first collection leg (initialization) and the charges to be levied. The collection confirmation request should be sent after this event is received</td></tr><tr><td>transaction.processing</td><td>PROCESSING</td><td>Notifies the merchant that transaction processing is ongoing. Returns mobile money instructions (when payment method is MOBILE_MONEY) or the card payment URL (when payment method is CARD)</td></tr><tr><td>transaction.failed</td><td>FAILED</td><td>The gateway tried to process the transaction but failed e.g. insufficient balance on the account. Failure reason is part of the callback payload</td></tr><tr><td>transaction.completed</td><td>COMPLETED</td><td>The gateway processed the collection successfully.</td></tr></tbody></table>

### Invoices (Generation and Payment)

<table><thead><tr><th width="180">Event</th><th>Description</th></tr></thead><tbody><tr><td>request.failed</td><td>Request failed and could not be processed by the gateway. Failure reason is part of the callback payload</td></tr><tr><td>invoice.generated</td><td>Event sent at successful generation of the invoice. Details of the invoice would be shared as part of the payload</td></tr><tr><td>invoice.paid</td><td>Event sent when the invoice is paid for. Part of the payload is the reference for the collection that fulfilled the invoice</td></tr></tbody></table>

### Payouts/Disbursements

<table><thead><tr><th width="218">Event</th><th width="175">Transaction Status</th><th>Description</th></tr></thead><tbody><tr><td>request.failed</td><td>FAILED</td><td>Request failed and could not be processed by the gateway. Failure reason is part of the callback payload</td></tr><tr><td>transaction.failed</td><td>FAILED</td><td>The gateway tried to process the transaction but failed. Failure reason is part of the callback payload</td></tr><tr><td>transaction.completed</td><td>COMPLETED</td><td>The gateway processed the payout successfully.</td></tr></tbody></table>
