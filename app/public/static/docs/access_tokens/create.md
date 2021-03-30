---
title: Create Token
---

### Create Access Tokens
> This page focuses on creating an access token using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=warden#/Access%20Tokens/postTokens); **however**, you will need to create your initial access token using the website which can then be used to access the REST API.

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
- `Token Expiry Date` - sets the date when the access token will expire. Leaving this out will generate token that will never expire.

&nbsp;  
After filling out the detail, go ahead and press the `CREATE` button, if everything is successful you will see the following form:

&nbsp;  
[![](/static/docs/access_tokens/tokens_post_create.png)](/static/docs/access_tokens/tokens_post_create.png)

As per instruction, do copy the token value since after closing the form the token value will no longer be accessible.

If you lose the token, delete the old and generate a new one. 