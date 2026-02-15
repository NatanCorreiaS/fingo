package utils

import "testing"

func TestConvertToInt(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		wantError bool
	}{
		{"18.91 -> 1891", 18.91, false},
		{"-631.2 -> ERROR", -631.2, true},
		{"0 -> 0", 0, false},
	}

	var m Money
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.ConvertToInt(tt.value)
			if (err != nil) != tt.wantError {
				t.Errorf("Money.ConvertToInt(%v) error = %v, wantError = %v\n", tt.value, err, tt.wantError)
			}

		})
	}
}

func TestConvertToFloat(t *testing.T) {
	tests := []struct {
		name     string
		value    Money
		expected float64
	}{
		{"431987 -> 4319.87", 431987, 4319.87},
		{"1999 -> 19.99", 1999, 19.99},
	}

	var m Money

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m = tt.value
			if m.ConvertToFloat() != tt.expected {
				t.Errorf("Money.ConvertToFloat() value = %v, expected = %v\n", m.ConvertToFloat(), tt.expected)
			}
		})
	}
}
