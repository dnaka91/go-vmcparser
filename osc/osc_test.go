package osc_test

import (
	"testing"

	"github.com/dnaka91/go-vmcparser/osc"
	"github.com/stretchr/testify/assert"
)

func assertPacket(t *testing.T, input []byte, want *osc.Packet) {
	t.Helper()

	got, buf, err := osc.ReadPacket(input)
	assert.NoError(t, err)
	assert.Empty(t, buf)
	assert.Equal(t, want, got)
}

func TestOscillatorSample(t *testing.T) {
	input := []byte("/oscillator/4/frequency\x00,f\x00\x00\x43\xdc\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:   []byte("/oscillator/4/frequency"),
			TypeTags:  []byte("f"),
			Arguments: []interface{}{float32(440)},
			Raw:       input,
		},
	})
}

func TestFooSample(t *testing.T) {
	input := []byte("/foo\x00\x00\x00\x00,iisff\x00\x00\x00\x00\x03\xe8\xff\xff\xff\xffhello\x00\x00\x00\x3f\x9d\xf3\xb6\x40\xb5\xb2\x2d")

	assertPacket(t, input, &osc.Packet{
		Message: &osc.Message{
			Address:  []byte("/foo"),
			TypeTags: []byte("iisff"),
			Arguments: []interface{}{
				int32(1000),
				int32(-1),
				[]byte("hello"),
				float32(1.234),
				float32(5.678),
			},
			Raw: input,
		},
	})
}

func TestBundleWithMessage(t *testing.T) {
	input := []byte("#bundle\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/\x00\x00\x00,s\x00\x00hi\x00\x00")

	assertPacket(t, input, &osc.Packet{
		Bundle: &osc.Bundle{
			TimeTag: 1,
			Contents: []osc.Packet{{
				Message: &osc.Message{
					Address:   []byte("/"),
					TypeTags:  []byte("s"),
					Arguments: []interface{}{[]byte("hi")},
					Raw:       []byte("/\x00\x00\x00,s\x00\x00hi\x00\x00"),
				},
			}},
		},
	})
}

func TestPacketIterator(t *testing.T) {
	input := []byte("#bundle\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/\x00\x00\x00,s\x00\x00hi\x00\x00")
	want := &osc.Message{
		Address:   []byte("/"),
		TypeTags:  []byte("s"),
		Arguments: []interface{}{[]byte("hi")},
		Raw:       []byte("/\x00\x00\x00,s\x00\x00hi\x00\x00"),
	}

	packet, buf, err := osc.ReadPacket(input)
	assert.NoError(t, err)
	assert.Empty(t, buf)

	count := 0
	got := &osc.Message{}

	err = packet.Iterate(func(m *osc.Message) error {
		count++
		got = m
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, want, got)
}

func TestPacketToMessages(t *testing.T) {
	input := []byte("#bundle\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0c/\x00\x00\x00,s\x00\x00hi\x00\x00")
	want := []*osc.Message{{
		Address:   []byte("/"),
		TypeTags:  []byte("s"),
		Arguments: []interface{}{[]byte("hi")},
		Raw:       []byte("/\x00\x00\x00,s\x00\x00hi\x00\x00"),
	}}

	packet, buf, err := osc.ReadPacket(input)
	assert.NoError(t, err)
	assert.Empty(t, buf)

	got := packet.ToMessages()
	assert.Equal(t, want, got)
}
