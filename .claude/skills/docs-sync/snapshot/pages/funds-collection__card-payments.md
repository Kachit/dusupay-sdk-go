
# Card Payments

The DusuPay API supports card payment as described in the Supported Countries section. The API allows for direct charge (sending card details directly) or the hosted page that's powered by DusuPay.

## Prerequisites

1. The merchant needs to have an approved merchant account on the live platform. This is only achieved after the go-live requirements are met (being on sandbox and clearing the UAT tests).
2. The card payment options need to be configured on the designated merchant account. The support teams will be available to assist with this.

## Collection Workflow Selection

The API allows the merchant to either use the hosted payment page (recommended) or build their own user experience based on the API workflow. With the hosted payment page, DusuPay handles all the different variations of authorizing card payments e.g. OTP, PIN, address verification etc, and this seamless workflow creates a good user experience. Some merchants however wish to create custom experiences using their brand guidelines and, in this situation, the direct charge can be used.

The table below briefly describes the two card payment workflows to select from.

<table><thead><tr><th width="522">Scenario</th><th>Resource</th></tr></thead><tbody><tr><td>You wish to use a managed user experience that handles all card workflow variations without any input from you, even when issuer policies change from time to time to require other parameters. <em><strong>(RECOMMENDED)</strong></em></td><td><a href="/funds-collection/card-payments/hosted-payment-page.md">Hosted Payment Page</a></td></tr><tr><td>You wish to have control of the user experience, page branding and directly call all the API endpoints yourself. (this requires a thorough understanding of the card API workflow and the variations in authorization modes)</td><td><a href="/funds-collection/card-payments/direct-card-payment-s2s.md">Direct Card Payment (S2S)</a></td></tr></tbody></table>
