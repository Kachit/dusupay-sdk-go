
# Handling Notifications/Callbacks

Whenever the status of the transaction changes (to processing, failure or success), we notify your service via callbacks. Below are the hints to take note of as regards these payment notifications.

* All callbacks are sent to the callback URL you provide during registration.
* A secure callback URL is required and notifications will ONLY be sent to **https** URLs. If you don't receive these notifications, ensure to check this.
* The callback request is a POST with JSON serialized data.
* We expect that your service will acknowledge the callback request by responding with HTTP code **200** otherwise we will retry a few more times (if the response code is below the 500 range) and then stop.
* The callback request includes two extra headers; **rsa-signature** (an RSA generated signature) and **hmac-signature**, and these can optionally be used to verify that the callback request originated from our servers. Verification can be done for any of these two signatures, NOT necessarily both of them. If the signature is valid, you can proceed with your business logic.
* The HMAC signature can be verified using [this](/callbacks/hmac-signature-verification.md) procedure while the RSA (DusuPay) signature can be verified using [this](/callbacks/rsa-signature-verification.md) procedure. The merchant can choose which procedure to work with.

{% hint style="info" %}
All **production** notifications/callbacks will originate from this IP address below. Whenever more IPs are enabled, you'll be notified, and the document updated accordingly. We recommend that the IP (s) are whitelisted on your platform and only accept callback requests originating from them.

**161.35.164.139**
{% endhint %}

{% hint style="warning" %}
If you would like to equally whitelist IPs for the **sandbox** notifications/callbacks, consider the IP **165.227.128.244** for that as all sandbox callback requests will originate from that IP
{% endhint %}
