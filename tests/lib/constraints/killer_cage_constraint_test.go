package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewKillerCageConstraint(t *testing.T) {
	tests := []struct {
		name      string
		cells     []int
		targetSum int
		shouldErr bool
	}{
		{"valid cage", []int{0, 1, 9}, 15, false},
		{"valid single cell", []int{0}, 5, false},
		{"empty cells", []int{}, 10, true},
		{"invalid cell index negative", []int{0, -1, 2}, 10, true},
		{"invalid cell index too large", []int{0, 81, 2}, 10, true},
		{"target sum too small", []int{0, 1, 2}, 0, true},
		{"target sum too large", []int{0, 1, 2}, 46, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc, err := constraints.NewKillerCageConstraint(tt.cells, tt.targetSum)
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
			if kc == nil {
				t.Errorf("expected constraint but got nil")
			}
		})
	}
}

func TestKillerCageConstraintIsValid(t *testing.T) {
	tests := []struct {
		name      string
		cells     []int
		targetSum int
		values    []int
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "empty cells",
			cells:     []int{0, 1, 2},
			targetSum: 15,
			values:    []int{0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid complete cage sum",
			cells:     []int{0, 1, 2},
			targetSum: 6,
			values:    []int{1, 2, 3},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid partial cage under sum",
			cells:     []int{0, 1, 2},
			targetSum: 15,
			values:    []int{5, 6, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "invalid complete cage wrong sum",
			cells:     []int{0, 1, 2},
			targetSum: 10,
			values:    []int{1, 2, 3},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "invalid partial cage exceeds sum",
			cells:     []int{0, 1, 2},
			targetSum: 10,
			values:    []int{8, 9, 0},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "duplicate values in cage",
			cells:     []int{0, 1, 2},
			targetSum: 15,
			values:    []int{5, 5, 5},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "valid two-cell cage",
			cells:     []int{0, 1},
			targetSum: 9,
			values:    []int{4, 5},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "partial duplicate values",
			cells:     []int{0, 1, 2},
			targetSum: 15,
			values:    []int{5, 5, 0},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc, err := constraints.NewKillerCageConstraint(tt.cells, tt.targetSum)
			if err != nil {
				t.Fatalf("failed to create constraint: %v", err)
			}

			board := lib.NewBoard()
			for i, cellIdx := range tt.cells {
				row := cellIdx / 9
				col := cellIdx % 9
				err := board.Set(row, col, tt.values[i])
				if err != nil {
					t.Fatalf("Set(%d, %d, %d) failed: %v", row, col, tt.values[i], err)
				}
			}

			valid, err := kc.IsValid(board)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v for values %v (sum=%d)", valid, tt.wantValid, tt.values, tt.targetSum)
			}
		})
	}
}

func TestKillerCageConstraintIsValidNilBoard(t *testing.T) {
	kc, err := constraints.NewKillerCageConstraint([]int{0, 1, 2}, 15)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	valid, err := kc.IsValid(nil)
	if err == nil {
		t.Error("expected error for nil board, got none")
	}
	if valid {
		t.Error("expected invalid result for nil board")
	}
}
