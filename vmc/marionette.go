package vmc

import (
	"fmt"
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
	CalibrationState *CalibrationState
	CalibrationMode  *CalibrationMode
	TrackingStatus   *bool
}

func (a *Available) isMessage() {}

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

func parseAvailable(tags, data []byte) (*Available, error) {
	const (
		typeTagsV1   = "i"
		typeTagsV2_5 = "iii"
		typeTagsV2_7 = "iiii"
	)

	if string(tags) != typeTagsV1 &&
		string(tags) != typeTagsV2_5 &&
		string(tags) != typeTagsV2_7 {
		return nil, InvalidTypeTagsError{
			Found:    tags,
			Expected: []string{typeTagsV1, typeTagsV2_5, typeTagsV2_7},
		}
	}

	if len(data) != 4 && len(data) != 12 && len(data) != 16 {
		return nil, InvalidBufferLengthError{
			Length:   len(data),
			Expected: []int{4, 8, 12},
		}
	}

	value := &Available{
		Loaded:           getInt32(data[0:4]) == 1,
		CalibrationState: nil,
		CalibrationMode:  nil,
		TrackingStatus:   nil,
	}

	if len(data) == 12 || len(data) == 16 {
		rawValue := getInt32(data[4:8])
		calibrationState := CalibrationState(rawValue)
		if !calibrationState.isValid() {
			return nil, InvalidEnumValueError{
				Name:  "calibration state",
				Value: rawValue,
			}
		}

		rawValue = getInt32(data[8:12])
		calibrationMode := CalibrationMode(rawValue)
		if !calibrationMode.isValid() {
			return nil, InvalidEnumValueError{
				Name:  "calibration mode",
				Value: rawValue,
			}
		}

		value.CalibrationState = &calibrationState
		value.CalibrationMode = &calibrationMode
	}

	if len(data) == 16 {
		trackingStatus := getInt32(data[12:16]) == 1
		value.TrackingStatus = &trackingStatus
	}

	return value, nil
}

type RelativeTime struct {
	Time float32
}

func (r *RelativeTime) isMessage() {}

func parseRelativeTime(tags, data []byte) (*RelativeTime, error) {
	if string(tags) != "f" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"f"}}
	}

	if len(data) != 4 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{4}}
	}

	return &RelativeTime{
		Time: getFloat32(data[0:4]),
	}, nil
}

type RootTransform struct {
	Name       []byte
	Position   Vec3
	Quaternion Vec4
	Scale      *Vec3
	Offset     *Vec3
}

func (r *RootTransform) isMessage() {}

func parseRootTransform(tags, data []byte) (*RootTransform, error) {
	const (
		typeTagsV2_0 = "sfffffff"
		typeTagsV2_1 = "sfffffffffffff"
	)
	if string(tags) != typeTagsV2_0 &&
		string(tags) != typeTagsV2_1 {
		return nil, InvalidTypeTagsError{
			Found:    tags,
			Expected: []string{typeTagsV2_0, typeTagsV2_1},
		}
	}

	name, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 28 && len(data) != 52 {
		return nil, InvalidBufferLengthError{
			Length:   len(data),
			Expected: []int{28, 52},
		}
	}

	value := &RootTransform{
		Name:       name,
		Position:   getVec3(data[0:12]),
		Quaternion: getVec4(data[12:28]),
		Scale:      nil,
		Offset:     nil,
	}

	if len(data) == 52 {
		scale := getVec3(data[28:40])
		offset := getVec3(data[40:52])
		value.Scale = &scale
		value.Offset = &offset
	}

	return value, nil
}

type BoneTransform struct {
	Name       []byte
	Position   Vec3
	Quaternion Vec4
}

func (b *BoneTransform) isMessage() {}

func parseBoneTransform(tags, data []byte) (*BoneTransform, error) {
	if string(tags) != "sfffffff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"sfffffff"}}
	}

	name, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 28 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{28}}
	}

	return &BoneTransform{
		Name:       name,
		Position:   getVec3(data[0:12]),
		Quaternion: getVec4(data[12:28]),
	}, nil
}

type BlendShapeProxyValue struct {
	Name  []byte
	Value float32
}

func (b *BlendShapeProxyValue) isMessage() {}

func parseBlendShapeProxyValue(tags, data []byte) (*BlendShapeProxyValue, error) {
	if string(tags) != "sf" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"sf"}}
	}

	name, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 4 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{4}}
	}

	return &BlendShapeProxyValue{
		Name:  name,
		Value: getFloat32(data[0:4]),
	}, nil
}

type BlendShapeProxyApply struct{}

func (b *BlendShapeProxyApply) isMessage() {}

func parseBlendShapeProxyApply(tags, data []byte) (*BlendShapeProxyApply, error) {
	if string(tags) != "" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: nil}
	}

	if len(data) != 0 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: nil}
	}

	return &BlendShapeProxyApply{}, nil
}

