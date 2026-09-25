package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero    = errors.New("cannot divide by zero")
	ErrNegativeSqrtInput = errors.New("cannot calculate square root of a negative number")
	ErrOverflow          = errors.New("calculation resulted in overflow or invalid number")
)

// Add calculates a + b.
func Add(a, b float64) (float64, error) {
	res := a + b
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Subtract calculates a - b.
func Subtract(a, b float64) (float64, error) {
	res := a - b
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Multiply calculates a * b.
func Multiply(a, b float64) (float64, error) {
	res := a * b
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Divide calculates a / b and guards against division by zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	res := a / b
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Power calculates a ^ b.
func Power(a, b float64) (float64, error) {
	res := math.Pow(a, b)
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Sqrt calculates square root of a and guards against negative inputs.
func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrtInput
	}
	res := math.Sqrt(a)
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}

// Percentage calculates (a * b) / 100.
func Percentage(a, b float64) (float64, error) {
	res := (a * b) / 100.0
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, ErrOverflow
	}
	return res, nil
}
