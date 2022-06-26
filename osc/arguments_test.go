package osc_test

import (
	"testing"

	"github.com/dnaka91/go-vmcparser/osc"
)

func TestParseInt(t *testing.T) {
	input := []byte("/\x00\x00\x00,i\x00\x00\x00\x00\x00\x05")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "i",
			Arguments: []interface{}{int32(5)},
			Raw:       input,
		},
	})
}

func TestParseFloat(t *testing.T) {
	input := []byte("/\x00\x00\x00,f\x00\x00\x40\xa0\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "f",
			Arguments: []interface{}{float32(5)},
			Raw:       input,
		},
	})
}

func TestParseString(t *testing.T) {
	input := []byte("/\x00\x00\x00,s\x00\x00tst\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "s",
			Arguments: []interface{}{"tst"},
			Raw:       input,
		},
	})
}

func TestParseBlob(t *testing.T) {
	input := []byte("/\x00\x00\x00,b\x00\x00\x00\x00\x00\x03\x01\x02\x03\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "b",
			Arguments: []interface{}{[]byte{1, 2, 3}},
			Raw:       input,
		},
	})
}

func TestParseInt64(t *testing.T) {
	input := []byte("/\x00\x00\x00,h\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "h",
			Arguments: []interface{}{int64(5)},
			Raw:       input,
		},
	})
}

func TestParseTimeTag(t *testing.T) {
	input := []byte("/\x00\x00\x00,t\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "t",
			Arguments: []interface{}{int64(5)},
			Raw:       input,
		},
	})
}

func TestParseDouble(t *testing.T) {
	input := []byte("/\x00\x00\x00,d\x00\x00\x40\x14\x00\x00\x00\x00\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "d",
			Arguments: []interface{}{float64(5)},
			Raw:       input,
		},
	})
}

func TestParseChar(t *testing.T) {
	input := []byte("/\x00\x00\x00,c\x00\x00\x00\x00\x00a")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "c",
			Arguments: []interface{}{rune('a')},
			Raw:       input,
		},
	})
}

func TestParseRgba(t *testing.T) {
	input := []byte("/\x00\x00\x00,r\x00\x00\x01\x02\x03\x04")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "r",
			Arguments: []interface{}{[4]byte{1, 2, 3, 4}},
			Raw:       input,
		},
	})
}

func TestParseMidi(t *testing.T) {
	input := []byte("/\x00\x00\x00,m\x00\x00\x01\x02\x03\x04")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "m",
			Arguments: []interface{}{[4]byte{1, 2, 3, 4}},
			Raw:       input,
		},
	})
}

func TestParseTrue(t *testing.T) {
	input := []byte("/\x00\x00\x00,T\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "T",
			Arguments: []interface{}{true},
			Raw:       input,
		},
	})
}

func TestParseFalse(t *testing.T) {
	input := []byte("/\x00\x00\x00,F\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "F",
			Arguments: []interface{}{false},
			Raw:       input,
		},
	})
}

func TestParseNil(t *testing.T) {
	input := []byte("/\x00\x00\x00,N\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "N",
			Arguments: []interface{}{nil},
			Raw:       input,
		},
	})
}

func TestParseInfinitum(t *testing.T) {
	input := []byte("/\x00\x00\x00,|\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   "/",
			TypeTags:  "|",
			Arguments: []interface{}{nil},
			Raw:       input,
		},
	})
}
