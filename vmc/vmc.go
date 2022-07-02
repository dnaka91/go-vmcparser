// Package vmc implements parsing of "Virtual Motion Capture" messages.
package vmc

import (
	"errors"
	"fmt"
)

// ErrUnknownAddress can happen during ParseMessage, if the message address describes either an
// unknown/unsupported VMC message, or describes a non-VMC message.
var ErrUnknownAddress = errors.New("unknown address")

// ErrFiltered happens when a list of address filters is passed to ParseMessage, and the given raw
// messsage didn't match any of the addresses. That means, it is not an actual error, but instead
// signals, that (any further) parsing of the content was cancelled due to the user-provided filter.
var ErrFiltered = errors.New("address filtered by user")

// InvalidTypeTagsError happens during VMC message parsing, if the message doesn't contain the
// expected arguments (and therefore defined type tags) for that particular message.
type InvalidTypeTagsError struct {
	Found    []byte
	Expected []string
}

var _ error = (*InvalidTypeTagsError)(nil)

func (e InvalidTypeTagsError) Error() string {
	switch len(e.Expected) {
	case 0:
		return fmt.Sprintf("invalid type tags `%v`, expected none", string(e.Found))
	case 1:
		return fmt.Sprintf("invalid type tags `%v`, expected `%v`", string(e.Found), e.Expected[0])
	default:
		return fmt.Sprintf("invalid type tags `%v`, expected one of %v", string(e.Found), e.Expected)
	}
}

// InvalidEnumValueError if an enum value of a VMC message is none of the known possible values.
type InvalidEnumValueError struct {
	Name  string
	Value int32
}

var _ error = (*InvalidEnumValueError)(nil)

func (e InvalidEnumValueError) Error() string {
	return fmt.Sprintf("invalid value for %v: %v", e.Name, e.Value)
}

type InvalidBufferLengthError struct {
	Length   int
	Expected []int
}

var _ error = (*InvalidBufferLengthError)(nil)

func (e InvalidBufferLengthError) Error() string {
	switch len(e.Expected) {
	case 0:
		return fmt.Sprintf("got %d bytes of data, expected none", e.Length)
	case 1:
		return fmt.Sprintf("got %d bytes of data, expected exactly %d", e.Length, e.Expected[0])
	default:
		return fmt.Sprintf("got %d bytes of data, expected one of %v", e.Length, e.Expected)
	}
}

// Message is a marker for any type that is considered a VMC message.
type Message interface {
	isMessage()
}

// ParseMessage takes a generic OSC message, and tries to parse it into one of the known VMC
// messages.
//
// The list of supported messages is not complete yet, but most of the "marionette" messages are
// implemented.
//
// Parsing of a message can be limited with optional address filters. If passed, and the message's
// address didn't match any of the filters, then all further processing is stopped and a ErrFiltered
// error is returned.
func ParseMessage(data []byte, addressFilters ...string) (Message, error) {
	address, newData, err := getString(data)
	if err != nil {
		return nil, err
	}
	data = newData

	if !filterAddress(address, addressFilters) {
		return nil, ErrFiltered
	}

	tags, newData, err := getTypeTags(data)
	if err != nil {
		return nil, err
	}
	data = newData

	switch string(address) {
	case AddressAvailable:
		return parseAvailable(tags, data)
	case AddressRelativeTime:
		return parseRelativeTime(tags, data)
	case AddressRootTransform:
		return parseRootTransform(tags, data)
	case AddressBoneTransform:
		return parseBoneTransform(tags, data)
	case AddressBlendShapeProxyValue:
		return parseBlendShapeProxyValue(tags, data)
	case AddressBlendShapeProxyApply:
		return parseBlendShapeProxyApply(tags, data)
	case AddressCameraTransform:
		return parseCameraTransform(tags, data)
	case AddressControllerInput:
		return parseControllerInput(tags, data)
	case AddressKeyboardInput:
		return parseKeyboardInput(tags, data)
	case AddressMidiNoteInput:
		return parseMidiNoteInput(tags, data)
	case AddressMidiCCValueInput:
		return parseMidiCCValueInput(tags, data)
	case AddressMidiCCButtonInput:
		return parseMidiCCButtonInput(tags, data)
	case AddressDeviceTransformHmd,
		AddressDeviceTransformCon,
		AddressDeviceTransformTra,
		AddressDeviceTransformHmdLocal,
		AddressDeviceTransformConLocal,
		AddressDeviceTransformTraLocal:
		return parseDeviceTransform(tags, data)
	case AddressReceiveEnable:
		return parseReceiveEnable(tags, data)
	case AddressDirectionalLight:
		return parseDirectionalLight(tags, data)
	case AddressLocalVrm:
		return parseLocalVrm(tags, data)
	case AddressRemoteVrm:
		return parseRemoteVrm(tags, data)
	case AddressOptionString:
		return parseOptionString(tags, data)
	case AddressBackgroundColor:
		return parseBackgroundColor(tags, data)
	case AddressWindowAttribute:
		return parseWindowAttribute(tags, data)
	case AddressLoadedSettingPath:
		return parseLoadedSettingPath(tags, data)
	default:
		return nil, ErrUnknownAddress
	}
}

func filterAddress(address []byte, filters []string) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		if string(address) == filter {
			return true
		}
	}

	return false
}
