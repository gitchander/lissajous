package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/envoker/gotk3/cairo"
)

type Canvas struct {
	mutex   sync.Mutex
	config  Config
	surface *cairo.Surface
	context *cairo.Context
	size    Size
}

func NewCanvas(config Config) *Canvas {
	return &Canvas{config: config}
}

func (c *Canvas) Config() Config {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.config
}

func (c *Canvas) SetConfig(config Config) error {

	if err := config.Error(); err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.config = config

	return nil
}

func (c *Canvas) Resize(Width, Height int) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.surface != nil {
		var (
			w = c.surface.GetWidth()
			h = c.surface.GetHeight()
		)
		if (Width <= w) && (Height <= h) {
			c.size = Size{Width, Height}
			return nil
		}
	}

	var (
		w = ceilPowerOfTwo(Width)
		h = ceilPowerOfTwo(Height)
	)
	surface, err := cairo.NewSurface(cairo.FORMAT_ARGB32, w, h)
	if err != nil {
		return err
	}

	fmt.Println("resize:", w, h)

	c.surface = surface
	c.context = cairo.Create(c.surface)
	c.size = Size{Width, Height}

	return nil
}

func (c *Canvas) Draw(context *cairo.Context) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.surface != nil {
		context.SetSourceSurface(c.surface, 0, 0)
		context.Paint()
	}
}

func (c *Canvas) Render(deltaT time.Duration) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	var (
		width  = c.size.Width
		height = c.size.Height
	)

	bg := Color{0, 0, 0}
	fg := Color{0, 1, 0}

	if c.surface == nil {
		return
	}

	context := c.context

	context.SetSourceRGB(bg.R, bg.G, bg.B)
	context.Rectangle(0, 0, float64(width), float64(height))
	context.Fill()

	context.SetSourceRGB(fg.R, fg.G, fg.B)

	center := Point{
		X: float64(width),
		Y: float64(height),
	}.DivScalar(2)

	t := deltaT.Seconds()

	tailDuration := c.config.TailDuration

	count := c.config.TailSegments + 1
	dt := tailDuration / float64(count)

	amplitude := 0.7 * (float64(min(width, height)) * 0.5)

	AmplitudeX := amplitude
	AmplitudeY := amplitude

	FrequencyX := c.config.FrequencyX
	FrequencyY := c.config.FrequencyY

	xw := 2 * math.Pi * FrequencyX
	yw := 2 * math.Pi * FrequencyY

	phaseShift := c.config.PhaseShift

	curr := Point{
		X: AmplitudeX * math.Sin(xw*t+phaseShift),
		Y: AmplitudeY * math.Sin(yw*t),
	}.Add(center)

	prev := curr

	for i := 0; i < count; i++ {

		curr = Point{
			X: AmplitudeX * math.Sin(xw*t+phaseShift),
			Y: AmplitudeY * math.Sin(yw*t),
		}.Add(center)

		cl := Clerp(bg, fg, float64(i)/float64(count-1))
		context.SetSourceRGB(cl.R, cl.G, cl.B)

		context.MoveTo(prev.X, prev.Y)
		context.LineTo(curr.X, curr.Y)
		context.Stroke()

		prev = curr

		t += dt
	}
}
