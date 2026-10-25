
# Introduction

The DusuPay platform enables merchants (businesses and online platforms) to process payments across multiple payment options. The document describes the current version of the API.

Welcome to the DusuPay API. This API enables merchants (businesses and other online platforms) to integrate online payments into their business processes and support payments for goods and services. Our developer friendly API gives you access to mobile money, bank, card and crypto integrations across a number of countries. More payment options to be enabled in due course and the document will be updated accordingly. With this platform, merchants are able to request payments from valid mobile money subscribers, card holders, bank accounts and crypto wallets. Merchants can equally send funds via mobile money, bank accounts and soon to crypto wallets

## Platform Features

### Sandbox Environment

This is a utility platform built to allow merchants simulate the workflows of the production/live environment. The sandbox environment offers an API that the merchant integrates first to test the platform and perform UATs. It offers test accounts that the merchant can use to try out the different operations of the API without affecting customer funds in the real world. When all tests are done to the satisfaction of the merchant and the DusuPay support staff, the merchant can then be approved to go live.

{% hint style="info" %}
We intend to make your migration from sandbox to production very smooth. You'll strictly need to change the API base URL and obtain a new set of credentials (security keys). Everything else will work as expected.
{% endhint %}

### Production Environment

Under this environment, real funds can be requested from real customers and real funds sent to real customers. In order for the merchant to be accepted to use the production environment, they will have completed the compliance/KYC phase and also integrated the sandbox API for all the required tests.

### Async Processing

The DusuPay API is asynchronous; meaning the merchant doesn't need to wait on the same request thread to obtain the new status of the transaction. The confirmation of the transaction status will be via a callback/webhook/IPN POST request. We will therefore require that you share the **https** URLs for your callback endpoints and these will be configured on the merchant account accordingly. Different callback URLs are configured for collections and payouts.

### Transparency

DusuPay provides a management portal to all merchants. This is the central point to manage your business. It allows for creation of multiple merchant accounts and provides reports for all the activity on the different accounts. It also offers accounts statements to reflect how the merchant balances have changed over time. For every merchant account, API settings can be configured to support the integration.

### API Endpoints

The DusuPay API can be accessed using the following base URLs. Both URLs will only be accessible when the merchant account is created and approved on either environment.

<table><thead><tr><th width="151.77069091796875">Environment</th><th>API Base URL</th></tr></thead><tbody><tr><td>Sandbox/Test</td><td>https://sandboxapi.dusupay.com/</td></tr><tr><td>Production/Live</td><td>https://payments.dusupay.com/</td></tr></tbody></table>

### Postman Collection

The Postman Collection is a great way to quickly explore the DusuPay APIs without writing any code. The collection contains pre-formatted requests for the different transaction methods/options. Download the collection and add your API keys in order to start interacting with the sandbox APIs.

{% embed url="<https://postman.dusupay.com>" %}
