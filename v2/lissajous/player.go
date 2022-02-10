package main

import (
	"errors"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Player struct {
	guard   sync.Mutex
	c       *Core
	da      *gtk.DrawingArea
	stop    chan struct{}
	durCh   chan time.Duration
	started bool

	td time.Duration
}

func NewPlayer(c *Core, da *gtk.DrawingArea) *Player {

	// glib.IdleAdd(func() bool {
	// 	da.QueueDraw()
	// 	return true // continue execute this func
	// 	//return false // if nead to stop this function
	// })

	return &Player{
		c:     c,
		da:    da,
		stop:  make(chan struct{}),
		durCh: make(chan time.Duration),
	}
}

func (p *Player) Started() bool {
	p.guard.Lock()
	defer p.guard.Unlock()
	return p.started
}

func (p *Player) Stopped() bool {
	p.guard.Lock()
	defer p.guard.Unlock()
	return !(p.started)
}

func (p *Player) Start() error {

	p.guard.Lock()
	defer p.guard.Unlock()

	if p.c == nil {
		return errors.New("Curve is not set")
	}
	// if p.da == nil {
	// 	return errors.New("DrawingArea is not set")
	// }

	if p.started {
		return errors.New("is started")
	}

	render := func(d time.Duration) {
		p.c.Render(d)
		glib.IdleAdd(p.da.QueueDraw)
	}

	go timeRender(p.stop, p.durCh, p.td, render)

	p.started = true

	return nil
}

func (p *Player) Stop() error {

	p.guard.Lock()
	defer p.guard.Unlock()

	if !p.started {
		return nil
	}

	p.stop <- struct{}{}
	p.td = <-p.durCh

	p.started = false

	return nil
}

// type Renderer interface {
// 	Render(time.Duration)
// }

func timeRender(stop <-chan struct{}, durCh chan<- time.Duration, d time.Duration, render func(time.Duration)) {

	framesPerSecond := 30 // frames per second

	dt := time.Second / time.Duration(framesPerSecond)
	t := time.Now()
	t0 := t.Add(-d)

	for {
		select {
		case <-stop:
			durCh <- d
			return
		default:
		}

		render(d)

		now := time.Now()
		d = now.Sub(t0)

		t = t.Add(dt)
		dSleep := t.Sub(now)
		if dSleep > 0 {
			time.Sleep(dSleep)
		}
	}
}
