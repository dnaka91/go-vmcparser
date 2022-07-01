package osc_test

import (
	"fmt"

	"github.com/dnaka91/go-vmcparser/osc"
)

func ExampleReadPacket() {
	raw := []byte("/hi\x00,s\x00\x00hello\x00\x00\x00")

	packet, _, err := osc.ReadPacket(raw)
	if err != nil {
		panic(err)
	}

	fmt.Println(packet)
	// Output: Packet { Message "/hi" "s" [[104 101 108 108 111]] }
}

func ExamplePacket_Iterate() {
	raw := []byte("#bundle\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/a\x00\x00,i\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/b\x00\x00,i\x00\x00\x00\x00\x00\x02")

	packet, _, err := osc.ReadPacket(raw)
	if err != nil {
		panic(err)
	}

	err = packet.Iterate(func(m *osc.Message) error {
		fmt.Println(m)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	// Message "/a" "i" [1]
	// Message "/b" "i" [2]
}

func ExamplePacket_ToMessages() {
	raw := []byte("#bundle\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/a\x00\x00,i\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/b\x00\x00,i\x00\x00\x00\x00\x00\x02")

	packet, _, err := osc.ReadPacket(raw)
	if err != nil {
		panic(err)
	}

	for _, msg := range packet.ToMessages() {
		fmt.Println(msg)
	}

	// Output:
	// Message "/a" "i" [1]
	// Message "/b" "i" [2]
}

func ExampleMessage_useArguments() {
	// The following shows how to easily cast arguments to the right type.
	// The assertions won't panic as the type tag guarantees the proper type
	// was assigned during the parsing process.

	raw := []byte("/a\x00\x00,iT\x00\x00\x00\x00\x05")

	packet, _, err := osc.ReadPacket(raw)
	if err != nil {
		panic(err)
	}
	if packet.Message == nil {
		panic("message must be present")
	}

	msg := packet.Message
	if string(msg.TypeTags) != "iT" {
		panic("unexpected type tags")
	}

	fmt.Println("arg 1:", msg.Arguments[0].(int32))
	fmt.Println("arg 2:", msg.Arguments[1].(bool))

	// Output:
	// arg 1: 5
	// arg 2: true
}
