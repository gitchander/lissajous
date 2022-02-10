package main

import (
	"fmt"
)

func clamp(x float64, min, max float64) float64 {
	if min > max {
		err := fmt.Errorf("clamp: (min=%f) > (max=%f)", min, max)
		panic(err)
	}
	if x < min {
		x = min
	}
	if x > max {
		x = max
	}
	return x
}

func clamp01(x float64) float64 {
	return clamp(x, 0, 1)
}

// Lerp - Linear interpolation
// t = [0..1]
// (t == 0) => v0
// (t == 1) => v1
func lerp(v0, v1 float64, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ceilPowerOfTwo(x int) int {
	d := 1
	for d < x {
		d *= 2
	}
	return d
}

func isPowerOfTwo(x int) bool {
	return ((x > 0) && ((x & (x - 1)) == 0))
}
