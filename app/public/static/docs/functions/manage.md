---
title: Manage Function
---

### Manage Function
> This page focuses on managing the function using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api)

Function management page can be accessed by navigating to [Functions](/app/functions) and selecting one of your functions from the list by pressing on the function ID

[![](/static/docs/functions/function_create_navbar_location.png)](/static/docs/functions/function_create_navbar_location.png)

&nbsp;  
When the page loads, select the function from the list you would like to manage.

[![](/static/docs/functions/function_manage_list_select.png)](/static/docs/functions/function_manage_list_select.png)
&nbsp;  

#### Invocation URLs
Once the page is loaded you will see a lot of different information. At the top of the page you will see a couple of links labeled `Sync URL` and `Async URL`.

[![](/static/docs/functions/function_manage_view_urls.png)](/static/docs/functions/function_manage_view_urls.png)

`Sync URL` is used for synchronous invocations and `Async URL` is used for asynchronous invocations.

#### Update and Delete
On the right hand side you will see two buttons: `DELETE` and `UPDATE FUNCTION`. When delete button is pressed a modal will appear to confirm the deletion, once you confirm the function will be scheduler for deletion and you will be redirected to function list view.
The documentation for updating the function can be found [here](/docs/functions/update) and the documentation for deleting a function can be found [here](/docs/functions/delete)


#### Overview
Right below that you will see different tabs: `DETAILS`, `TIMELINES`, `LOGS`, `METRICS`. On page load the default is always `DETAILS`

[![](/static/docs/functions/function_manage_view_tabs.png)](/static/docs/functions/function_manage_view_tabs.png)


#### Details

The details tab which is visible by default has three components to it. The details of the function can be updating using function update page. This documentation about the process can be found [here](/docs/function/update)

##### Function Information

[![](/static/docs/functions/function_manage_view_details_info.png)](/static/docs/functions/function_manage_view_details_info.png)

The function info card contains configuration fields as well as the status of the function which is either `AVAILABLE` or `UNAVAILABLE`.

`AVAILABLE` indicates that everything is green and ready to go. `UNAVAILABLE` indicates something is not right, possible scenarios are:
- function just deployed and it takes some time to become ready
- function is in the middle of an update
- internal Eywa error


The rest of the fields are:

- `ID` - the id (UUID) of the function
- `Name` - the name of the function given during deployment
- `Image ID` - the id (UUID) of the image which is deployed. Clicking on it will take you to the image details page.
- `Image Name` - the name of the image
- `Available Replicas` - available replicas at the current time. This number will vary depending on configuration and load
- `Min Replicas` - minimum replicas the function is allowed to scale down to
- `Max Replicas` - maximum replicas the function is allowed to scale up to
- `Scaling Factor` - the percentage step of `Max Replicas` by which the function will be scaled down after load
- `Per Instance Concurrency` - maximum allowed concurrent requests that can be handled by each function replica
- `Debug Mode` - enables `stderr` and `stdout` logging. Can only be updated during [function update process](/docs/function/update)
- `Write Timeout` - write timeout is the amount of time in seconds the function will attempt to write the response for before cancelling it
- `Read Timeout` - read timeout is the amount of time in seconds the function will attempt to read the request for before cancelling it


##### Environment Variables

[![](/static/docs/functions/function_manage_view_details_env.png)](/static/docs/functions/function_manage_view_details_env.png)

The environment variables card contains all the environment variables which are associated to your deployed functions. At the very minimum there will be at least one environment variable called `mongodb_host` which is used to access `MongoDB` database. More details on the database access can be found [here](/docs/database/access)

The rest of the environment variables will depend on how you have configured your function deployment


##### Mounted Secrets

[![](/static/docs/functions/function_manage_view_details_secrets.png)](/static/docs/functions/function_manage_view_details_secrets.png)

The mounted secrets card card contains all the secrets which have been made available as files to your deployed functions. At the very minimum there will be at least one mounted secret called `{prefix}-mongodb-credentials` which is used to access `MongoDB` database. More details on the database access can be found [here](/docs/database/access)

The rest of the secrets will depend on how you have configured your function deployment. Documentation on accessing mounted secrets can be found [here](/docs/secrets/access)

#### Timelines

The execution timelines are the same as the ones globally available except they are scoped to the function. Documentation on executions can be found [here](/docs/executions/overview)
[![](/static/docs/functions/function_manage_view_timelines.png)](/static/docs/functions/function_manage_view_timelines.png)

#### Logs

The logs are the same as the ones globally available except they are scoped to the function. Documentation on logs can be found [here](/docs/logs/overview)
[![](/static/docs/functions/function_manage_view_logs.png)](/static/docs/functions/function_manage_view_logs.png)

#### Metrics

More details on the function metrics can be found [here](/docs/functions/metrics)

[![](/static/docs/functions/function_manage_view_metrics.png)](/static/docs/functions/function_manage_view_metrics.png)

