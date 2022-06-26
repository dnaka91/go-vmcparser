// Package vmc implements parsing of "Virtual Motion Capture" messages.
package vmc

import (
	"errors"
	"fmt"

	"github.com/dnaka91/go-vmcparser/osc"
)

// ErrUnknownAddress can happen during ParseMessage, if the message address describes either an
// unknown/unsupported VMC message, or describes a non-VMC message.
var ErrUnknownAddress = errors.New("unknown address")

// InvalidTypeTagsError happens during VMC message parsing, if the message doesn't contain the
// expected arguments (and therefore defined type tags) for that particular message.
type InvalidTypeTagsError struct {
	Found    string
	Expected string
}

var _ error = (*InvalidTypeTagsError)(nil)

func (e InvalidTypeTagsError) Error() string {
	if e.Expected == "" {
		return fmt.Sprintf("invalid type tags `%v`, expected none", e.Found)
	}

	return fmt.Sprintf("invalid type tags `%v`, expected `%v`", e.Found, e.Expected)
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

// Message is a marker for any type that is considered a VMC message.
type Message interface {
	isMessage()

	// Address returns the OSC address belonging to the particular VMC message.
	Address() string
}

// ParseMessage takes a generic OSC message, and tries to parse it into one of the known VMC
// messages.
//
// The list of supported messages is not complete yet, but most of the "marionette" messages are
// implemented.
func ParseMessage(msg *osc.Message) (Message, error) {
	switch msg.Address {
	case AddressAvailable:
		return parseAvailable(msg)
	case AddressRelativeTime:
		return parseRelativeTime(msg)
	case AddressRootTransform:
		return parseRootTransform(msg)
	case AddressBoneTransform:
		return parseBoneTransform(msg)
	case AddressBlendShapeProxyValue:
		return parseBlendShapeProxyValue(msg)
	case AddressBlendShapeProxyApply:
		return parseBlendShapeProxyApply(msg)
	case AddressCameraTransform:
		return parseCameraTransform(msg)
	case AddressControllerInput:
		return parseControllerInput(msg)
	case AddressKeyboardInput:
		return parseKeyboardInput(msg)
	case AddressMidiNoteInput:
		return parseMidiNoteInput(msg)
	case AddressMidiCCValueInput:
		return parseMidiCCValueInput(msg)
	case AddressMidiCCButtonInput:
		return parseMidiCCButtonInput(msg)
	case AddressDeviceTransformHmd,
		AddressDeviceTransformCon,
		AddressDeviceTransformTra,
		AddressDeviceTransformHmdLocal,
		AddressDeviceTransformConLocal,
		AddressDeviceTransformTraLocal:
		return parseDeviceTransform(msg)
	case AddressReceiveEnable:
		return parseReceiveEnable(msg)
	case AddressDirectionalLight:
		return parseDirectionalLight(msg)
	case AddressLocalVrm:
		return parseLocalVrm(msg)
	case AddressRemoteVrm:
		return parseRemoteVrm(msg)
	case AddressOptionString:
		return parseOptionString(msg)
	case AddressBackgroundColor:
		return parseBackgroundColor(msg)
	case AddressWindowAttribute:
		return parseWindowAttribute(msg)
	case AddressLoadedSettingPath:
		return parseLoadedSettingPath(msg)
	default:
		return nil, ErrUnknownAddress
	}
}
