---
title: Getting Started
---

### Getting Started

This page will take you from fresh start to having a basic function deployed and ready to receive requests.


### 1. Creating an image
> In depth documentation about image creation can be found [here](/docs/images/create)

In order to create an image you will first need a handler. In this example we will use `Go` to create our basic handler. If you wish to use some other runtime you can take a look at [create images documentation](/docs/images/create).

Make sure you have `Go` installed. You can find documentation on how to set it up [here](https://golang.org/).

Once you have setup your `Go` environment, create a new file called `handler.go` and add the following code:
```go
package function

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Body: %s", string(input))))
}
```

The code reads request body if present and echoes it back as a response. Not very exciting but it is enough as your first function.

Zip the `handler.go` into a zip file in a way that the handler file remains at the top level i.e.
```
my-handler.zip
│   handler.go
```

This is **very important**, if you fail to do so, during the build process you will get an error stating that the `top-level file handler.go` is missing.
If you wish to dive into more details about `Go` runtime, you can find the documentation [here](/docs/images/go)

Once you have your zip file ready, go ahead and navigate to the images pages. Image creation page can be accessed by clicking on [Images](/app/images) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/images/create)

&nbsp;  
[![](/static/docs/images/images_navbar_location.png)](/static/docs/images/images_navbar_location.png)

&nbsp;  
Once you are on the page, at the top right of the page you should see a `CREATE IMAGE` button:

&nbsp;  
[![](/static/docs/images/images_create_location.png)](/static/docs/images/images_create_location.png)

&nbsp;  
Click the button and a page with the following form will appear:

&nbsp;  
[![](/static/docs/images/images_create_form.png)](/static/docs/images/create_form.png)

##### Enter Name
Go ahead and enter the name of your image. It is a required field and must match the following criteria:
- contain at most 63 characters
- contain only lowercase alphanumeric characters or '-'
- start with an alphanumeric character
- end with an alphanumeric character

&nbsp;  
##### Enter Version
Choose a semantic version of your image, it's perfectly fine to leave as a default `0.1.0`. The version field must match the following criteria:
- a normal version number must take the form X.Y.Z where X, Y, and Z are non-negative integers- must not contain leading zeroes. 
- x is the major version, Y is the minor version, and Z is the patch version.
- each element must increase numerically. For instance: 1.9.0 -> 1.10.0 -> 1.11.0.

&nbsp;  
##### Choose a Runtime

Since in this example we are trying to create an image that uses `Go` runtime, you don't have to do anything, `Go` is selected by default. If you decided to use a different runtime, change to appropriate one instead.


&nbsp;  
##### Upload Zip File

In order to upload the the zipped `handler.go`, either drag the zip file into the file drop zone or click on it and select it from file menu.

Once you have selected the zip file, go ahead and press the `CREATE` button. You will be redirected to the image build page which shows the build progress.

&nbsp;  
##### Observing Build logs

Once creation is successful, you will be redirected to an image build page that outputs build logs. The contents of the page will look similar to this:

&nbsp;  
[![](/static/docs/overview/image_build_logs.png)](/static/docs/overview/image_build_logs.png)

At the beginning of the logs you will see `BUILD QUEUED` message that contains some general information about your image.

Once the build process starts, you will see `BUILD STARTED` message print out and the build logs will start appearing. The logs are just `buildkit` output. 

Assuming everything goes well, at the end of the build process you should see `BUILD FINISHED` message. This means that the build has successfully finished:

&nbsp;  
[![](/static/docs/overview/image_build_logs_success.png)](/static/docs/overview/image_build_logs_success.png)

&nbsp;  
Go back to images list page and you should see your image available. 

&nbsp;
### 2. Deploying a Function
> In depth documentation about function deployment can be found [here](/docs/functions/create)

Now that you have your first image ready. Go ahead and navigate to the functions creation page.

Function creation page can be accessed by navigating to [functions](/app/functions) and clicking `CREATE FUNCTION` button at the top of the page ([a direct link](/app/functions/create)).

