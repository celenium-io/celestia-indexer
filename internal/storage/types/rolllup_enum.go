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
	// RollupCategoryUncategorized is a RollupCategory of type uncategorized.
	RollupCategoryUncategorized RollupCategory = "uncategorized"
	// RollupCategoryFinance is a RollupCategory of type finance.
	RollupCategoryFinance RollupCategory = "finance"
	// RollupCategoryGaming is a RollupCategory of type gaming.
	RollupCategoryGaming RollupCategory = "gaming"
	// RollupCategoryNft is a RollupCategory of type nft.
	RollupCategoryNft RollupCategory = "nft"
)

var ErrInvalidRollupCategory = fmt.Errorf("not a valid RollupCategory, try [%s]", strings.Join(_RollupCategoryNames, ", "))

var _RollupCategoryNames = []string{
	string(RollupCategoryUncategorized),
	string(RollupCategoryFinance),
	string(RollupCategoryGaming),
	string(RollupCategoryNft),
}

// RollupCategoryNames returns a list of possible string values of RollupCategory.
func RollupCategoryNames() []string {
	tmp := make([]string, len(_RollupCategoryNames))
	copy(tmp, _RollupCategoryNames)
	return tmp
}

// RollupCategoryValues returns a list of the values for RollupCategory
func RollupCategoryValues() []RollupCategory {
	return []RollupCategory{
		RollupCategoryUncategorized,
		RollupCategoryFinance,
		RollupCategoryGaming,
		RollupCategoryNft,
	}
}

// String implements the Stringer interface.
func (x RollupCategory) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RollupCategory) IsValid() bool {
	_, err := ParseRollupCategory(string(x))
	return err == nil
}

var _RollupCategoryValue = map[string]RollupCategory{
	"uncategorized": RollupCategoryUncategorized,
	"finance":       RollupCategoryFinance,
	"gaming":        RollupCategoryGaming,
	"nft":           RollupCategoryNft,
}

// ParseRollupCategory attempts to convert a string to a RollupCategory.
func ParseRollupCategory(name string) (RollupCategory, error) {
	if x, ok := _RollupCategoryValue[name]; ok {
		return x, nil
	}
	return RollupCategory(""), fmt.Errorf("%s is %w", name, ErrInvalidRollupCategory)
}

// MarshalText implements the text marshaller method.
func (x RollupCategory) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *RollupCategory) UnmarshalText(text []byte) error {
	tmp, err := ParseRollupCategory(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errRollupCategoryNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *RollupCategory) Scan(value interface{}) (err error) {
	if value == nil {
		*x = RollupCategory("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseRollupCategory(v)
	case []byte:
		*x, err = ParseRollupCategory(string(v))
	case RollupCategory:
		*x = v
	case *RollupCategory:
		if v == nil {
			return errRollupCategoryNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errRollupCategoryNilPtr
		}
		*x, err = ParseRollupCategory(*v)
	default:
		return errors.New("invalid type for RollupCategory")
	}

	return
}

// Value implements the driver Valuer interface.
func (x RollupCategory) Value() (driver.Value, error) {
	return x.String(), nil
}

const (
	// RollupTypeSovereign is a RollupType of type sovereign.
	RollupTypeSovereign RollupType = "sovereign"
	// RollupTypeSettled is a RollupType of type settled.
	RollupTypeSettled RollupType = "settled"
)

var ErrInvalidRollupType = fmt.Errorf("not a valid RollupType, try [%s]", strings.Join(_RollupTypeNames, ", "))

var _RollupTypeNames = []string{
	string(RollupTypeSovereign),
	string(RollupTypeSettled),
}

// RollupTypeNames returns a list of possible string values of RollupType.
func RollupTypeNames() []string {
	tmp := make([]string, len(_RollupTypeNames))
	copy(tmp, _RollupTypeNames)
	return tmp
}

// RollupTypeValues returns a list of the values for RollupType
func RollupTypeValues() []RollupType {
	return []RollupType{
		RollupTypeSovereign,
		RollupTypeSettled,
	}
}

// String implements the Stringer interface.
func (x RollupType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RollupType) IsValid() bool {
	_, err := ParseRollupType(string(x))
	return err == nil
}

var _RollupTypeValue = map[string]RollupType{
	"sovereign": RollupTypeSovereign,
	"settled":   RollupTypeSettled,
}

// ParseRollupType attempts to convert a string to a RollupType.
func ParseRollupType(name string) (RollupType, error) {
	if x, ok := _RollupTypeValue[name]; ok {
		return x, nil
	}
	return RollupType(""), fmt.Errorf("%s is %w", name, ErrInvalidRollupType)
}

// MarshalText implements the text marshaller method.
func (x RollupType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *RollupType) UnmarshalText(text []byte) error {
	tmp, err := ParseRollupType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

var errRollupTypeNilPtr = errors.New("value pointer is nil") // one per type for package clashes

// Scan implements the Scanner interface.
func (x *RollupType) Scan(value interface{}) (err error) {
	if value == nil {
		*x = RollupType("")
		return
	}

	// A wider range of scannable types.
	// driver.Value values at the top of the list for expediency
	switch v := value.(type) {
	case string:
		*x, err = ParseRollupType(v)
	case []byte:
		*x, err = ParseRollupType(string(v))
	case RollupType:
		*x = v
	case *RollupType:
		if v == nil {
			return errRollupTypeNilPtr
		}
		*x = *v
	case *string:
		if v == nil {
			return errRollupTypeNilPtr
		}
		*x, err = ParseRollupType(*v)
	default:
		return errors.New("invalid type for RollupType")
	}

	return
}

// Value implements the driver Valuer interface.
func (x RollupType) Value() (driver.Value, error) {
	return x.String(), nil
}
