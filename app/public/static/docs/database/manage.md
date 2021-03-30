---
title: Database Management
---

### Manage database
> This page focuses on managing the database using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=tugrik#/)

You can access the database management page  by clicking on the [Database](/app/database) directive on the navigation bar. Alternatively you can access it via a direct [link](/app/database)


[![](/static/docs/database/database_navbar_location.png)](/static/docs/database/database_navbar_location.png)

&nbsp;  
Once the page loads you should see the following (number of cards will depend on how many collection your database has):

&nbsp;  
[![](/static/docs/database/database_details.png)](/static/docs/database/database_details.png)

&nbsp;  
On the details page you will see at least one card (in this example two):

#### Database Overview

[![](/static/docs/database/database_details_overview.png)](/static/docs/database/database_details_overview.png)

This card contains basic overview of the database. The fields are as such:
- `Total Size` - the total size of your database. Combined `Storage Size` and `Index Size`
- `Storage Size` - the metric is equal to the size (in bytes) of all the data extents in the database. This number is larger than `Data Size` because it includes yet-unused space (in data extents) and space vacated by deleted or moved documents within extents. The `Storage Size` does not decrease as you remove or shrink documents
- `Index Size` - the size of your indices including `_id`
- `Data Size` - the metric is the sum of the sizes (in bytes) of all the documents and padding stored in the database. While `Data Size` does decrease when you delete documents, `Data Size` does not decrease when documents shrink because the space used by the original document has already been allocated (to that particular document) and cannot be used by other documents. Alternatively, if a user updates a document with more data, `Data Size` will remain the same as long as the new document fits within its originally padded pre-allocated space
- `Average Object Size` - the average size of objects in your database
- `Collection Count` - the number of collections in your database
- `Index Count` - the number of indices in your database


#### Collection Information

> Will only be visible if you actually have a collection in your database

[![](/static/docs/database/database_details_collection.png)](/static/docs/database/database_details_collection.png)

Per each collection in your database you will see a card displayed information about that collection. The fields are as such:
- `Total Size` - the total size of your database. Combined `Storage Size` and `Index Size`
- `Storage Size` - the metric is equal to the size (in bytes) of all the data extents in the database. This number is larger than `Data Size` because it includes yet-unused space (in data extents) and space vacated by deleted or moved documents within extents. The `Storage Size` does not decrease as you remove or shrink documents
- `Total Index Size` - the size of your indices including `_id`
- `Average Object Size` - the average size of objects in your database
- `Index Count` - the number of indices in your database

Below that you will see a list of all indices in that collection as `Index Name` and `Size` pairs


#### Delete Collection

In order to delete the collection you should press on the delete icon in the top right corner of your collection card you wish to delete:

&nbsp;  
[![](/static/docs/database/database_delete_collection_location.png)](/static/docs/database/database_delete_location_details.png)

&nbsp;  
Once you press the button the following modal will pop up:

&nbsp;  
[![](/static/docs/database/database_delete_modal.png)](/static/docs/database/database_delete_modal.png)

&nbsp;  
You must enter the exact collection name in the text field in order to confirm (in this case `my_collection`). If the value matches exactly, delete button will be enabled which you can then press in order to confirm deletion.

After a successful deletion the collection card will be removed from the database information page.