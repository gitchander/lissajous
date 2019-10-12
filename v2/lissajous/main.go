package main

import (
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"

	//"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {

	// config := Config{
	// 	FrequencyX:   200,
	// 	FrequencyY:   300.1,
	// 	PhaseShift:   0,
	// 	TailDuration: 0.01,
	// 	TailSegments: 50,
	// }

	config := configSamples[2]

	core := NewCore(config)

	// canvas := NewCanvas(configSamples[2])

	// player := NewPlayer()
	// player.SetCurve(canvas)

	gtk.Init(nil)

	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	checkError(err)

	da, err := makeDrawingArea(core)
	checkError(err)

	// redraw := func() {
	// 	glib.IdleAdd(func() bool {
	// 		da.QueueDraw()
	// 		//return true // continue execute this func
	// 		return false // if nead to stop this function
	// 	})
	// }

	player := NewPlayer(core, da)
	player.Start()

	w.Add(da)

	// box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	// addControls(box, canvas, player)
	// addDrawingArea(box, canvas, player)
	//w.Add(box)

	w.Connect("destroy", func() {
		gtk.MainQuit()
	})

	w.SetTitle("Lissajous curves")
	w.SetSizeRequest(610, 377)
	w.SetPosition(gtk.WIN_POS_CENTER)

	w.ShowAll()

	gtk.Main()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
