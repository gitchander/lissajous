package main

import (
	"log"

	"github.com/envoker/gotk3/cairo"
	"github.com/envoker/gotk3/gdk"
	//"github.com/envoker/gotk3/glib"
	"github.com/envoker/gotk3/gtk"
)

func main() {

	canvas := NewCanvas(configSamples[2])

	player := NewPlayer()
	player.SetCurve(canvas)

	gtk.Init(nil)

	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	addControls(box, canvas, player)
	addDrawingArea(box, canvas, player)

	window.Add(box)

	window.Connect("destroy", gtk.MainQuit)

	window.SetTitle("Lissajous curves")
	window.SetSizeRequest(610, 377)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.ShowAll()

	gtk.Main()
}

func addControls(box *gtk.Box, canvas *Canvas, player *Player) {

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	grid, _ := gtk.GridNew()

	grid.SetBorderWidth(10)
	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(15)

	configEntrys := NewConfigEntrys()
	configEntrys.SetConfig(canvas.Config())

	// Frequency X
	grid.Attach(newLabel("Frequency X"), 0, 0, 1, 1)
	grid.Attach(configEntrys.entryFrequencyX, 1, 0, 1, 1)

	// Frequency Y
	grid.Attach(newLabel("Frequency Y"), 0, 1, 1, 1)
	grid.Attach(configEntrys.entryFrequencyY, 1, 1, 1, 1)

	// PhaseShift
	grid.Attach(newLabel("Phase Shift"), 0, 2, 1, 1)
	grid.Attach(configEntrys.entryPhaseShift, 1, 2, 1, 1)

	// TailDuration
	grid.Attach(newLabel("Tail Duration (sec)"), 0, 3, 1, 1)
	grid.Attach(configEntrys.entryTailDuration, 1, 3, 1, 1)

	// TailSegments
	grid.Attach(newLabel("Tail Segments"), 0, 4, 1, 1)
	grid.Attach(configEntrys.entryResolution, 1, 4, 1, 1)

	vbox.Add(grid)

	// start stop panel
	{
		hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)

		hbox.SetBorderWidth(10)

		// buttonApply
		buttonApply, _ := gtk.ButtonNewWithLabel("Apply")
		buttonApply.Connect("clicked", func() {
			var (
				config Config
				err    error
			)
			if err := configEntrys.GetConfig(&config); err != nil {
				errorMessageBox(err.Error())
				return
			}
			if err = canvas.SetConfig(config); err != nil {
				errorMessageBox(err.Error())
				return
			}
			/*
				if player.Stoped() {
					canvas.Render(0)
				}
			*/
		})
		hbox.PackStart(buttonApply, true, true, 0)

		// buttonStart
		buttonStart, _ := gtk.ButtonNewWithLabel("Start")
		buttonStart.Connect("clicked", func() {
			if err := player.Start(); err != nil {
				log.Println(err)
			}
		})
		hbox.PackStart(buttonStart, true, true, 0)

		// buttonStop
		buttonStop, _ := gtk.ButtonNewWithLabel("Stop")
		buttonStop.Connect("clicked", func() {
			if err := player.Stop(); err != nil {
				log.Println(err)
			}
		})
		hbox.PackStart(buttonStop, true, true, 0)

		vbox.PackEnd(hbox, false, false, 0)
	}

	box.Add(vbox)
}

func addDrawingArea(box *gtk.Box, canvas *Canvas, player *Player) {

	drawingArea, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal(err.Error())
	}

	drawingArea.Connect("configure-event", func(da *gtk.DrawingArea, event *gdk.Event) {

		var (
			width  = da.GetAllocatedWidth()
			height = da.GetAllocatedHeight()
		)

		if err := canvas.Resize(width, height); err != nil {
			log.Fatal(err)
		}

		if player.Stoped() {
			//glib.IdleAdd(da.QueueDraw)
			canvas.Render(0)
		}
	})

	drawingArea.Connect("draw", func(da *gtk.DrawingArea, c *cairo.Context) {
		canvas.Draw(c)
	})

	box.PackEnd(drawingArea, true, true, 0)

	player.SetDrawingArea(drawingArea)

	//player.Start()
}

func newLabel(name string) *gtk.Label {
	p, _ := gtk.LabelNew(name)
	return p
}

func newEntry() *gtk.Entry {
	entry, _ := gtk.EntryNew()
	return entry
}

func errorMessageBox(format string, a ...interface{}) {
	dialog := gtk.MessageDialogNew(nil,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		format, a...,
	)
	dialog.Run()
	dialog.Destroy()
}
