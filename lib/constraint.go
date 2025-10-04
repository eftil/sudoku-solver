package lib

import (
	"github.com/eftil/sudoku-solver.git/lib/logger"
	"github.com/eftil/sudoku-solver.git/lib/observer"
	"github.com/eftil/sudoku-solver.git/lib/utils"
)

// Constraint is the interface that all constraint types must implement
// Constraints now act as observers of their cells
type Constraint interface {
	observer.CellObserver // Constraints observe their cells

	// GetCells returns the cell indices (0-80) that are part of this constraint
	GetCells() []int

	// IsValid checks if the constraint is currently satisfied on the board
	IsValid(board *Board) (bool, error)

	// GetName returns a human-readable name for this constraint
	GetName() string

	// GetDescription returns a detailed description of the constraint
	GetDescription() string

	// PropagateValueChange propagates constraint effects when a cell value changes
	// This is called automatically via the observer pattern
	PropagateValueChange(row, col, value int)

	// ApplyPencilMarkConstraints applies advanced techniques like naked pairs/triples
	// Returns true if any candidates were eliminated
	ApplyPencilMarkConstraints(board *Board) bool

	// RequiresUniqueness returns true if this constraint enforces uniqueness
	// (used to determine if pencil mark techniques apply)
	RequiresUniqueness() bool
}

// BaseConstraint provides common functionality for all constraints
type BaseConstraint struct {
	Cells []int
	Name  string
	Board *Board // Exported so embedded constraints can access it
}

func (bc *BaseConstraint) GetCells() []int {
	return bc.Cells
}

func (bc *BaseConstraint) GetName() string {
	return bc.Name
}

// SetBoard sets the board reference for this constraint
func (bc *BaseConstraint) SetBoard(board *Board) {
	bc.Board = board
}

// PropagateValueChange is called when a cell value changes (via observer pattern)
// Subclasses should override this to implement specific propagation logic
func (bc *BaseConstraint) PropagateValueChange(row, col, value int) {
	// Base implementation does nothing
	logger.Debug("BaseConstraint: PropagateValueChange called for R%dC%d = %d", row+1, col+1, value)
}

// OnCellSolved is called when a cell is solved (observer interface)
func (bc *BaseConstraint) OnCellSolved(row, col, value int) {
	bc.PropagateValueChange(row, col, value)
}

// OnSingleCandidate is called when a cell has only one candidate (observer interface)
func (bc *BaseConstraint) OnSingleCandidate(row, col, candidate int) {
	// Base implementation doesn't need to do anything
}

// OnCandidateEliminated is called when a candidate is eliminated (observer interface)
func (bc *BaseConstraint) OnCandidateEliminated(row, col, candidate, remainingCount int) {
	// Base implementation doesn't need to do anything
}

func (bc *BaseConstraint) RequiresUniqueness() bool {
	// Base constraint doesn't require uniqueness by default
	return false
}

func (bc *BaseConstraint) ApplyPencilMarkConstraints(board *Board) bool {
	// Base implementation does nothing
	return false
}

// HasUniqueNonZeros checks if all non-zero values in a slice are unique
// Deprecated: Use utils.HasUniqueNonZeros instead
func HasUniqueNonZeros(values []int) bool {
	return utils.HasUniqueNonZeros(values)
}

// ApplyNakedSubsets implements the naked pairs/triples/quads technique
// When n unsolved cells in a constraint collectively have exactly n candidates,
// those candidates can be eliminated from all other cells in the constraint
func ApplyNakedSubsets(board *Board, cellIndices []int, maxSubsetSize int) bool {
	if board == nil || len(cellIndices) == 0 {
		return false
	}

	logger.Debug("Applying naked subsets (max size: %d) to %d cells", maxSubsetSize, len(cellIndices))
	changed := false

	// Get all unsolved cells
	unsolvedCells := make([]*Cell, 0)
	for _, idx := range cellIndices {
		cell := board.GetCell(idx)
		if cell != nil && !cell.IsSolved() {
			unsolvedCells = append(unsolvedCells, cell)
		}
	}

	if len(unsolvedCells) < 2 {
		return false
	}

	// Try subset sizes from 2 up to maxSubsetSize (or number of unsolved cells)
	maxSize := maxSubsetSize
	if len(unsolvedCells) < maxSize {
		maxSize = len(unsolvedCells)
	}

	for subsetSize := 2; subsetSize <= maxSize; subsetSize++ {
		// Generate all combinations of subsetSize cells
		combinations := utils.GenerateCombinations(len(unsolvedCells), subsetSize)

		for _, combo := range combinations {
			// Get the union of candidates for this subset
			candidateUnion := make(map[int]bool)
			subsetCells := make([]*Cell, 0, subsetSize)

			for _, idx := range combo {
				cell := unsolvedCells[idx]
				subsetCells = append(subsetCells, cell)
				candidates := cell.GetCandidates()
				for candidate := range candidates {
					candidateUnion[candidate] = true
				}
			}

			// If the union has exactly subsetSize candidates, we found a naked subset
			if len(candidateUnion) == subsetSize {
				logger.Debug("Found naked subset of size %d with candidates: %v",
					subsetSize, utils.GetCandidatesAsSlice(candidateUnion))

				// Remove these candidates from all cells NOT in the subset
				eliminatedCount := 0
				for _, cell := range unsolvedCells {
					if !contains(subsetCells, cell) {
						for candidate := range candidateUnion {
							if cell.HasCandidate(candidate) {
								cell.RemoveCandidate(candidate)
								changed = true
								eliminatedCount++
							}
						}
					}
				}

				if eliminatedCount > 0 {
					logger.Info("Naked subset eliminated %d candidate(s)", eliminatedCount)
				}
			}
		}
	}

	return changed
}

