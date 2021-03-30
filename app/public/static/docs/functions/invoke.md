---
title: Invoke Function
---

### Invoke a function

You can invoke the function using `curl` or any other HTTP framework you like. 

In order to call a function you will need a couple of things:
1. A URL to the function. This can be either `Synchronous` or `Asynchronous`. Documentation on how to get the url can be found [here](/docs/functions/manage#invocation-urls)
2. An Access token. Documentation on how to create an access token can be found [here](/docs/access_tokens/create)

The retrieved access token **MUST** be passed to the HTTP call as a header `X-Eywa-Token`.

Each invocation generates a unique `request id` that is returned as a header `X-Request-Id`. Use the value of from the header in order to look up for logs or execution timelines. Documentation on logs can be found [here](/docs/logs/overview) and documentation on execution timeline can be found [here](/docs/executions/overview). Simply passing the request id to the search bar will filter out logs and executions to that particular request.


#### Synchronous invocation

Synchronous executions are in a sense atomic meaning they either happen or not. Once invoked you will immediately receive the result of the execution. Supported HTTP methods are:
- `POST`
- `PUT`
- `PATCH`
- `DELETE`
- `GET`

#### Asynchronous invocation

Asynchronous invocations are similar to synchronous in a sense that you still **need** the token header `X-Eywa-Token` and it is still an HTTP request.
However, **you can only make a `POST` request**, every other method will return `405 Method Not Allowed`. Since asynchronous invocations are queued up by design, other than `POST` type does not make much sense. 

A successful async request will return `202 Accepted` which indicates that the request is successful and the execution has been queued. You can use the returned `X-Request-Id` header value to track your request.

When making asynchronous requests you can also pass a callback URL via header `X-Callback-Url`. If the function execution is successful and the callback url is present, the function result will be sent over to the url via a `POST` request. 

#### Curl Examples

For the sake of the documentation lets assume our URL is the following:
```
https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/
```

##### Basic GET request
```bash
curl https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/ \
-H 'X-Eywa-Token: {your_token_here}' 
```

##### GET request with a path
```bash
curl https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/some/path/my/function/can/handle \
-H 'X-Eywa-Token: {your_token_here}' 
```

##### GET request with params
```bash
curl https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/foo?bar=baz&baz=foo \
-H 'X-Eywa-Token: {your_token_here}'
```

##### POST request 
```bash
curl https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/ \
-H 'X-Eywa-Token: {your_token_here}' \
-d '{"key1": "value1", "key2": "value2"}' \
-X POST
```