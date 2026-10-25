
# NGN Bank Transfers

The DusuPay API supports for funds collections from NGN bank accounts. Customers complete payments by transferring funds to a specified bank account. The guide describes the two supported modes

## Overview

Before starting, your merchant account needs to have been approved enabled to do funds collections from Nigerian banks. We equally recommend that you go through the [Getting Started](/funds-collection/getting-started.md) section to have a high-level understanding of the funds collection process.

## Collection Workflow Selection

The API allows the merchant to either use the hosted payment page (recommended) or build their own user experience based on the API workflow. With the hosted payment page, DusuPay handles the end to end customer journey and this seamless workflow creates a good user experience. Some merchants however wish to create custom experiences using their brand guidelines and, in this situation, the direct charge can be used.

The table below briefly describes the two collection workflows to select from.

<table><thead><tr><th width="532.6217041015625">Scenario</th><th>Resource</th></tr></thead><tbody><tr><td>You wish to use a managed user experience that handles the customer experience without any input from you, and you don't mind generic brand display <em><strong>(RECOMMENDED)</strong></em></td><td><a href="/funds-collection/ngn-bank-transfers/ngn-bank-transfers-hosted-page.md">Hosted Page Workflow</a></td></tr><tr><td>You wish to have control of the user experience, page branding and directly call all the API endpoints yourself. (this requires a thorough understanding of the NGN Bank Transfer collection workflow as described on the page)</td><td><a href="/funds-collection/ngn-bank-transfers/ngn-bank-transfers-direct-charge.md">Direct Charge Workflow</a></td></tr></tbody></table>

{% hint style="info" %}
We strongly recommend that when the merchant chooses to use the direct charge workflow, they should first study the hosted page workflow before building their custom experiences. This ensures the customers have a smooth payment process
{% endhint %}
