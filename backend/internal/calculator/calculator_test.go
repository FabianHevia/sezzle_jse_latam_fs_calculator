package calculator_test

import (
	"errors"
	"math"
	"testing"

	"github.com/takehome/calculator/backend/internal/calculator"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 10, 5, 15, nil},
		{"negative numbers", -10, -5, -15, nil},
		{"mixed signs", -10, 15, 5, nil},
		{"floating point", 1.5, 2.3, 3.8, nil},
		{"zero", 0, 0, 0, nil},
		{"overflow to inf", math.MaxFloat64, math.MaxFloat64, 0, calculator.ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Add(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Add(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 10, 5, 5, nil},
		{"negative result", 5, 10, -5, nil},
		{"subtract negative", 10, -5, 15, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Subtract(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Subtract(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 4, 5, 20, nil},
		{"zero factor", 0, 100, 0, nil},
		{"negative factor", -4, 5, -20, nil},
		{"overflow", math.MaxFloat64, 2, 0, calculator.ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Multiply(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Multiply(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"exact division", 10, 2, 5, nil},
		{"decimal division", 7, 2, 3.5, nil},
		{"division by zero", 10, 0, 0, calculator.ErrDivisionByZero},
		{"zero divided by number", 0, 5, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Divide(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Divide(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"2 power 3", 2, 3, 8, nil},
		{"10 power 0", 10, 0, 1, nil},
		{"power overflow", 10, 1000, 0, calculator.ErrOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Power(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Power(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{"sqrt of 16", 16, 4, nil},
		{"sqrt of 0", 0, 0, nil},
		{"negative sqrt input", -4, 0, calculator.ErrNegativeSqrtInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Sqrt(tt.a)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Sqrt(%v) error = %v, wantErr %v", tt.a, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.a, got, tt.want)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"20 percent of 100", 20, 100, 20, nil},
		{"50 percent of 80", 50, 80, 40, nil},
		{"0 percent of 50", 0, 50, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Percentage(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Percentage(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
			}
			if err == nil && math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
