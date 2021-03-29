#### Delete Function
> This page focuses on deleting a function using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Functions/deleteFunctionsFunctionid)

In order to delete the function you will need to access the function management page. Documentation on how to access function management page can be found [here](/docs/functions/manage)

Once you are on the function management page you will see a `DELETE` button at the top right of the page:

[![](/static/docs/functions/function_delete_location.png)](/static/docs/functions/function_delete_location.png)

Once you press the button, the following modal will pop up:

&nbsp;  
[![](/static/docs/functions/function_delete_modal.png)](/static/docs/functions/function_delete_modal.png)

&nbsp;  
You must enter the exact function name in the text field in order to confirm (in this case `custom-runner`). If the value matches exactly, delete button will be enabled which you can then press in order to confirm deletion.

&nbsp;  
[![](/static/docs/functions/function_delete_modal_confirm.png)](/static/docs/functions/function_delete_modal_confirm.png)

&nbsp;  
Once you press the `delete` button you will be redirected to the function list view and the function you've just deleted will have status set to `TERMINATING`. Depending on the function, termination process can take a bit of time so keep refreshing and it eventually will disappear from the list.

&nbsp;  
[![](/static/docs/functions/function_delete_terminating.png)](/static/docs/functions/function_delete_terminating.png)