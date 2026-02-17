package utils

import (
	"errors"
	"math"
)

type Money uint64

// Represents monetary values in cents (e.g., 1050 for R$ 10.50)
func (m *Money) ConvertToFloat() float64 {
	return float64(*m) / 100
}

func (m *Money) ConvertToInt(value float64) error {
	if value < 0 {
		return errors.New("value must be greater than 0")
	}
	*m = Money(math.Round(value * 100))
	return nil

}

// Placeholder for additional Money methods
