
# Authentication

The DusuPay API uses a key pair (Public and Secret Keys) for authentication. To make an authenticated request to our API, you need to pass custom headers as described in this section.

## Obtain the Keys

The Public and Secret keys can be copied from the merchant dashboard under the specific merchant account details as described [here](/getting-started/merchant-account-credentials.md). These keys should be sent as headers with the header names; **public-key** and **secret-key**

{% hint style="info" %}
Please ensure that your keys are stored safely and not shared with the public. In the event your keys are compromised, please contact us immediately for assistance. You an alternatively regenerate the secret key from within your merchant account dashboard
{% endhint %}

## Request Headers

For most of the API requests, only the public-key header is required. Scenarios that require the secret-key header will be described accordingly.

<table><thead><tr><th width="154.33333333333331">Header Name</th><th width="487">Value</th><th>Required</th></tr></thead><tbody><tr><td>Content-Type</td><td>application/json</td><td>YES</td></tr><tr><td>x-api-version</td><td>API version (currently 1)</td><td>YES</td></tr><tr><td>public-key</td><td>The Public Key as obtained from the merchant account settings</td><td>YES</td></tr><tr><td>secret-key</td><td>The Secret Key as obtained from the merchant account settings. Required for specific situations as the document will guide</td><td>NO</td></tr></tbody></table>
