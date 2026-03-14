package filter

import (
	"testing"
)

func TestNewOptions_Defaults(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	if opts.Limit() != 10 {
		t.Errorf("expected limit 10, got %d", opts.Limit())
	}
	if opts.Page() != 1 {
		t.Errorf("expected page 1, got %d", opts.Page())
	}
	if opts.SortBy() != "id" {
		t.Errorf("expected sortBy 'id', got %q", opts.SortBy())
	}
	if opts.Desc() != false {
		t.Errorf("expected desc false, got %v", opts.Desc())
	}
}

func TestNewOptions_Descending(t *testing.T) {
	opts := NewOptions(20, 3, "name", true)
	if opts.Desc() != true {
		t.Error("expected desc true")
	}
	if opts.Limit() != 20 {
		t.Errorf("expected limit 20, got %d", opts.Limit())
	}
	if opts.Page() != 3 {
		t.Errorf("expected page 3, got %d", opts.Page())
	}
}

func TestNewOptions_NegativeLimit(t *testing.T) {
	opts := NewOptions(-5, 1, "id", false)
	if opts.Limit() != 1 {
		t.Errorf("expected limit 1 for negative input, got %d", opts.Limit())
	}
}

func TestNewOptions_ZeroLimit(t *testing.T) {
	opts := NewOptions(0, 1, "id", false)
	if opts.Limit() != 1 {
		t.Errorf("expected limit 1 for zero input, got %d", opts.Limit())
	}
}

func TestNewOptions_NegativePage(t *testing.T) {
	opts := NewOptions(10, -1, "id", false)
	if opts.Page() != 1 {
		t.Errorf("expected page 1 for negative input, got %d", opts.Page())
	}
}

func TestNewOptions_ZeroPage(t *testing.T) {
	opts := NewOptions(10, 0, "id", false)
	if opts.Page() != 1 {
		t.Errorf("expected page 1 for zero input, got %d", opts.Page())
	}
}

func TestAddField_Valid(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.AddField("name", OperatorEq, "John", DataTypeStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := opts.Fields()
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].Name != "name" || fields[0].Operator != OperatorEq || fields[0].Value != "John" || fields[0].Type != DataTypeStr {
		t.Errorf("unexpected field: %+v", fields[0])
	}
}

func TestAddField_InvalidOperator(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.AddField("name", "invalid", "John", DataTypeStr)
	if err == nil {
		t.Fatal("expected error for invalid operator, got nil")
	}
}

func TestAddField_InvalidDataType(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.AddField("name", OperatorEq, "John", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid data type, got nil")
	}
}

func TestAddField_Multiple(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", OperatorEq, "John", DataTypeStr)
	_ = opts.AddField("age", OperatorGreaterThan, "18", DataTypeInt)
	if len(opts.Fields()) != 2 {
		t.Errorf("expected 2 fields, got %d", len(opts.Fields()))
	}
}

func TestFieldByName_Found(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", OperatorEq, "John", DataTypeStr)
	f, found := opts.FieldByName("name")
	if !found {
		t.Fatal("expected field to be found")
	}
	if f.Value != "John" {
		t.Errorf("expected value 'John', got %q", f.Value)
	}
}

func TestFieldByName_NotFound(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	_, found := opts.FieldByName("nonexistent")
	if found {
		t.Fatal("expected field not to be found")
	}
}

func TestUpdateFromQueryParams_Basic(t *testing.T) {
	type Filter struct {
		Name   string `form:"name" type:"string"`
		Age    string `form:"age" type:"int"`
		Active string `form:"active" type:"bool"`
	}
	f := Filter{Name: "John", Age: "25", Active: "true"}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Fields()) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(opts.Fields()))
	}
}

func TestUpdateFromQueryParams_WithOperator(t *testing.T) {
	type Filter struct {
		Age string `form:"age" type:"int"`
	}
	f := Filter{Age: "gte#18"}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := opts.Fields()
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].Operator != OperatorGreaterThanEq {
		t.Errorf("expected operator %q, got %q", OperatorGreaterThanEq, fields[0].Operator)
	}
	if fields[0].Value != "18" {
		t.Errorf("expected value '18', got %q", fields[0].Value)
	}
}

func TestUpdateFromQueryParams_EmptyValueSkipped(t *testing.T) {
	type Filter struct {
		Name string `form:"name" type:"string"`
	}
	f := Filter{Name: ""}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Fields()) != 0 {
		t.Errorf("expected 0 fields for empty value, got %d", len(opts.Fields()))
	}
}

func TestUpdateFromQueryParams_MissingTypeTag(t *testing.T) {
	type Filter struct {
		Name string `form:"name"`
	}
	f := Filter{Name: "John"}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err == nil {
		t.Fatal("expected error for missing type tag, got nil")
	}
}

func TestUpdateFromQueryParams_NonStringField(t *testing.T) {
	type Filter struct {
		Count int `form:"count" type:"int"`
	}
	f := Filter{Count: 5}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err == nil {
		t.Fatal("expected error for non-string field, got nil")
	}
}

func TestUpdateFromQueryParams_Pointer(t *testing.T) {
	type Filter struct {
		Name string `form:"name" type:"string"`
	}
	f := &Filter{Name: "John"}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Fields()) != 1 {
		t.Errorf("expected 1 field, got %d", len(opts.Fields()))
	}
}

func TestUpdateFromQueryParams_NonStruct(t *testing.T) {
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams("not a struct", "form")
	if err == nil {
		t.Fatal("expected error for non-struct, got nil")
	}
}

func TestUpdateFromQueryParams_NoMatchingTag(t *testing.T) {
	type Filter struct {
		Name string `json:"name" type:"string"`
	}
	f := Filter{Name: "John"}
	opts := NewOptions(10, 1, "id", false)
	err := opts.UpdateFromQueryParams(f, "form")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Fields()) != 0 {
		t.Errorf("expected 0 fields when tag doesn't match, got %d", len(opts.Fields()))
	}
}
