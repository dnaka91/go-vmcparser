package vmc

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/dnaka91/go-vmc/osc"
)

// Various pre-defined and recommended buffer sizes to use with NewUDPServer.
const (
	// BufSizeMaxMTU is the maximum transmission unit for Ethernet II.
	//
	// This is an optimized value that can be used as buffer size, when creating a NewUDPServer. It
	// gives the best possible performance as it allows the OS to read the buffer in the biggest
	// possible single packet.
	BufSizeMaxMTU = 1536

	// BufSizeLarge defines a fairly large buffer size of 16kB, in case BufSizeMaxMTU is not
	// sufficient.
	BufSizeLarge = 16384

	// BufSizeHuge defines an even large buffer size of 64kB (the maximum value of an uint16). Only
	// use this if you know you'll need it, as it allocates a lot of data.
	BufSizeHuge = 65535
)

// UDPHandler is a function that handles a single VMC message. The given network address is the
// origin where the message was received from. The raw value is a slice of the internal UPDServer's
// buffer with the original unparsed payload.
//
// The handler can return an error, cancelling any further processing of the received payload. That
// is, in case the read raw data contained more than a single VMC message, or the payload is an OSC
// bundle (which can hold multiple messages).
//
// Warning: Don't keep the `raw` byte slice around, as it is a view into the the internal UDPServer
// buffer. It'll be overwritten with new data on the next call to Read. If you need to keep the data
// for a longer time, copy the content with the built-in `copy` function.
type UDPHandler = func(addr net.Addr, raw []byte, message Message) error

// UDPServer is a VMC message reader over a UDP connection.
//
// This server keeps and re-uses an internal buffer to read messages, reducing the amount of
// allocations required to read messages.
type UDPServer struct {
	conn net.PacketConn // UDP server connection.
	buf  []byte         // Internal, re-usable buffer, to avoid allocations.
}

// Read tries to receive a new message. The handler might be called multiple times, for each message
// received in a single read. For example, the payload could contain multiple message at once.
func (s *UDPServer) Read(timeout time.Duration, handler UDPHandler) error {
	if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("failed to set read deadline on the connection: %w", err)
	}

	n, addr, err := s.conn.ReadFrom(s.buf)
	if err != nil {
		return fmt.Errorf("failed to read from the UDP connection: %w", err)
	}

	buf := s.buf[:n]

	for len(buf) > 0 {
		packet, newBuf, err := osc.ReadPacket(buf)
		if err != nil {
			return fmt.Errorf("failed to parse OSC packet: %w", err)
		}
		buf = newBuf

		err = packet.Iterate(func(msg *osc.Message) error {
			message, err := ParseMessage(msg)
			if errors.Is(err, ErrUnknownAddress) {
				// skip any unknown VMC messages.
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to parse VMC message: %w", err)
			}

			return handler(addr, msg.Raw, message)
		})
		if err != nil {
			// no error wrapping, this is just the inner error from `Iterate`.
			return fmt.Errorf("failed handling VMC message: %w", err)
		}
	}

	return nil
}

// NewUDPServer create a new, simple UDP server that can read VMC messages. Prefer to use one of the
// buffer size constants for the right buffer size (but a custom size is suitable as well).
func NewUDPServer(conn net.PacketConn, bufSize int) UDPServer {
	return UDPServer{
		conn: conn,
		buf:  make([]byte, bufSize),
	}
}
