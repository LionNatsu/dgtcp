# DGTCP

## Overview

**Datagram Transmission Control Protocol** (DGTCP) is an TCP-based datagram-style protocol. It provide
a datagram model over TCP (which is stream model).

DGTCP can handle packets not larger than 65536 bytes.

## Protocol Details

For sending a packet, DGTCP sends two bytes firstly, to show the length of the following data, in *network order* i.e. big-endian byte order.

Here is an example of 2 packets sent in DGTCP.

	[ packet #1: 00 01 02 ...(0x1FF bytes) ]
	[ packet #2: 23 12 AB ]
	...

The following chart shows how these two packets was sent
in the underlying TCP stream.

	<len> <          data         >   <len> < data >   ...
	01 FF 00 01 02 ...(0x1FF bytes)   00 03 23 12 AB   ...


## Golang Package

[![GoDoc](https://godoc.org/github.com/LionNatsu/dgtcp?status.svg)](https://godoc.org/github.com/LionNatsu/dgtcp)

Package dgtcp overrides Read() and Write() methods of `net.TCPConn`.
### Usage & Example

Client:

```go
package main

import (
	"net"

	"github.com/LionNatsu/dgtcp"
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
package main

import (
	"fmt"
	"net"

	"github.com/LionNatsu/dgtcp"
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
