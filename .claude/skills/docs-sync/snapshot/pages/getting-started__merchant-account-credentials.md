
# Merchant Account Credentials

For every approved merchant account, a set of credentials is automatically generated and stored on the record. This section describes how the merchant can manage credentials after the account creation

As described [here](/getting-started/authentication.md), the API requires that certain headers are sent for the various requests. To be specific, the `public-key`header is mandatory for all the API requests and the `secret-key` is required only in a few situations as the documentation will advise along the way. The other critical variable is the **signing key**, a randomly generated string stored on the merchant account record, whose sole purpose is to support the generation of the HMAC Hash signature sent alongside the merchant callbacks.

When the merchant account is created and approved by the admins, all the above 3 mentioned values are generated and securely stored. The public-key (as its name suggests) is displayed among the merchant account details in plain text. The signing key too is displayed in plain text for the merchant to view it. The secret key on the other hand is stored as a strongly hashed value and the gateway team never gets to see the plain text version of this secret key. The merchant should therefore generate a new secret key from their dashboard in order to temporarily see the secret key in plain text and get a chance to store is somewhere safe.

The merchant is at liberty to re-generate just the secret key OR all the keys mentioned above that way in case of compromise, the merchant has full control of the keys change without needing the gateway support team. The table below advises further.

<table><thead><tr><th width="537">Scenario/Use case</th><th>Resource</th></tr></thead><tbody><tr><td>If you wish to generate only the secret key and the other keys remain the same</td><td><a href="/getting-started/merchant-account-credentials/generate-secret-key.md">Generate Secret Key</a></td></tr><tr><td>If you wish to generate a fresh set of credentials and replace all the existing ones</td><td><a href="/getting-started/merchant-account-credentials/regenerate-security-keys.md">Regenerate Security Keys</a></td></tr></tbody></table>

{% hint style="info" %}
We strongly recommend routine change of the credentials as a security measure especially in situations where integration work is outsourced to parties external to the merchant/organization.
{% endhint %}