// ApplyHiddenSubsets implements the hidden pairs/triples/quads technique
// When n candidates appear in exactly n cells (and nowhere else in the constraint),
// those cells can't contain any other candidates
func ApplyHiddenSubsets(board *Board, cellIndices []int, maxSubsetSize int) bool {
	if board == nil || len(cellIndices) == 0 {
		return false
	}

	logger.Debug("Applying hidden subsets (max size: %d) to %d cells", maxSubsetSize, len(cellIndices))
	changed := false

	// Get all unsolved cells
	unsolvedCells := make([]*Cell, 0)
	for _, idx := range cellIndices {
		cell := board.GetCell(idx)
		if cell != nil && !cell.IsSolved() {
			unsolvedCells = append(unsolvedCells, cell)
		}
	}

	if len(unsolvedCells) < 2 {
		return false
	}

	// Build a map of candidate -> cells that have it
	candidateLocations := make(map[int][]*Cell)
	for candidate := 1; candidate <= 9; candidate++ {
		candidateLocations[candidate] = make([]*Cell, 0)
	}

	for _, cell := range unsolvedCells {
		candidates := cell.GetCandidates()
		for candidate := range candidates {
			candidateLocations[candidate] = append(candidateLocations[candidate], cell)
		}
	}

	// Get candidates that appear in at least 2 cells
	activeCandidates := make([]int, 0)
	for candidate := 1; candidate <= 9; candidate++ {
		if len(candidateLocations[candidate]) >= 2 {
			activeCandidates = append(activeCandidates, candidate)
		}
	}

	if len(activeCandidates) < 2 {
		return false
	}

	maxSize := maxSubsetSize
	if len(activeCandidates) < maxSize {
		maxSize = len(activeCandidates)
	}

	// Try subset sizes from 2 up to maxSize
	for subsetSize := 2; subsetSize <= maxSize; subsetSize++ {
		combinations := utils.GenerateCombinations(len(activeCandidates), subsetSize)

		for _, combo := range combinations {
			// Get the candidates in this subset
			subsetCandidates := make([]int, 0, subsetSize)
			for _, idx := range combo {
				subsetCandidates = append(subsetCandidates, activeCandidates[idx])
			}

			// Get the union of cells that contain any of these candidates
			cellUnion := make(map[*Cell]bool)
			for _, candidate := range subsetCandidates {
				for _, cell := range candidateLocations[candidate] {
					cellUnion[cell] = true
				}
			}

			// If exactly subsetSize cells contain these candidates, it's a hidden subset
			if len(cellUnion) == subsetSize {
				logger.Debug("Found hidden subset of size %d with candidates: %v",
					subsetSize, subsetCandidates)

				// These cells can only contain these candidates
				eliminatedCount := 0
				for cell := range cellUnion {
					for candidate := 1; candidate <= 9; candidate++ {
						if !utils.ContainsInt(subsetCandidates, candidate) {
							if cell.HasCandidate(candidate) {
								cell.RemoveCandidate(candidate)
								changed = true
								eliminatedCount++
							}
						}
					}
				}

				if eliminatedCount > 0 {
					logger.Info("Hidden subset eliminated %d candidate(s)", eliminatedCount)
				}
			}
		}
	}

	return changed
}

// Helper function to check if a cell is in a slice of cells
func contains(cells []*Cell, target *Cell) bool {
	for _, cell := range cells {
		if cell == target {
			return true
		}
	}
	return false
}
