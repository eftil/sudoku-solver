package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// RenbanConstraint ensures values form a consecutive set (no duplicates, consecutive when sorted)
type RenbanConstraint struct {
	lib.BaseConstraint
}

func NewRenbanConstraint(cells []int) (*RenbanConstraint, error) {
	if len(cells) == 0 {
		return nil, fmt.Errorf("renban constraint must have at least one cell")
	}

	for _, cell := range cells {
		if cell < 0 || cell > 80 {
			return nil, fmt.Errorf("invalid cell index: %d (must be 0-80)", cell)
		}
	}

	return &RenbanConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  "Renban Line",
		},
	}, nil
}

func (rc *RenbanConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	cells := rc.GetCells()
	values := make([]int, len(cells))
	hasEmpty := false

	for i, cellIdx := range cells {
		row := cellIdx / 9
		col := cellIdx % 9
		values[i] = board.Get(row, col)
		if values[i] == 0 {
			hasEmpty = true
		}
	}

	// If any cells are empty, we can't fully validate yet (return true for partial boards)
	if hasEmpty {
		return lib.HasUniqueNonZeros(values), nil
	}

	// Check uniqueness first
	if !lib.HasUniqueNonZeros(values) {
		return false, nil
	}

	// Sort and check if consecutive
	sorted := make([]int, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1]+1 {
			return false, nil
		}
	}

	return true, nil
}

func (rc *RenbanConstraint) GetDescription() string {
	return fmt.Sprintf("Renban line with %d cells - values must form a consecutive set with no gaps or repeats", len(rc.GetCells()))
}

// PropagateValueChange propagates the value change to other cells in the Renban line
// This is called automatically via the observer pattern when a cell is solved
func (rc *RenbanConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if rc.Board == nil {
		return
	}

	cells := rc.GetCells()
	cellIndex := row*9 + col

	// First: Remove the set value from all other cells (uniqueness constraint)
	for _, otherIndex := range cells {
		if otherIndex != cellIndex {
			otherRow, otherCol := otherIndex/9, otherIndex%9
			otherCell := rc.Board.GetCellAt(otherRow, otherCol)
			if otherCell != nil && !otherCell.IsSolved() {
				otherCell.RemoveCandidate(value)
			}
		}
	}

	// Second: Get all current values in the constraint
	currentValues := make([]int, 0, len(cells))
	for _, idx := range cells {
		r, c := idx/9, idx%9
		otherCell := rc.Board.GetCellAt(r, c)
		if otherCell != nil && otherCell.GetValue() != 0 {
			currentValues = append(currentValues, otherCell.GetValue())
		}
	}

	// Determine valid range for consecutive values
	minVal := value
	maxVal := value
	for _, val := range currentValues {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	totalCells := len(cells)

	// Update candidates for empty cells based on consecutive constraint
	for _, idx := range cells {
		r, c := idx/9, idx%9
		otherCell := rc.Board.GetCellAt(r, c)
		if otherCell != nil && otherCell.GetValue() == 0 {
			// Remove candidates that would break consecutive constraint
			for candidate := 1; candidate <= 9; candidate++ {
				// Check if this candidate would maintain consecutive property
				newMin := minVal
				newMax := maxVal
				if candidate < newMin {
					newMin = candidate
				}
				if candidate > newMax {
					newMax = candidate
				}

				// Check if the range would be too large for remaining cells
				rangeSize := newMax - newMin + 1
				if rangeSize > totalCells {
					otherCell.RemoveCandidate(candidate)
				}
			}
		}
	}
}

func (rc *RenbanConstraint) RequiresUniqueness() bool {
	return true
}

func (rc *RenbanConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// Apply both naked and hidden subset techniques
	// Use smaller max size since renban constraints are often smaller than 9 cells
	maxSize := 4
	if len(rc.Cells) < maxSize {
		maxSize = len(rc.Cells)
	}

	changed := false
	changed = lib.ApplyNakedSubsets(board, rc.Cells, maxSize) || changed
	changed = lib.ApplyHiddenSubsets(board, rc.Cells, maxSize) || changed
	return changed
}
