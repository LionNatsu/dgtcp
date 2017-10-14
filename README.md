# DGTCP

## Datagram Transmission Control Protocol

DGTCP is an TCP-based datagram-style protocol. It provide
a datagram model over TCP (which is stream model).

## Protocol Details
Here is an example of 2 packets sent in DGTCP.

	[ packet #1: 00 01 02 ...(0x1FF bytes) ]
	[ packet #2: 23 12 AB ]
	...

The following chart shows how these two packets was sent
in the underlying TCP stream.

	<len> <          data         >   <len> < data >   ...
	01 FF 00 01 02 ...(0x1FF bytes)   00 03 23 12 AB   ...

## Golang Package

Package dgtcp overrides Read() and Write() methods of `net.TCPConn`.
### Usage & Example

Client:

```go
import (
	"net"

	"dgtcp"
)

func main() {
	connRaw, _ := net.Dial("tcp", "[::1]:50000")
	conn := dgtcp.New(connRaw) // <=
	buf := make([]byte, 20000)
	conn.Write(buf)
	conn.Close()
}
```

Server:

```go
import (
	"fmt"
	"net"

	"dgtcp"
)

func main() {
	listener, _ := net.Listen("tcp", "[::1]:50000")
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go process(conn)
	}
}

func process(connRaw net.Conn) {
	conn := dgtcp.New(connRaw)
	buf := make([]byte, 65536)
	size, _ := conn.Read(buf)
	fmt.Println("recieved:", buf[:size], "from", conn.RemoteAddr())
	conn.Close()
}
```
