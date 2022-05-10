package main

import (
	"image"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	colorBarBackground = color.NRGBA{R: 211, G: 211, B: 211, A: 255}
	colorBarTextWon    = color.NRGBA{R: 8, G: 229, A: 255}
	colorBarTextLost   = color.NRGBA{R: 202, G: 44, A: 255}
)

type Bar struct {
	th *material.Theme

	restart widget.Clickable
}

func NewBar(th *material.Theme) *Bar {
	return &Bar{
		th: th,
	}
}

func (b *Bar) Layout(gtx C) D {
	// Handle events.
	if b.restart.Clicked() {
		if state.Running || state.Finished {
			ResetGame()
		} else {
			StartGame()
		}
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Max.Y = 80

			const inset = 3
			min := image.Pt(inset, inset)
			max := gtx.Constraints.Max
			max.X -= inset
			max.Y -= inset
			rect := clip.Rect{Min: min, Max: max}
			paint.FillShape(gtx.Ops, colorBarBackground, rect.Op())

			return D{Size: gtx.Constraints.Max}
		}),
		layout.Expanded(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: unit.Dp(16),
					}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return material.H4(b.th, "⏲️ "+strconv.FormatUint(uint64(state.SecondsLeft), 10)).Layout(gtx)
						})
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						if state.Running || !state.Finished {
							return D{}
						}

						margin := layout.Inset{Left: unit.Dp(40)}

						color := colorBarTextWon
						if state.Lost {
							color = colorBarTextLost
						}

						label := material.H4(b.th, "You")
						label.Alignment = text.End
						label.Color = color

						return margin.Layout(gtx, label.Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						t := "Start"
						if state.Running || state.Finished {
							t = "Reset"
						}

						return material.Button(b.th, &b.restart, t).Layout(gtx)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						if state.Running || !state.Finished {
							return D{}
						}

						var (
							t     = "Won"
							color = colorBarTextWon
						)
						if state.Lost {
							t = "Lost"
							color = colorBarTextLost
						}

						margin := layout.Inset{Right: unit.Dp(40)}

						label := material.H4(b.th, t)
						label.Color = color

						return margin.Layout(gtx, label.Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Right: unit.Dp(16),
					}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return material.H4(b.th, strconv.FormatUint(uint64(state.selectFlagsLeft()), 10)+" 💣").Layout(gtx)
						})
					})
				}),
			)
		}),
	)
}
