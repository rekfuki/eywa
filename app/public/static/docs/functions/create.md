---
title: Create Function
---

### Creating Functions
> This page focuses on creating a function using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Functions/postFunctions)

Function creation page can be accessed by clicking on [Functions](/app/functions) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/functions/create)

&nbsp;
[![](/static/docs/functions/function_create_navbar_location.png)](/static/docs/functions/function_create_navbar_location.png)

&nbsp;
When the page with a list of function loads. You will see a `CREATE FUNCTION` button at the top right of the page

&nbsp;
[![](/static/docs/functions/function_create_button_location.png)](/static/docs/functions/function_create_button_location.png)

&nbsp;
Press the button and a page with the following form will appear

&nbsp;
[![](/static/docs/functions/function_create_form_details.png)](/static/docs/functions/function_create_form_details.png)

The form will have the following fields:

- `Function Name` - is required and must follow these rules: 
    - match `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$` regular expression.
    - be at least 5 characters long

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

- `Mounted Secrets` - select from the dropdown which secrets you would like to attach to your function. Documentation on how to create secrets can be found [here](/docs/secrets/create)


Once you have filled out the necessary fields, click `NEXT` to see environment variables configuration

[![](/static/docs/functions/function_create_form_env.png)](/static/docs/functions/function_create_form_env.png)

You can add more environment variables by clicking the `ADD` button which will insert additional `KEY:VALUE` pair field

[![](/static/docs/functions/function_create_form_env_entry.png)](/static/docs/functions/function_create_form_env_entry.png)

Both `KEY` and `VALUE` fields must follow the following rules:
   - minimum character length is 1
   - maximum character length is 255 

If you wish to remove environment variable entry, press the X icon next to it.

Once you press the `CREATE` button you should be redirected to your function management page. Documentation describing function management page can be found [here](/docs/functions/manage)