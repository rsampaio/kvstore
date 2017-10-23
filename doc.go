/*
Package kvstore is a simple key-value store that listens on TCP and TLS sockets and is intended to exercise software development and demonstrate the design and the use of documentation, tests and benchmarks in Go projects.

I started by defining the project structure which helped me separate each aspect of the project into well-defined packages. I decided to implement the project with the command SET, GET, DELETE first and unlimited capacity and later evolved adding the access and modify time tracking for the replacement policy and to implement the STREAM command that returns keys and values in last modified order.

The use of a test-driven approach made the transition of unlimited capacity with untracked actions to the final version much easier since the changes could be made with confidence and the success could be measured. The ability to generate detailed documentation with godoc was also very positive and can give the reader a great overview of how the project is implemented in greater detail than this README.

Project structure

The project is organized into independent packages (protocol, store and server) that are used kvserver which is the only main package under cmd.

Protocol

The protocol package defines a Parser interface and a Protocol struct that implements the parser, the Protocol parser will fill up its fields Command and Args when the input line is successfully parsed, the Protocol struct also has the boolean field ReceiveValue that indicates when a command should read the next line as an input value.

In this implementation the parser uses simple string split to indentify the keywords that compose the commands and its basic structure and even though it works for a small set of commands it is definitely not easily expansible and can be difficult to maintain, the use of a lexical parser would work better in future iterations of this package.

Store

The store package defines the Store interface with functions to store key value pairs, retrieve its capacity and a list of last modified pairs.

The store package also defines a MemoryStore struct with internal fields to track access and modify time of each key and their mutexes and a capacity counter also guarded by a mutex, this struct implements the Store interface the initilizer NewMemoryStore receives the desidered capacity for the MemoryStore.

The tracking of access time is used to implement an LRU (Least-Recently-Used) replacement policy and the modify time is used to implement the STREAM the results ordered by the last modified key.

The heavy use of mutexes in this package is limitation and lock contention causes additional complexy making all the operations at least O(N).

Server

The server package defines Handlers, helper functions for new Listeners that can be TCP or TLS and a Commander that is responsible for reading lines from the client, parse it using the protocol package and execute handlers appropriate for the command returned by the protocol parser.

The handler.go file also defines a variable DefaultHandler that is a map initialized with the default handlers for the commands GET, SET, DELETE and STREAM.

The implementation of this package was tricky and I ended up facing interesting issues with connection used in bufio Readers and re-used later for direct IO operations with different results due to buffered nature of the bufio. Once I realized that that I should perform Read operations on the buffer the implementation got simpler.

*/
package kvstore
