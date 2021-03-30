---
title: Create Secret
---

### Create Secrets
> This page focuses on creating a secret using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=gateway-api#/Secrets/postSecrets)

Secret details page can be accessed by clicking on the [Secrets](/app/secrets) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/secrets)

[![](/static/docs/secrets/secrets_navbar_location.png)](/static/docs/secrets/secrets_navbar_location.png)

&nbsp;  
When the page loads, at the top right corner of the  page you will see a `CREATE SECRET` button.

&nbsp;  
[![](/static/docs/secrets/secrets_create_location.png)](/static/docs/secrets/secrets_create_location.png)

&nbsp;  
Once you press the button the following form will appear:

&nbsp;  
[![](/static/docs/secrets/secrets_create_form.png)](/static/docs/secrets/secrets_create_form.png)

&nbsp;  
The form will have the following fields:
- `Secret Name` - This is the name of your secret and it must match the following requirements:
    - must be at least 5 characters long
    - must match `^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$` regex
- `Secret fields` - these fields are `KEY:VALUE` pairs and you can have as many as you want. However, the total secret size **CANNOT EXCEED 1MB IN TOTAL SIZE**. In order to add extra fields you can click on the `ADD` button. In order to remove the field you can click on the `X` button on the right hand side of the field.
    - `KEY` must match the following rules:
        - must be at least 1 characters long
        - must bet at most 255 characters long
    - `VALUE` must match the following rules:
        - must be at least 1 characters long
        - must bet at most 2000 characters long

Once you are done filling out the fields, press `CREATE` button in order to create the secret. If creation process is successful you will be redirected to secrets details page, where you should see your newly created secret. Documentation on secret information page can be found [here](/docs/secrets/overview)


