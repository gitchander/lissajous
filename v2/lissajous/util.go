package main

func crop(x float64) float64 {
	if x < 0 {
		x = 0
	}
	if x > 1 {
		x = 1
	}
	return x
}

// Lerp - Linear interpolation
// t= [0, 1]
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
