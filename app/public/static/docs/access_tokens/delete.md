### Delete Access Token

> This page focuses on delete an access token using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=warden#/Access%20Tokens/deleteTokensTokenid)

Access token management page can be accessed by navigating clicking on the [Tokens](/app/tokens) directive in the navigation bar. Alternatively you can access it via a direct [link](/app/tokens)

[![](/static/docs/access_tokens/tokens_navbar_location.png)](/static/docs/access_tokens/tokens_navbar_location.png)

&nbsp;  
When the page loads, you will see a list of tokens. Each of the rows have a red garbage can icon that indicates deletion.

&nbsp;  
[![](/static/docs/access_tokens/tokens_delete_location.png)](/static/docs/access_tokens/tokens_delete_location.png)

&nbsp;  
Once you press the garbage can button the following modal will pop up:

&nbsp;  
[![](/static/docs/access_tokens/tokens_delete_modal.png)](/static/docs/access_tokens/tokens_delete_modal.png)

&nbsp;  
You must enter the exact token name in the text field in order to confirm (in this case `87840366-demo-token`). If the value matches exactly, delete button will be enabled which you can then press in order to confirm deletion.

After a deletion is successful the access token will be removed rom the list

&nbsp;  
[![](/static/docs/access_tokens/tokens_delete_confirm.png)](/static/docs/access_tokens/tokens_delete_confirm.png)
