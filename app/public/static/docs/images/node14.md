## NodeJS 14 Runtime Image

This page describes the layout and naming requirements of `NodeJS` runtime which you need to follow in order to deploy your source code that supposed to run on `NodeJS 14` runtime.

#### File Requirements

When you submit your code for the build process, the builder will first try to validate the file structure inside provided zip file. In order for `NodeJS` runtime to be bootstrapped, the builder **requires the following files** to appear inside the submitted zip file:
```
my-code.zip
│   handler.js
|   package.json
```
Your entroypoint point logic should go to `handler.js` and any other dependencies required should be described in the `package.json` file.

#### Handler.js Requirements

Since `handler.js` is the main entry point to your custom source code, it has to meet certain criteria in order for the build process to succeed.

Here is the most basic template for `handler.js`
```js
'use strict'

module.exports = async (event, context) => {
  const result = {
    'body': JSON.stringify(event.body),
    'content-type': event.headers["content-type"]
  }

  return context
    .status(200)
    .succeed(result)
}
```


The **key strict requirements** are:
- must be a module export
- must implement `async (event, context)`
- must write some response back (how else would you know if it succeeded)


#### Package.json Requirements

Here is an example of a `package.json` file:

```json
{
  "name": "my-nodejs-function",
  "version": "0.1.0",
  "description": "Some JS function",
  "main": "handler.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 0"
  },
  "keywords": [],
  "author": "Me :)",
  "license": "MIT"
}
```

The **key strict requirements** are:
- must contain `"main": "handler.js"` entry
- must conform to npms [specifics](https://docs.npmjs.com/cli/v7/configuring-npm/package-json)