package main

import (
	"flag"
	"log"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// https://en.wikipedia.org/wiki/Lissajous_curve

func main() {

	var config Config
	var tailDuration string

	flag.Float64Var(&(config.FreqA), "freqA", 2, "frequency A")
	flag.Float64Var(&(config.FreqB), "freqB", 3, "frequency B")
	flag.Float64Var(&(config.Phase), "phase", 0, "phase difference [0..1]")
	flag.StringVar(&(tailDuration), "dur", "1s", "tail duration")
	flag.IntVar(&(config.TailSegments), "seg", 100, "tail segments")

	flag.Parse()

	d, err := time.ParseDuration(tailDuration)
	checkError(err)
	config.TailDuration = d

	checkError(run(config))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run(config Config) error {

	core, err := NewCore(config)
	if err != nil {
		return err
	}

	gtk.Init(nil)

	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return err
	}

	da, err := makeDrawingArea(core)
	if err != nil {
		return err
	}

	player := NewPlayer(core, da)
	player.Start()

	w.Add(da)

	w.Connect("destroy", func() {
		gtk.MainQuit()
	})

	w.SetTitle("Lissajous curve")
	w.SetSizeRequest(720, 480)
	w.SetPosition(gtk.WIN_POS_CENTER)

	w.ShowAll()

	gtk.Main()

	return nil
}

func makeDrawingArea(core *Core) (*gtk.DrawingArea, error) {

	da, err := gtk.DrawingAreaNew()
	if err != nil {
		return nil, err
	}

	da.Connect("configure-event", func(da *gtk.DrawingArea, event *gdk.Event) {
		var (
			w = da.GetAllocatedWidth()
			h = da.GetAllocatedHeight()
		)
		core.Resize(w, h)
	})

	da.Connect("draw", func(da *gtk.DrawingArea, c *cairo.Context) {
		core.Draw(c)
	})

	return da, nil
}
