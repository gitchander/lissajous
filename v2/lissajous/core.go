package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gotk3/gotk3/cairo"
)

type Core struct {
	guard     sync.Mutex
	config    Config
	surface   *cairo.Surface
	context   *cairo.Context
	size      Size
	allocSize Size
}

func NewCore(config Config) (*Core, error) {
	c := new(Core)
	err := c.SetConfig(config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Core) Config() Config {
	c.guard.Lock()
	defer c.guard.Unlock()
	return c.config
}

func (c *Core) SetConfig(config Config) error {

	err := checkConfig(config)
	if err != nil {
		return err
	}

	c.guard.Lock()
	defer c.guard.Unlock()

	c.config = config

	return nil
}

func (c *Core) Resize(Width, Height int) error {

	c.guard.Lock()
	defer c.guard.Unlock()

	if c.surface != nil {

		if (Width <= c.allocSize.Width) && (Height <= c.allocSize.Height) {

			c.size = Size{Width, Height}

			fmt.Println(c.size, c.allocSize)
			return nil
		}
	}

	allocSize := Size{
		Width:  ceilPowerOfTwo(Width),
		Height: ceilPowerOfTwo(Height),
	}

	surface := cairo.CreateImageSurface(cairo.FORMAT_ARGB32, allocSize.Width, allocSize.Height)

	//fmt.Println("alloc size:", allocSize)

	c.surface = surface
	c.context = cairo.Create(c.surface)
	c.size = Size{Width, Height}
	c.allocSize = allocSize

	return nil
}

func (c *Core) Draw(context *cairo.Context) {
	c.guard.Lock()
	{
		if c.surface != nil {
			context.SetSourceSurface(c.surface, 0, 0)
			context.Paint()
		}
	}
	c.guard.Unlock()
}

func (c *Core) Render(deltaT time.Duration) {

	c.guard.Lock()
	defer c.guard.Unlock()

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

	center := Point2f{
		X: float64(width),
		Y: float64(height),
	}.DivScalar(2)

	t := deltaT.Seconds()

	tailDuration := c.config.TailDuration.Seconds()

	count := c.config.TailSegments + 1
	dt := tailDuration / float64(count)

	amplitude := 0.7 * (float64(minInt(width, height)) * 0.5)

	var (
		AmplitudeX = amplitude
		AmplitudeY = amplitude

		freqA = 2 * math.Pi * c.config.FreqA
		freqB = 2 * math.Pi * c.config.FreqB
		phase = 2 * math.Pi * c.config.Phase
	)

	curr := Point2f{
		X: AmplitudeX * math.Sin(freqA*t+phase),
		Y: AmplitudeY * math.Sin(freqB*t),
	}.Add(center)

	prev := curr

	for i := 0; i < count; i++ {

		curr = Point2f{
			X: AmplitudeX * math.Sin(freqA*t+phase),
			Y: AmplitudeY * math.Sin(freqB*t),
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
