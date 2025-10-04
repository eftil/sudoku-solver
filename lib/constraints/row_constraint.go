package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// RowConstraint ensures all values in a row are unique
type RowConstraint struct {
	lib.BaseConstraint
	row int
}

func NewRowConstraint(row int) (*RowConstraint, error) {
	if row < 0 || row > 8 {
		return nil, fmt.Errorf("row must be between 0 and 8, got %d", row)
	}

	cells := make([]int, 9)
	for col := 0; col < 9; col++ {
		cells[col] = row*9 + col
	}

	return &RowConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  fmt.Sprintf("Row %d", row+1),
		},
		row: row,
	}, nil
}

func (rc *RowConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	rowData := board.GetRow(rc.row)
	return lib.HasUniqueNonZeros(rowData[:]), nil
}

func (rc *RowConstraint) GetDescription() string {
	return fmt.Sprintf("All values in row %d must be unique (1-9)", rc.row+1)
}

// PropagateValueChange propagates the value change to other cells in the row
// This is called automatically via the observer pattern when a cell is solved
func (rc *RowConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if rc.Board == nil {
		return
	}

	// Remove the value from candidates of all other cells in this row
	for _, cellIndex := range rc.Cells {
		otherRow, otherCol := cellIndex/9, cellIndex%9
		if otherRow != row || otherCol != col {
			otherCell := rc.Board.GetCellAt(otherRow, otherCol)
			if otherCell != nil && !otherCell.IsSolved() {
				otherCell.RemoveCandidate(value)
			}
		}
	}
}

func (rc *RowConstraint) RequiresUniqueness() bool {
	return true
}

func (rc *RowConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// Apply both naked and hidden subset techniques up to quads (size 4)
	changed := false
	changed = lib.ApplyNakedSubsets(board, rc.Cells, 4) || changed
	changed = lib.ApplyHiddenSubsets(board, rc.Cells, 4) || changed
	return changed
}
