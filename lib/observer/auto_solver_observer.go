package observer

import (
	"fmt"
)

// AutoSolverObserver automatically sets cell values when only one candidate remains
type AutoSolverObserver struct {
	enabled       bool
	cellsToSolve  map[string]int // Map of "row,col" -> value
	solutionCount int
}

// NewAutoSolverObserver creates a new auto-solver observer
func NewAutoSolverObserver() *AutoSolverObserver {
	return &AutoSolverObserver{
		enabled:       true,
		cellsToSolve:  make(map[string]int),
		solutionCount: 0,
	}
}

// OnSingleCandidate is called when a cell has only one candidate remaining
func (aso *AutoSolverObserver) OnSingleCandidate(row, col, candidate int) {
	if !aso.enabled {
		return
	}

	key := fmt.Sprintf("%d,%d", row, col)
	aso.cellsToSolve[key] = candidate
	fmt.Printf("📝 Observer detected: Cell R%dC%d can be auto-solved with value %d\n",
		row+1, col+1, candidate)
}

// OnCellSolved is called when a cell's value is set
func (aso *AutoSolverObserver) OnCellSolved(row, col, value int) {
	if !aso.enabled {
		return
	}

	aso.solutionCount++
	fmt.Printf("✓ Cell R%dC%d solved with value %d (Total solved: %d)\n",
		row+1, col+1, value, aso.solutionCount)

	// Remove from cellsToSolve if it was there
	key := fmt.Sprintf("%d,%d", row, col)
	delete(aso.cellsToSolve, key)
}

// OnCandidateEliminated is called when a candidate is removed from a cell
func (aso *AutoSolverObserver) OnCandidateEliminated(row, col, candidate, remainingCount int) {
	// This observer doesn't need to do anything for general candidate elimination
	// (OnSingleCandidate will be called when count reaches 1)
}

// GetCellsToSolve returns the cells that have been identified as solvable
func (aso *AutoSolverObserver) GetCellsToSolve() map[string]int {
	return aso.cellsToSolve
}

// ClearCellsToSolve clears the list of cells to solve
func (aso *AutoSolverObserver) ClearCellsToSolve() {
	aso.cellsToSolve = make(map[string]int)
}

// GetSolutionCount returns the total number of cells solved
func (aso *AutoSolverObserver) GetSolutionCount() int {
	return aso.solutionCount
}

// Enable enables the observer
func (aso *AutoSolverObserver) Enable() {
	aso.enabled = true
}

// Disable disables the observer
func (aso *AutoSolverObserver) Disable() {
	aso.enabled = false
}

// IsEnabled returns whether the observer is enabled
func (aso *AutoSolverObserver) IsEnabled() bool {
	return aso.enabled
}
