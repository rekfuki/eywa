---
title: Delete Secret
---

### Delete Secrets
> This page focuses on deleting a secret using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Secrets/deleteSecretsSecretid)

In order to delete a secret you will need to access your secret details page. In order to access the details page you will need to first go to the secrets list page which can be done by clicking on the [Secrets](/app/secrets) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/secrets)

[![](/static/docs/secrets/secrets_navbar_location.png)](/static/docs/secrets/secrets_navbar_location.png)

&nbsp;  
When the page loads, select the secret from the list you would like to delete.

&nbsp;  
[![](/static/docs/secrets/secrets_delete_select.png)](/static/docs/secrets/secrets_delete_select.png)

&nbsp;  
Once you are on the secret details page (about which you can read [here](/docs/secret/overview)), at the top right corner you will see a `DELETE` button: 

&nbsp;  
[![](/static/docs/secrets/secrets_delete_location.png)](/static/docs/secrets/secrets_delete_location.png)

Click on the `DELETE` button and the following modal will pop up:

&nbsp;  
[![](/static/docs/secrets/secrets_delete_modal.png)](/static/docs/secrets/secrets_delete_modal.png)

&nbsp;  
You must enter the exact token name in the text field in order to confirm (in this case `899c7b3a-demo-secret`). If the value matches exactly, delete button will be enabled which you can then press in order to confirm deletion.

&nbsp;  
[![](/static/docs/secrets/secrets_delete_confirm.png)](/static/docs/secrets/secrets_delete_confirm.png)

After a deletion is successful you will be redirected to the secrets list page and the deleted secret will be gone.

