package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// ColumnConstraint ensures all values in a column are unique
type ColumnConstraint struct {
	lib.BaseConstraint
	col int
}

func NewColumnConstraint(col int) (*ColumnConstraint, error) {
	if col < 0 || col > 8 {
		return nil, fmt.Errorf("column must be between 0 and 8, got %d", col)
	}

	cells := make([]int, 9)
	for row := 0; row < 9; row++ {
		cells[row] = row*9 + col
	}

	return &ColumnConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  fmt.Sprintf("Column %d", col+1),
		},
		col: col,
	}, nil
}

func (cc *ColumnConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	colData := board.GetColumn(cc.col)
	return lib.HasUniqueNonZeros(colData[:]), nil
}

func (cc *ColumnConstraint) GetDescription() string {
	return fmt.Sprintf("All values in column %d must be unique (1-9)", cc.col+1)
}

// PropagateValueChange propagates the value change to other cells in the column
// This is called automatically via the observer pattern when a cell is solved
func (cc *ColumnConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if cc.Board == nil {
		return
	}

	// Remove the value from candidates of all other cells in this column
	for _, cellIndex := range cc.Cells {
		otherRow, otherCol := cellIndex/9, cellIndex%9
		if otherRow != row || otherCol != col {
			otherCell := cc.Board.GetCellAt(otherRow, otherCol)
			if otherCell != nil && !otherCell.IsSolved() {
				otherCell.RemoveCandidate(value)
			}
		}
	}
}

func (cc *ColumnConstraint) RequiresUniqueness() bool {
	return true
}

func (cc *ColumnConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// Apply both naked and hidden subset techniques up to quads (size 4)
	changed := false
	changed = lib.ApplyNakedSubsets(board, cc.Cells, 4) || changed
	changed = lib.ApplyHiddenSubsets(board, cc.Cells, 4) || changed
	return changed
}
