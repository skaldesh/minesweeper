package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Grid struct {
	th *material.Theme

	rowChildren []layout.FlexChild
}

func NewGrid(th *material.Theme) *Grid {
	// Create the grid.
	g := &Grid{
		th: th,
	}

	// Create the cells themselves and their layout children.
	//g.cells = make([]*GridCell, 0, g.rows*g.cols)
	g.rowChildren = make([]layout.FlexChild, state.Rows)
	for row := 0; row < int(state.Rows); row++ {
		colChildren := make([]layout.FlexChild, state.Cols)
		for col := 0; col < int(state.Cols); col++ {
			cell := NewGridCell(th, (row*int(state.Cols))+col)
			//g.cells = append(g.cells, cell)
			colChildren[col] = layout.Flexed(1, cell.Layout)
		}

		g.rowChildren[row] = layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceEvenly,
			}.Layout(gtx, colChildren...)
		})
	}

	return g
}

func (g *Grid) Layout(gtx C) D {
	if !state.Running {
		// Disable events if game is not running.
		gtx = gtx.Disabled()
	}

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEvenly,
	}.Layout(gtx, g.rowChildren...)
}
