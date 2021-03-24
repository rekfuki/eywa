## Creating Images
> This page focuses on creating an image using the website. Documentation for the REST API can be found [here](/api-docs/?urls.primaryName=registry#/Images/postImages)

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

- `Runtime` - select one of runtimes based on your requirements. In order for Eywa to successfully bootstrap a runtime to your provided function code, each of the runtimes have certain layout and naming requirements.

    **IMPORTANT**: Make sure to match the exact file layout and naming conventions of your chosen runtime.

    Currently supported layouts are the following:
    - `Go` - [documentation](/docs/images/go),
    - `NodeJS 14` - [documentation](/docs/images/nodejs14),
    - `Python 3` - [documentation](/docs/images/python3),
    - `Ruby` - [documentation](/docs/images/ruby),
    - `Custom` - [documentation](/docs/images/custom),


- `Executable path` (only applicable to `custom` runtimes) - path to the executable (relative to your zip file).    

  Imagine the following zip structure:
   ```
      my-code.zip
      │   my-executable
      │
      └───templates
      │   │   index.html
    ```
  The `Executable path` would be `my-executable`
  
  Or your executable could be somewhere in another folder instead of being at the root of the zip file:
   ```
      my-code.zip
      │
      └───templates
      │   │   index.html
      └───runner
      │   │   my-executable
    ```
  In this case the `Executable path` would be `runner/my-executable`
- `Zip File` - either click or drag to upload the zip file that contains your code (runtime in case of `custom`). Only one zip file is accepted which means the last uploaded one will be uploaded for image creation. 

  Currently, the source code **cannot exeed 50MB (Megabytes)** in size


**IMPORTANT**: Combination of `name of the image`, `version` and `runtime` must be unique. For example if you already have an image named `my-awesome-image` version `0.1.0` runtime `go` and you tried to create another exact one, you would get an error.


Upon a successful creation you will be redirected to an image build progress page which displays the state of of the image and any logs produced during the building process.