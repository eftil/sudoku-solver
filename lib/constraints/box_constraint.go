package constraints

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib"
)

// BoxConstraint ensures all values in a 3x3 box are unique
type BoxConstraint struct {
	lib.BaseConstraint
	box int
}

func NewBoxConstraint(box int) (*BoxConstraint, error) {
	if box < 0 || box > 8 {
		return nil, fmt.Errorf("box must be between 0 and 8, got %d", box)
	}

	cells := make([]int, 9)
	boxRow := (box / 3) * 3
	boxCol := (box % 3) * 3

	idx := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			cells[idx] = (boxRow+r)*9 + (boxCol + c)
			idx++
		}
	}

	return &BoxConstraint{
		BaseConstraint: lib.BaseConstraint{
			Cells: cells,
			Name:  fmt.Sprintf("Box %d", box+1),
		},
		box: box,
	}, nil
}

func (bc *BoxConstraint) IsValid(board *lib.Board) (bool, error) {
	if board == nil {
		return false, fmt.Errorf("board cannot be nil")
	}

	boxData := board.GetBox(bc.box)
	return lib.HasUniqueNonZeros(boxData[:]), nil
}

func (bc *BoxConstraint) GetDescription() string {
	return fmt.Sprintf("All values in 3x3 box %d must be unique (1-9)", bc.box+1)
}

// PropagateValueChange propagates the value change to other cells in the box
// This is called automatically via the observer pattern when a cell is solved
func (bc *BoxConstraint) PropagateValueChange(row, col, value int) {
	if value == 0 {
		return // No value set, nothing to propagate
	}

	// Get the board from the base constraint
	if bc.Board == nil {
		return
	}

	// Remove the value from candidates of all other cells in this box
	for _, cellIndex := range bc.Cells {
		otherRow, otherCol := cellIndex/9, cellIndex%9
		if otherRow != row || otherCol != col {
			otherCell := bc.Board.GetCellAt(otherRow, otherCol)
			if otherCell != nil && !otherCell.IsSolved() {
				otherCell.RemoveCandidate(value)
			}
		}
	}
}

func (bc *BoxConstraint) RequiresUniqueness() bool {
	return true
}

func (bc *BoxConstraint) ApplyPencilMarkConstraints(board *lib.Board) bool {
	// Apply both naked and hidden subset techniques up to quads (size 4)
	changed := false
	changed = lib.ApplyNakedSubsets(board, bc.Cells, 4) || changed
	changed = lib.ApplyHiddenSubsets(board, bc.Cells, 4) || changed
	return changed
}
