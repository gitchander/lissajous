package main

type Size struct {
	Width, Height int
}

func (x Size) IsZero() bool {
	return (x.Width == 0) && (x.Height == 0)
}
