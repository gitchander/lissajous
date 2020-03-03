package main

type Point2f struct {
	X, Y float64
}

func Pt2f(a, b float64) Point2f {
	return Point2f{
		X: a,
		Y: b,
	}
}

func (p Point2f) Add(q Point2f) Point2f {
	return Point2f{
		X: p.X + q.X,
		Y: p.Y + q.Y,
	}
}

func (p Point2f) Sub(q Point2f) Point2f {
	return Point2f{
		X: p.X - q.X,
		Y: p.Y - q.Y,
	}
}

func (p Point2f) MulScalar(scalar float64) Point2f {
	return Point2f{
		X: p.X * scalar,
		Y: p.Y * scalar,
	}
}

func (p Point2f) DivScalar(scalar float64) Point2f {
	return Point2f{
		X: p.X / scalar,
		Y: p.Y / scalar,
	}
}
