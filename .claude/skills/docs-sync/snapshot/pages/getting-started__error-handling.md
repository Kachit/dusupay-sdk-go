
# Error Handling

The DusuPay API uses HTTP status codes to indicate the success or failure of the API requests. The section describes the summary of these codes and what they mean in our context.

## Summary

* Codes in the 2XX range mean that the API request was processed successfully.
* Codes in the 4XX range mean that something was wrong with the data that you sent. For example, you might have missed some required parameters/headers, or you might be using the wrong API credentials.
* Codes in the 5XX range indicate an error in processing on our end

## Common HTTP Codes

<table><thead><tr><th width="132.57142857142856">Code</th><th>Description</th></tr></thead><tbody><tr><td>200</td><td>OK - Request was successful</td></tr><tr><td>202</td><td>Accepted - Request has been accepted and acknowledged. We will now go ahead to process the request and notify you of the status afterwards.</td></tr><tr><td>400</td><td>Bad Request - Malformed request or missing required parameters</td></tr><tr><td>401</td><td>Unauthorized - Missing required headers, wrong Public or Secret Key etc</td></tr><tr><td>403</td><td>Forbidden - You are trying to access a resource for which you don't have proper access rights.</td></tr><tr><td>404</td><td>Not Found - You are trying to access a resource that does not exist</td></tr><tr><td>422</td><td>Unprocessable Entity - You provided all the required parameters, but they are not proper for the request</td></tr><tr><td>500</td><td>Internal Server Error - We had a glitch in our servers. Retry the request in a little while or contact support</td></tr><tr><td>503</td><td>Service Unavailable – We are temporarily offline for maintenance. Please try again later.</td></tr></tbody></table>
