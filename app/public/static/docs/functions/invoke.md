### Invoke a function

In order to call a function you will need a couple of things:
1. A URL to the function. This can be either `Synchronous` or `Asynchronous`. Documentation on how to get the url can be found [here](/docs/invocation-urls)
2. An Access token. Documentation on how to create an access token can be found [here](/docs/access_tokens/create)

Each invocation generates a unique `request id` that is returned as a header `X-Request-Id`. Use the value of from the header in order to look up for logs or execution timelines. Documentation on logs can be found [here](/docs/logs/overview) and documentation on execution timeline can be found [here](/docs/executions/overview). Simply passing the request id to the search bar will filter out logs and executions to that particular request.

Once you have the url and the token you can invoke the function as such
> Note: These are just `curl` examples, you can use anything equivalent to make the requests

#### Basic GET request
```bash
curl -H 'X-Eywa-Token: {{your_generated_access_token}}' {{your_function_url}}
```

#### GET request with a path
```bash
curl -H 'X-Eywa-Token: {{your_generated_access_token}}' {{your_function_url}}/some/path/that/my/function/can/handle
```