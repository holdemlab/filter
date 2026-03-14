package adapter

import (
	"context"
	"errors"
	"testing"

	"github.com/holdemlab/filter"
)

func TestMongoQueryD_NilOptions(t *testing.T) {
	_, _, err := MongoQueryD(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil options, got nil")
	}
	if !errors.Is(err, filter.ErrNilOptions) {
		t.Errorf("expected ErrNilOptions, got: %v", err)
	}
}

func TestMongoQueryD_NoFields(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	bsonD, findOpts, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bsonD) != 0 {
		t.Errorf("expected empty bson.D, got %d elements", len(bsonD))
	}
	if findOpts == nil {
		t.Fatal("expected non-nil FindOptions")
	}
}

func TestMongoQueryD_EqField(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", filter.OperatorEq, "John", filter.DataTypeStr)
	bsonD, _, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bsonD) != 1 {
		t.Errorf("expected 1 element, got %d", len(bsonD))
	}
}

func TestMongoQueryD_MultipleFields(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", filter.OperatorEq, "John", filter.DataTypeStr)
	_ = opts.AddField("age", filter.OperatorGreaterThan, "18", filter.DataTypeInt)
	bsonD, _, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bsonD) != 2 {
		t.Errorf("expected 2 elements, got %d", len(bsonD))
	}
}

func TestMongoQueryD_AllOperators(t *testing.T) {
	operators := []struct {
		op    string
		value string
		dType string
	}{
		{filter.OperatorEq, "test", filter.DataTypeStr},
		{filter.OperatorNotEq, "test", filter.DataTypeStr},
		{filter.OperatorGreaterThan, "10", filter.DataTypeInt},
		{filter.OperatorGreaterThanEq, "10", filter.DataTypeInt},
		{filter.OperatorLowerThan, "10", filter.DataTypeInt},
		{filter.OperatorLowerThanEq, "10", filter.DataTypeInt},
		{filter.OperatorLike, "test", filter.DataTypeStr},
		{filter.OperatorBetween, "2024-01-01*2024-12-31", filter.DataTypeDate},
	}
	for _, tc := range operators {
		opts := filter.NewOptions(10, 1, "id", false)
		_ = opts.AddField("field", tc.op, tc.value, tc.dType)
		_, _, err := MongoQueryD(context.Background(), opts)
		if err != nil {
			t.Errorf("unexpected error for operator %q: %v", tc.op, err)
		}
	}
}

func TestMongoQueryD_InOperator(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("tags", filter.OperatorIn, `["a","b"]`, filter.DataTypeList)
	bsonD, _, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bsonD) != 1 {
		t.Errorf("expected 1 element, got %d", len(bsonD))
	}
}

func TestMongoQueryD_SortAsc(t *testing.T) {
	opts := filter.NewOptions(10, 1, "name", false)
	_, findOpts, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findOpts == nil {
		t.Fatal("expected non-nil FindOptions")
	}
}

func TestMongoQueryD_SortDesc(t *testing.T) {
	opts := filter.NewOptions(10, 1, "name", true)
	_, findOpts, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findOpts == nil {
		t.Fatal("expected non-nil FindOptions")
	}
}

func TestMongoQueryD_EmptySortBy(t *testing.T) {
	opts := filter.NewOptions(10, 1, "", false)
	_, findOpts, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findOpts == nil {
		t.Fatal("expected non-nil FindOptions")
	}
}

func TestMongoQueryD_InvalidFieldValue(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("age", filter.OperatorEq, "notanumber", filter.DataTypeInt)
	_, _, err := MongoQueryD(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error for invalid field value, got nil")
	}
	var ce *filter.ConversionError
	if !errors.As(err, &ce) {
		t.Errorf("expected *filter.ConversionError in chain, got: %T: %v", err, err)
	}
}

func TestMongoQueryD_Pagination(t *testing.T) {
	opts := filter.NewOptions(10, 3, "id", false)
	_, findOpts, err := MongoQueryD(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findOpts == nil {
		t.Fatal("expected non-nil FindOptions")
	}
}
