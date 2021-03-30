---
title: Ruby Runtime
---

### Ruby Runtime Image

This page describes the layout and naming requirements of `Ruby` runtime which you need to follow in order to deploy your source code written in `ruby`.

#### File Requirements

When you submit your code for the build process, the builder will first try to validate the file structure inside provided zip file. In order for `Ruby` runtime to be bootstrapped, the builder **requires the following files** to appear inside the submitted zip file:
```
my-code.zip
│   handler.rb
```

> Note: You may also bundle together a `Gemfile` (at the same level as `handler.rb`) if your code requires additional dependencies

Your entrypoint point logic should go to `handler.rb` and any other dependencies required should be described in the `Gemfile` file.


#### Handler.rb Requirements

Since `handler.rb` is the main entry point to your custom source code, it has to meet certain criteria in order for the build process to succeed.

Here is the most basic template for `handler.rb`
```rb
class Handler
  def run(body, headers)
    status_code = 200 # Optional status code, defaults to 200
    response_headers = {"content-type" => "text/plain"}
    body = "Hello world from the Ruby template"

    return body, response_headers, status_code
  end
end
```


The **key strict requirements** are:
- must implement `Class Handler`
- the class `Handler` must implement a method `def run(body, header)`
- must write some response back (how else would you know if it succeeded)


#### Gemfile Requirements

Here is an example of a basic `Gemfile` file:

```Gemfile
source 'https://rubygems.org'

gem 'rspec'
```

The **key strict requirements** are:
- must conform to typical ruby Gemfile standards [documentation](https://bundler.io/gemfile.html)
