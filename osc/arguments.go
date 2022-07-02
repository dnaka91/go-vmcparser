package osc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Possible errors while reading basic OSC types.
var (
	ErrIntTooShort             = errors.New("content is too short for an int")
	ErrFloatTooShort           = errors.New("content is too short for a float")
	ErrStringMissingTerminator = errors.New("string missing 0 terminator")
	ErrBlobTooShort            = errors.New("content is too short for a blob")
	ErrInt64TooShort           = errors.New("content is too short for an int64")
	ErrTimeTagTooShort         = errors.New("content is too short for an time tag")
	ErrDoubleTooShort          = errors.New("content is too short for a double")
	ErrCharTooShort            = errors.New("content is too short for a char")
	ErrRgbaTooShort            = errors.New("content is too short for rgba")
	ErrMidiTooShort            = errors.New("content is too short for midi")
	ErrNegativeLength          = errors.New("invalid negative length")
)

// Fixed lengths for the different OSC types.
const (
	lenInt     = 4
	lenFloat   = 4
	lenInt64   = 8
	lenTimeTag = 8
	lenDouble  = 8
	lenChar    = 4
	lenRgba    = 4
	lenMidi    = 4
)

func readInt(buf []byte) (int32, []byte, error) {
	if len(buf) < lenInt {
		return 0, nil, ErrIntTooShort
	}

	return int32(binary.BigEndian.Uint32(buf[:lenInt])), buf[lenInt:], nil
}

func readFloat(buf []byte) (float32, []byte, error) {
	if len(buf) < lenFloat {
		return 0, nil, ErrFloatTooShort
	}

	return math.Float32frombits(binary.BigEndian.Uint32(buf[:lenFloat])), buf[lenFloat:], nil
}

// ReadString reads the string content from the given raw OSC encoded content, and returns it
// together with the advanced buffer and potential error if decoding failed.
func ReadString(buf []byte) ([]byte, []byte, error) {
	pos := bytes.IndexByte(buf, 0)
	if pos == -1 {
		return nil, nil, ErrStringMissingTerminator
	}

	value := buf[:pos]

	return value, buf[len(value)+pad(len(value)):], nil
}

func readLength(buf []byte) (int, []byte, error) {
	length, newBuf, err := readInt(buf)
	if err != nil {
		return 0, nil, fmt.Errorf("failed reading length: %w", err)
	}

	if length < 0 {
		return 0, nil, ErrNegativeLength
	}

	return int(length), newBuf, nil
}

func readBlob(buf []byte) ([]byte, []byte, error) {
	length, newBuf, err := readLength(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading blob length: %w", err)
	}
	buf = newBuf

	if len(buf) < length {
		return nil, nil, ErrBlobTooShort
	}

	return buf[:length], buf[length+pad(length):], nil
}

func pad(length int) int {
	const padSize = 4

	return padSize - length%padSize
}

func readInt64(buf []byte) (int64, []byte, error) {
	if len(buf) < lenInt64 {
		return 0, nil, ErrInt64TooShort
	}

	return int64(binary.BigEndian.Uint64(buf[:lenInt64])), buf[lenInt64:], nil
}

func readTimeTag(buf []byte) (int64, []byte, error) {
	if len(buf) < lenTimeTag {
		return 0, nil, ErrTimeTagTooShort
	}

	return int64(binary.BigEndian.Uint64(buf[:lenTimeTag])), buf[lenTimeTag:], nil
}

func readDouble(buf []byte) (float64, []byte, error) {
	if len(buf) < lenDouble {
		return 0, nil, ErrDoubleTooShort
	}

	return math.Float64frombits(binary.BigEndian.Uint64(buf[:lenDouble])), buf[lenDouble:], nil
}

func readChar(buf []byte) (rune, []byte, error) {
	if len(buf) < lenChar {
		return 0, nil, ErrCharTooShort
	}

	return rune(binary.BigEndian.Uint32(buf[:lenChar])), buf[lenChar:], nil
}

func readRgba(buf []byte) ([4]byte, []byte, error) {
	if len(buf) < lenRgba {
		return [lenRgba]byte{}, nil, ErrRgbaTooShort
	}

	return [lenRgba]byte{buf[0], buf[1], buf[2], buf[3]}, buf[lenRgba:], nil
}

func readMidi(buf []byte) ([4]byte, []byte, error) {
	if len(buf) < lenMidi {
		return [lenMidi]byte{}, nil, ErrMidiTooShort
	}

	return [lenMidi]byte{buf[0], buf[1], buf[2], buf[3]}, buf[lenMidi:], nil
}
