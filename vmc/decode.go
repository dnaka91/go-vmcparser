package vmc

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/dnaka91/go-vmcparser/osc"
)

func getString(buf []byte) ([]byte, []byte, error) {
	value, newBuf, err := osc.ReadString(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading string: %w", err)
	}

	return value, newBuf, nil
}

func getTypeTags(buf []byte) ([]byte, []byte, error) {
	tags, newBuf, err := osc.ReadTypeTags(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading type tags: %w", err)
	}

	return tags, newBuf, nil
}

func getInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf[0:4]))
}

func getFloat32(buf []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(buf[0:4]))
}

func getVec3(buf []byte) Vec3 {
	return Vec3{
		X: math.Float32frombits(binary.BigEndian.Uint32(buf[0:4])),
		Y: math.Float32frombits(binary.BigEndian.Uint32(buf[4:8])),
		Z: math.Float32frombits(binary.BigEndian.Uint32(buf[8:12])),
	}
}

func getVec4(buf []byte) Vec4 {
	return Vec4{
		X: math.Float32frombits(binary.BigEndian.Uint32(buf[0:4])),
		Y: math.Float32frombits(binary.BigEndian.Uint32(buf[4:8])),
		Z: math.Float32frombits(binary.BigEndian.Uint32(buf[8:12])),
		W: math.Float32frombits(binary.BigEndian.Uint32(buf[12:16])),
	}
}
