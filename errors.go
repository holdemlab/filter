package filter

import (
	"errors"
	"fmt"
)

// Sentinel errors for common nil / missing checks.
var (
	// ErrNilOptions is returned when filter Options are nil.
	ErrNilOptions = errors.New("filter options are nil")

	// ErrNilDataset is returned when the query dataset is nil.
	ErrNilDataset = errors.New("dataset is nil")

	// ErrEmptyTypeTag is returned when a struct field has no "type" tag.
	ErrEmptyTypeTag = errors.New("empty data type in struct annotation")

	// ErrNotStruct is returned when a non-struct value is passed to UpdateFromQueryParams.
	ErrNotStruct = errors.New("value must be a struct")
)

// ValidationError is returned when an operator, data type, or their combination is invalid.
type ValidationError struct {
	// Kind identifies what was being validated: "operator", "data_type", "operator_data_type".
	Kind string
	// Values contains the invalid values that caused the error.
	Values []string
}

func (e *ValidationError) Error() string {
	switch e.Kind {
	case "operator":
		return fmt.Sprintf("undefined operator [%s]", e.Values[0])
	case "data_type":
		return fmt.Sprintf("undefined data type [%s]", e.Values[0])
	case "operator_data_type":
		return fmt.Sprintf("operator %q is not compatible with data type %q", e.Values[0], e.Values[1])
	default:
		return fmt.Sprintf("validation error: %v", e.Values)
	}
}

// ConversionError is returned when a field value cannot be converted to its target data type.
type ConversionError struct {
	// DataType is the target data type (e.g. "int", "bool", "uuid").
	DataType string
	// Value is the original string value that failed to convert.
	Value string
	// Err is the underlying conversion error.
	Err error
}

func (e *ConversionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cannot convert %q to %s: %v", e.Value, e.DataType, e.Err)
	}
	return fmt.Sprintf("cannot convert %q to %s", e.Value, e.DataType)
}

func (e *ConversionError) Unwrap() error {
	return e.Err
}

// ParseError is returned when parsing a structured value (e.g. a between range) fails.
type ParseError struct {
	// Value is the raw input that failed to parse.
	Value string
	// Message provides a human-readable description.
	Message string
	// Err is the underlying parse error, if any.
	Err error
}

func (e *ParseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
