---
title: Update Function
---

### Update Function
> This page focuses on updating a function using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Functions/putFunctionsFunctionid)

In order to update the function you will need to access the function update page. The page can be accessed from the function management page. Documentation on how to access function management page can be found [here](/docs/functions/manage)

Once you are on the function management page you will see an `EDIT` button at the top right of the page:

[![](/static/docs/functions/function_update_location.png)](/static/docs/functions/function_update_location.png)

Once you press the button, a page with the following form will be displayed:

&nbsp;  
[![](/static/docs/functions/function_update_form.png)](/static/docs/functions/function_update_form.png)

&nbsp;  
The form contains the following fields:

- `Image` - select one of the available images. Images that are building or their build resulted in an error will be grayed out.
- `Minimum Replicas` - the minimum instances of your code running. Setting it to zero will despawn all function instances when idle. The field must follow these rules:
    - minimum allowed value is 0
    - maximum allowed value is 100
    - value must bet lesser than or equal to `Maximum Replicas`
- `Maximum Replicas` - the maximum instances of your code running. The field must follow these rules:
    - minimum allowed value is 1
    - maximum allowed value is 100
    - value must bet greater than or equal to `Minimum Replicas`
- `Scaling Factor` - the percentage of `Maximum Replicas` to be used as a step when downscaling function instances. The field must follow these rules:
    - minimum allowed value is 0. Setting the value to zero will disable gradual downscaling and all replicas will be despawned at once
    - maximum allowed value is 100
- `Per Instance Concurrency` - limit the maximum number of requests in flight. Setting the value to zero means that each instance of your function can handle at most one request at the time. The field must follow these rules:
    - minimum allowed value is 0 (disabled concurrency)
- `Read Timeout` - the amount of seconds the function instance will wait while reading the request. After the timeout is reached it will terminate reading the request body and return an error. The field must follow these rules:
    - minimum allowed value is 0
- `Write Timeout` - the amount of seconds the function instance will try to write the response for. After the timeout is reached it will terminate writing the response and return an error. The field must follow these rules:
    - minimum allowed value is 0
- `Environment Variables` - key value pairs that can be added ar removed. In order to add additional entries, press the `ADD` button. If you wish to remove an entry, press the `X` button on the right hand side. Environment variables must follow these rules:
   - minimum character length is 1
   - maximum character length is 255 
- `Mounted Secrets` - select from the dropdown which secrets you would like to attach to your function. Documentation on how to create secrets can be found [here](/docs/secrets/create)
- `Write Debug` - enabled debug output. If the function has any `stderr` or `stdout` logging, enabling this option will write those outputs to the logs. Documentation on function scoped logs can be found [here](/docs/functions/manage#metrics) and global logs can be found [here](/docs/logs/overview)

&nbsp;  
Once you hit `UPDATE` your changes will be registered and you will be redirected back to [function management page](/docs/functions/manage)