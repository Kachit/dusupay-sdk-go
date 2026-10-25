
# RSA Signature Verification

The section describes how the RSA signature sent in the callback header can be verified. The signature is generated using RSA Signing. For verification to succeed, the public key is required.

## Obtain the Public Key - Based on environment

Copy the public key for the environment you're working with from [here](/getting-started/dusupay-public-keys.md). This document assumes that the public key would be stored somewhere on your server under the name `dusupay.public.key.pem`

## Next Steps

Below is the sample callback data for this demonstration

```json
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
```

1. Obtain the value of the `rsa-signature` header (if callback) OR the value of the `rsa_signature` query parameter (if redirect).
2. Form the string payload to be used in signature verification. This is obtained by concatenating values of the callback/redirect data in the format; `event:merchant_reference:internal_reference:transaction_type:transaction_status` and these values are obtained from the callback/redirect data. The string payload would therefore be `transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED`
3. Use the public key obtained above to verify the signature as described in the sample source codes below.

{% tabs %}
{% tab title="PHP" %}

```php
<?php

public function isValidSignature() {
    $file = "path-to-file/dusupay.public.key.pem";
    $keyContent = file_get_contents($file);
    $publicKey = openssl_get_publickey($keyContent);
    $strPayload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED";
    $signature = base64_decode("value-of-rsa-signature");

    /*true or false*/
    return openssl_verify($strPayload, $signature, $publicKey, "sha256") == 1;
}

?>
```

{% endtab %}

{% tab title="NodeJS" %}

```javascript
const crypto = require('crypto');
const fs = require('fs');

function isValidSignature() {
    const strPayload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED";
    const signature = "value-of-rsa-signature";
    const publicKeyFile = "path-to-file/dusupay.public.key.pem";
    const publicKey = fs.readFileSync(publicKeyFile).toString().replace(/\\n/g, '\n');

    const verify = crypto.createVerify("SHA256");
    verify.write(strPayload);
    verify.end();

    /*true or false*/
    return verify.verify(publicKey, signature, 'base64');
}
```

{% endtab %}

{% tab title="Java" %}

```java
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.PublicKey;
import java.security.Signature;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

public class SignatureVerifier {

    public boolean isValidSignature() throws Exception {
        // Read public key from file
        Path pathToFile = Paths.get("path-to-file/dusupay.public.key.pem");
        byte[] keyBytes = Files.readAllBytes(pathToFile);

        // Decode public key
        X509EncodedKeySpec keySpec = new X509EncodedKeySpec(keyBytes);
        KeyFactory keyFactory = KeyFactory.getInstance("RSA");
        PublicKey publicKey = keyFactory.generatePublic(keySpec);

        // Signature and payload
        String strPayload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED";
        byte[] signatureBytes = Base64.getDecoder().decode("value-of-rsa-signature");

        // Verify signature
        Signature signature = Signature.getInstance("SHA256withRSA");
        signature.initVerify(publicKey);
        signature.update(strPayload.getBytes());
        return signature.verify(signatureBytes);
    }

    public static void main(String[] args) throws Exception {
        SignatureVerifier verifier = new SignatureVerifier();
        System.out.println(verifier.isValidSignature());
    }
}

```

{% endtab %}

{% tab title="C#" %}

```csharp
using System;
using System.IO;
using System.Security.Cryptography;
using System.Text;

public class SignatureVerifier
{
    public bool IsValidSignature()
    {
        // Read public key from file
        string pathToFile = @"path-to-file\dusupay.public.key.pem";
        string keyContent = File.ReadAllText(pathToFile);

        // Decode public key
        RSA rsa = RSA.Create();
        rsa.ImportFromPem(keyContent);

        // Signature and payload
        string strPayload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED";
        byte[] signatureBytes = Convert.FromBase64String("value-of-rsa-signature");

        // Verify signature
        var verifier = new RSAPKCS1SignatureDeformatter(rsa);
        verifier.SetHashAlgorithm("SHA256");
        byte[] payloadBytes = Encoding.UTF8.GetBytes(strPayload);
        byte[] sha256Hash;

        using (SHA256 sha256 = SHA256.Create())
        {
            sha256Hash = sha256.ComputeHash(payloadBytes);
        }

        return verifier.VerifySignature(sha256Hash, signatureBytes);
    }

    public static void Main(string[] args)
    {
        SignatureVerifier verifier = new SignatureVerifier();
        Console.WriteLine(verifier.IsValidSignature());
    }
}

```

