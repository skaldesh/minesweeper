package main

import (
	"log"
	"math/rand"
	"sync/atomic"
)

func reducer(a Action) {
	// Copy state.
	s := state
	s.Cells = make([]CellState, len(state.Cells))
	copy(s.Cells, state.Cells)

	switch a.Type {
	case startGame:
		s = startGameReducer(s)
	case clickedCell:
		s = clickedCellReducer(s, a.Data.(clickedCellData))
	case resetGame:
		s = resetGameReducer(s)
	case decrementTimer:
		s = decrementTimerReducer(s)
	default:
		log.Fatal("unhandled action type", a)
	}

	// Update global state.
	state = s
}

func startGameReducer(s State) State {
	if s.Running || s.Finished {
		return s
	}

	s.Running = true

	atomic.StoreUint32(s.timerActive, 1)
	return s
}

func clickedCellReducer(s State, data clickedCellData) State {
	// Do nothing, if game is not running.
	if !s.Running || s.Finished {
		return s
	}

	cell := s.Cells[data.idx]

	if data.rightClick {
		// Handle right-click.
		if cell.IsRevealed || (!cell.IsFlagged && s.selectFlagsLeft() == 0) {
			return s
		}

		s.Cells[data.idx].IsFlagged = !cell.IsFlagged
	} else {
		// Ignore no-op left-click.
		if cell.IsFlagged || cell.IsRevealed {
			return s
		}

		// Left-click.
		s.Cells[data.idx].IsRevealed = true

		// Game lost, if it was a mine.
		if cell.IsMine {
			s.Running = false
			s.Finished = true
			s.Lost = true
			return s
		}

		// If the cell is empty, reveal all its neighbours.
		if cell.IsEmpty() {
			// Reveal all its neighbours.
			queue := make([]CellState, 0)
			visited := make(map[CellState]struct{})

			// Root of search, add immediately to both visited set and queue.
			queue = append(queue, cell)
			visited[cell] = struct{}{}

			for len(queue) > 0 {
				// Dequeue next element.
				cell, queue = queue[0], queue[1:]

				// Iterate over its neighbours and enqueue the not visited ones,
				// but only if they are hidden and not a mine.
				for _, idx := range neighbours(s, cell) {
					nc := s.Cells[idx]
					if _, ok := visited[nc]; !ok && !nc.IsRevealed && !nc.IsMine && !nc.IsFlagged {
						// Mark visited.
						visited[nc] = struct{}{}

						// Only if the cell is empty visit its neighbours.
						if nc.IsEmpty() {
							queue = append(queue, nc)
						}

						// Reveal the cell.
						s.Cells[idx].IsRevealed = true
					}
				}
			}
		}
	}

	// Check win condition.
	if s.selectFlagsLeft() == 0 && s.selectRemainingCells() == 0 {
		s.Running = false
		s.Finished = true
		s.Won = true
		return s
	}

	return s
}

func resetGameReducer(s State) State {
	s.Running = false
	s.Finished = false
	s.Won = false
	s.Lost = false
	s.NumMines = defaultMines
	s.SecondsLeft = defaultSecondsLeft

	atomic.StoreUint32(s.timerActive, 0)

	// Reset grid cells.
	s.Cells = make([]CellState, int(s.Cols)*int(s.Rows))
	for row := 0; row < int(s.Rows); row++ {
		for col := 0; col < int(s.Cols); col++ {
			idx := (row * int(s.Cols)) + col
			s.Cells[idx].Row, s.Cells[idx].Col = row, col
		}
	}

	// Setup game board.
	// Distribute the mines across all cells.
	numMines := s.NumMines
	for numMines != 0 {
		// Choose random index of a cell that is not yet a mine.
		idx := rand.Intn(int(s.Rows) * int(s.Cols))
		if s.Cells[idx].IsMine {
			continue
		}

		s.Cells[idx].IsMine = true
		numMines--
	}

	// Calculate the number of adjacent mines for each cell.
	for i := 0; i < len(s.Cells); i++ {
		if !s.Cells[i].IsMine {
			continue
		}

		for _, idx := range neighbours(s, s.Cells[i]) {
			if !s.Cells[idx].IsMine {
				s.Cells[idx].NumAdjacentMines++
			}
		}
	}
	return s
}

func neighbours(s State, c CellState) (idxs []int) {
	// Allocate space for up to 8 neighbours.
	idxs = make([]int, 0, 8)
	row, col := c.Row, c.Col

	if row != 0 {
		idxs = append(idxs, s.toIdx(row-1, col)) // Top.
		if col != 0 {
			idxs = append(idxs, s.toIdx(row-1, col-1)) // Top-left.
		}
		if col != int(s.Cols)-1 {
			idxs = append(idxs, s.toIdx(row-1, col+1)) // Top-right.
		}
	}

	if col != 0 {
		idxs = append(idxs, s.toIdx(row, col-1)) // Left.
	}
	if col != int(s.Cols)-1 {
		idxs = append(idxs, s.toIdx(row, col+1)) // Right.
	}

	if row != int(s.Rows)-1 {
		idxs = append(idxs, s.toIdx(row+1, col)) // Bottom.
		if col != 0 {
			idxs = append(idxs, s.toIdx(row+1, col-1)) // Bottom-left.
		}
		if col != int(s.Cols)-1 {
			idxs = append(idxs, s.toIdx(row+1, col+1)) // Bottom-right.
		}
	}

	return
}

func (s *State) toIdx(row, col int) int {
	return (row * int(s.Cols)) + col
}

func decrementTimerReducer(s State) State {
	if s.Finished || !s.Running {
		return s
	}

	s.SecondsLeft--
	if s.SecondsLeft == 0 {
		s.Running = false
		s.Finished = true
		s.Lost = true
		atomic.StoreUint32(s.timerActive, 0)
	}
	return s
}
