package filter

import (
	"testing"
	"time"
)

func TestParseBetweenOperator_ValidDate(t *testing.T) {
	result, err := ParseBetweenOperator(formatDate, "2024-01-01*2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	start := result.Start().(time.Time)
	end := result.End().(time.Time)
	if start != time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("unexpected start: %v", start)
	}
	if end != time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC) {
		t.Errorf("unexpected end: %v", end)
	}
}

func TestParseBetweenOperator_ValidDateTime(t *testing.T) {
	result, err := ParseBetweenOperator(formatDateTime, "2024-01-01T10:00:00*2024-12-31T23:59:59")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	start := result.Start().(time.Time)
	end := result.End().(time.Time)
	if start != time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC) {
		t.Errorf("unexpected start: %v", start)
	}
	if end != time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC) {
		t.Errorf("unexpected end: %v", end)
	}
}

func TestParseBetweenOperator_MissingSeparator(t *testing.T) {
	_, err := ParseBetweenOperator(formatDate, "2024-01-01")
	if err == nil {
		t.Fatal("expected error for missing separator, got nil")
	}
}

func TestParseBetweenOperator_TooManySeparators(t *testing.T) {
	_, err := ParseBetweenOperator(formatDate, "2024-01-01*2024-06-15*2024-12-31")
	if err == nil {
		t.Fatal("expected error for too many separators, got nil")
	}
}

func TestParseBetweenOperator_EmptyValue(t *testing.T) {
	_, err := ParseBetweenOperator(formatDate, "")
	if err == nil {
		t.Fatal("expected error for empty value, got nil")
	}
}

func TestParseBetweenOperator_InvalidStartDate(t *testing.T) {
	_, err := ParseBetweenOperator(formatDate, "invalid*2024-12-31")
	if err == nil {
		t.Fatal("expected error for invalid start date, got nil")
	}
}

func TestParseBetweenOperator_InvalidEndDate(t *testing.T) {
	_, err := ParseBetweenOperator(formatDate, "2024-01-01*invalid")
	if err == nil {
		t.Fatal("expected error for invalid end date, got nil")
	}
}
