package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// GermanWhispersConstraint ensures adjacent values differ by at least 5
type GermanWhispersConstraint struct {
	lib.BaseConstraint
}

func NewGermanWhispersConstraint(cells []int) (*GermanWhispersConstraint, error) {
	if len(cells) < 2 {
		return nil, fmt.Errorf("german whispers constraint must have at least two cells")
	}

	for _, cell := range cells {
		if cell < 0 || cell > 80 {
			return nil, fmt.Errorf("invalid cell index: %d (must be 0-80)", cell)
		}
	}

	return &GermanWhispersConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  "German Whispers",
		},
	}, nil
}

func (gw *GermanWhispersConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	cells := gw.GetCells()
	for i := 0; i < len(cells)-1; i++ {
		cellIdx1 := cells[i]
		cellIdx2 := cells[i+1]

		row1, col1 := cellIdx1/9, cellIdx1%9
		row2, col2 := cellIdx2/9, cellIdx2%9

		val1 := board.Get(row1, col1)
		val2 := board.Get(row2, col2)

		// Skip if either cell is empty
		if val1 == 0 || val2 == 0 {
			continue
		}

		diff := val1 - val2
		if diff < 0 {
			diff = -diff
		}

		if diff < 5 {
			return false, nil
		}
	}

	return true, nil
}

func (gw *GermanWhispersConstraint) GetDescription() string {
	return fmt.Sprintf("German whispers line with %d cells - adjacent values must differ by at least 5", len(gw.GetCells()))
}

// PropagateValueChange propagates the value change to adjacent cells in the German Whispers line
// This is called automatically via the observer pattern when a cell is solved
func (gw *GermanWhispersConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if gw.Board == nil {
		return
	}

	cells := gw.GetCells()
	cellIndex := row*9 + col

	// Find the position of this cell in the constraint
	var pos int = -1
	for i, idx := range cells {
		if idx == cellIndex {
			pos = i
			break
		}
	}

	if pos == -1 {
		return // Cell not in this constraint
	}

	// Update previous cell if it exists
	if pos > 0 {
		prevIdx := cells[pos-1]
		prevRow, prevCol := prevIdx/9, prevIdx%9
		prevCell := gw.Board.GetCellAt(prevRow, prevCol)
		if prevCell != nil && !prevCell.IsSolved() {
			// Remove candidates that would violate the constraint
			for candidate := 1; candidate <= 9; candidate++ {
				diff := candidate - value
				if diff < 0 {
					diff = -diff
				}
				if diff < 5 {
					prevCell.RemoveCandidate(candidate)
				}
			}
		}
	}

	// Update next cell if it exists
	if pos < len(cells)-1 {
		nextIdx := cells[pos+1]
		nextRow, nextCol := nextIdx/9, nextIdx%9
		nextCell := gw.Board.GetCellAt(nextRow, nextCol)
		if nextCell != nil && !nextCell.IsSolved() {
			// Remove candidates that would violate the constraint
			for candidate := 1; candidate <= 9; candidate++ {
				diff := candidate - value
				if diff < 0 {
					diff = -diff
				}
				if diff < 5 {
					nextCell.RemoveCandidate(candidate)
				}
			}
		}
	}
}

func (gw *GermanWhispersConstraint) RequiresUniqueness() bool {
	// German Whispers doesn't enforce uniqueness by itself
	return false
}

func (gw *GermanWhispersConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// German Whispers doesn't enforce uniqueness, so pencil mark techniques don't apply
	return false
}
