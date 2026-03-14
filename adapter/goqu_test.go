package adapter

import (
	"context"
	"errors"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/holdemlab/filter"
)

func TestGoquQuery_NilOptions(t *testing.T) {
	dataset := goqu.From("test").Select("*")
	_, err := GoquQuery(context.Background(), nil, dataset)
	if err == nil {
		t.Fatal("expected error for nil options, got nil")
	}
	if !errors.Is(err, filter.ErrNilOptions) {
		t.Errorf("expected ErrNilOptions, got: %v", err)
	}
}

func TestGoquQuery_NilDataset(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_, err := GoquQuery(context.Background(), opts, nil)
	if err == nil {
		t.Fatal("expected error for nil dataset, got nil")
	}
	if !errors.Is(err, filter.ErrNilDataset) {
		t.Errorf("expected ErrNilDataset, got: %v", err)
	}
}

func TestGoquQuery_NoFields(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGoquQuery_EqField(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", filter.OperatorEq, "John", filter.DataTypeStr)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_LikeField(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", filter.OperatorLike, "John", filter.DataTypeStr)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_BetweenField(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("created", filter.OperatorBetween, "2024-01-01*2024-12-31", filter.DataTypeDate)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_SortAsc(t *testing.T) {
	opts := filter.NewOptions(10, 1, "name", false)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_SortDesc(t *testing.T) {
	opts := filter.NewOptions(10, 1, "name", true)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_EmptySortBy(t *testing.T) {
	opts := filter.NewOptions(10, 1, "", false)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sql, _, err := result.ToSQL()
	if err != nil {
		t.Fatalf("error generating SQL: %v", err)
	}
	if sql == "" {
		t.Fatal("expected non-empty SQL")
	}
}

func TestGoquQuery_MultipleFields(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("name", filter.OperatorEq, "John", filter.DataTypeStr)
	_ = opts.AddField("age", filter.OperatorGreaterThan, "18", filter.DataTypeInt)
	dataset := goqu.From("test").Select("*")
	result, err := GoquQuery(context.Background(), opts, dataset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGoquQuery_InvalidFieldValue(t *testing.T) {
	opts := filter.NewOptions(10, 1, "id", false)
	_ = opts.AddField("age", filter.OperatorEq, "notanumber", filter.DataTypeInt)
	dataset := goqu.From("test").Select("*")
	_, err := GoquQuery(context.Background(), opts, dataset)
	if err == nil {
		t.Fatal("expected error for invalid field value, got nil")
	}
	var ce *filter.ConversionError
	if !errors.As(err, &ce) {
		t.Errorf("expected *filter.ConversionError in chain, got: %T: %v", err, err)
	}
}
