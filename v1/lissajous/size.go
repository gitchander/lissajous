package main

type Size struct {
	Width, Height int
}

func (s Size) IsZero() bool {
	return (s.Width == 0) && (s.Height == 0)
}
