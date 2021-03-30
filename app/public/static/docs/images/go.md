---
title: Go Runtime
---

### Go Runtime Image

This page describes the layout and naming requirements of `Go` runtime which you need to follow in order to deploy your source code written in `Go`.

#### File Requirements

When you submit your code for the build process, the builder will first try to validate the file structure inside provided zip file. In order for `Go` runtime to be bootstrapped, the builder **requires the following files** to appear inside the submitted zip file:
```
my-code.zip
│   handler.go
```

That's correct, only one file is mandatory for the `Go` runtime and that is `handler.go` which has to placed at the root of the zip directory.
No `go.mod` or `go.sum` files are necessary (and are in fact ignored) during the build process. All the dependencies are automatically downloaded.


#### Handler.go Requirements

Since `handler.go` is the main entry point to your custom source code, it has to meet certain criteria in order for the build process to succeed.

Here is the most basic template for `handler.go`
``` go
package function

import (
    "net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Hello from Eywa!"`))
}
```


As you can see it's as minimal as it gets; however, it still has a few **strict requirements**:
- must be implemented for `package main`
- must implement `go func Handle(w http.ResponseWriter, r *http.Request)`
- must write some response back (how else would you know if it succeeded)


#### Examples


Example writing a JSON response
```go
package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		// read request payload
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		input = reqBody
		}
	}

	// log to stdout
	fmt.Printf("request body: %s", string(input))

	response := struct {
		Payload     string              `json:"payload"`
		Headers     map[string][]string `json:"headers"`
		Environment []string            `json:"environment"`
	}{
		Payload:     string(input),
		Headers:     r.Header,
		Environment: os.Environ(),
	}

	resBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // write result
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
```

Example persistent database connection pool between function calls:

```go
package function

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mc *mongo.Client
var db *mongo.Database

// Handle handles request
func Handle(w http.ResponseWriter, r *http.Request) {
	if mc == nil {
		username, err := getSecretValue("8a85b30f-mongodb-credentials", "username") // MongoDB credentials name will vary from user to user
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		password, err := getSecretValue("8a85b30f-mongodb-credentials", "password")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		mongoDBHost := os.Getenv("mongodb_host")
		opts := options.Client().
			SetHosts([]string{"mongodb.mongodb:27017"}).
			SetAuth(options.Credential{
				Username:      string(username),
				Password:      string(password),
				AuthSource:    mongoDBHost,
				AuthMechanism: "SCRAM-SHA-1",
			})
		opts = opts.SetConnectTimeout(time.Second * 10)
		mc, err = mongo.Connect(context.TODO(), opts)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		db = mc.Database(string(username))
	}

	_, err := db.Collection("foo").InsertOne(context.Background(), bson.M{"some-test-value": "yes"})
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func getSecretValue(secretName, key string) ([]byte, error) {
	path := filepath.Join("/var/faas/secrets", secretName, key)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
```