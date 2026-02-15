package utils

import (
	"errors"
	"math"
)

type Money uint64

// 1050 for R$ 10,50
func (m *Money) ConvertToFloat() float64 {
	return float64(*m) / 100
}

func (m *Money) ConvertToInt(value float64) error {
	if value < 0 {
		return errors.New("Value must be greater than 0")
	}
	*m = Money(math.Round(value * 100))
	return nil

}

// func (m *Money)