&nbsp;
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

&nbsp;
##### Enter Name
Go ahead and enter the name of your function. It is a required field and must match the following criteria:
- match `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$` regular expression.
- be at least 5 characters long

&nbsp;
##### Select Image
Click on the image dropdown and you should see your built image in the list, select it.

&nbsp;
##### Further Configuration

Bellow image selection, you will see the following fields:

- `Minimum Replicas`
- `Maximum Replicas`
- `Scaling Factor`
- `Per Instance Concurrency`
- `Read Timeout`
- `Write Timeout`
- `Mounted Secrets` 

In this guide we will not touch these fields and will leave them set to their defaults. If you wish to play with the parameters, more details about them can be found [here](/docs/images/create)

&nbsp;
##### Environment Variables

Once you have filled out the necessary fields, click on the `NEXT` button which will bring you to the `Environment Variables` portion. 

&nbsp;
[![](/static/docs/functions/function_create_form_env.png)](/static/docs/functions/function_create_form_env.png)

Again, we will not be adding any extra environment variables in this guide, but if you wish to you can do so. Documentation on environment variables can be found [here](/docs/functions/create)

&nbsp;
##### Creating Function

Once you are ready, on the `Environment Variables` page press `CREATE` button at the bottom right of the form. This will redirect you to the function details page. 

&nbsp;
### 3. Accessing Invocation URLs
> In depth documentation about function management can be found [here](/docs/functions/manage)

After creating your function, you will be redirected to the page that holds all the information about your function. At the top of the page you will see the invocation URLs

&nbsp;
[![](/static/docs/functions/function_manage_view_urls.png)](/static/docs/functions/function_manage_view_urls.png)

`Sync URL` is used for synchronous invocations and `Async URL` is used for asynchronous invocations.

In this example we are going to use `Sync URL`, copy it. The URL should look something like this:
```
https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/
```
> The function UUID will be different

&nbsp;
### 4. Generating Access Token
> In depth documentation about access token creation can be found [here](/docs/access_tokens/create)

Now that you have the URL to your function. You will need an access token in order for you to call it.

Access token management page can be accessed by clicking on the [Tokens](/app/tokens) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/tokens)

[![](/static/docs/access_tokens/tokens_navbar_location.png)](/static/docs/access_tokens/tokens_navbar_location.png)

&nbsp;  
When the page loads, at the top right corner of the  page you will see a `CREATE TOKEN` button.

&nbsp;  
[![](/static/docs/access_tokens/tokens_create_location.png)](/static/docs/access_tokens/tokens_create_location.png)

&nbsp;  
Once you press the button the following modal will pop up:

&nbsp;  
[![](/static/docs/access_tokens/tokens_create_modal.png)](/static/docs/access_tokens/tokens_create_modal.png)

&nbsp;  
It will have a couple of fields:
- `Name of the Token` - this is simply a name of your access token. It must follow these requirements:
    - must be at least 5 characters long
    - must be at most 63 characters long
- `Token Expiry Date` - sets the date when the access token will expire. Leaving it empty.

&nbsp;  
After filling out the detail, go ahead and press the `CREATE` button, if everything is successful you will see the following form:

&nbsp;  
[![](/static/docs/access_tokens/tokens_post_create.png)](/static/docs/access_tokens/tokens_post_create.png)

As per instruction, do copy the token value since after closing the form the token value will no longer be accessible.

If you lose the token, delete the old and generate a new one. 

&nbsp;
### 5. Invoking Function
> In depth documentation about function invocations can be found [here](/docs/functions/invoke)

At this point you have everything you need in order to call your function. 

You can use whichever HTTP client you wish to. In this example we will be using `curl`.

In order to make the request you will need the URL and the access token acquired in previous steps. The access token must be sent via a header called `X-Eywa-Token`. 

For the sake of the example lets assume our URL is the following (**replace with your own**):
```
https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/
```

