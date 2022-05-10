package main

func (s State) selectFlagsLeft() (left uint16) {
	left = s.NumMines
	for _, c := range s.Cells {
		if c.IsFlagged {
			left--
		}
	}
	return
}

func (s State) selectRemainingCells() (hidden uint16) {
	for _, c := range s.Cells {
		if !c.IsRevealed && !c.IsMine {
			hidden++
		}
	}
	return
}
