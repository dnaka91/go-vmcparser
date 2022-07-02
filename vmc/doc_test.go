package vmc_test

import (
	"fmt"
	"net"
	"time"

	"github.com/dnaka91/go-vmcparser/vmc"
)

func ExampleParseMessage() {
	// A VMC "Available" message in raw form.
	raw := []byte("/VMC/Ext/OK\x00,iii\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x03\x00\x00\x00\x01")

	// Parse the raw content into a VMC message.
	message, err := vmc.ParseMessage(raw)
	if err != nil {
		panic(err)
	}

	switch m := message.(type) {
	case *vmc.Available:
		// Let's just print it out.
		fmt.Println(m)
	default:
		panic("message must be the `Available` VMC message")
	}

	// Output: &{true Calibrated MrNormal <nil>}
}

func ExampleUDPServer_Read() {
	// Create a new UDP listener at the VMC default port.
	conn, err := net.ListenPacket("udp", ":39539")
	if err != nil {
		panic(err)
	}

	// Create a new VMC server over UDP with default buffer size.
	server := vmc.NewUDPServer(conn, vmc.BufSizeMaxMTU)

	// Try to read some incoming data, within the next 10 seconds.
	// Our callback might be called multiple times (in case of an OSC bundle).
	err = server.Read(10*time.Second, func(addr net.Addr, raw []byte, message vmc.Message) error {
		fmt.Printf("got message from %v: %v\n", addr, message)

		// Do something with the message ...

		return nil
	})
	if err != nil {
		panic(err)
	}
}
