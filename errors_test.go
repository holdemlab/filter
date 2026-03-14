package filter

import (
	"errors"
	"testing"
)

// --- ValidationError ---

func TestValidationError_OperatorKind(t *testing.T) {
	err := validateOperator("bad_op")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Kind != "operator" {
		t.Errorf("expected Kind %q, got %q", "operator", ve.Kind)
	}
	if len(ve.Values) != 1 || ve.Values[0] != "bad_op" {
		t.Errorf("expected Values [bad_op], got %v", ve.Values)
	}
}

func TestValidationError_DataTypeKind(t *testing.T) {
	err := validateDataType("float")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Kind != "data_type" {
		t.Errorf("expected Kind %q, got %q", "data_type", ve.Kind)
	}
}

func TestValidationError_OperatorDataTypeKind(t *testing.T) {
	err := validateOperatorDataType(OperatorLike, DataTypeInt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Kind != "operator_data_type" {
		t.Errorf("expected Kind %q, got %q", "operator_data_type", ve.Kind)
	}
	if len(ve.Values) != 2 || ve.Values[0] != OperatorLike || ve.Values[1] != DataTypeInt {
		t.Errorf("unexpected Values: %v", ve.Values)
	}
}

func TestValidationError_UnwrapFromAddField(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.AddField("x", "bad_op", "v", DataTypeStr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected errors.As(*ValidationError) to succeed for wrapped error, got %T: %v", err, err)
	}
	if ve.Kind != "operator" {
		t.Errorf("expected Kind %q, got %q", "operator", ve.Kind)
	}
}

// --- ConversionError ---

func TestConversionError_IntField(t *testing.T) {
	f := Field{Name: "age", Operator: OperatorEq, Value: "abc", Type: DataTypeInt}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *ConversionError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConversionError, got %T", err)
	}
	if ce.DataType != DataTypeInt {
		t.Errorf("expected DataType %q, got %q", DataTypeInt, ce.DataType)
	}
	if ce.Value != "abc" {
		t.Errorf("expected Value %q, got %q", "abc", ce.Value)
	}
	if ce.Unwrap() == nil {
		t.Error("expected non-nil underlying error")
	}
}

func TestConversionError_BoolField(t *testing.T) {
	f := Field{Name: "active", Operator: OperatorEq, Value: "notbool", Type: DataTypeBool}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *ConversionError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConversionError, got %T", err)
	}
	if ce.DataType != DataTypeBool {
		t.Errorf("expected DataType %q, got %q", DataTypeBool, ce.DataType)
	}
}

func TestConversionError_UUIDField(t *testing.T) {
	f := Field{Name: "id", Operator: OperatorEq, Value: "not-a-uuid", Type: DataTypeUUID}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *ConversionError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConversionError, got %T", err)
	}
	if ce.DataType != DataTypeUUID {
		t.Errorf("expected DataType %q, got %q", DataTypeUUID, ce.DataType)
	}
}

func TestConversionError_ListField(t *testing.T) {
	f := Field{Name: "tags", Operator: OperatorEq, Value: "not-json", Type: DataTypeList}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *ConversionError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConversionError, got %T", err)
	}
	if ce.DataType != DataTypeList {
		t.Errorf("expected DataType %q, got %q", DataTypeList, ce.DataType)
	}
}

func TestConversionError_DateTimestampField(t *testing.T) {
	f := Field{Name: "created", Operator: OperatorEq, Value: "notanumber", Type: DataTypeDate}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *ConversionError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConversionError, got %T", err)
	}
	if ce.DataType != DataTypeDate {
		t.Errorf("expected DataType %q, got %q", DataTypeDate, ce.DataType)
	}
}

// --- ParseError ---

func TestParseError_InvalidBetweenFormat(t *testing.T) {
	_, err := ParseBetweenOperator("2006-01-02", "nodash")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Value != "nodash" {
		t.Errorf("expected Value %q, got %q", "nodash", pe.Value)
	}
	if pe.Err != nil {
		t.Errorf("expected nil Err for format error, got %v", pe.Err)
	}
}

func TestParseError_InvalidStartDate(t *testing.T) {
	_, err := ParseBetweenOperator("2006-01-02", "bad*2024-01-01")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Err == nil {
		t.Error("expected non-nil underlying Err for date parse failure")
	}
}

func TestParseError_InvalidEndDate(t *testing.T) {
	_, err := ParseBetweenOperator("2006-01-02", "2024-01-01*bad")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T", err)
	}
	if pe.Err == nil {
		t.Error("expected non-nil underlying Err for date parse failure")
	}
}

// --- Sentinel errors ---

func TestErrNotStruct_FromUpdateFromQueryParams(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams("not a struct", "filter")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotStruct) {
		t.Errorf("expected errors.Is(err, ErrNotStruct) to be true, got false; err: %v", err)
	}
}

func TestErrEmptyTypeTag_FromUpdateFromQueryParams(t *testing.T) {
	type bad struct {
		Name string `filter:"name"`
	}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(bad{Name: "test"}, "filter")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrEmptyTypeTag) {
		t.Errorf("expected errors.Is(err, ErrEmptyTypeTag) to be true, got false; err: %v", err)
	}
}

// --- ValidationError from GetValue with unknown type ---

func TestValidationError_UnknownDataTypeInGetValue(t *testing.T) {
	f := Field{Name: "x", Operator: OperatorEq, Value: "v", Type: "unknown"}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if ve.Kind != "data_type" {
		t.Errorf("expected Kind %q, got %q", "data_type", ve.Kind)
	}
}

// --- Error message formatting ---

func TestValidationError_ErrorMessage(t *testing.T) {
	tests := []struct {
		kind     string
		values   []string
		contains string
	}{
		{"operator", []string{"bad"}, "undefined operator [bad]"},
		{"data_type", []string{"float"}, "undefined data type [float]"},
		{"operator_data_type", []string{"like", "int"}, `operator "like" is not compatible with data type "int"`},
		{"other", []string{"x"}, "validation error"},
	}
	for _, tc := range tests {
		ve := &ValidationError{Kind: tc.kind, Values: tc.values}
		msg := ve.Error()
		if msg != tc.contains && !contains(msg, tc.contains) {
			t.Errorf("Kind=%q: expected message to contain %q, got %q", tc.kind, tc.contains, msg)
		}
	}
}

func TestConversionError_ErrorMessage(t *testing.T) {
	ce := &ConversionError{DataType: "int", Value: "abc"}
	msg := ce.Error()
	if !contains(msg, "abc") || !contains(msg, "int") {
		t.Errorf("expected message to contain value and type, got %q", msg)
	}
}

func TestParseError_ErrorMessage(t *testing.T) {
	pe := &ParseError{Value: "x", Message: "bad format"}
	if pe.Error() != "bad format" {
		t.Errorf("expected %q, got %q", "bad format", pe.Error())
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && len(sub) > 0 && stringContains(s, sub)))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
