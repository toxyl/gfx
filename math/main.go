package math

import (
	"fmt"
	"math"
)

const (
	Pi                     = math.Pi
	Tau                    = 2 * Pi
	Phi                    = math.Phi
	MaxFloat64             = math.MaxFloat64
	SmallestNonzeroFloat64 = math.SmallestNonzeroFloat64
)

type Number interface {
	int64 | int32 | int16 | int8 | int | uint64 | uint32 | uint16 | uint8 | uint | float32 | float64
}

type Float interface {
	float32 | float64
}

func Cos[N Number](x N) N                 { return N(math.Cos(float64(x))) }
func Sin[N Number](x N) N                 { return N(math.Sin(float64(x))) }
func Exp[N Number](x N) N                 { return N(math.Exp(float64(x))) }
func Round[N Number](x N) N               { return N(math.Round(float64(x))) }
func Abs[N Number](x N) N                 { return N(math.Abs(float64(x))) }
func Sqrt[N Number](x N) N                { return N(math.Sqrt(float64(x))) }
func Add[N Number](x, y N) N              { return N(x + y) }
func Sub[N Number](x, y N) N              { return N(x - y) }
func Div[N Number](x, y N) N              { return N(x / y) }
func Mul[N Number](x, y N) N              { return N(x * y) }
func Pow[N Number](x, y N) N              { return N(math.Pow(float64(x), float64(y))) }
func Max[N Number](x, y N) N              { return N(math.Max(float64(x), float64(y))) }
func Min[N Number](x, y N) N              { return N(math.Min(float64(x), float64(y))) }
func Mod[N Number](x, y N) N              { return N(math.Mod(float64(x), float64(y))) }
func Blend[N Number](x, y N, p float64) N { return N(float64(x)*(1-p) + float64(y)*p) }

func Avg[N Number](x ...N) N {
	res := 0.0
	for i, n := range x {
		if i == 0 {
			res = float64(n)
			continue
		}
		res = (res + float64(n)) / 2.0
	}
	return N(res)
}

func Clamp[N Number](x, min, max N) N {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func Delta[N Number](x, y N) N {
	max := Max(x, y)
	min := Min(x, y)
	if max >= 0 && min <= 0 {
		return N(Abs(min) + max)
	}
	return Abs(max - min)
}

func Wrap[N Number](val, min, max N) N {
	if min > max {
		min, max = max, min
	}
	rangeSize := max - min
	val = Mod(val-min, rangeSize)
	if val < 0 {
		val += rangeSize
	}
	return val + min
}

func FormatNumber[N Number](number N, decimals int) string {
	suffix := " "
	divisor := 1.0
	n := float64(number)

	switch {
	case n >= 1_000_000_000_000 || n <= -1_000_000_000_000:
		suffix = "t"
		divisor = 1_000_000_000_000
	case n >= 1_000_000_000 || n <= -1_000_000_000:
		suffix = "b"
		divisor = 1_000_000_000
	case n >= 1_000_000 || n <= -1_000_000:
		suffix = "m"
		divisor = 1_000_000
	case n >= 1_000 || n <= -1_000:
		suffix = "k"
		divisor = 1_000
	}

	return fmt.Sprintf("%*.*f", 5+decimals, decimals, n/divisor) + suffix
}
