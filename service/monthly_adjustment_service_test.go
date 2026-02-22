package service

import (
	"testing"
)

func TestMonthsBetween_EmptyFrom(t *testing.T) {
	months, err := monthsBetween("", "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 1 {
		t.Fatalf("expected 1 month, got %d", len(months))
	}

	if months[0] != "2025-07" {
		t.Errorf("expected %q, got %q", "2025-07", months[0])
	}
}

func TestMonthsBetween_ConsecutiveMonths(t *testing.T) {
	months, err := monthsBetween("2025-06", "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 1 {
		t.Fatalf("expected 1 month, got %d", len(months))
	}

	if months[0] != "2025-07" {
		t.Errorf("expected %q, got %q", "2025-07", months[0])
	}
}

func TestMonthsBetween_SameMonth(t *testing.T) {
	months, err := monthsBetween("2025-07", "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 0 {
		t.Errorf("expected 0 months when from == to, got %d: %v", len(months), months)
	}
}

func TestMonthsBetween_MultipleMonths(t *testing.T) {
	months, err := monthsBetween("2025-03", "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"2025-04", "2025-05", "2025-06", "2025-07"}
	if len(months) != len(expected) {
		t.Fatalf("expected %d months, got %d: %v", len(expected), len(months), months)
	}

	for i, m := range months {
		if m != expected[i] {
			t.Errorf("month[%d]: expected %q, got %q", i, expected[i], m)
		}
	}
}

func TestMonthsBetween_CrossYearBoundary(t *testing.T) {
	months, err := monthsBetween("2024-11", "2025-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"2024-12", "2025-01", "2025-02"}
	if len(months) != len(expected) {
		t.Fatalf("expected %d months, got %d: %v", len(expected), len(months), months)
	}

	for i, m := range months {
		if m != expected[i] {
			t.Errorf("month[%d]: expected %q, got %q", i, expected[i], m)
		}
	}
}

func TestMonthsBetween_FullYearGap(t *testing.T) {
	months, err := monthsBetween("2024-01", "2025-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 12 {
		t.Fatalf("expected 12 months for a full year gap, got %d: %v", len(months), months)
	}

	if months[0] != "2024-02" {
		t.Errorf("first month: expected %q, got %q", "2024-02", months[0])
	}
	if months[11] != "2025-01" {
		t.Errorf("last month: expected %q, got %q", "2025-01", months[11])
	}
}

func TestMonthsBetween_FromAfterTo(t *testing.T) {
	months, err := monthsBetween("2025-07", "2025-03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 0 {
		t.Errorf("expected 0 months when from > to, got %d: %v", len(months), months)
	}
}

func TestMonthsBetween_InvalidFromFormat(t *testing.T) {
	_, err := monthsBetween("invalid", "2025-07")
	if err == nil {
		t.Fatal("expected error for invalid 'from' format, got nil")
	}
}

func TestMonthsBetween_InvalidToFormat(t *testing.T) {
	_, err := monthsBetween("2025-07", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid 'to' format, got nil")
	}
}

func TestMonthsBetween_DecemberToJanuary(t *testing.T) {
	months, err := monthsBetween("2025-12", "2026-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 1 {
		t.Fatalf("expected 1 month, got %d: %v", len(months), months)
	}

	if months[0] != "2026-01" {
		t.Errorf("expected %q, got %q", "2026-01", months[0])
	}
}

func TestMonthsBetween_MultiYearGap(t *testing.T) {
	months, err := monthsBetween("2023-06", "2025-06")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2023-07 through 2025-06 = 24 months
	if len(months) != 24 {
		t.Fatalf("expected 24 months for a 2-year gap, got %d", len(months))
	}

	if months[0] != "2023-07" {
		t.Errorf("first month: expected %q, got %q", "2023-07", months[0])
	}
	if months[23] != "2025-06" {
		t.Errorf("last month: expected %q, got %q", "2025-06", months[23])
	}
}

func TestCurrentYearMonth_Format(t *testing.T) {
	ym := currentYearMonth()

	if len(ym) != 7 {
		t.Errorf("expected year-month format YYYY-MM (length 7), got %q (length %d)", ym, len(ym))
	}

	if ym[4] != '-' {
		t.Errorf("expected dash at position 4, got %q", string(ym[4]))
	}
}