type CameraTransform struct {
	Name       []byte
	Position   Vec3
	Quaternion Vec4
	FOV        float32
}

func (c *CameraTransform) isMessage() {}

func parseCameraTransform(tags, data []byte) (*CameraTransform, error) {
	if string(tags) != "sffffffff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"sffffffff"}}
	}

	name, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 32 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{32}}
	}

	return &CameraTransform{
		Name:       name,
		Position:   getVec3(data[0:12]),
		Quaternion: getVec4(data[12:28]),
		FOV:        getFloat32(data[28:32]),
	}, nil
}

type ControllerInput struct {
	Active  ControllerActive
	Name    []byte
	IsLeft  bool
	IsTouch bool
	IsAxis  bool
	Axis    Vec3
}

func (c *ControllerInput) isMessage() {}

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

func parseControllerInput(tags, data []byte) (*ControllerInput, error) {
	if string(tags) != "isiiifff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"isiiifff"}}
	}

	if len(data) <= 4 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{4}}
	}

	rawValue := getInt32(data[0:4])
	active := ControllerActive(rawValue)
	if !active.isValid() {
		return nil, InvalidEnumValueError{
			Name:  "active (controller)",
			Value: rawValue,
		}
	}

	name, newData, err := getString(data[4:])
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 24 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{24}}
	}

	return &ControllerInput{
		Active:  active,
		Name:    name,
		IsLeft:  getInt32(data[0:4]) == 1,
		IsTouch: getInt32(data[4:8]) == 1,
		IsAxis:  getInt32(data[8:12]) == 1,
		Axis:    getVec3(data[12:24]),
	}, nil
}

type KeyboardInput struct {
	Active  bool
	Name    []byte
	KeyCode int32
}

func (k *KeyboardInput) isMessage() {}

func parseKeyboardInput(tags, data []byte) (*KeyboardInput, error) {
	if string(tags) != "isi" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"isi"}}
	}

	if len(data) <= 4 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{4}}
	}

	active := getInt32(data[0:4]) == 1
	name, newData, err := getString(data[4:])
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 4 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{4}}
	}

	return &KeyboardInput{
		Active:  active,
		Name:    name,
		KeyCode: getInt32(data[0:4]),
	}, nil
}

type MidiNoteInput struct {
	Active   bool
	Channel  int32
	Note     int32
	Velocity float32
}

func (m *MidiNoteInput) isMessage() {}

func parseMidiNoteInput(tags, data []byte) (*MidiNoteInput, error) {
	if string(tags) != "iiif" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"iiif"}}
	}

	if len(data) != 16 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{16}}
	}

	return &MidiNoteInput{
		Active:   getInt32(data[0:4]) == 1,
		Channel:  getInt32(data[4:8]),
		Note:     getInt32(data[8:12]),
		Velocity: getFloat32(data[12:16]),
	}, nil
}

type MidiCCValueInput struct {
	Knob  int32
	Value float32
}

func (m *MidiCCValueInput) isMessage() {}

func parseMidiCCValueInput(tags, data []byte) (*MidiCCValueInput, error) {
	if string(tags) != "if" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"if"}}
	}

	if len(data) != 8 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{8}}
	}

	return &MidiCCValueInput{
		Knob:  getInt32(data[0:4]),
		Value: getFloat32(data[4:8]),
	}, nil
}

type MidiCCButtonInput struct {
	Knob   int32
	Active bool
}

func (m *MidiCCButtonInput) isMessage() {}

func parseMidiCCButtonInput(tags, data []byte) (*MidiCCButtonInput, error) {
	if string(tags) != "ii" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"ii"}}
	}

	if len(data) != 8 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{8}}
	}

	return &MidiCCButtonInput{
		Knob:   getInt32(data[0:4]),
		Active: getInt32(data[4:8]) == 1,
	}, nil
}

type DeviceTransform struct {
	Serial     []byte
	Position   Vec3
	Quaternion Vec4
}

func (d *DeviceTransform) isMessage() {}

func parseDeviceTransform(tags, data []byte) (*DeviceTransform, error) {
	if string(tags) != "sfffffff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"sfffffff"}}
	}

	serial, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 28 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{28}}
	}

	return &DeviceTransform{
		Serial:     serial,
		Position:   getVec3(data[0:12]),
		Quaternion: getVec4(data[12:28]),
	}, nil
}

type ReceiveEnable struct {
	Enable    bool
	Port      int32
	IPAddress *[]byte
}

func (r *ReceiveEnable) isMessage() {}

