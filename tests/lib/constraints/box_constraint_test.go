package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewBoxConstraint(t *testing.T) {
	tests := []struct {
		name      string
		box       int
		shouldErr bool
	}{
		{"valid box 0", 0, false},
		{"valid box 4", 4, false},
		{"valid box 8", 8, false},
		{"invalid box -1", -1, true},
		{"invalid box 9", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc, err := constraints.NewBoxConstraint(tt.box)
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
			if bc == nil {
				t.Errorf("expected constraint but got nil")
				return
			}

			// Check cells are correct
			cells := bc.GetCells()
			if len(cells) != 9 {
				t.Errorf("expected 9 cells, got %d", len(cells))
			}
		})
	}
}

func TestBoxConstraintIsValid(t *testing.T) {
	tests := []struct {
		name      string
		box       int
		values    [9]int
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "empty box",
			box:       0,
			values:    [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid partial box",
			box:       0,
			values:    [9]int{1, 2, 3, 0, 0, 0, 0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid complete box",
			box:       0,
			values:    [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "duplicate in box",
			box:       0,
			values:    [9]int{1, 2, 3, 1, 0, 0, 0, 0, 0},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc, err := constraints.NewBoxConstraint(tt.box)
			if err != nil {
				t.Fatalf("failed to create constraint: %v", err)
			}

			board := lib.NewBoard()
			boxRow := (tt.box / 3) * 3
			boxCol := (tt.box % 3) * 3

			idx := 0
			for r := 0; r < 3; r++ {
				for c := 0; c < 3; c++ {
					err := board.Set(boxRow+r, boxCol+c, tt.values[idx])
					if err != nil {
						t.Fatalf("Set(%d, %d, %d) failed: %v", boxRow+r, boxCol+c, tt.values[idx], err)
					}
					idx++
				}
			}

			valid, err := bc.IsValid(board)
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

func TestBoxConstraintIsValidNilBoard(t *testing.T) {
	bc, err := constraints.NewBoxConstraint(0)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	valid, err := bc.IsValid(nil)
	if err == nil {
		t.Error("expected error for nil board, got none")
	}
	if valid {
		t.Error("expected invalid result for nil board")
	}
}
