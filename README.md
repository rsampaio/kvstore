# KVStore

KVStore is a simple key value store that listens on TCP and TLS sockets and is intended to excercise software development and demonstrate the design and the use of documentation, tests and benchmarks in Go projects.

I started by defining the project structure which helped me separate each aspect of the project into well defined packages. I decided to implement the project with the command SET, GET, DELETE first and unlimited capacity and later evolved adding the access and modify time tracking for the replacment policy and to implement the STREAM command that returns keys and values in last modified order.

The use of a test driven approach made the transition of unlimited capacity with untracked actions to the final version much easier since the changes could be made with confidence and the success could be measured. The ability to generate detailed documentation with `godoc` was also very positive and can give the reader a great overview of how the project is implemented in greater detail than this README.

## Project structure

The project is organized into independent packages (`protocol`, `store` and `server`) that are used kvserver which is the only main package under `cmd`.

### protocol

The `protocol` package defines a `Parser` interface and a `Protocol` struct that implements the parser, the `Protocol` parser will fill up its fields `Command` and `Args` when the input line is succesfully parsed, the `Protocol` struct also has the boolean field `ReceiveValue` that indicates when a command should read the next line as an input value.

In this implementation the parser uses simple string split to indentify the keywords that compose the commands and its basic structure and even though it works for a small set of commands it is definitely not easily expansible and can be difficult to maintain, the use of a lexical parser would work better in future iterations of this package.

### store

The `store` package defines the `Store` interface with functions to store key value pairs, retrieve its capacity and a list of last modified pairs.

The `store` package also defines a `MemoryStore` struct with internal fields to track access and modify time of each key and their mutexes and a capacity counter also guarded by a mutex, this struct implements the `Store` interface the initilizer `NewMemoryStore` receives the desidered capacity for the `MemoryStore`.

The tracking of access time is used to implement an LRU (Least-Recently-Used) replacement policy and the modify time is used to implement the STREAM the results ordered by the last modified key.

The heavy use of mutexes in this package is limitation and lock contention causes additional complexy making all the operations at least O(N).

### server

The `server` package defines `Handlers`, helper functions for new `Listeners` that can be TCP or TLS and a `Commander` that is responsible for reading lines from the client, parse it using the `protocol` package and execute handlers appropriate for the command returned by the protocol parser.

The `handler.go` file also defines a variable `DefaultHandler` that is a map initialized with the default handlers for the commands GET, SET, DELETE and STREAM.

The implementation of this package was tricky and I ended up facing interesting issues with connection used in `bufio` Readers and re-used later for direct IO operations with different results due to buffered nature of the bufio. Once I realized that that I should peform Read operations on the buffer the implementation go easier.

## Build, Test and Execution

To build this project you can use:

```
go get -u github.com/rsampaio/kvstore/cmd/kvserver
```

This will clone the repository into `$GOPATH/src/github.com/rsampaio/kvstore` directory, compile the `kvserver` binary and install it in your `$GOBIN` directory (usually `$GOPATH/bin`).

To run tests after the project is installed change to this directory and run:

```
make test
```

Benchmarks for GET and SET operations are also available and can be run with:

```
make benchmarks
```

To run `kvserver` with TLS enabled run:

```
make run-tls
```

This will compile and install the binary in `$GOBIN`, create `cert.pem` and `cert.key` with the openssl command and start `kvserver` with TLS as well as TCP enabled.

This `openssl` command can be used to connect to the TLS port: `openssl s_client -host localhost -port 2021`.

A detailed documentation of the code can be viewed with `godoc` and for easy access can be started with `make godoc` inside the cloned repository. The command will show the URL assuming the repository was cloned inside your `$GOPATH`.

The parameters for `kvserver` binary are:

```
λ kvserver -h
Usage of kvserver:
  -capacity-bytes int
        Max capacity in bytes (default 1000)
  -enable-tls
        Enables TLS server (requires --tls-cert and --tls-key)
  -tcp-listen string
        TCP server listen address (default ":2020")
  -tls-cert string
        PEM certificate file
  -tls-key string
        Cerficate key file
  -tls-listen string
        TLS server listen address (default ":2021")
```
