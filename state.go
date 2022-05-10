package main

import (
	"sync/atomic"
	"time"
)

const (
	defaultMines       = 20
	defaultSecondsLeft = 180
	defaultGridRows    = 20
	defaultGridCols    = 15
)

var state State

func init() {
	state.timerActive = new(uint32)
	state.Rows = defaultGridRows
	state.Cols = defaultGridCols
	ResetGame()

	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			// Wait for next tick.
			<-t.C

			// Check, if timer should stop.
			if atomic.LoadUint32(state.timerActive) == 0 {
				continue
			}

			// Decrement our timer and invalidate the window to trigger a redraw.
			window.Run(DecrementTimer)
			window.Invalidate()
		}
	}()
}

type State struct {
	Running  bool
	Finished bool
	Won      bool
	Lost     bool

	NumMines    uint16
	SecondsLeft uint16

	Rows  uint8
	Cols  uint8
	Cells []CellState

	// Atomic
	timerActive *uint32
}

func dispatch(a Action) {
	reducer(a)
}

type ActionType int

const (
	startGame ActionType = iota
	clickedCell
	resetGame
	decrementTimer
)

type Action struct {
	Type ActionType
	Data interface{}
}

func StartGame() {
	dispatch(Action{Type: startGame})
}

type clickedCellData struct {
	idx        int
	rightClick bool
}

func ClickedCell(idx int, rightClick bool) {
	dispatch(Action{
		Type: clickedCell,
		Data: clickedCellData{idx: idx, rightClick: rightClick},
	})
}

func ResetGame() {
	dispatch(Action{Type: resetGame})
}

func DecrementTimer() {
	dispatch(Action{Type: decrementTimer})
}
