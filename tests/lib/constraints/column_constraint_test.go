package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewColumnConstraint(t *testing.T) {
	tests := []struct {
		name      string
		col       int
		shouldErr bool
	}{
		{"valid column 0", 0, false},
		{"valid column 4", 4, false},
		{"valid column 8", 8, false},
		{"invalid column -1", -1, true},
		{"invalid column 9", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := constraints.NewColumnConstraint(tt.col)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if cc == nil {
				t.Errorf("expected constraint but got nil")
				return
			}

			// Check cells are correct
			cells := cc.GetCells()
			if len(cells) != 9 {
				t.Errorf("expected 9 cells, got %d", len(cells))
			}

			// Verify cells correspond to the correct column
			for row := 0; row < 9; row++ {
				expectedCell := row*9 + tt.col
				if cells[row] != expectedCell {
					t.Errorf("cell[%d] = %d, want %d", row, cells[row], expectedCell)
				}
			}
		})
	}
}

func TestColumnConstraintIsValid(t *testing.T) {
	tests := []struct {
		name      string
		col       int
		values    [9]int
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "empty column",
			col:       0,
			values:    [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid partial column",
			col:       0,
			values:    [9]int{1, 2, 3, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid complete column",
			col:       0,
			values:    [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "duplicate in column",
			col:       0,
			values:    [9]int{1, 2, 3, 1, 0, 0, 0, 0, 0},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := constraints.NewColumnConstraint(tt.col)
			if err != nil {
				t.Fatalf("failed to create constraint: %v", err)
			}

			board := lib.NewBoard()
			for row := 0; row < 9; row++ {
				err := board.Set(row, tt.col, tt.values[row])
				if err != nil {
					t.Fatalf("Set(%d, %d, %d) failed: %v", row, tt.col, tt.values[row], err)
				}
			}

			valid, err := cc.IsValid(board)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestColumnConstraintIsValidNilBoard(t *testing.T) {
	cc, err := constraints.NewColumnConstraint(0)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	valid, err := cc.IsValid(nil)
	if err == nil {
		t.Error("expected error for nil board, got none")
	}
	if valid {
		t.Error("expected invalid result for nil board")
	}
}
