package main

type CellState struct {
	Row int
	Col int

	IsRevealed       bool
	IsMine           bool
	IsFlagged        bool
	NumAdjacentMines uint8
}

func (c *CellState) IsEmpty() bool {
	return !c.IsMine && c.NumAdjacentMines == 0
}
