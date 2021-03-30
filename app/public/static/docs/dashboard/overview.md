---
title: Dashboard Overview
---

### Dashboard Overview

Dashboard page can be accessed by clicking on the [Dashboard](/app/dashboard) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/dashboard)

[![](/static/docs/dashboard/dashboard_navbar_location.png)](/static/docs/dashboard/dashboard_navbar_location.png)

&nbsp;  
Once the page loads, you should see the following (subject to your own data):

&nbsp;  
[![](/static/docs/dashboard/dashboard_details.png)](/static/docs/dashboard/dashboard_details.png)

&nbsp;  
Dashboard page shows what happened in the last 15 minutes. The page has four components to it: `Counters`, `Metrics`, `Logs`, `Executions`.


#### Counters

&nbsp;  
[![](/static/docs/dashboard/dashboard_counters.png)](/static/docs/dashboard/dashboard_counter.png)

&nbsp;  
The counters contain generalised information about the functions and the database:
- `Requests Made Since Epoch` - total requests made across all functions since account was created
- `Requests Made` - total request made in the last 15 minutes
- `Image Count` - total image count. Clicking on it will take you to the images page
- `Function Count` - total function count. Clicking on it will take you to the functions page
- `Database Collection Count` - total collection count in your database
- `Database size` - total database size
- `Database Average Object Size` - average object size in your database
- `Database Object Count` - total object count in your database

#### Metrics

&nbsp;  
[![](/static/docs/dashboard/dashboard_graphs.png)](/static/docs/dashboard/dashboard_graphs.png)

&nbsp;  
Dashboard contains a couple of metrics panels:
- `Request and Error Rates (per second)` - shows `2XX`, `4XX`, `5XX` request rates relative to total requests made across all of your functions. The equivalent `Prometheus` queries would look like:

```rust
// All requests
sum(rate(gateway_function_invocation_total{user_id="{{your_user_id}}"}[5000ms])) by(user_id)
// 2XX requests
sum(rate(gateway_function_invocation_total{code=~"2..",user_id="{{your_user_id}}"}[5000ms])) by(user_id)
// 4XX requests
sum(rate(gateway_function_invocation_total{code=~"4..",user_id="{{your_user_id}}"}[5000ms])) by(user_id)
// 5XX requests
sum(rate(gateway_function_invocation_total{code=~"5..",user_id="{{your_user_id}}"}[5000ms])) by(user_id)
```
- `Top 5 API calls (by path)` - shows the top five url paths that were used to call all your functions. The equivalent `Prometheus` query would look like:
```rust
topk(5, sum(rate(gateway_function_invocation_total{user_id="{{your_user_id}}"}[5000ms] by(path))
```

#### Logs

&nbsp;  
[![](/static/docs/dashboard/dashboard_logs.png)](/static/docs/dashboard/dashboard_logs.png)

&nbsp;  
Logs panel shows last five log entries that have been produced in the last 15 minutes. In order to see the rest of the logs you can head to the [Logs](/app/logs) page or you can click on the `SEE ALL` button on the bottom right of the `Logs` card

#### Executions

&nbsp;  
[![](/static/docs/dashboard/dashboard_executions.png)](/static/docs/dashboard/dashboard_executions.png)

&nbsp;  
Executions panel shows last five execution entries that have been produced in the last 15 minutes. In order to see the rest of the executions you can head to the [Executions](/app/executions) page or you can click on the `SEE ALL` button on the bottom right of the `Executions` card