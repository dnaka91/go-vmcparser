package vmc

import (
	"fmt"

	"github.com/dnaka91/go-vmc/osc"
)

// OSC message addresses for the VMC marionette messages.
const (
	AddressAvailable               = "/VMC/Ext/OK"
	AddressRelativeTime            = "/VMC/Ext/T"
	AddressRootTransform           = "/VMC/Ext/Root/Pos"
	AddressBoneTransform           = "/VMC/Ext/Bone/Pos"
	AddressBlendShapeProxyValue    = "/VMC/Ext/Blend/Val"
	AddressBlendShapeProxyApply    = "/VMC/Ext/Blend/Apply"
	AddressCameraTransform         = "/VMC/Ext/Cam"
	AddressControllerInput         = "/VMC/Ext/Con"
	AddressKeyboardInput           = "/VMC/Ext/Key"
	AddressMidiNoteInput           = "/VMC/Ext/Midi/Note"
	AddressMidiCCValueInput        = "/VMC/Ext/Midi/CC/Val"
	AddressMidiCCButtonInput       = "/VMC/Ext/Midi/CC/Bit"
	AddressDeviceTransformHmd      = "/VMC/Ext/Hmd/Pos"
	AddressDeviceTransformCon      = "/VMC/Ext/Con/Pos"
	AddressDeviceTransformTra      = "/VMC/Ext/Tra/Pos"
	AddressDeviceTransformHmdLocal = "/VMC/Ext/Hmd/Pos/Local"
	AddressDeviceTransformConLocal = "/VMC/Ext/Con/Pos/Local"
	AddressDeviceTransformTraLocal = "/VMC/Ext/Tra/Pos/Local"
	AddressReceiveEnable           = "/VMC/Ext/Rcv"
	AddressDirectionalLight        = "/VMC/Ext/Light"
	AddressLocalVrm                = "/VMC/Ext/VRM"
	AddressRemoteVrm               = "/VMC/Ext/Remote"
	AddressOptionString            = "/VMC/Ext/Opt"
	AddressBackgroundColor         = "/VMC/Ext/Setting/Color"
	AddressWindowAttribute         = "/VMC/Ext/Setting/Win"
	AddressLoadedSettingPath       = "/VMC/Ext/Config"
)

type Available struct {
	Loaded           bool
	CalibrationState CalibrationState
	CalibrationMode  CalibrationMode
}

func (a *Available) isMessage() {}

func (a *Available) Address() string {
	return AddressAvailable
}

type CalibrationState uint8

// Possible values for the calibration state.
const (
	CalibrationStateUncalibrated CalibrationState = iota
	CalibrationStateWaitingForCalibration
	CalibrationStateCalibrating
	CalibrationStateCalibrated
)

func (s CalibrationState) isValid() bool {
	return s <= CalibrationStateCalibrated
}

func (s CalibrationState) String() string {
	switch s {
	case CalibrationStateUncalibrated:
		return "Uncalibrated"
	case CalibrationStateWaitingForCalibration:
		return "WaitingForCalibration"
	case CalibrationStateCalibrating:
		return "Calibrating"
	case CalibrationStateCalibrated:
		return "Calibrated"
	default:
		return fmt.Sprintf("Unknown(%d)", uint8(s))
	}
}

type CalibrationMode uint8

// Possible values for the calibration mode.
const (
	CalibrationModeNormal CalibrationMode = iota
	CalibrationModeMrNormal
	CalibrationModeMrFloorFix
)

func (m CalibrationMode) isValid() bool {
	return m <= CalibrationModeMrFloorFix
}

func (m CalibrationMode) String() string {
	switch m {
	case CalibrationModeNormal:
		return "Normal"
	case CalibrationModeMrNormal:
		return "MrNormal"
	case CalibrationModeMrFloorFix:
		return "MrFloorFix"
	default:
		return fmt.Sprintf("Unknown(%d)", uint8(m))
	}
}