Since the function we have deployed echoes back the request body, lets make a `POST` request with a simple `JSON` body:

**Request**:
```bash
curl https://eywa.rekfuki.dev/eywa/api/functions/sync/a72edb6d-ee2b-5b91-9fa2-bd173c1eb269/ \
-H 'X-Eywa-Token: {your_token_here}' \
-d '{"foo": "bar"}' \
-X POST
```

You should receive the following response:

**Response**:
```bash
Body: {"foo": "bar"}
```

If you enable verbose mode you will also see that it returns a header called `X-Request-Id`. This is a unique UUID of a single request. You can use it to look at logs and executions. 

&nbsp;
### 6. Checking out Executions
> In depth documentation about executions can be found [here](/docs/executions/overview)

Executions are like timelines of each function invocation. Whenever you make a request to your function either synchronously or asynchronously, the entire journey of your request is tracked throughout the system until the result is returned. 

In order to access the execution timelines of the invocation, you can head to the the function management page of your function you just invoked and clicking on the `Timelines` tab. It is the same page as the one you went to grab invocation URLs.

&nbsp;
[![](/static/docs/functions/function_manage_view_timelines.png)](/static/docs/functions/function_manage_view_timelines.png)

&nbsp;
You should should see all the executions from your invocations. More details on how to parse timelines can be found [here](/docs/executions/overview).

&nbsp;
### 7. Checking out Logs
> In depth documentation about logs can be found [here](/docs/logs/overview)

Each invocation made generates logs. In order to access the logs of the invocation, you can head to the function management page of your function you just invoked and clicking on the `Logs` tab. It is the same page as the one you went to grab invocation URLs.

&nbsp;
[![](/static/docs/functions/function_manage_view_logs.png)](/static/docs/functions/function_manage_view_logs.png)

&nbsp;
You should should see all the logs from your invocations. More details on how to parse logs can be found [here](/docs/logs/overview).

&nbsp;
### 8. Checking out Metrics
> In depth documentation about metrics can be found [here](/docs/functions/metrics)

**Important**: Due to the way `Prometheus` works you may need to make several invocations before any metrics appear.

Each invocation made generates certain metrics. In order to access the metrics of your function, you can head to the function management page of your function you just invoked and clicking on the `Metrics` tab. It is the same page as the one you went to grab invocation URLs.

&nbsp;
[![](/static/docs/functions/function_manage_view_metrics.png)](/static/docs/functions/function_manage_view_metrics.png)


You can control the time ranges and refresh intervals using provided controls just above the individual metric panes. The controls going left to right are: 
- `Time Range` - a dropdown that lets you chose the time range scope of your metrics, i.e. last 5 minutes, last 15 minutes, last hour, etc...
- `Time`- an end time of your metrics data. In combination with `Time Range` the metrics returned would be in the range of `[time - time_range; time]`. By default this field is set to current time meaning you will get up to date metrics ranging back to whatever your `Time Range` is set to
- `Update Interval` - the interval of time to fetch your metrics at

Just below the time controls you can see 6 panels of different metrics. Going row by row the different metrics panels are:
- `Total Requests Made` - this counter shows total requests made to your function since it's deployment
- `Errors %` - is the rate of errors relative to total requests made. Any responses that are either `4XX` or `5XX` will be considered as errors
- `Request and Error Rates (per second)` - shows `2XX`, `4XX`, `5XX` request rates relative to total requests made
- `Request Rate by Method (per second)` - shows requests per second grouped by HTTP method
- `Request Duration (%)` - shows the percentage of requests duration grouped by `< 10ms`, `< 50ms`, `< 100ms`, `< 500ms` and `> 500ms`. For example if there are total of ten request made where two of them took `< 100ms` to respond and the rest took `< 10ms`, the pane would show two separate graphs: one at `20%` for the `< 100ms` and the other at `80%` for the `< 10ms`
- `Top 3 API calls (by path)` - shows the top url paths that were used to call your function
