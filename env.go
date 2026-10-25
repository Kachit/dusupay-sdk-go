package dusupay

// ProdAPIUrl production ENV API url
const ProdAPIUrl = "https://payments.dusupay.com"

// SandboxAPIUrl sandbox ENV API url
const SandboxAPIUrl = "https://sandboxapi.dusupay.com"

// SandboxPublicKey DusuPay sandbox public key (RSA signature verification and card data encryption).
// See https://developer.dusupay.com/getting-started/dusupay-public-keys.md
const SandboxPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAzg8OSczAV7IYVsb5q7Ke
1ipnVkbL0fQdb8tS+0+JsZFhgZwwQopvuAgtA3XQh8Mlq1ul7iIl99qII8WHD96a
HRg1w8JxHkMP5glCFAIrRp/rRn72W3MC2O0utKiKyihEaW6ObB1eQKzuYfXAP1HO
dTnal+SzBjB2LyMFf4M7eE7Yksa9rdjFvIpfUkWSBlG028AxRrUI5fvem07TMK5i
I3GOv7xH/a3siaT2rQ+dZ83G+k7nCL/jNlekjD7QBPhRCExHGp6CKici74nWPZB+
cbCCObnBYfoJI+5AOwRjvimd8lBc/zBEL4Q5ElUG5UhPuiVMmjiZ4oEzbnIIswFx
ujFNuvqd8W3nv86OuTwUTpDszan1EHjoKY4wJ/wKJg38+3Gtfy0jzuaE0F71kvAy
j6yhM05hm+oybQeWkQFr/G1CiefG3uGHgE07lkU3PegdkRUHCDBhhkjgoaguHzna
Q8ReJDnRPlezncsW9VQyN10vuWIwLaYJvUqEBZdsp7/rIH4P81m6lI0vH2DG9ns0
acFA6hK+uIjzwiEY/n1NfdjRHpvYBSA493NcVbKGZ3yVHsvPvVSxkE/jX+P/0zVZ
DGu2LLv3IQilVElvjhT5qNeFBHG+2E9hS7Cqd+KsxNYpuRlow7NCiAflKx7/AeO2
P2DOy3sG/3E+u3OngnED5rUCAwEAAQ==
-----END PUBLIC KEY-----
`

// ProdPublicKey DusuPay production public key (RSA signature verification and card data encryption).
// See https://developer.dusupay.com/getting-started/dusupay-public-keys.md
const ProdPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0hPd57qbKYvZLwyqfZ3H
teqJ6VGniXjw/QWB4Ja4SRul7h7ShMwdDAyOZ1nZoFEJnIs95agoJdjh7sQfBOrB
YznHAK6tR/a8HP0z4VjN3V1d0wNDTLnx942XB+SnKWzQ+8atqKqBIsiD8ATWhryV
/AmLC3e8w38JP30qzV3xLVE469hY+Qh0rpbHNX1X/kpX2647G94Yj9rlPi7vggWz
qd3icwzjiEqTUr6LytUYLbgUEHfm9DDOEe96DUzZGCP2St0OhrBlAEMvD0md4X6Z
6GbXKrTz0Z8RF3dKm3R6k4VBBd390/qbd5Qt0BsAThK5qaFohhAtUF/DZCoA2bfl
K0pwQ43C3vfXcmew0+CzOAtowNbJ0ELqkIqoT1TD5OaVz6rMhLtBuSsfg233Bzup
iDqzltjWCqNfbNOZ8JKQQj9RQTwQpT/+izEdHjkM+HmMH4j//5/nuqzMYFcFkS00
mRDFhdSH1kORqs/Jj5pmBUbxeMnBXeKMh37maQBub5QygSkFDRmBz0pEqLynM/Yk
6YbgDqYDI7GWVFxh3c7TyLaixQbISZToTDEtnaQwFQOhNn6cVxFaITWJyftMFwSX
1oqrLjwnhi13k2eposPQibZgeY/s1xmWK7yYsC/GHjoZkNxXNvGafsIFuYD+793M
CNLBHhKqgIfU3YUiOKA1ne8CAwEAAQ==
-----END PUBLIC KEY-----
`
