---
title: Access Secret
---

### Access your secret

In order to access your secret from the function you will have to read it as a file. But before you do that you need to go to your secrets details page which you can find the details on how to do that [here](/docs/secrets/details)

Once you are on the function details page, look for `Mount Path` on the `Function Info` card. It will look something like this `/var/faas/secrets/{name_of_your_secret}`. 

All secrets follow the same mount path structure: `/var/faas/secrets/{the_name_of_your_function}/{your_field_you_want_to_access}`. For example let's say you have a secret that is named `899c7b3a-demo-secret` and has the following fields:
- `KEY1` : `VALUE1`
- `KEY2` : `VALUE2`

In order to access each of those fields values you will have to read the following files:
- `/var/faas/secrets/899c7b3a-demo-secret/KEY1`
- `/var/faas/secrets/899c7b3a-demo-secret/KEY2`