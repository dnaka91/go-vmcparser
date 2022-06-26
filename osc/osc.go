// Package osc implements parsing of "Open Sound Control" packets, in a read-only fashion.
//
// Only reading/parsing is implemented and allocations are avoided wherever possible. The use case
// is to only inspect OSC packets and then pass them on to some other application for handling.
package osc

import (
	"errors"
	"fmt"
)

// Possible errors while reading OSC packets.
var (
	ErrInputEmpty              = errors.New("input data is empty")
	ErrInvalidPacket           = errors.New("invalid packet (neither message nor bundle)")
	ErrTypeTagsStartMissing    = errors.New("expected start of type tags")
	ErrArraysNotSupported      = errors.New("arrays not supported")
	ErrInvalidBundleIdentifier = errors.New("invalid bundle identifier")
	ErrElementTooShort         = errors.New("element content is too short")
)

// UnknownTypeTagError occurs when an unknown type tag was discovered during parsing.
type UnknownTypeTagError struct {
	Tag rune // Tag is the unexpected type tag.
}

var _ error = (*UnknownTypeTagError)(nil)

func (e UnknownTypeTagError) Error() string {
	return fmt.Sprintf("unknown type tag `%c`", e.Tag)
}

// Standard OSC type tags.
const (
	TypeTagInt    = 'i' // 32-bit integer.
	TypeTagFloat  = 'f' // 32-bit floating point number.
	TypeTagString = 's' // OSC string.
	TypeTagBlob   = 'b' // OSC blob.
)

// Extended (non-standard) type tags.
const (
	TypeTagInt64      = 'h' // 64-bit integer.
	TypeTagTimeTag    = 't' // OSC time tag.
	TypeTagDouble     = 'd' // 64-bit floating point number.
	TypeTagSymbol     = 'S' // Alternate OSC string (symbol).
	TypeTagChar       = 'c' // 32-bit character.
	TypeTagRgba       = 'r' // 32-bit RGBA color.
	TypeTagMidi       = 'm' // 4 byte MIDI message.
	TypeTagTrue       = 'T' // Boolean true, no argument data.
	TypeTagFalse      = 'F' // Boolean false, no argument data.
	TypeTagNil        = 'N' // Nil value, no argument data.
	TypeTagInfinitum  = '|' // Infinitum value, no argument data.
	TypeTagArrayStart = '[' // Start indicator of an array.
	TypeTagArrayEnd   = ']' // End indicator of an array.
)

// Packet is a complete OSC packet, that is either a message or a bundle.
//
// Only one of the fields is set at any time, and one field is always set, after successfully
// parsing a packet with ReadPacket. Otherwise, the packet is considered invalid.
type Packet struct {
	Message *Message
	Bundle  *Bundle
}

var _ fmt.Stringer = (*Packet)(nil)

func (p Packet) String() string {
	if p.Message != nil {
		return fmt.Sprintf("Packet { %v }", p.Message)
	}

	if p.Bundle != nil {
		return fmt.Sprintf("Packet { %v }", p.Bundle)
	}

	return "Packet { <invalid> }"
}

// Iterate unpacks the packet into individual messages and calls the given handler for each. In case
// the handler returns an error, it is returned from this method.
func (p Packet) Iterate(handler func(msg *Message) error) error {
	if p.Message != nil {
		return handler(p.Message)
	}

	if p.Bundle != nil {
		for _, packet := range p.Bundle.Contents {
			if err := packet.Iterate(handler); err != nil {
				return err
			}
		}
	}

	return nil
}

// ToMessages unpacks the packet into individual messages. If the packet is a message, it'll just
// return a single element slice. If it is a bundle, it'll recursively iterate over the contents and
// extract all messages into a single slice.
//
// Note: This may cause several slice allocations, depending on the level of nesting in the bundles.
// The Iterate method is a more lightweight alternative.
func (p Packet) ToMessages() []*Message {
	if p.Message != nil {
		return []*Message{p.Message}
	}

	if p.Bundle != nil {
		messages := make([]*Message, 0, len(p.Bundle.Contents))
		for _, packet := range p.Bundle.Contents {
			messages = append(messages, packet.ToMessages()...)
		}

		return messages
	}

	return nil
}

// ReadPacket reads and parses a raw byte slice into a OSC packet. The remaining bytes (if any) are
// returned for further processing by the user, as well.
func ReadPacket(buf []byte) (*Packet, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, ErrInputEmpty
	}

	switch buf[0] {
	case '/':
		message, newBuf, err := readMessage(buf)
		if err != nil {
			return nil, nil, err
		}

		return &Packet{
			Message: message,
			Bundle:  nil,
		}, newBuf, nil
	case '#':
		bundle, newBuf, err := readBundle(buf)
		if err != nil {
			return nil, nil, err
		}

		return &Packet{
			Message: nil,
			Bundle:  bundle,
		}, newBuf, nil
	default:
		return nil, nil, ErrInvalidPacket
	}
}

