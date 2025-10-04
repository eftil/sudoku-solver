package lib

import (
	"github.com/eftil/sudoku-solver.git/lib/logger"
	"github.com/eftil/sudoku-solver.git/lib/observer"
	"github.com/eftil/sudoku-solver.git/lib/utils"
)

type Cell struct {
	row        int
	col        int
	index      int
	value      int
	board      *Board
	candidates map[int]bool
	notifier   *observer.CellNotifier
}

func NewCell(row, col int, board *Board) *Cell {
	candidates := make(map[int]bool)
	for i := 1; i <= 9; i++ {
		candidates[i] = true
	}

	logger.DebugCell(row, col, "Cell created with all candidates available")

	return &Cell{
		row:        row,
		col:        col,
		index:      row*9 + col,
		board:      board,
		candidates: candidates,
		value:      0,
		notifier:   observer.NewCellNotifier(),
	}
}

func (c *Cell) GetIndex() int {
	return c.index
}

func (c *Cell) GetRow() int {
	return c.row
}

func (c *Cell) GetCol() int {
	return c.col
}

func (c *Cell) GetValue() int {
	return c.value
}

func (c *Cell) GetBoard() *Board {
	return c.board
}

func (c *Cell) SetValue(value int) error {
	if value < 0 || value > 9 {
		logger.Error("Cell R%dC%d: Invalid value %d (must be 0-9)", c.row+1, c.col+1, value)
		return &BoardError{Message: "value must be between 0 and 9"}
	}

	oldValue := c.value
	c.value = value

	if value != 0 {
		logger.InfoCell(c.row, c.col, "Value set to %d (previous: %d)", value, oldValue)

		// Clear candidates when a value is set
		c.candidates = make(map[int]bool)

		// Notify observers that cell is solved (including constraints!)
		// This automatically propagates to all constraints via the observer pattern
		if c.notifier != nil {
			c.notifier.NotifyCellSolved(c.row, c.col, value)
		}

		logger.DebugCell(c.row, c.col, "Notified observers about value %d", value)
	} else {
		logger.DebugCell(c.row, c.col, "Value cleared (was: %d)", oldValue)
	}

	return nil
}

// Note: AddConstraint and GetConstraints removed!
// Constraints are now observers and don't need to be tracked separately

// GetCandidates returns the current candidates for this cell
func (c *Cell) GetCandidates() map[int]bool {
	if c.value != 0 {
		return make(map[int]bool) // No candidates if value is set
	}
	return c.candidates
}

// RemoveCandidate removes a candidate from this cell
func (c *Cell) RemoveCandidate(candidate int) {
	if c.value == 0 && c.candidates[candidate] {
		delete(c.candidates, candidate)
		remainingCount := len(c.candidates)

		logger.DebugCell(c.row, c.col, "Removed candidate %d (remaining: %v)",
			candidate, utils.GetCandidatesAsSlice(c.candidates))

		// Notify observers
		if c.notifier != nil {
			c.notifier.NotifyCandidateEliminated(c.row, c.col, candidate, remainingCount)

			// If only one candidate remains, notify that too
			if remainingCount == 1 {
				lastCandidate := utils.GetCandidatesAsSlice(c.candidates)[0]
				logger.InfoCell(c.row, c.col, "Only one candidate remains: %d", lastCandidate)
				c.notifier.NotifySingleCandidate(c.row, c.col, lastCandidate)
			}
		}
	}
}

// AddCandidate adds a candidate to this cell
func (c *Cell) AddCandidate(candidate int) {
	if c.value == 0 && candidate >= 1 && candidate <= 9 {
		if !c.candidates[candidate] {
			c.candidates[candidate] = true
			logger.DebugCell(c.row, c.col, "Added candidate %d (total: %v)",
				candidate, utils.GetCandidatesAsSlice(c.candidates))
		}
	}
}

// HasCandidate checks if a candidate is available for this cell
func (c *Cell) HasCandidate(candidate int) bool {
	if c.value != 0 {
		return false
	}
	return c.candidates[candidate]
}

// IsSolved returns true if the cell has a value set
func (c *Cell) IsSolved() bool {
	return c.value != 0
}

// CandidateCount returns the number of candidates for this cell
func (c *Cell) CandidateCount() int {
	if c.value != 0 {
		return 0
	}
	return len(c.candidates)
}

// GetNotifier returns the cell's notifier for adding observers
func (c *Cell) GetNotifier() *observer.CellNotifier {
	return c.notifier
}

// AddObserver adds an observer to this cell's notifier
func (c *Cell) AddObserver(obs observer.CellObserver) {
	if c.notifier != nil {
		c.notifier.AddObserver(obs)
	}
}