func parseAvailable(msg *osc.Message) (*Available, error) {
	if msg.TypeTags != "iii" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "iii"}
	}

	calibrationState := CalibrationState(msg.Arguments[1].(int32))
	if !calibrationState.isValid() {
		return nil, InvalidEnumValueError{
			Name:  "calibration state",
			Value: msg.Arguments[1].(int32),
		}
	}

	calibrationMode := CalibrationMode(msg.Arguments[2].(int32))
	if !calibrationMode.isValid() {
		return nil, InvalidEnumValueError{
			Name:  "calibration mode",
			Value: msg.Arguments[2].(int32),
		}
	}

	return &Available{
		Loaded:           msg.Arguments[0].(int32) == 1,
		CalibrationState: calibrationState,
		CalibrationMode:  calibrationMode,
	}, nil
}

type RelativeTime struct {
	Time float32
}

func (r *RelativeTime) isMessage() {}

func (r *RelativeTime) Address() string {
	return AddressRelativeTime
}

func parseRelativeTime(msg *osc.Message) (*RelativeTime, error) {
	if msg.TypeTags != "f" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "f"}
	}

	return &RelativeTime{
		Time: msg.Arguments[0].(float32),
	}, nil
}

type RootTransform struct {
	Name       string
	Position   Vec3
	Quaternion Vec4
	Scale      Vec3
	Offset     Vec3
}

func (r *RootTransform) isMessage() {}

func (r *RootTransform) Address() string {
	return AddressRootTransform
}

func parseRootTransform(msg *osc.Message) (*RootTransform, error) {
	if msg.TypeTags != "sfffffffffffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sfffffffffffff"}
	}

	return &RootTransform{
		Name: msg.Arguments[0].(string),
		Position: Vec3{
			X: msg.Arguments[1].(float32),
			Y: msg.Arguments[2].(float32),
			Z: msg.Arguments[3].(float32),
		},
		Quaternion: Vec4{
			X: msg.Arguments[4].(float32),
			Y: msg.Arguments[5].(float32),
			Z: msg.Arguments[6].(float32),
			W: msg.Arguments[7].(float32),
		},
		Scale: Vec3{
			X: msg.Arguments[8].(float32),
			Y: msg.Arguments[9].(float32),
			Z: msg.Arguments[10].(float32),
		},
		Offset: Vec3{
			X: msg.Arguments[11].(float32),
			Y: msg.Arguments[12].(float32),
			Z: msg.Arguments[13].(float32),
		},
	}, nil
}

type BoneTransform struct {
	Name       string
	Position   Vec3
	Quaternion Vec4
}

func (b *BoneTransform) isMessage() {}

func (b *BoneTransform) Address() string {
	return AddressBoneTransform
}

func parseBoneTransform(msg *osc.Message) (*BoneTransform, error) {
	if msg.TypeTags != "sfffffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sfffffff"}
	}

	return &BoneTransform{
		Name: msg.Arguments[0].(string),
		Position: Vec3{
			X: msg.Arguments[1].(float32),
			Y: msg.Arguments[2].(float32),
			Z: msg.Arguments[3].(float32),
		},
		Quaternion: Vec4{
			X: msg.Arguments[4].(float32),
			Y: msg.Arguments[5].(float32),
			Z: msg.Arguments[6].(float32),
			W: msg.Arguments[7].(float32),
		},
	}, nil
}

type BlendShapeProxyValue struct {
	Name  string
	Value float32
}

func (b *BlendShapeProxyValue) isMessage() {}

func (b *BlendShapeProxyValue) Address() string {
	return AddressBlendShapeProxyValue
}

func parseBlendShapeProxyValue(msg *osc.Message) (*BlendShapeProxyValue, error) {
	if msg.TypeTags != "sf" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sf"}
	}

	return &BlendShapeProxyValue{
		Name:  msg.Arguments[0].(string),
		Value: msg.Arguments[1].(float32),
	}, nil
}

type BlendShapeProxyApply struct{}

func (b *BlendShapeProxyApply) isMessage() {}

func (b *BlendShapeProxyApply) Address() string {
	return AddressBlendShapeProxyApply
}

func parseBlendShapeProxyApply(msg *osc.Message) (*BlendShapeProxyApply, error) {
	if msg.TypeTags != "" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: ""}
	}

	return &BlendShapeProxyApply{}, nil
}

