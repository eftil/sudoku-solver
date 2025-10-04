package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewRowConstraint(t *testing.T) {
	tests := []struct {
		name      string
		row       int
		shouldErr bool
	}{
		{"valid row 0", 0, false},
		{"valid row 4", 4, false},
		{"valid row 8", 8, false},
		{"invalid row -1", -1, true},
		{"invalid row 9", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := constraints.NewRowConstraint(tt.row)
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
			if rc == nil {
				t.Errorf("expected constraint but got nil")
				return
			}

			// Check cells are correct
			cells := rc.GetCells()
			if len(cells) != 9 {
				t.Errorf("expected 9 cells, got %d", len(cells))
			}

			// Verify cells correspond to the correct row
			for col := 0; col < 9; col++ {
				expectedCell := tt.row*9 + col
				if cells[col] != expectedCell {
					t.Errorf("cell[%d] = %d, want %d", col, cells[col], expectedCell)
				}
			}
		})
	}
}

func TestRowConstraintIsValid(t *testing.T) {
	tests := []struct {
		name      string
		row       int
		values    [9]int
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "empty row",
			row:       0,
			values:    [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid partial row",
			row:       0,
			values:    [9]int{1, 2, 3, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid complete row",
			row:       0,
			values:    [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "duplicate in row",
			row:       0,
			values:    [9]int{1, 2, 3, 1, 0, 0, 0, 0, 0},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := constraints.NewRowConstraint(tt.row)
			if err != nil {
				t.Fatalf("failed to create constraint: %v", err)
			}

			board := lib.NewBoard()
			for col := 0; col < 9; col++ {
				err := board.Set(tt.row, col, tt.values[col])
				if err != nil {
					t.Fatalf("Set(%d, %d, %d) failed: %v", tt.row, col, tt.values[col], err)
				}
			}

			valid, err := rc.IsValid(board)
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

func TestRowConstraintIsValidNilBoard(t *testing.T) {
	rc, err := constraints.NewRowConstraint(0)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	valid, err := rc.IsValid(nil)
	if err == nil {
		t.Error("expected error for nil board, got none")
	}
	if valid {
		t.Error("expected invalid result for nil board")
	}
}
