package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// KillerCageConstraint ensures values sum to a target and are unique
type KillerCageConstraint struct {
	lib.BaseConstraint
	targetSum int
}

func NewKillerCageConstraint(cells []int, targetSum int) (*KillerCageConstraint, error) {
	if len(cells) == 0 {
		return nil, fmt.Errorf("killer cage must have at least one cell")
	}

	for _, cell := range cells {
		if cell < 0 || cell > 80 {
			return nil, fmt.Errorf("invalid cell index: %d (must be 0-80)", cell)
		}
	}

	if targetSum < 1 || targetSum > 45 {
		return nil, fmt.Errorf("target sum must be between 1 and 45, got %d", targetSum)
	}

	return &KillerCageConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  fmt.Sprintf("Killer Cage (%d)", targetSum),
		},
		targetSum: targetSum,
	}, nil
}

func (kc *KillerCageConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	cells := kc.GetCells()
	values := make([]int, len(cells))
	sum := 0
	hasEmpty := false

	for i, cellIdx := range cells {
		row := cellIdx / 9
		col := cellIdx % 9
		values[i] = board.Get(row, col)
		if values[i] == 0 {
			hasEmpty = true
		} else {
			sum += values[i]
		}
	}

	// Check uniqueness
	if !lib.HasUniqueNonZeros(values) {
		return false, nil
	}

	// If cage is complete, check sum
	if !hasEmpty {
		return sum == kc.targetSum, nil
	}

	// If incomplete, sum shouldn't exceed target
	return sum <= kc.targetSum, nil
}

func (kc *KillerCageConstraint) GetDescription() string {
	return fmt.Sprintf("Killer cage with %d cells - values must sum to %d and be unique", len(kc.GetCells()), kc.targetSum)
}

// PropagateValueChange propagates the value change to other cells in the killer cage
// This is called automatically via the observer pattern when a cell is solved
func (kc *KillerCageConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if kc.Board == nil {
		return
	}

	cells := kc.GetCells()
	cellIndex := row*9 + col

	// First: Remove the set value from all other cells (uniqueness constraint)
	for _, otherIndex := range cells {
		if otherIndex != cellIndex {
			otherRow, otherCol := otherIndex/9, otherIndex%9
			otherCell := kc.Board.GetCellAt(otherRow, otherCol)
			if otherCell != nil && !otherCell.IsSolved() {
				otherCell.RemoveCandidate(value)
			}
		}
	}

	// Second: Calculate current sum and apply sum constraints
	currentSum := 0
	filledCount := 0
	for _, idx := range cells {
		r, c := idx/9, idx%9
		otherCell := kc.Board.GetCellAt(r, c)
		if otherCell != nil && otherCell.GetValue() != 0 {
			currentSum += otherCell.GetValue()
			filledCount++
		}
	}

	remainingCells := len(cells) - filledCount
	remainingSum := kc.targetSum - currentSum

	// Update candidates for empty cells based on sum constraints
	for _, idx := range cells {
		r, c := idx/9, idx%9
		otherCell := kc.Board.GetCellAt(r, c)
		if otherCell != nil && otherCell.GetValue() == 0 {
			// Remove candidates that would violate sum constraint
			for candidate := 1; candidate <= 9; candidate++ {
				// Check if this candidate would make the sum impossible
				if remainingCells == 1 {
					// Last cell must equal remaining sum
					if candidate != remainingSum {
						otherCell.RemoveCandidate(candidate)
					}
				} else {
					// Check if remaining sum is achievable with remaining cells
					minPossibleSum := remainingCells - 1
					maxPossibleSum := (remainingCells - 1) * 9
					if remainingSum-candidate < minPossibleSum || remainingSum-candidate > maxPossibleSum {
						otherCell.RemoveCandidate(candidate)
					}
				}
			}
		}
	}
}

func (kc *KillerCageConstraint) RequiresUniqueness() bool {
	return true
}

func (kc *KillerCageConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// Apply both naked and hidden subset techniques
	// Use smaller max size for killer cages since they're often smaller than 9 cells
	maxSize := 4
	if len(kc.Cells) < maxSize {
		maxSize = len(kc.Cells)
	}

	changed := false
	changed = lib.ApplyNakedSubsets(board, kc.Cells, maxSize) || changed
	changed = lib.ApplyHiddenSubsets(board, kc.Cells, maxSize) || changed
	return changed
}
