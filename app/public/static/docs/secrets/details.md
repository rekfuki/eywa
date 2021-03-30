---
title: Secret Details
---

### Secret Details
> This page focuses on getting secret details using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Secrets/getSecretsSecretid)

In order to get access to a secret details page you will first need to navigate to the secrets list page which can be done by clicking on the [Secrets](/app/secrets) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/secrets)

[![](/static/docs/secrets/secrets_navbar_location.png)](/static/docs/secrets/secrets_navbar_location.png)

&nbsp;  
When the page loads, select the secret from the list you would like to access details of:
[![](/static/docs/secrets/secrets_delete_select.png)](/static/docs/secrets/secrets_delete_select.png)

&nbsp;  
Once the page loads you will see the following page (subject to your secret configuration):
&nbsp;  
[![](/static/docs/secrets/secrets_details.png)](/static/docs/secrets/secrets_details.png)

&nbsp;  
On the page itself you will see three cards:

&nbsp;  
#### Secret Info

[![](/static/docs/secrets/secrets_details_info.png)](/static/docs/secrets/secrets_details_info.png)

The secret info card contains general information about the secret such as:
- `ID` - the id (UUID) of the secret
- `Name` - the name of the secret given during creation plus a random prefix attached by the system
- `Mount Path` - the mount path when the secret is mounted to the function. Whenever you attached a secret to a function, that secret is considered as `mounted` meaning its contents are available on that `mount path`. In the given example above, the mount path is `/var/faas/secrets/899c7b3a-demo-secret` and the secret has two fields: `key1` and `key2`. In order to access those fields from within your deployed function, you would read the the following files:
   - `/var/faas/secrets/899c7b3a-demo-secret/key1` which would contain `value1`
   - `/var/faas/secrets/899c7b3a-demo-secret/key2` which would contain `value2`
- `Total Mounts` - the number of functions that have this secret mounted
- `Updated At` - when the secret was last updated
- `Created At` - when the secret was created


&nbsp;  
##### Mounted Functions

[![](/static/docs/secrets/secrets_details_mounted.png)](/static/docs/secrets/secrets_details_mounted.png)

The mounted functions card contains information about which functions have that particular secret mounted (loaded):
- `ID` - the id (UUID) of the function. Clicking on it will redirect to the function management page
- `Name of the function` - the name of the function


&nbsp;  
##### Secret Fields

[![](/static/docs/secrets/secrets_details_fields.png)](/static/docs/secrets/secrets_details_fields.png)

The secret fields card contains the fields you've configured for your secret. The values of each field will always be hidden (they are actually not there).
In order to add additional secret field press the `ADD` button. If you wish to delete a field, press the `X` button on the right hand side of each entry.
The `UPDATE` button will be grayed out unless a change is detected.
