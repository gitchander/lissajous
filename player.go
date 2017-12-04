package main

import (
	"errors"
	"sync"
	"time"

	"github.com/envoker/gotk3/glib"
	"github.com/envoker/gotk3/gtk"
)

type Player struct {
	mutex   sync.Mutex
	da      *gtk.DrawingArea
	c       *Canvas
	quit    chan bool
	started bool
}

func NewPlayer() *Player {
	return &Player{
		quit: make(chan bool),
	}
}

func (p *Player) SetCurve(c *Canvas) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.c = c
}

func (p *Player) SetDrawingArea(da *gtk.DrawingArea) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.da = da
}

func (p *Player) Start() error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.c == nil {
		return errors.New("Curve is not set")
	}
	if p.da == nil {
		return errors.New("DrawingArea is not set")
	}

	if p.started {
		return errors.New("is started")
	}

	go func() {
		d := time.Second / 50
		startTime := time.Now()
		for {
			select {
			case <-p.quit:
				p.quit <- true
				return
			case <-time.After(d):
				p.c.Render(time.Now().Sub(startTime))
				glib.IdleAdd(p.da.QueueDraw)
			}
		}
	}()

	p.started = true

	return nil
}

func (p *Player) Stop() error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.started {
		return errors.New("is stopped")
	}

	p.quit <- true
	<-p.quit

	p.started = false

	return nil
}

func (p *Player) Stoped() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return !p.started
}