type CameraTransform struct {
	Name       string
	Position   Vec3
	Quaternion Vec4
	FOV        float32
}

func (c *CameraTransform) isMessage() {}

func (c *CameraTransform) Address() string {
	return AddressCameraTransform
}

func parseCameraTransform(msg *osc.Message) (*CameraTransform, error) {
	if msg.TypeTags != "sffffffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sffffffff"}
	}

	return &CameraTransform{
		Name: msg.Arguments[0].(string),
		Position: Vec3{
			X: msg.Arguments[1].(float32),
			Y: msg.Arguments[2].(float32),
			Z: msg.Arguments[3].(float32),
		},
		Quaternion: Vec4{
			X: msg.Arguments[4].(float32),
			Y: msg.Arguments[5].(float32),
			Z: msg.Arguments[6].(float32),
			W: msg.Arguments[7].(float32),
		},
		FOV: msg.Arguments[8].(float32),
	}, nil
}

type ControllerInput struct {
	Active  ControllerActive
	Name    string
	IsLeft  bool
	IsTouch bool
	IsAxis  bool
	Axis    Vec3
}

func (c *ControllerInput) isMessage() {}

func (c *ControllerInput) Address() string {
	return AddressControllerInput
}

type ControllerActive uint8

// Possible values for the controller active state.
const (
	ControllerActiveRelease ControllerActive = iota
	ControllerActivePress
	ControllerActiveChangeAxis
)

func (a ControllerActive) isValid() bool {
	return a <= ControllerActiveChangeAxis
}

func parseControllerInput(msg *osc.Message) (*ControllerInput, error) {
	if msg.TypeTags != "isiiifff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "isiiifff"}
	}

	active := ControllerActive(msg.Arguments[0].(int32))
	if !active.isValid() {
		return nil, InvalidEnumValueError{
			Name:  "active (controller)",
			Value: msg.Arguments[0].(int32),
		}
	}

	return &ControllerInput{
		Active:  active,
		Name:    msg.Arguments[1].(string),
		IsLeft:  msg.Arguments[2].(int32) == 1,
		IsTouch: msg.Arguments[3].(int32) == 1,
		IsAxis:  msg.Arguments[4].(int32) == 1,
		Axis: Vec3{
			X: msg.Arguments[5].(float32),
			Y: msg.Arguments[6].(float32),
			Z: msg.Arguments[7].(float32),
		},
	}, nil
}

type KeyboardInput struct {
	Active  bool
	Name    string
	KeyCode int32
}

func (k *KeyboardInput) isMessage() {}

func (k *KeyboardInput) Address() string {
	return AddressKeyboardInput
}

func parseKeyboardInput(msg *osc.Message) (*KeyboardInput, error) {
	if msg.TypeTags != "isi" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "isi"}
	}

	return &KeyboardInput{
		Active:  msg.Arguments[0].(int32) == 1,
		Name:    msg.Arguments[1].(string),
		KeyCode: msg.Arguments[2].(int32),
	}, nil
}

type MidiNodeInput struct {
	Active   bool
	Channel  int32
	Note     int32
	Velocity float32
}

func (m *MidiNodeInput) isMessage() {}

func (m *MidiNodeInput) Address() string {
	return AddressMidiNoteInput
}

func parseMidiNoteInput(msg *osc.Message) (*MidiNodeInput, error) {
	if msg.TypeTags != "iiif" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "iiif"}
	}

	return &MidiNodeInput{
		Active:   msg.Arguments[0].(int32) == 1,
		Channel:  msg.Arguments[1].(int32),
		Note:     msg.Arguments[2].(int32),
		Velocity: msg.Arguments[3].(float32),
	}, nil
}

type MidiCCValueInput struct {
	Knob  int32
	Value float32
}

func (m *MidiCCValueInput) isMessage() {}

func (m *MidiCCValueInput) Address() string {
	return AddressMidiCCValueInput
}

