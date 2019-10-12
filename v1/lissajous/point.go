package main

type Point struct {
	X, Y float64
}

func (a Point) Add(b Point) (c Point) {
	c.X = a.X + b.X
	c.Y = a.Y + b.Y
	return
}

func (a Point) Sub(b Point) (c Point) {
	c.X = a.X - b.X
	c.Y = a.Y - b.Y
	return
}

func (a Point) MulScalar(s float64) (b Point) {
	b.X = s * a.X
	b.Y = s * a.Y
	return
}

func (a Point) DivScalar(s float64) (b Point) {
	b.X = a.X / s
	b.Y = a.Y / s
	return
}
