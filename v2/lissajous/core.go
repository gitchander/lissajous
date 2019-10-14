package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gotk3/gotk3/cairo"
)

type Core struct {
	mutex     sync.Mutex
	config    Config
	surface   *cairo.Surface
	context   *cairo.Context
	size      Size
	allocSize Size
}

func NewCore(config Config) *Core {
	return &Core{config: config}
}

func (c *Core) Config() Config {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.config
}

func (c *Core) SetConfig(config Config) error {

	if err := config.Error(); err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.config = config

	return nil
}

func (c *Core) Resize(Width, Height int) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

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

	fmt.Println("alloc size:", allocSize)

	c.surface = surface
	c.context = cairo.Create(c.surface)
	c.size = Size{Width, Height}
	c.allocSize = allocSize

	return nil
}

func (c *Core) Draw(context *cairo.Context) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.surface != nil {
		context.SetSourceSurface(c.surface, 0, 0)
		context.Paint()
	}
}

func (c *Core) Render(deltaT time.Duration) {

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

	center := Point2f{
		X: float64(width),
		Y: float64(height),
	}.DivScalar(2)

	t := deltaT.Seconds()

	tailDuration := c.config.TailDuration

	count := c.config.TailSegments + 1
	dt := tailDuration / float64(count)

	amplitude := 0.7 * (float64(minInt(width, height)) * 0.5)

	AmplitudeX := amplitude
	AmplitudeY := amplitude

	FrequencyX := c.config.FrequencyX
	FrequencyY := c.config.FrequencyY

	xw := 2 * math.Pi * FrequencyX
	yw := 2 * math.Pi * FrequencyY

	phaseShift := c.config.PhaseShift

	curr := Point2f{
		X: AmplitudeX * math.Sin(xw*t+phaseShift),
		Y: AmplitudeY * math.Sin(yw*t),
	}.Add(center)

	prev := curr

	for i := 0; i < count; i++ {

		curr = Point2f{
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
