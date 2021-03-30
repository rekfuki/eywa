---
title: Database Access
---

### Access the Database

In order to connect to the database you will need to find the secret that contains your database credentials which have been generated during user registration process and the database address

#### Getting the database address

Chose a function you want to connect to the database and head over to its management page. Documentation on how to access function management page can be found [here](/docs/functions/manage)

Once on the function management page, on the `Environment Variables` card look for an environment variable called `mongodb_host`. That is the variable you will need to read from your function. Alternative you can just hardcode the value of `mongodb_host` variable inside your function code.

#### Getting the secret

The secret will have a name of `mongodb-credentials` prefixed with a random string `{prefix}-mongodb-credentials`. You can find the secret by heading to the secrets page which you can access either by clicking on the [Secrets](/app/secrets) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/secrets.

&nbsp;  
[![](/static/docs/secrets/secrets_navbar_location.png)](/static/docs/secrets/secrets_navbar_location.png)

Once the page is loaded look for a secret that has a name `{prefix}-mongodb-credentials`. If you have loads of secrets try searching for it using the search bar above the list.

&nbsp;  
[![](/static/docs/database/database_mongo_secret_list.png)](/static/docs/database/database_mongo_secret_list.png)

&nbsp;  
Navigate to the secrets details page by clicking on the `Secret ID` UUID. 

&nbsp;  
[![](/static/docs/database/database_mongo_secret_click_id.png)](/static/docs/database/database_mongo_secret_click_id.png)

&nbsp;  
Once the page loads you will see something similar (the full name of the secret is subject to change due to random prefix)

&nbsp;  
[![](/static/docs/database/database_mongo_secret_details.png)](/static/docs/database/database_mongo_secret_details.png)

&nbsp;  
You should see the same `Secret Fields` card:

&nbsp;  
[![](/static/docs/database/database_mongo_secret_fields.png)](/static/docs/database/database_mongo_secret_fields.png)

&nbsp;  
With the following fields:
- `username`
- `password`
- `database`


#### Connecting to the database

In order to connect to the database from your function you will need to have access to the `mongodb_host` which is passed as an environment variable.

You will also need to read the secret fields to get access to the `username`, `password` and `database`.

You can read more about reading the secrets from your function [here](/docs/secrets/access)

For the sake of this example, lets assume our mongo database secret is called `4f1b6641-mongodb-credentials` (just like in the screenshots). As already mentioned in the [secret access docs](/docs/secrets/access), you will need to read all three secret fields as files in order to get access to the values. The following file paths would be read:
- `/var/faas/secrets/4f1b6641-mongodb-credentials/username` (username field)
- `/var/faas/secrets/4f1b6641-mongodb-credentials/password` (password field)
- `/var/faas/secrets/4f1b6641-mongodb-credentials/database` (database field)

The actual establishment of the connection depends entirely on the language and the library of your choice. More details on forming a MongoDB URI can be found [here](https://docs.mongodb.com/manual/reference/connection-string/)