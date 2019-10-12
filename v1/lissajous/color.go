package main

type Color struct {
	R, G, B float64
}

func (a Color) Norm() Color {
	return Color{
		R: crop(a.R),
		G: crop(a.G),
		B: crop(a.B),
	}
}

func Clerp(a, b Color, t float64) Color {
	return Color{
		R: lerp(a.R, b.R, t),
		G: lerp(a.G, b.G, t),
		B: lerp(a.B, b.B, t),
	}
}
