package vmc_test

import (
	"testing"

	"github.com/dnaka91/go-vmc/osc"
	"github.com/dnaka91/go-vmc/vmc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertMessage(t *testing.T, input []byte, want vmc.Message) {
	t.Helper()

	packet, buf, err := osc.ReadPacket(input)
	require.NoError(t, err)
	require.Empty(t, buf)
	require.NotNil(t, packet.Message)

	got, err := vmc.ParseMessage(packet.Message)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseAvailable(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/OK\x00,iii\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x03\x00\x00\x00\x01"),
		&vmc.Available{
			Loaded:           true,
			CalibrationState: vmc.CalibrationStateCalibrated,
			CalibrationMode:  vmc.CalibrationModeMrNormal,
		},
	)
}

func TestParseRelativeTime(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/T\x00\x00,f\x00\x00\x40\xa0\x00\x00"),
		&vmc.RelativeTime{
			Time: 5,
		},
	)
}

func TestParseRootTransform(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Root/Pos\x00\x00\x00,sfffffffffffff\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a\x40\x46\x66\x66\x40\x4c\xcc\xcd\x40\x53\x33\x33\x40\x83\x33\x33\x40\x86\x66\x66\x40\x89\x99\x9a"),
		&vmc.RootTransform{
			Name:       "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
			Scale:      vmc.Vec3{X: 3.1, Y: 3.2, Z: 3.3},
			Offset:     vmc.Vec3{X: 4.1, Y: 4.2, Z: 4.3},
		},
	)
}

func TestParseBoneTransform(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Bone/Pos\x00\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.BoneTransform{
			Name:       "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
}

func TestParseBlendShapeProxyValue(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Blend/Val\x00\x00,sf\x00tst\x00\x40\xa0\x00\x00"),
		&vmc.BlendShapeProxyValue{
			Name:  "tst",
			Value: 5,
		},
	)
}

func TestParseBlendShapeProxyApply(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Blend/Apply\x00\x00\x00\x00,\x00\x00\x00"),
		&vmc.BlendShapeProxyApply{},
	)
}

func TestParseCameraTransform(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Cam\x00\x00\x00\x00,sffffffff\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a\x40\xa0\x00\x00"),
		&vmc.CameraTransform{
			Name:       "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
			FOV:        5,
		},
	)
}

func TestParseControllerInput(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Con\x00\x00\x00\x00,isiiifff\x00\x00\x00\x00\x00\x00\x01tst\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66"),
		&vmc.ControllerInput{
			Active:  vmc.ControllerActivePress,
			Name:    "tst",
			IsLeft:  true,
			IsTouch: true,
			IsAxis:  false,
			Axis:    vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
		},
	)
}

func TestParseKeyboardInput(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Key\x00\x00\x00\x00,isi\x00\x00\x00\x00\x00\x00\x00\x01tst\x00\x00\x00\x00\x05"),
		&vmc.KeyboardInput{
			Active:  true,
			Name:    "tst",
			KeyCode: 5,
		},
	)
}

func TestParseMidiNodeInput(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Midi/Note\x00\x00,iiif\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x02\x3f\x8c\xcc\xcd"),
		&vmc.MidiNodeInput{
			Active:   true,
			Channel:  1,
			Note:     2,
			Velocity: 1.1,
		},
	)
}

func TestParseMidiCCValueInput(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Midi/CC/Val\x00\x00\x00\x00,if\x00\x00\x00\x00\x01\x3f\x8c\xcc\xcd"),
		&vmc.MidiCCValueInput{
			Knob:  1,
			Value: 1.1,
		},
	)
}

func TestParseMidiCCButtonInput(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Midi/CC/Bit\x00\x00\x00\x00,ii\x00\x00\x00\x00\x01\x00\x00\x00\x01"),
		&vmc.MidiCCButtonInput{
			Knob:   1,
			Active: true,
		},
	)
}

func TestParseDeviceTransform(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Hmd/Pos\x00\x00\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Hmd",
			Local:      false,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
	assertMessage(
		t,
		[]byte("/VMC/Ext/Con/Pos\x00\x00\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Con",
			Local:      false,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
	assertMessage(
		t,
		[]byte("/VMC/Ext/Tra/Pos\x00\x00\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Tra",
			Local:      false,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)

	assertMessage(
		t,
		[]byte("/VMC/Ext/Hmd/Pos/Local\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Hmd",
			Local:      true,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
	assertMessage(
		t,
		[]byte("/VMC/Ext/Con/Pos/Local\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Con",
			Local:      true,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
	assertMessage(
		t,
		[]byte("/VMC/Ext/Tra/Pos/Local\x00\x00,sfffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.DeviceTransform{
			Device:     "Tra",
			Local:      true,
			Serial:     "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
}

func TestParseReceiveEnable(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Rcv\x00\x00\x00\x00,ii\x00\x00\x00\x00\x01\x00\x00\x1f\x90"),
		&vmc.ReceiveEnable{
			Enable: true,
			Port:   8080,
		},
	)
}

func TestParseDirectionalLight(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Light\x00\x00,sfffffffffff\x00\x00\x00tst\x00\x3f\x8c\xcc\xcd\x3f\x99\x99\x9a\x3f\xa6\x66\x66\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a\x40\x46\x66\x66\x40\x4c\xcc\xcd\x40\x53\x33\x33\x40\x59\x99\x9a"),
		&vmc.DirectionalLight{
			Name:       "tst",
			Position:   vmc.Vec3{X: 1.1, Y: 1.2, Z: 1.3},
			Quaternion: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
			Color:      vmc.Vec4{X: 3.1, Y: 3.2, Z: 3.3, W: 3.4},
		},
	)
}

func TestParseLocalVrm(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/VRM\x00\x00\x00\x00,ss\x00t01\x00t02\x00"),
		&vmc.LocalVrm{
			Path:  "t01",
			Title: "t02",
		},
	)
}

func TestParseRemoteVrm(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Remote\x00,ss\x00tst\x00{}\x00\x00"),
		&vmc.RemoteVrm{
			Service: "tst",
			JSON:    "{}",
		},
	)
}

func TestParseOptionString(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Opt\x00\x00\x00\x00,s\x00\x00tst\x00"),
		&vmc.OptionString{
			Option: "tst",
		},
	)
}

func TestParseBackgroundColor(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Setting/Color\x00\x00,ffff\x00\x00\x00\x40\x06\x66\x66\x40\x0c\xcc\xcd\x40\x13\x33\x33\x40\x19\x99\x9a"),
		&vmc.BackgroundColor{
			Color: vmc.Vec4{X: 2.1, Y: 2.2, Z: 2.3, W: 2.4},
		},
	)
}

func TestParseWindowAttribute(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Setting/Win\x00\x00\x00\x00,iiii\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x01"),
		&vmc.WindowAttribute{
			IsTopMost:          true,
			IsTransparent:      true,
			WindowClickThrough: true,
			HideBorder:         true,
		},
	)
}

func TestParseLoadedSettingPath(t *testing.T) {
	assertMessage(
		t,
		[]byte("/VMC/Ext/Config\x00,s\x00\x00tst\x00"),
		&vmc.LoadedSettingPath{
			Path: "tst",
		},
	)
}
