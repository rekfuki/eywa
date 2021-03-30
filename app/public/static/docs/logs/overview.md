---
title: Logs Overview
---

### Logs

This page describes how logs are generated and how they can be accessed.

Whenever you make a call to your function, each stage of the execution generates appropriate logs.

In order to access the logs page you can either click on the `Logs` in the side navigation bar or by navigating there via a direct [link](/app/logs).

[![](/static/docs/logs/logs_navbar_location.png)](/static/docs/executions/logs_navbar_location.png)

&nbsp;  
Once the page is loaded you should see a page with the following content

&nbsp;  
[![](/static/docs/logs/logs_list_view.png)](/static/docs/logs/logs_list_view.png)

&nbsp;  
Just above the list of logs you can see the basic filtering controls:
- `Search box` - allows to filter logs by `Status`, `Request ID`, `Function ID`, `Function Name`, `Message` and `Timestamp`
- `Start Time` - The lower boundary of time at which the logs have been generated at
- `End Time` - The upper boundary of time at which the logs have been generated at
- `Show only errors` - when ticked, it will show only logs that resulted in errors

Each log entry row has the following columns: 
- `Status` - the status of the stage
- `Function ID` - the id (UUID) of the function that was invoked. Clicking on it will redirect you to the function management page (information about management can be found [here](/docs/functions/manage))
- `Function Name` - the name of the function
- `Request ID` - the id (UUID) of the request made. Clicking on it will redirect you to the that particular requests execution timeline page (information about execution timelines can be found [here](/docs/executions/overview))
- `Message` - summary of the log message. `Synchronous` execution messages will start with `SYNC EXECUTION...`, whereas `Asynchronous` execution messages will start with `ASYNC EXECUTION...`
- `Timestamp` - the timestamp of the log

By clicking on the chevron on the very left of the row you can expand the log message to see more details about the execution stage.

Since there are two types of invocations: `Synchronous` and `Asynchronous` and since they both differ in the amount of stages generated (read more about it [here](/docs/executions/overview)) the amount of logs will differ as well.

> **NOTE**: logs will appear in chronological order with latest log entry at the top


#### Synchronous execution logs

`Synchronous` invocation generates two logs:
1. The start of the execution. This log entry holds the request information such as `body`, `query`, `path` and `headers`
2. The end of the execution. This log entry holds the response information such as `body`, `headers` and `status`. If the `DEBUG` mode is enabled and your function logs any `stderr` or `stdout` messages, those messages will appear in the logs details under `stderr` and `stdout` respectively. Debug mode can be enabled by updating the function and turning on debug (read more about it [here](/docs/functions/update))

An example of a `Synchronous` execution log entries:

&nbsp;  
[![](/static/docs/logs/logs_list_sync.png)](/static/docs/logs/logs_list_sync.png)

As you can see there are only two entries as mention before. Both of the entry messages start with `SYNC EXECUTION ...`

Here is how the log entry for `STARTED` looks like once expanded:

&nbsp;  
[![](/static/docs/logs/logs_sync_started_expanded.png)](/static/docs/logs/logs_sync_started_expanded.png)

And here is how the log entry for `ENDED` looks like once expanded:

&nbsp;  
[![](/static/docs/logs/logs_sync_ended_expanded.png)](/static/docs/logs/logs_sync_ended_expanded.png)

#### Asynchronous execution logs

`Asynchronous` invocation generates at least three logs:
1. Execution queued up. This log entry holds basic information informing that the async execution request has been accepted and has been successfully queued up
2. Execution Attempt #1. This log entry holds the execution attempt information such as `body`, `query`, `path` and `headers`
3. Execution Attempt #1 Result. This log entry holds the execution attempt #1 result information such as `body`, `query`, `path` and `headers`. If the `DEBUG` mode is enabled and your function logs any `stderr` or `stdout` messages, those messages will appear in the logs details under `stderr` and `stdout` respectively. Debug mode can be enabled by updating the function and turning on debug (read more about it [here](/docs/functions/update))
4. Execution Attempt #2 (**only if Attempt #1 failed**)
5. Execution Attempt #2 Result (**only if Attempt #1 failed**)
6. Execution Attempt #3 (**only if Attempt #2 failed**)
7. Execution Attempt #3 Result (**only if Attempt #2 failed**)


&nbsp;  
An example of a `Asynchronous` execution log entries:

&nbsp;  
[![](/static/docs/logs/logs_list_async.png)](/static/docs/logs/logs_list_async.png)

As you can see there are three entries as mention before. The first Entry starts with `QUEUED...`, the second one with `ASYNC ATTEMPT #1 STARTED` and the third one with `ASYNC ATTEMPT #1 FINISHED`. These are the three minimum entries you will always see with asynchronous executions. If the first attempt is unsuccessful you may see additional entries.

Here is how the log entry for `QUEUED` looks like once expanded:

&nbsp;  
[![](/static/docs/logs/logs_async_queued_expanded.png)](/static/docs/logs/logs_async_queued_expanded.png)

Here is how the log entry for `ASYNC ATTEMPT #1 STARTED` looks like once expanded:

&nbsp;  
[![](/static/docs/logs/logs_async_started_expanded.png)](/static/docs/logs/logs_async_started_expanded.png)

And here is how the log entry for `ASYNC ATTEMPT #1 FINISHED` looks like once expanded:

&nbsp;  
[![](/static/docs/logs/logs_async_finished_expanded.png)](/static/docs/logs/logs_async_finished_expanded.png)
