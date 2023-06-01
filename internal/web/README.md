# Package web

Package `web` provides "plumbing" primitives for building web applications.
It is not a framework but rather a set of simple utilities that can be used independently of each other.


which provides higher-level abstraction on top of this package.

Nevertheless, nothing prevents you from grabbing the parts you need to create a fury application from this repository.
️The only caveat is that you should be aware for example how to chain this repository's middleware or even registering the ping handler.

## Features

- HTTP router built upon [chi](https://github.com/go-chi/chi) with API grouping capabilities.
- Data binding and validation for JSON.
- Handy functions to create custom errors.
- Centralized HTTP error handling.

## Quick Example

```go
package main

import (
    "log"
    "net"
    "net/http"

    "github.com/rromero96/roro-lib/cmd/web"
)

func main() {
    // Instantiate the router that will be used for registering HTTP handlers.
    w := web.New()

    // Register a simple handler that always returns 200 OK
    w.Get("/", func(w http.ResponseWriter, r *http.Request) error {
        w.WriteHeader(http.StatusOK)
        return nil
    })

    // Create the listener that will be pass in to the underlying http.Server for attending incoming requests
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }

    // Run blocks and listen for an interrupt or terminate signal from the OS.
    // After signalling, a graceful shutdown is triggered without affecting any live connections/clients connected to the server.
    // It will complete executing all the active/live requests before shutting down.
    if err := web.Run(ln, web.DefaultTimeouts, w); err != nil {
        log.Fatal(err)
    }
}
```

