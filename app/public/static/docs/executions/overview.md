#### Executions
> This page focuses on looking at executions using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=execution-tracker#/Timeline/getTimeline)

Executions are like timelines of each function invocation. Whenever you make a request to your function either synchronously or asynchronously, the entire journey of your request is tracked throughout the system until the result is returned.

The execution page can be access by clicking on the `Executions` in the side navigation bar or by navigating there via a direct [link](/app/timelines).

[![](/static/docs/executions/executions_navbar.png)](/static/docs/executions/executions_navbar.png)

&nbsp;  
Once the page is loaded you should see a page with the following content

&nbsp;  
[![](/static/docs/executions/executions_list.png)](/static/docs/executions/executions_list.png)

Each timeline entry row has the following columns: 
- `Request ID` - the id (UUID) of the request made. These are generated every time a new invocation is made
- `Function ID` - the id (UUID) of the function that was invoked
- `Function Name` - the name of the function
- `Status` - the status of the execution represented in HTTP status codes
- `Duration` - how long it took to execute the invocation. Asynchronous invocations will naturally have longer durations due to queueing
- `Created at` - the timestamp of the invocation

Timelines can be filtered by `Request ID`, `Function ID`, `Function Name`, `Status`, `Created at`. They can also be ordered by age (`Created At`) in both orders.

By clicking on the request id of one of the timeline records you will be taken to execution details page which shows the timeline of that particular request.

&nbsp;  
[![](/static/docs/executions/executions_list_click.png)](/static/docs/executions/executions_list_click.png)

When opened up you will see the details of an execution timeline. 


#### Execution Details

##### Timeline stages

Each timeline is composed of stages, a single stage will look like the following

[![](/static/docs/executions/executions_details_stage.png)](/static/docs/executions/executions_details_stage.png)

Each stage holds four pieces of information (reference image above):
- time it took to execute the stage. Visible on the left hand side
- the status code returned by the stage execution, the green `200` in the middle. Based on execution stage the status code may be different in each stage. If something is wrong the status code will be red
- the name of the stage. Visible on the right hand side (reference image above `custom-runner`). The name of stage will depend on the position in the execution chain
- below name of the stage, the timestamp of when the event was created

&nbsp;  
##### Synchronous Execution Timeline

Synchronous executions at most have a single timeline stage. Since they are synchronous they either succeed or fail, unlike asynchronous, there are no queues, dwell times and retries.
Since it only has a single stage, the stage will have the name of the function set as its name.

An example of a synchronous execution:

[![](/static/docs/executions/executions_details_sync.png)](/static/docs/executions/executions_details_sync.png)


&nbsp;  
##### Async Execution Timeline

Asynchronous execution has multiple stages throughout the entire timeline. It will at least consist of three stages:
1. Creation stage. The stage will be named after the function
2. Dwell Time. This stage indicates how long has the function execution request has been in the queue
3. Attempt #1. This show the information about the first execution attempt.

If the first attempt fails, the function will be executed again up to three times. Each failed attempt will add additional three minutes to the wait timer.

An example of a asynchronous execution:

[![](/static/docs/executions/executions_details.png)](/static/docs/executions/executions_details.png)

#### Execution Logs

Execution logs of the timeline can be accessed by clicking on the logs tab:

[![](/static/docs/executions/executions_details_access_logs.png)](/static/docs/executions/executions_details_access_logs.png)

&nbsp;  
The logs are the same as as globally available logs except they are `scoped to the request id`. Documentation on logs can be found [here](/docs/logs/overview)