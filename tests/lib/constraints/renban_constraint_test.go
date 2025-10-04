package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewRenbanConstraint(t *testing.T) {
	tests := []struct {
		name      string
		cells     []int
		shouldErr bool
	}{
		{"valid single cell", []int{0}, false},
		{"valid multiple cells", []int{0, 1, 2}, false},
		{"empty cells", []int{}, true},
		{"invalid cell index negative", []int{0, -1, 2}, true},
		{"invalid cell index too large", []int{0, 81, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := constraints.NewRenbanConstraint(tt.cells)
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
			}
		})
	}
}

func TestRenbanConstraintIsValid(t *testing.T) {
	tests := []struct {
		name      string
		cells     []int
		values    []int
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "empty cells",
			cells:     []int{0, 1, 2},
			values:    []int{0, 0, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid consecutive sequence 123",
			cells:     []int{0, 1, 2},
			values:    []int{1, 2, 3},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid consecutive sequence 567",
			cells:     []int{0, 1, 2},
			values:    []int{5, 6, 7},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid consecutive sequence out of order",
			cells:     []int{0, 1, 2},
			values:    []int{3, 1, 2},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "partial sequence valid",
			cells:     []int{0, 1, 2},
			values:    []int{1, 2, 0},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "duplicate values",
			cells:     []int{0, 1, 2},
			values:    []int{1, 1, 2},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "non-consecutive sequence",
			cells:     []int{0, 1, 2},
			values:    []int{1, 2, 4},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "gap in sequence",
			cells:     []int{0, 1, 2, 3},
			values:    []int{1, 3, 4, 5},
			wantValid: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc, err := constraints.NewRenbanConstraint(tt.cells)
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

			valid, err := rc.IsValid(board)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v for values %v", valid, tt.wantValid, tt.values)
			}
		})
	}
}

func TestRenbanConstraintIsValidNilBoard(t *testing.T) {
	rc, err := constraints.NewRenbanConstraint([]int{0, 1, 2})
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
