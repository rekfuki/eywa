### Function Metrics

Function metrics can be accessed by navigating to your specific function page and clicking on the `Metrics` tab. More details on how to access your specific function page can be found [here](/docs/functions/manage)

[![](/static/docs/functions/function_manage_view_metrics.png)](/static/docs/functions/function_manage_view_metrics.png)

You can control the time ranges and refresh intervals using provided controls just above the individual metric panes. The controls going left to right are: 
- `Time Range` - a dropdown that lets you chose the time range scope of your metrics, i.e. last 5 minutes, last 15 minutes, last hour, etc...
- `Time`- an end time of your metrics data. In combination with `Time Range` the metrics returned would be in the range of `[time - time_range; time]`. By default this field is set to current time meaning you will get up to date metrics ranging back to whatever your `Time Range` is set to
- `Update Interval` - the interval of time to fetch your metrics at


&nbsp;  
Just below the time controls you can see 6 panels of different metrics. Going row by row the different metrics panels are:
- `Total Requests Made` - this counter shows total requests made to your function since it's deployment. The equivalent `Prometheus` queries would look like:
```rust
sum(gateway_function_invocation_total{function_id="{{your_function_id}}",user_id="{{your_user_id}}"}) by(function_id)
```
- `Errors %` - is the rate of errors relative to total requests made. Any responses that are either `4XX` or `5XX` will be considered as errors. The equivalent `Prometheus` queries would look like:
```rust
((sum by(function_id) (rate(gateway_function_invocation_total{function_id="{{your_function_id}}",code=~"4..|5..",user_id="{{your_user_id}}"}[5000ms]))) / (sum by(function_id) (rate(gateway_function_invocation_total[5000ms])))) * 100
```
- `Request and Error Rates (per second)` - shows `2XX`, `4XX`, `5XX` request rates relative to total requests made. The equivalent `Prometheus` queries would look like:

```rust
// All requests
sum(rate(gateway_function_invocation_total{function_id="{{your_function_id}}",user_id="{{your_user_id}}"}[5000ms])) by(function_id)
// 2XX requests
sum(rate(gateway_function_invocation_total{function_id="{{your_function_id}}",code=~"2..",user_id="{{your_user_id}}"}[5000ms])) by(function_id)
// 4XX requests
sum(rate(gateway_function_invocation_total{function_id="{{your_function_id}}",code=~"4..",user_id="{{your_user_id}}"}[5000ms])) by(function_id)
// 5XX requests
sum(rate(gateway_function_invocation_total{function_id="{{your_function_id}}",code=~"5..",user_id="{{your_user_id}}"}[5000ms])) by(function_id)
```
- `Request Rate by Method (per second)` - shows requests per second grouped by HTTP method. The equivalent `Prometheus` queries would look like:
```rust
sum(rate(gateway_function_invocation_started{function_id="{{your_function_id}}",user_id="{{your_user_id}}"}[5000ms])) by(method)
```
- `Request Duration (%)` - shows the percentage of requests duration grouped by `< 10ms`, `< 50ms`, `< 100ms`, `< 500ms` and `> 500ms`. For example if there are total of ten request made where two of them took `< 100ms` to respond and the rest took `< 10ms`, the pane would show two separate graphs: one at `20%` for the `< 100ms` and the other at `80%` for the `< 10ms`. The equivalent `Prometheus` queries would look like:
```rust
// Sub 10ms
(sum(rate(gateway_function_duration_milliseconds_bucket{function_id="{{your_function_id}}",le="10",user_id="{{your_user_id}}"}[5000ms)) / sum(rate(gateway_function_duration_milliseconds_count{function_id="{{your_function_id}}"}[5000ms]))) * 100

// Sub 50ms
(sum(rate(gateway_function_duration_milliseconds_bucket{function_id="{{your_function_id}}",le="50",user_id="{{your_user_id}}"}[5000ms)) / sum(rate(gateway_function_duration_milliseconds_count{function_id="{{your_function_id}}"}[5000ms]))) * 100

// Sub 100ms
(sum(rate(gateway_function_duration_milliseconds_bucket{function_id="{{your_function_id}}",le="100",user_id="{{your_user_id}}"}[5000ms)) / sum(rate(gateway_function_duration_milliseconds_count{function_id="{{your_function_id}}"}[5000ms]))) * 100

// Sub 500ms
(sum(rate(gateway_function_duration_milliseconds_bucket{function_id="{{your_function_id}}",le="500",user_id="{{your_user_id}}"}[5000ms)) / sum(rate(gateway_function_duration_milliseconds_count{function_id="{{your_function_id}}"}[5000ms]))) * 100

// Above 500ms
((gateway_function_duration_milliseconds_bucket{function_id="{{your_function_id}}",le="500",user_id="{{your_user_id}}"}) / ignoring (le) gateway_function_duration_milliseconds_count{function_id="{{your_function_id}}"}) * 100
```
- `Top 3 API calls (by path)` - shows the top three paths url paths that were used to call your function. The equivalent `Prometheus` queries would look like:
```rust
topk(3, sum(rate(gateway_function_invocation_total{function_id="{{your_function_id}}",user_id="{{your_user_id}}"}[5000ms]by(path))
```