func parseMidiCCValueInput(msg *osc.Message) (*MidiCCValueInput, error) {
	if msg.TypeTags != "if" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "if"}
	}

	return &MidiCCValueInput{
		Knob:  msg.Arguments[0].(int32),
		Value: msg.Arguments[1].(float32),
	}, nil
}

type MidiCCButtonInput struct {
	Knob   int32
	Active bool
}

func (m *MidiCCButtonInput) isMessage() {}

func (m *MidiCCButtonInput) Address() string {
	return AddressMidiCCButtonInput
}

func parseMidiCCButtonInput(msg *osc.Message) (*MidiCCButtonInput, error) {
	if msg.TypeTags != "ii" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "ii"}
	}

	return &MidiCCButtonInput{
		Knob:   msg.Arguments[0].(int32),
		Active: msg.Arguments[1].(int32) == 1,
	}, nil
}

type DeviceTransform struct {
	Device Device
	Local  bool // Local means device's raw scale, otherwise avatar scale.

	Serial     string
	Position   Vec3
	Quaternion Vec4
}

func (d *DeviceTransform) isMessage() {}

func (d *DeviceTransform) Address() string {
	if d.Local {
		switch d.Device {
		case DeviceHmd:
			return AddressDeviceTransformHmdLocal
		case DeviceCon:
			return AddressDeviceTransformConLocal
		case DeviceTra:
			return AddressDeviceTransformTraLocal
		}
	} else {
		switch d.Device {
		case DeviceHmd:
			return AddressDeviceTransformHmd
		case DeviceCon:
			return AddressDeviceTransformCon
		case DeviceTra:
			return AddressDeviceTransformTra
		}
	}

	return "" // Should never happen on valid instances.
}

type Device string

// Possible values for the device.
const (
	DeviceHmd Device = "Hmd"
	DeviceCon Device = "Con"
	DeviceTra Device = "Tra"
)

type InvalidDeviceError struct {
	Address string
}

func (e InvalidDeviceError) Error() string {
	return fmt.Sprintf("invalid device address `%s`", e.Address)
}

func parseDeviceTransform(msg *osc.Message) (*DeviceTransform, error) {
	if msg.TypeTags != "sfffffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sfffffff"}
	}

	var device Device
	switch msg.Address {
	case AddressDeviceTransformHmd, AddressDeviceTransformHmdLocal:
		device = DeviceHmd
	case AddressDeviceTransformCon, AddressDeviceTransformConLocal:
		device = DeviceCon
	case AddressDeviceTransformTra, AddressDeviceTransformTraLocal:
		device = DeviceTra
	default:
		return nil, InvalidDeviceError{Address: msg.Address}
	}

	return &DeviceTransform{
		Device: device,
		Local: msg.Address == AddressDeviceTransformHmdLocal ||
			msg.Address == AddressDeviceTransformConLocal ||
			msg.Address == AddressDeviceTransformTraLocal,
		Serial: msg.Arguments[0].(string),
		Position: Vec3{
			X: msg.Arguments[1].(float32),
			Y: msg.Arguments[2].(float32),
			Z: msg.Arguments[3].(float32),
		},
		Quaternion: Vec4{
			X: msg.Arguments[4].(float32),
			Y: msg.Arguments[5].(float32),
			Z: msg.Arguments[6].(float32),
			W: msg.Arguments[7].(float32),
		},
	}, nil
}

type ReceiveEnable struct {
	Enable bool
	Port   int32
}

func (r *ReceiveEnable) isMessage() {}

func (r *ReceiveEnable) Address() string {
	return AddressReceiveEnable
}

func parseReceiveEnable(msg *osc.Message) (*ReceiveEnable, error) {
	if msg.TypeTags != "ii" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "ii"}
	}

	return &ReceiveEnable{
		Enable: msg.Arguments[0].(int32) == 1,
		Port:   msg.Arguments[1].(int32),
	}, nil
}

type DirectionalLight struct {
	Name       string
	Position   Vec3
	Quaternion Vec4
	Color      Vec4
}

func (d *DirectionalLight) isMessage() {}

func (d *DirectionalLight) Address() string {
	return AddressDirectionalLight
}