// Message is a single OSC message, that contains an address to identify its type, type tags to
// describe the argument types, and the list of arguments.
//
// The type tags can be used to infer the argument type of the generic arguments field, allowing for
// easy casting of the type. For example, with type tag `sf`, it is guaranteed that the first
// argument is a Go `string` and the second argument is a Go `float32`.
//
// Mapping from type tag to Go types is as follows:
//
//     # Standard OSC type tags.
//
//     i -> int32
//     f -> float32
//     s -> string
//     b -> []byte
//
//     # Extended (non-standard) type tags.
//
//     h -> int64
//     t -> int64
//     d -> float64
//     S -> string
//     c -> rune
//     r -> [4]byte
//     m -> [4]byte
//     T -> bool
//     F -> bool
//     N -> nil
//     | -> nil
//     [ -> not supported!
//     ] -> not supported!
//
// The Raw field can be used to forward the original message to any real VMC server, to handle it.
// This is especially helpful, when the message is part of a bundle, and only some of them are
// supposed to be forwarded the target, instead of the dropping the whole bundle.
type Message struct {
	Address   string        // Address is the message address
	TypeTags  string        // TypeTags contains OSC type tags for each argument
	Arguments []interface{} // Arguments contains all parsed arguments
	Raw       []byte        // Raw is the un-parsed message content
}

var _ fmt.Stringer = (*Message)(nil)

func (m Message) String() string {
	return fmt.Sprintf("Message \"%v\" \"%v\" %v", m.Address, m.TypeTags, m.Arguments)
}

func readMessage(buf []byte) (*Message, []byte, error) {
	raw := buf

	address, newBuf, err := readString(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading address: %w", err)
	}
	buf = newBuf

	if len(buf) == 0 || buf[0] != ',' {
		return nil, nil, ErrTypeTagsStartMissing
	}

	typeTags, newBuf, err := readString(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading type tags: %w", err)
	}
	buf = newBuf

	arguments := make([]interface{}, len(typeTags)-1)

	for idx, tag := range typeTags[1:] {
		switch tag {
		case TypeTagInt:
			v, b, err := readInt(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagFloat:
			v, b, err := readFloat(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagString, TypeTagSymbol:
			v, b, err := readString(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagBlob:
			v, b, err := readBlob(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagInt64:
			v, b, err := readInt64(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagTimeTag:
			v, b, err := readTimeTag(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagDouble:
			v, b, err := readDouble(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagChar:
			v, b, err := readChar(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagRgba:
			v, b, err := readRgba(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagMidi:
			v, b, err := readMidi(buf)
			if err != nil {
				return nil, nil, err
			}
			buf = b
			arguments[idx] = v
		case TypeTagTrue:
			arguments[idx] = true
		case TypeTagFalse:
			arguments[idx] = false
		case TypeTagNil, TypeTagInfinitum:
			arguments[idx] = nil
		case TypeTagArrayStart, TypeTagArrayEnd:
			return nil, nil, ErrArraysNotSupported
		default:
			return nil, nil, UnknownTypeTagError{Tag: tag}
		}
	}

	return &Message{
		Address:   address,
		TypeTags:  typeTags[1:],
		Arguments: arguments,
		Raw:       raw,
	}, buf, nil
}

// Bundle is a single OCS bundle, which is in turn a collection of packets, that are either messages
// or bundles themselves.
//
// The time tag is kept for reference and only important when executing contained commands, but not
// interesting for inspection of the content.
type Bundle struct {
	TimeTag  int64
	Contents []Packet
}

var _ fmt.Stringer = (*Bundle)(nil)

func (b Bundle) String() string {
	return fmt.Sprintf("Bundle %v %v", b.TimeTag, b.Contents)
}

func readBundle(buf []byte) (*Bundle, []byte, error) {
	ident, newBuf, err := readString(buf)
	if err != nil {
		return nil, nil, err
	}
	buf = newBuf

	if ident != "#bundle" {
		return nil, nil, ErrInvalidBundleIdentifier
	}

	timeTag, newBuf, err := readTimeTag(buf)
	if err != nil {
		return nil, nil, err
	}
	buf = newBuf

	contents := []Packet{}

	for len(buf) > 4 {
		length, newBuf, err := readLength(buf)
		if err != nil {
			return nil, nil, err
		}
		buf = newBuf

		if len(buf) < length {
			return nil, nil, ErrElementTooShort
		}

		packet, newBuf, err := ReadPacket(buf)
		if err != nil {
			return nil, nil, err
		}
		buf = newBuf

		contents = append(contents, *packet)
	}

	return &Bundle{
		TimeTag:  timeTag,
		Contents: contents,
	}, buf, nil
}
