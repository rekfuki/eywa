## Custom Runtime Image

This page describes the layout and naming requirements of a `Custom` runtime which you need to follow in order to deploy your source code written in language of your choice.

#### File Requirements

Unlike the other runtimes, `custom` does not adhere to any strict file requirements. Since `custom` runtime gives the most freedom to the user, we do not want to restrict you to specific layout or naming conventions and therefore allow you to provide the path to your executable (see more at [image create documentation](/docs/images/create)).


That being said, there are stil a couple of rules you **must follow**:
- the executable must be **statically compilled** (or alternatively entirely self-contained). This means that it should not have any dynamic dependencies loaded during execution. There are checks in place which will get tripped if the condition is not satisfied.
- the executable must implement a **web server** that listens either on `0.0.0.0:8082` or `127.0.0.2:8082`.


When your executable is deployed, any requests will be proxied by the system to the port `:8082`. The handling of the request is entirely up to you as long as it returns some sort of a response.


#### Examples

Here is an example of a custom runtime written in `Rust` (don't forget to **compile it statically**):


##### main.rs
```rust
extern crate actix;
extern crate actix_web;
extern crate env_logger;
#[macro_use]
extern crate tera;

use std::collections::HashMap;

use actix_web::{error, http, middleware, server, App, Error, HttpResponse, Query, State};

struct AppState {
    template: tera::Tera, // <- store tera template in application state
}

fn index(
    (state, _query): (State<AppState>, Query<HashMap<String, String>>),
) -> Result<HttpResponse, Error> {
    let s = state
        .template
        .render("index.html", &tera::Context::new())
        .map_err(|_| error::ErrorInternalServerError("Template error"))?;

    let paths = fs::read_dir("./").unwrap();

    Ok(HttpResponse::Ok().content_type("text/html").body(s))
}

fn main() {
    ::std::env::set_var("RUST_LOG", "eywa-rust=info");
    env_logger::init();
    let sys = actix::System::new("eywa-rust");

    server::new(|| {
        let tera = compile_templates!("./templates/**/*");

        App::with_state(AppState { template: tera })
            // enable logger
            .middleware(middleware::Logger::default())
            .resource("/", |r| r.method(http::Method::GET).with(index))
    })
    .bind("0.0.0.0:8082")
    .unwrap()
    .start();

    println!("Started http server: 0.0.0.0:8082");
    let _ = sys.run();
}
```

##### templates/index.html

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>eywa-rust</title>
</head>
<body>
	<h1>Welcome to eywa-rust!</h1>
</body>
</html>

```