package filter

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

func TestGetValue_Int(t *testing.T) {
	f := &Field{Name: "age", Operator: OperatorEq, Value: "42", Type: DataTypeInt}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestGetValue_IntEmpty(t *testing.T) {
	f := &Field{Name: "age", Operator: OperatorEq, Value: "", Type: DataTypeInt}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 0 {
		t.Errorf("expected 0, got %v", val)
	}
}

func TestGetValue_IntInvalid(t *testing.T) {
	f := &Field{Name: "age", Operator: OperatorEq, Value: "abc", Type: DataTypeInt}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestGetValue_IntIn(t *testing.T) {
	f := &Field{Name: "id", Operator: OperatorIn, Value: "1,2,3", Type: DataTypeInt}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ints, ok := val.([]int)
	if !ok {
		t.Fatalf("expected []int, got %T", val)
	}
	if len(ints) != 3 || ints[0] != 1 || ints[1] != 2 || ints[2] != 3 {
		t.Errorf("expected [1 2 3], got %v", ints)
	}
}

func TestGetValue_IntInWithSpaces(t *testing.T) {
	f := &Field{Name: "id", Operator: OperatorIn, Value: "1, 2, 3", Type: DataTypeInt}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ints := val.([]int)
	if len(ints) != 3 {
		t.Errorf("expected 3 elements, got %d", len(ints))
	}
}

func TestGetValue_IntInInvalid(t *testing.T) {
	f := &Field{Name: "id", Operator: OperatorIn, Value: "1,abc,3", Type: DataTypeInt}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid int in list, got nil")
	}
}

func TestGetValue_String(t *testing.T) {
	f := &Field{Name: "name", Operator: OperatorEq, Value: "hello", Type: DataTypeStr}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}
}

func TestGetValue_StringEmpty(t *testing.T) {
	f := &Field{Name: "name", Operator: OperatorEq, Value: "", Type: DataTypeStr}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string, got %v", val)
	}
}

func TestGetValue_BoolTrue(t *testing.T) {
	f := &Field{Name: "active", Operator: OperatorEq, Value: "true", Type: DataTypeBool}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != true {
		t.Errorf("expected true, got %v", val)
	}
}

func TestGetValue_BoolFalse(t *testing.T) {
	f := &Field{Name: "active", Operator: OperatorEq, Value: "false", Type: DataTypeBool}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != false {
		t.Errorf("expected false, got %v", val)
	}
}

func TestGetValue_BoolInvalid(t *testing.T) {
	f := &Field{Name: "active", Operator: OperatorEq, Value: "yes", Type: DataTypeBool}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid bool, got nil")
	}
}

func TestGetValue_DateString(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorEq, Value: "2024-06-15", Type: DataTypeDate}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if val != expected {
		t.Errorf("expected %v, got %v", expected, val)
	}
}

func TestGetValue_DateTimestamp(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorEq, Value: "1718409600", Type: DataTypeDate}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := val.(time.Time)
	if ts.Unix() != 1718409600 {
		t.Errorf("expected timestamp 1718409600, got %d", ts.Unix())
	}
}

func TestGetValue_DateBetween(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorBetween, Value: "2024-01-01*2024-12-31", Type: DataTypeDate}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	between, ok := val.(*Between)
	if !ok {
		t.Fatalf("expected *Between, got %T", val)
	}
	start := between.Start().(time.Time)
	if start.Year() != 2024 || start.Month() != 1 || start.Day() != 1 {
		t.Errorf("unexpected start: %v", start)
	}
}

func TestGetValue_DateInvalidTimestamp(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorEq, Value: "notanumber", Type: DataTypeDate}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid date timestamp, got nil")
	}
}

func TestGetValue_DateDefaultBranch(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorLike, Value: "1718409600", Type: DataTypeDate}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := val.(time.Time)
	if ts.Unix() != 1718409600 {
		t.Errorf("expected timestamp 1718409600, got %d", ts.Unix())
	}
}

func TestGetValue_DateDefaultBranchInvalid(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorLike, Value: "invalid", Type: DataTypeDate}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetValue_DateTimeString(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorEq, Value: "2024-06-15T10:30:00", Type: DataTypeDateTime}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	if val != expected {
		t.Errorf("expected %v, got %v", expected, val)
	}
}

func TestGetValue_DateTimeTimestamp(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorEq, Value: "1718409600", Type: DataTypeDateTime}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := val.(time.Time)
	if ts.Unix() != 1718409600 {
		t.Errorf("expected timestamp 1718409600, got %d", ts.Unix())
	}
}

func TestGetValue_DateTimeBetween(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorBetween, Value: "2024-01-01T00:00:00*2024-12-31T23:59:59", Type: DataTypeDateTime}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := val.(*Between)
	if !ok {
		t.Fatalf("expected *Between, got %T", val)
	}
}

func TestGetValue_DateTimeDefaultBranch(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorLike, Value: "1718409600", Type: DataTypeDateTime}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ts := val.(time.Time)
	if ts.Unix() != 1718409600 {
		t.Errorf("expected timestamp 1718409600, got %d", ts.Unix())
	}
}

func TestGetValue_DateTimeDefaultBranchInvalid(t *testing.T) {
	f := &Field{Name: "created", Operator: OperatorLike, Value: "invalid", Type: DataTypeDateTime}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetValue_UUID(t *testing.T) {
	id := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	f := &Field{Name: "user_id", Operator: OperatorEq, Value: id, Type: DataTypeUUID}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	uid, ok := val.(uuid.UUID)
	if !ok {
		t.Fatalf("expected uuid.UUID, got %T", val)
	}
	if uid.String() != id {
		t.Errorf("expected %s, got %s", id, uid.String())
	}
}

func TestGetValue_UUIDInvalid(t *testing.T) {
	f := &Field{Name: "user_id", Operator: OperatorEq, Value: "not-a-uuid", Type: DataTypeUUID}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid uuid, got nil")
	}
}

func TestGetValue_UUIDIn(t *testing.T) {
	id1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	id2 := "6ba7b811-9dad-11d1-80b4-00c04fd430c8"
	f := &Field{Name: "user_id", Operator: OperatorIn, Value: id1 + "," + id2, Type: DataTypeUUID}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	uids, ok := val.([]uuid.UUID)
	if !ok {
		t.Fatalf("expected []uuid.UUID, got %T", val)
	}
	if len(uids) != 2 {
		t.Errorf("expected 2 uuids, got %d", len(uids))
	}
}

func TestGetValue_UUIDInInvalid(t *testing.T) {
	f := &Field{Name: "user_id", Operator: OperatorIn, Value: "6ba7b810-9dad-11d1-80b4-00c04fd430c8,invalid", Type: DataTypeUUID}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid uuid in list, got nil")
	}
}

func TestGetValue_List(t *testing.T) {
	f := &Field{Name: "tags", Operator: OperatorEq, Value: `["a","b","c"]`, Type: DataTypeList}
	val, err := f.GetValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := val.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", val)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
}

func TestGetValue_ListInvalid(t *testing.T) {
	f := &Field{Name: "tags", Operator: OperatorEq, Value: "not json", Type: DataTypeList}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid list, got nil")
	}
}

func TestGetValue_InvalidDataType(t *testing.T) {
	f := &Field{Name: "x", Operator: OperatorEq, Value: "v", Type: "unknown"}
	_, err := f.GetValue()
	if err == nil {
		t.Fatal("expected error for invalid data type, got nil")
	}
}