func parseReceiveEnable(tags, data []byte) (*ReceiveEnable, error) {
	const (
		typeTagsV2_4 = "ii"
		typeTagsV2_7 = "iis"
	)

	if string(tags) != typeTagsV2_4 &&
		string(tags) != typeTagsV2_7 {
		return nil, InvalidTypeTagsError{
			Found:    tags,
			Expected: []string{typeTagsV2_4, typeTagsV2_7},
		}
	}

	if len(data) < 8 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{8}}
	}

	value := &ReceiveEnable{
		Enable:    getInt32(data[0:4]) == 1,
		Port:      getInt32(data[4:8]),
		IPAddress: nil,
	}

	if len(data) > 8 {
		ipAddress, _, err := getString(data[8:])
		if err != nil {
			return nil, err
		}
		value.IPAddress = &ipAddress
	}

	return value, nil
}

type DirectionalLight struct {
	Name       []byte
	Position   Vec3
	Quaternion Vec4
	Color      Vec4
}

func (d *DirectionalLight) isMessage() {}

func parseDirectionalLight(tags, data []byte) (*DirectionalLight, error) {
	if string(tags) != "sfffffffffff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"sfffffffffff"}}
	}

	name, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if len(data) != 44 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{44}}
	}

	return &DirectionalLight{
		Name:       name,
		Position:   getVec3(data[0:12]),
		Quaternion: getVec4(data[12:28]),
		Color:      getVec4(data[28:44]),
	}, nil
}

type LocalVrm struct {
	Path  []byte
	Title []byte
	Hash  *[]byte
}

func (l *LocalVrm) isMessage() {}

func parseLocalVrm(tags, data []byte) (*LocalVrm, error) {
	const (
		typeTagsV2_4 = "ss"
		typeTagsV2_7 = "sss"
	)

	if string(tags) != typeTagsV2_4 &&
		string(tags) != typeTagsV2_7 {
		return nil, InvalidTypeTagsError{
			Found:    tags,
			Expected: []string{typeTagsV2_4, typeTagsV2_7},
		}
	}

	if len(data) == 0 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{1}}
	}

	path, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	title, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	value := &LocalVrm{
		Path:  path,
		Title: title,
		Hash:  nil,
	}

	if len(data) > 0 {
		hash, _, err := getString(data)
		if err != nil {
			return nil, err
		}
		value.Hash = &hash
	}

	return value, nil
}

type RemoteVrm struct {
	Service []byte
	JSON    []byte
}

func (r *RemoteVrm) isMessage() {}

func parseRemoteVrm(tags, data []byte) (*RemoteVrm, error) {
	if string(tags) != "ss" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"ss"}}
	}

	if len(data) == 0 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{1}}
	}

	service, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	json, _, err := getString(data)
	if err != nil {
		return nil, err
	}

	return &RemoteVrm{
		Service: service,
		JSON:    json,
	}, nil
}

type OptionString struct {
	Option []byte
}

func (o *OptionString) isMessage() {}

func parseOptionString(tags, data []byte) (*OptionString, error) {
	if string(tags) != "s" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"s"}}
	}

	if len(data) == 0 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{1}}
	}

	option, _, err := getString(data)
	if err != nil {
		return nil, err
	}

	return &OptionString{
		Option: option,
	}, nil
}

type BackgroundColor struct {
	Color Vec4
}

func (b *BackgroundColor) isMessage() {}

func parseBackgroundColor(tags, data []byte) (*BackgroundColor, error) {
	if string(tags) != "ffff" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"ffff"}}
	}

	if len(data) != 16 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{16}}
	}

	return &BackgroundColor{
		Color: getVec4(data[0:16]),
	}, nil
}

type WindowAttribute struct {
	IsTopMost          bool
	IsTransparent      bool
	WindowClickThrough bool
	HideBorder         bool
}

func (w *WindowAttribute) isMessage() {}

func parseWindowAttribute(tags, data []byte) (*WindowAttribute, error) {
	if string(tags) != "iiii" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"iiii"}}
	}

	if len(data) != 16 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{16}}
	}

	return &WindowAttribute{
		IsTopMost:          getInt32(data[0:4]) == 1,
		IsTransparent:      getInt32(data[4:8]) == 1,
		WindowClickThrough: getInt32(data[8:12]) == 1,
		HideBorder:         getInt32(data[12:16]) == 1,
	}, nil
}

type LoadedSettingPath struct {
	Path []byte
}

func (l *LoadedSettingPath) isMessage() {}

func parseLoadedSettingPath(tags, data []byte) (*LoadedSettingPath, error) {
	if string(tags) != "s" {
		return nil, InvalidTypeTagsError{Found: tags, Expected: []string{"s"}}
	}

	if len(data) == 0 {
		return nil, InvalidBufferLengthError{Length: len(data), Expected: []int{1}}
	}

	path, _, err := getString(data)
	if err != nil {
		return nil, err
	}

	return &LoadedSettingPath{
		Path: path,
	}, nil
}
