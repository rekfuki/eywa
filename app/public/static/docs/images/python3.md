## Python3 Runtime Image

This page describes the layout and naming requirements of `Python3` runtime which you need to follow in order to deploy your source code written in `python`.

#### File Requirements

When you submit your code for the build process, the builder will first try to validate the file structure inside provided zip file. In order for `Python3` runtime to be bootstrapped, the builder **requires the following files** to appear inside the submitted zip file:
```
my-code.zip
│   handler.py
```
> Note: You may also bundle together a `requirements.txt` file (at the same level as `handler.py`) if your code requires additional dependencies

Your entroypoint point logic should go to `handler.py` and any other dependencies required should be described in the `requirements.txt` file.


#### Handler.py Requirements

Since `handler.py` is the main entry point to your custom source code, it has to meet certain criteria in order for the build process to succeed.

Here is the most basic template for `handler.py`
``` py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": "Hello from Eywa!"
    }
```


The **key strict requirements** are:
- must implement `def handle(event, context)`
- must write some response back (how else would you know if it succeeded)


#### Requirements.txt Requirements

Here is an example of a `requirements.txt` file:

```txt
pyOpenSSL==0.13.1
pyparsing==2.0.1
python-dateutil==1.5
pytz==2013.7
scipy==0.13.0b1
six==1.4.1
virtualenv==16.3.0
```

The **key strict requirements** are:
- must conform to typical python3 requirements.txt syntax

#### Examples


Successful response status code and JSON response body

```py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": {
            "key": "value"
        }
    }
```

Successful response status code and string response body

```py
def handle(event, context):
    return {
        "statusCode": 201,
        "body": "Object successfully created"
    }
```

Failure response status code and JSON error message

```py
def handle(event, context):
    return {
        "statusCode": 400,
        "body": {
            "error": "Bad request"
        }
    }
```

Setting custom response headers

```py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": {
            "key": "value"
        },
        "headers": {
            "Location": "https://www.example.com/"
        }   
    }
```

Accessing request body

```py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": "You said: " + str(event.body)
    }
```

Accessing request method

```py
def handle(event, context):
    if event.method == 'GET':
        return {
            "statusCode": 200,
            "body": "GET request"
        }
    else:
        return {
            "statusCode": 405,
            "body": "Method not allowed"
        }
```

Accessing request query string arguments

``` py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": {
            "name": event.query['name']
        }
    }
```

Accessing request headers

```py
def handle(event, context):
    return {
        "statusCode": 200,
        "body": {
            "content-type-received": event.headers.get('Content-Type')
        }
    }
```