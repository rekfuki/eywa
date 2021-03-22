## Creating Images
> This page focuses on creating an image using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=registry#/Images/postImages)

Eywa images are simple container images that are built from your uploaded source code which is combined with one of the available runtimes. Currently, the platfrom supports the following runtimes:

- `Go` 
- `NodeJS 14` 
- `Python 3`
- `Ruby`
- `C#` 
- `Custom`

In order for Eywa to successfully bootstrap a runtime to your provided function code, each of the runtimes have certain layout and naming requirements which can be found here:

[Go](/docs/images/go),
[NodeJS](/docs/images/nodejs),
[Python](/docs/images/python),
[Ruby](/docs/images/ruby),
[C#](/docs/images/csharp),
[Custom](/docs/images/custom)

**IMPORTANT**: Make sure to match the exact file layout and naming conventions of your chosen runtime

> NOTE: Eywa currently allows to upload source code that is no more than 50MB (Megabytes) in size.

### Using the website

Image creation page can be accessed by navigating to [images](/app/images) and clicking `CREATE IMAGE` button at the top of the page ([a direct link](/app/images/create)).

![](/static/docs/images/image_create_form.png "Create image form")

Once the page is loaded you should see a form with the following fields:

- `Name of the image` - is required and must match a [DNS Standard](https://tools.ietf.org/html/rfc1123). This means that the name must:
    - contain at most 63 characters
    - contain only lowercase alphanumeric characters or '-'
    - start with an alphanumeric character
    - end with an alphanumeric character

- `Version` - is required and must matched a normal [SemVer 2.0](https://semver.org/#spec-item-2) standard. This means that the version must:
    - a normal version number must take the form X.Y.Z where X, Y, and Z are non-negative integers- must not contain leading zeroes. 
    - x is the major version, Y is the minor version, and Z is the patch version.
    - each element must increase numerically. For instance: 1.9.0 -> 1.10.0 -> 1.11.0.

- `Runtime` - select one of runtimes based on your requirements.




### Using the API