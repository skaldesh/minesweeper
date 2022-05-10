package main

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	colorBackground = color.NRGBA{R: 165, G: 165, B: 165, A: 255}
	colorHidden     = color.NRGBA{R: 195, G: 195, B: 195, A: 255}
	colorFlagged    = color.NRGBA{G: 170, A: 255}
	colorMine       = color.NRGBA{R: 255, A: 255}

	colorsLabels = []color.NRGBA{
		{B: 255, A: 255},
		{G: 255, A: 255},
		{R: 255, A: 255},
		{B: 255, G: 255, A: 255},
		{B: 255, R: 255, A: 255},
		{G: 255, R: 255, A: 255},
		{B: 80, G: 130, R: 208, A: 255},
		{B: 25, G: 90, R: 110, A: 255},
	}
)

type GridCell struct {
	th *material.Theme

	button widget.Clickable

	idx int
}

func NewGridCell(th *material.Theme, idx int) *GridCell {
	return &GridCell{
		th:  th,
		idx: idx,
	}
}

func (c *GridCell) Layout(gtx C) D {
	// Handle events.
	for _, e := range gtx.Events(c) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				ClickedCell(c.idx, e.Buttons.Contain(pointer.ButtonSecondary))
			}
		}
	}

	// Draw.
	return c.draw(gtx)
}

func (c *GridCell) draw(gtx C) layout.Dimensions {
	// Retrieve our state data.
	cellState := state.Cells[c.idx]

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx C) D { // Background.
			// Build a grid by giving each cell a bottom right border of one pixel.
			max := gtx.Constraints.Max
			max.X--
			max.Y--
			rect := clip.Rect{Max: max}

			// Push onto op stack to confine the following actions to this area.
			rectArea := rect.Push(gtx.Ops)

			// Define input area.
			pointer.InputOp{
				Tag:   c,
				Types: pointer.Press,
			}.Add(gtx.Ops)

			// Define color.
			paint.ColorOp{Color: colorBackground}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Finish ops on this area.
			rectArea.Pop()

			// Draw inner area to distinct between hidden state.
			if !cellState.IsRevealed {
				const inset = 2
				min := image.Pt(inset, inset)
				max.X -= inset
				max.Y -= inset
				rect = clip.Rect{Min: min, Max: max}
				rectArea = rect.Push(gtx.Ops)
				var color = colorHidden
				if cellState.IsFlagged {
					color = colorFlagged
				}
				paint.ColorOp{Color: color}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				rectArea.Pop()
			}

			return D{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx C) D { // Text.
			return layout.Center.Layout(gtx, func(gtx C) D {
				if !cellState.IsRevealed || cellState.IsEmpty() || cellState.IsMine {
					return D{}
				}

				label := material.H4(c.th, strconv.Itoa(int(cellState.NumAdjacentMines)))
				label.Alignment = text.Middle
				label.Color = colorsLabels[cellState.NumAdjacentMines-1]
				return label.Layout(gtx)
			})
		}),
		layout.Stacked(func(gtx C) D { // Mine.
			if !cellState.IsMine || (!cellState.IsRevealed && !state.Lost) {
				return D{}
			}

			return layout.Center.Layout(gtx, func(gtx C) D {
				const inset = 8
				min := float32(math.Min(float64(gtx.Constraints.Max.X), float64(gtx.Constraints.Max.Y))) / 2.0

				circle := clip.Ellipse{
					Min: f32.Pt(-min+inset, -min+inset),
					Max: f32.Pt(min-inset, min-inset),
				}.Op(gtx.Ops)
				paint.FillShape(gtx.Ops, colorMine, circle)

				return D{}
			})
		}),
	)
}
