package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

func init() {
	rand.Seed(time.Now().UnixMicro())
}

var window *app.Window

func main() {
	go func() {
		// Create the window
		window = app.NewWindow(
			app.Title("Minesweeper"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)

		if err := draw(window); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func draw(w *app.Window) error {
	// Ops are the operations from the UI.
	var ops op.Ops

	// th defines the material design theme
	th := material.NewTheme(gofont.Collection())

	bar := NewBar(th)
	grid := NewGrid(th)

	// Handle window events.
	for e := range w.Events() {
		switch e := e.(type) {
		// This is sent when the application should re-render.
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return bar.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					return grid.Layout(gtx)
				}),
			)
			e.Frame(gtx.Ops)

		case system.DestroyEvent:
			return e.Err
		}
	}

	return nil
}
