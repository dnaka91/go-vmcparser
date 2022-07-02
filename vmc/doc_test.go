package vmc_test

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/dnaka91/go-vmcparser/osc"
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

func Example_udpServer() {
	// Create a new UDP listener at the VMC default port.
	conn, err := net.ListenPacket("udp", ":39539")
	if err != nil {
		panic(err)
	}

	// Create a new buffer to read UPD payloads into.
	//
	// The value 1536 is a common maximum size for Ethernet II, meaning we only need
	// a single (or few) system calls to get the content from the OS handlers.
	//
	// Pick any other value as you like ;-)
	buf := make([]byte, 1536)

	// Start an endless loop, trying to get messages until the end of time, or some error happens.
	for {
		// Read a new packet from the connection.
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatalf("failed to read from the UDP connection: %v", err)
		}

		// Fail if we got an OSC bundle, instead of a single OSC message.
		if osc.IsBundle(buf[:n]) {
			log.Fatal("got an OSC bundle, we don't handle them (yet)")
		}

		// Parse the message into a known VMC message.
		//
		// Here we pass some extra filters, so we'll only fully parse
		// the root and bone transform messages (for best possible performance).
		message, err := vmc.ParseMessage(
			buf,
			vmc.AddressRootTransform,
			vmc.AddressBoneTransform,
		)
		if errors.Is(err, vmc.ErrFiltered) {
			// Here we skip on filtered messages, could as well forward the raw buffer to
			// some real VMC handler.
			continue
		}
		if err != nil {
			log.Fatalf("failed to parse VMC message: %v", err)
		}

		// Finally we got our message parsed and ready. Now we can cast it into
		// one of the several defined messages and access their content as proper
		// Go structs.
		switch m := message.(type) {
		case *vmc.RootTransform:
			log.Printf("new root transformation: %v\n", m)
			// Do something actually specific with the message ...
		case *vmc.BoneTransform:
			log.Printf("new root transformation, named  %v with position %v\n", m.Name, m.Position)
		default:
			log.Printf("got message from %v: %v\n", addr, message)
		}
	}
}
