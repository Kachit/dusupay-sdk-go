
# Mobile Money Collection

The DusuPay API supports mobile money as described in the Supported Countries section. The API allows for direct charge or making use of the hosted page that's powered by DusuPay (recommended).

### Pre-requisites

1. The merchant needs to have an approved merchant account on the live platform. This is only achieved after the go-live requirements are met (being on sandbox and clearing the UAT tests).
2. The mobile money channels required by the merchant need to be configured on the designated merchant account. The support teams will be available to assist with this.

### Collection Workflow Selection

The API allows the merchant to either use the hosted payment page (recommended) or build their own user experience based on the API workflow. With the hosted payment page, DusuPay handles all the different variations of authorizing mobile money payments e.g. PIN entry, verifying phone numbers with OTP before completing transactions (e.g. for GHS Mobile Money collections.) and this seamless workflow creates a good user experience. Some merchants however wish to create custom experiences using their brand guidelines and, in this situation, the direct charge can be used.

The table below briefly describes the two Mobile Money workflows to select from.

<table><thead><tr><th width="529">Scenario</th><th>Resource</th></tr></thead><tbody><tr><td>You wish to use a managed user experience that handles all Mobile Money workflow variations without any input from you, even when telecom policies change from time to time to require other parameters. <em><strong>(RECOMMENDED)</strong></em></td><td><a href="/funds-collection/mobile-money-collection/mobile-money-hosted-page.md">Hosted Page Workflow</a></td></tr><tr><td>You wish to have control of the user experience, page branding and directly call all the API endpoints yourself. (this requires a thorough understanding of the Mobile Money API workflow and the variations in authorization modes)</td><td><a href="/funds-collection/mobile-money-collection/mobile-money-direct-charge.md">Direct Charge Workflow</a></td></tr></tbody></table>

{% hint style="info" %}
We strongly recommend that when the merchant chooses to use the direct charge workflow, they should first study the hosted page workflow before building their custom experiences. Merchants outside the mobile money processing countries are especially encouraged for try the hosted page first.
{% endhint %}