{% endtab %}

{% tab title="Python" %}

```python
import base64
import hashlib
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.serialization import load_pem_public_key

def is_valid_signature():
    # Read public key from file
    with open("path-to-file/dusupay.public.key.pem", "rb") as key_file:
        key_content = key_file.read()

    # Decode public key
    public_key = load_pem_public_key(key_content, backend=default_backend())

    # Signature and payload
    str_payload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED"
    signature_bytes = base64.b64decode("value-of-rsa-signature")

    # Verify signature
    verifier = public_key.verifier(
        signature_bytes,
        padding.PKCS1v15(),
        hashlib.sha256
    )
    verifier.update(str_payload.encode())

    try:
        verifier.verify()
        return True
    except Exception as e:
        print("Signature verification failed:", e)
        return False

```

{% endtab %}

{% tab title="Ruby" %}

```ruby
require 'openssl'
require 'base64'

def is_valid_signature
  # Read public key from file
  file = File.read("path-to-file/dusupay.public.key.pem")

  # Decode public key
  public_key = OpenSSL::PKey::RSA.new(file)

  # Signature and payload
  str_payload = "transaction.completed:MCTREFT2WMNWZ23SBN6Y:DUSUPAYRMGRXNNYBWATKJ:COLLECTION:COMPLETED"
  signature_bytes = Base64.decode64("value-of-rsa-signature")

  # Verify signature
  sha256 = OpenSSL::Digest::SHA256.new
  result = public_key.verify(sha256, signature_bytes, str_payload)

  result
end

```

{% endtab %}
{% endtabs %}

Below is an example of a generated RSA signature (rsa-signature) on the sandbox environment. Feel free to verify it using the sample payload above and the public key.

{% code overflow="wrap" %}

```
bX6aa5avegrpqS7I9eK6OZXD14RkuEuKCkRgWHIGnXsI2Vr2O/NrRXUHLzaTm0ujk48hLdGV6dVxzapkB3NBOtHT8cntyXLX7rw5PsYg2FDsUT/7l5uH0xugG7JzLuY1J72GVKwNBMRu6VEJo9c43OLSha5Lldket7NBNub5ZO6zD0PIN6fSblp/FGuycfasxbjilpKwdXH9x1EDOOW2Ds4gyHesZjON/jiQuvs/qeekXqBCRVh/5D0tKfGoAGaug5sa3hOhSyskeiHihslpPBs5Yw4xy7ajuHNKCRlq6w8PsTW6RXsbcTGJbowpGv1j1xJoUxjVCIr5vJtwHIKXU9v1S8pZLc4Z5nYV/ejCSZZrMYB6Fg5n0bCz4zgUTbSmES3G/0Cthb/40zbDX83reFxweLLqO1VRQVd+3CFijLSalKmtvzFkAP2SSN6NpC/9yJfvOnq7sXt3u4igMsIkChzSYF/OZksZlCyC4q7zTR5k+7XJmpuE4FeYKWDaOSFtqRYEfsdA7+VgPdHaMGLW9wwBKA867DBvNuZC0lP5wEOkbqCVOmVtcbbXDDKUo3nybNRMky9Lz7fptIrdB8m7ZVdVkLZ6OZIM3vq3qx2y8cb5jqO8lUmsNIf+BJtVxhlz7kwc2NQEC34ga92wQoDVQWWv9211tF7fCOSnkA8QSc0=
```

{% endcode %}
