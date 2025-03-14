# Internals of `ListenAndServe` in Go

In Go, **`ListenAndServe`** is a convenience function in the `net/http` package that starts up an HTTP server on a given address. Under the hood, it does three main things:

1. **Open a network listener** (typically TCP) on the specified address.
2. **Create a default `http.Server`** if one isn’t provided.
3. **Call the server’s `Serve` method**, which enters the main loop of accepting connections and handling requests.

Below is a high-level look at how `ListenAndServe` is implemented and what happens internally.

---

## 1. The `ListenAndServe` Function

If you look at the Go source code for `ListenAndServe` (simplified version shown), you’ll see something like this:

```go
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```

It just constructs a `Server` and then calls its `ListenAndServe` method:

```go
func (srv *Server) ListenAndServe() error {
    ln, err := net.Listen("tcp", srv.Addr)
    if err != nil {
        return err
    }
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

So the steps are:

1. **`net.Listen("tcp", srv.Addr)`** creates a TCP listener on the address (e.g., `":8080"`).
2. It **wraps the listener** in a `tcpKeepAliveListener`.
3. Finally, it calls **`srv.Serve(listener)`** to handle connections.

---

## 2. The `Serve` Method — The Core Loop

Inside `Serve`, you’ll find a loop that **accepts incoming connections** and then **spawns goroutines** to handle them. A simplified version of the code in `Serve` looks like this:

```go
func (srv *Server) Serve(l net.Listener) error {
    defer l.Close()
    var tempDelay time.Duration // for retrying on accept errors

    for {
        // 1. Accept a connection
        rw, err := l.Accept()
        if err != nil {
            // If it's a temporary error, back off and retry
            if ne, ok := err.(net.Error); ok && ne.Temporary() {
                if tempDelay == 0 {
                    tempDelay = 5 * time.Millisecond
                } else {
                    tempDelay *= 2
                }
                if max := 1 * time.Second; tempDelay > max {
                    tempDelay = max
                }
                time.Sleep(tempDelay)
                continue
            }
            return err
        }
        tempDelay = 0

        // 2. Handle the connection in a new goroutine
        go srv.handleConn(rw)
    }
}
```

### Key Observations

- **Accept Loop**: The server continuously calls `Accept()` on the listener to get new connections.
- **Temporary Errors**: If there’s a temporary network error, the code uses an **exponential backoff** (`tempDelay`) and retries.
- **Goroutine per Connection**: For each accepted connection `rw`, the server launches `go srv.handleConn(rw)` in its own goroutine.  
  This means each new client connection is handled **independently**, enabling concurrency.

---

## 3. Handling the Connection: `handleConn`

`handleConn` in turn creates a `conn` object (an internal struct in `net/http`) and calls `c.serve()`:

```go
func (srv *Server) handleConn(rw net.Conn) {
    c := srv.newConn(rw)
    c.setState(c.rwc, StateNew) // before Serve can return
    go c.serve()
}
```

The `conn` struct (found in `server.go`) holds information about:

- The raw network connection (`rwc`)
- The server reference
- Remote address, local address
- Connection state (e.g., active, closed)
- Buffers for reading/writing, etc.

---

## 4. Serving Requests: `c.serve()`

Inside `c.serve()`, the server **enters a loop** that reads requests from the connection:

```go
func (c *conn) serve() {
    // Defer a function to recover from panics, close connection, etc.
    defer func() {
        c.close()
    }()

    // 1. Read loop
    for {
        w, err := c.readRequest(ctx)
        if err != nil {
            // ... handle errors, break on EOF, etc.
            return
        }

        // 2. Handle the request
        serverHandler{c.server}.ServeHTTP(w, w.req)
        w.cancelCtx()

        // 3. If it's a keep-alive connection, keep reading more requests
        if someCondition {
            // read next request
            continue
        }
        
        // If not keep-alive or we need to close, break
        break
    }
}
```

### What Happens Here:

- **Parsing the Request**:  
  `c.readRequest(ctx)` reads from the TCP connection, parses the HTTP headers, method, URL, etc., and returns a `Request` plus a `ResponseWriter` (`w`).

- **Calling the Handler**:  
  The server wraps your handler (the one passed to `ListenAndServe` or the default `DefaultServeMux`) and calls `ServeHTTP(w, req)`. This is where your application logic runs.

- **Keep-Alive**:  
  If the client sets `Connection: keep-alive`, the server can handle multiple requests over the **same** TCP connection. The loop continues until the client or server closes the connection.

- **Cleanup**:  
  When the loop ends (client disconnected, error, or no keep-alive), it defers to `c.close()`, which closes resources.

---

## 5. Concurrency Model

> **`1 Connection → 1 Goroutine`**

Each accepted connection is processed in its **own** goroutine. That goroutine will sequentially read and serve multiple requests for **that** connection (if keep-alive is used).

- If a single client uses keep-alive and sends multiple requests on the same connection, they are handled **in series** in that single goroutine.
- Multiple clients (each with their own connections) will be handled **concurrently** because each connection gets its own goroutine.

---

## 6. Putting It All Together

1. **`ListenAndServe(addr, handler)`** → Creates a default `Server`.
2. **`net.Listen("tcp", addr)`** → Opens a TCP listener on the address.
3. **`Server.Serve(listener)`** → An infinite loop:
    1. `Accept()` a connection
    2. `go handleConn()` in a goroutine
4. **`handleConn(rw net.Conn)`** → Create a `conn` struct, then call `conn.serve()`.
5. **`conn.serve()`** → Read **multiple** requests (if keep-alive) in a loop:
    1. Parse HTTP request
    2. Call your handler (`ServeHTTP`)
    3. Write HTTP response
    4. Continue or close the connection

---