func parseDirectionalLight(msg *osc.Message) (*DirectionalLight, error) {
	if msg.TypeTags != "sfffffffffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "sfffffffffff"}
	}

	return &DirectionalLight{
		Name: msg.Arguments[0].(string),
		Position: Vec3{
			X: msg.Arguments[1].(float32),
			Y: msg.Arguments[2].(float32),
			Z: msg.Arguments[3].(float32),
		},
		Quaternion: Vec4{
			X: msg.Arguments[4].(float32),
			Y: msg.Arguments[5].(float32),
			Z: msg.Arguments[6].(float32),
			W: msg.Arguments[7].(float32),
		},
		Color: Vec4{
			X: msg.Arguments[8].(float32),
			Y: msg.Arguments[9].(float32),
			Z: msg.Arguments[10].(float32),
			W: msg.Arguments[11].(float32),
		},
	}, nil
}

type LocalVrm struct {
	Path  string
	Title string
}

func (l *LocalVrm) isMessage() {}

func (l *LocalVrm) Address() string {
	return AddressLocalVrm
}

func parseLocalVrm(msg *osc.Message) (*LocalVrm, error) {
	if msg.TypeTags != "ss" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "ss"}
	}

	return &LocalVrm{
		Path:  msg.Arguments[0].(string),
		Title: msg.Arguments[1].(string),
	}, nil
}

type RemoteVrm struct {
	Service string
	JSON    string
}

func (r *RemoteVrm) isMessage() {}

func (r *RemoteVrm) Address() string {
	return AddressRemoteVrm
}

func parseRemoteVrm(msg *osc.Message) (*RemoteVrm, error) {
	if msg.TypeTags != "ss" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "ss"}
	}

	return &RemoteVrm{
		Service: msg.Arguments[0].(string),
		JSON:    msg.Arguments[1].(string),
	}, nil
}

type OptionString struct {
	Option string
}

func (o *OptionString) isMessage() {}

func (o *OptionString) Address() string {
	return AddressOptionString
}

func parseOptionString(msg *osc.Message) (*OptionString, error) {
	if msg.TypeTags != "s" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "s"}
	}

	return &OptionString{
		Option: msg.Arguments[0].(string),
	}, nil
}

type BackgroundColor struct {
	Color Vec4
}

func (b *BackgroundColor) isMessage() {}

func (b *BackgroundColor) Address() string {
	return AddressBackgroundColor
}

func parseBackgroundColor(msg *osc.Message) (*BackgroundColor, error) {
	if msg.TypeTags != "ffff" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "ffff"}
	}

	return &BackgroundColor{
		Color: Vec4{
			X: msg.Arguments[0].(float32),
			Y: msg.Arguments[1].(float32),
			Z: msg.Arguments[2].(float32),
			W: msg.Arguments[3].(float32),
		},
	}, nil
}

type WindowAttribute struct {
	IsTopMost          bool
	IsTransparent      bool
	WindowClickThrough bool
	HideBorder         bool
}

func (w *WindowAttribute) isMessage() {}

func (w *WindowAttribute) Address() string {
	return AddressWindowAttribute
}

func parseWindowAttribute(msg *osc.Message) (*WindowAttribute, error) {
	if msg.TypeTags != "iiii" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "iiii"}
	}

	return &WindowAttribute{
		IsTopMost:          msg.Arguments[0].(int32) == 1,
		IsTransparent:      msg.Arguments[1].(int32) == 1,
		WindowClickThrough: msg.Arguments[2].(int32) == 1,
		HideBorder:         msg.Arguments[3].(int32) == 1,
	}, nil
}

type LoadedSettingPath struct {
	Path string
}

func (l *LoadedSettingPath) isMessage() {}

func (l *LoadedSettingPath) Address() string {
	return AddressLoadedSettingPath
}

func parseLoadedSettingPath(msg *osc.Message) (*LoadedSettingPath, error) {
	if msg.TypeTags != "s" {
		return nil, InvalidTypeTagsError{Found: msg.TypeTags, Expected: "s"}
	}

	return &LoadedSettingPath{
		Path: msg.Arguments[0].(string),
	}, nil
}
