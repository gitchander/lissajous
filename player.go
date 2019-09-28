package main

import (
	"errors"
	"sync"
	"time"

	"github.com/envoker/gotk3/glib"
	"github.com/envoker/gotk3/gtk"
)

type Player struct {
	guard   sync.Mutex
	da      *gtk.DrawingArea
	c       *Canvas
	quit    chan bool
	started bool

	td    time.Duration
	durCh chan time.Duration
}

func NewPlayer() *Player {
	return &Player{
		quit:  make(chan bool),
		durCh: make(chan time.Duration),
	}
}

func (p *Player) Started() bool {
	p.guard.Lock()
	ok := p.started
	p.guard.Unlock()
	return ok
}

func (p *Player) Stopped() bool {
	p.guard.Lock()
	ok := !(p.started)
	p.guard.Unlock()
	return ok
}

func (p *Player) SetCurve(c *Canvas) {
	p.guard.Lock()
	p.c = c
	p.guard.Unlock()
}

func (p *Player) SetDrawingArea(da *gtk.DrawingArea) {
	p.guard.Lock()
	p.da = da
	p.guard.Unlock()
}

func (p *Player) Start() error {

	p.guard.Lock()
	defer p.guard.Unlock()

	if p.c == nil {
		return errors.New("Curve is not set")
	}
	if p.da == nil {
		return errors.New("DrawingArea is not set")
	}

	if p.started {
		return errors.New("is started")
	}

	go func(dur time.Duration) {

		framesPerSecond := 30 // frames per second

		dt := time.Second / time.Duration(framesPerSecond)
		t := time.Now()
		t0 := t.Add(-dur)

		for {
			select {
			case <-p.quit:
				p.durCh <- dur
				return
			default:
			}

			p.c.Render(dur)
			glib.IdleAdd(p.da.QueueDraw)

			now := time.Now()
			dur = now.Sub(t0)

			t = t.Add(dt)
			d := t.Sub(now)
			if d > 0 {
				time.Sleep(d)
			}
		}
	}(p.td)

	p.started = true

	return nil
}

func (p *Player) Stop() error {

	p.guard.Lock()
	defer p.guard.Unlock()

	if !p.started {
		return errors.New("is stopped")
	}

	p.quit <- true
	p.td = <-p.durCh

	p.started = false

	return nil
}
