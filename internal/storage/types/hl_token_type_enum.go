// SPDX-FileCopyrightText: 2025 Bb Strategy Pte. Ltd. <celenium@baking-bad.org>
// SPDX-License-Identifier: MIT

// Code generated by go-enum DO NOT EDIT.
// Version: 0.5.7
// Revision: bf63e108589bbd2327b13ec2c5da532aad234029
// Build Date: 2023-07-25T23:27:55Z
// Built By: goreleaser

package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

const (
	// HLTokenTypeSynthetic is a HLTokenType of type synthetic.
	HLTokenTypeSynthetic HLTokenType = "synthetic"
	// HLTokenTypeCollateral is a HLTokenType of type collateral.
	HLTokenTypeCollateral HLTokenType = "collateral"
)

var ErrInvalidHLTokenType = fmt.Errorf("not a valid HLTokenType, try [%s]", strings.Join(_HLTokenTypeNames, ", "))

var _HLTokenTypeNames = []string{
	string(HLTokenTypeSynthetic),
	string(HLTokenTypeCollateral),
}

// HLTokenTypeNames returns a list of possible string values of HLTokenType.
func HLTokenTypeNames() []string {
	tmp := make([]string, len(_HLTokenTypeNames))
	copy(tmp, _HLTokenTypeNames)
	return tmp
}

// HLTokenTypeValues returns a list of the values for HLTokenType
func HLTokenTypeValues() []HLTokenType {
	return []HLTokenType{
		HLTokenTypeSynthetic,
		HLTokenTypeCollateral,
	}
}

// String implements the Stringer interface.
func (x HLTokenType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x HLTokenType) IsValid() bool {
	_, err := ParseHLTokenType(string(x))
	return err == nil
}

var _HLTokenTypeValue = map[string]HLTokenType{
	"synthetic":  HLTokenTypeSynthetic,
	"collateral": HLTokenTypeCollateral,
}

// ParseHLTokenType attempts to convert a string to a HLTokenType.
func ParseHLTokenType(name string) (HLTokenType, error) {
	if x, ok := _HLTokenTypeValue[name]; ok {
		return x, nil
	}
	return HLTokenType(""), fmt.Errorf("%s is %w", name, ErrInvalidHLTokenType)
}

// MarshalText implements the text marshaller method.
func (x HLTokenType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *HLTokenType) UnmarshalText(text []byte) error {
	tmp, err := ParseHLTokenType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errHLTokenTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *HLTokenType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = HLTokenType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseHLTokenType(v)
	case []byte:
		*x, err = ParseHLTokenType(string(v))
	case HLTokenType:
		*x = v
	case *HLTokenType:
		if v == nil {
			return errHLTokenTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errHLTokenTypeNilPtr
		}
		*x, err = ParseHLTokenType(*v)
	default:
		return errors.New("invalid type for HLTokenType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x HLTokenType) Value() (driver.Value, error) {
	return x.String(), nil
}

const (
	// HLTransferTypeSend is a HLTransferType of type send.
	HLTransferTypeSend HLTransferType = "send"
	// HLTransferTypeReceive is a HLTransferType of type receive.
	HLTransferTypeReceive HLTransferType = "receive"
)

var ErrInvalidHLTransferType = fmt.Errorf("not a valid HLTransferType, try [%s]", strings.Join(_HLTransferTypeNames, ", "))

var _HLTransferTypeNames = []string{
	string(HLTransferTypeSend),
	string(HLTransferTypeReceive),
}

// HLTransferTypeNames returns a list of possible string values of HLTransferType.
func HLTransferTypeNames() []string {
	tmp := make([]string, len(_HLTransferTypeNames))
	copy(tmp, _HLTransferTypeNames)
	return tmp
}

// HLTransferTypeValues returns a list of the values for HLTransferType
func HLTransferTypeValues() []HLTransferType {
	return []HLTransferType{
		HLTransferTypeSend,
		HLTransferTypeReceive,
	}
}

// String implements the Stringer interface.
func (x HLTransferType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x HLTransferType) IsValid() bool {
	_, err := ParseHLTransferType(string(x))
	return err == nil
}

var _HLTransferTypeValue = map[string]HLTransferType{
	"send":    HLTransferTypeSend,
	"receive": HLTransferTypeReceive,
}

// ParseHLTransferType attempts to convert a string to a HLTransferType.
func ParseHLTransferType(name string) (HLTransferType, error) {
	if x, ok := _HLTransferTypeValue[name]; ok {
		return x, nil
	}
	return HLTransferType(""), fmt.Errorf("%s is %w", name, ErrInvalidHLTransferType)
}

// MarshalText implements the text marshaller method.
func (x HLTransferType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *HLTransferType) UnmarshalText(text []byte) error {
	tmp, err := ParseHLTransferType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errHLTransferTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *HLTransferType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = HLTransferType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseHLTransferType(v)
	case []byte:
		*x, err = ParseHLTransferType(string(v))
	case HLTransferType:
		*x = v
	case *HLTransferType:
		if v == nil {
			return errHLTransferTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errHLTransferTypeNilPtr
		}
		*x, err = ParseHLTransferType(*v)
	default:
		return errors.New("invalid type for HLTransferType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x HLTransferType) Value() (driver.Value, error) {
	return x.String(), nil
}
