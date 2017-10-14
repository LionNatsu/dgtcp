// Package dgtcp provides a simple implementation of
// TCP-based datagram-style protocol.
// It overrides Read() and Write(). It can send and recieve not larger
// than 65536 bytes per packet.
//
// Read() and Write() are coroutine-safe.
//
// DGTCP Protocol:
//   [       packet #1             ]  [ packet #2  ] ...
//   <len> < data                  >  <len> < data > ...
//   01 FF 00 01 02 ...(0x1FF bytes)  00 03 23 12 AB ...
package dgtcp

import (
	"errors"
	"net"
	"sync"
)

// DGTCPConn is the instance of DGTCP connection.
// It has all methods of *net.TCPConn.
type DGTCPConn struct {
	*net.TCPConn
	readLen int
	wl, rl  sync.Locker
}

// If your buffer is not big enough to recieve(read), or your buffer is too
// big to send(write), ErrOverflow will be returned.
var ErrOverflow = errors.New("packet is too large to recieve or send")

func (c *DGTCPConn) fillBuf(buf []byte) error {
	var pos int
	for pos != len(buf) {
		size, err := c.TCPConn.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos += size
	}
	return nil
}

// Read recieves one packet and put it into buf.
func (c *DGTCPConn) Read(buf []byte) (int, error) {
	c.rl.Lock()
	defer c.rl.Unlock()

	var plen [2]byte
	var ilen int

	if c.readLen == 0 {
		err := c.fillBuf(plen[:])
		if err != nil {
			return 0, err
		}
		ilen = int(plen[0])<<8 + int(plen[1])
	} else {
		ilen = c.readLen
	}

	if ilen > len(buf) {
		c.readLen = ilen
		return 0, ErrOverflow
	}
	err := c.fillBuf(buf[:ilen])
	c.readLen = 0
	return ilen, err
}

// Write sends one packet.
func (c *DGTCPConn) Write(buf []byte) error {
	c.wl.Lock()
	defer c.wl.Unlock()

	var ilen = len(buf)
	if ilen > 65536 {
		return ErrOverflow
	}
	var plen [2]byte
	plen[0] = byte((ilen & 0xff00) >> 8)
	plen[1] = byte(ilen & 0xff)
	_, err := c.TCPConn.Write(plen[:])
	if err != nil {
		return err
	}
	_, err = c.TCPConn.Write(buf)
	return err
}

// New returns an instance of DGTCPConn. c must be a pointer of net.TCPConn.
func New(c net.Conn) *DGTCPConn {
	if tcp, ok := c.(*net.TCPConn); ok {
		return &DGTCPConn{tcp, 0, new(sync.Mutex), new(sync.Mutex)}
	}
	return nil
}
