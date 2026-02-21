package utils

import (
	"errors"
	"math"
)

// Money represents a monetary value stored in cents to avoid floating-point precision issues.
// For example, R$ 10.50 is stored as 1050.
type Money uint64

// ConvertToFloat converts the Money value from cents to a float64 representation.
func (m *Money) ConvertToFloat() float64 {
	return float64(*m) / 100
}

// ConvertToInt converts a float64 value to cents and stores it in the Money receiver.
// Returns an error if the value is negative.
func (m *Money) ConvertToInt(value float64) error {
	if value < 0 {
		return errors.New("value must be greater than 0")
	}
	*m = Money(math.Round(value * 100))
	return nil
}
