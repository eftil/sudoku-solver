package constraints_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestNewGermanWhispersConstraint(t *testing.T) {
	tests := []struct {
		name      string
		cells     []int
		shouldErr bool
	}{
		{"valid two cells", []int{0, 1}, false},
		{"valid multiple cells", []int{0, 1, 2, 3}, false},
		{"single cell", []int{0}, true},
		{"empty cells", []int{}, true},
		{"invalid cell index negative", []int{0, -1, 2}, true},
		{"invalid cell index too large", []int{0, 81, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw, err := constraints.NewGermanWhispersConstraint(tt.cells)
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
			if gw == nil {
				t.Errorf("expected constraint but got nil")
			}
		})
	}
}

func TestGermanWhispersConstraintIsValid(t *testing.T) {
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
			name:      "valid difference of 5",
			cells:     []int{0, 1},
			values:    []int{1, 6},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid difference of 6",
			cells:     []int{0, 1},
			values:    []int{3, 9},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid multiple pairs",
			cells:     []int{0, 1, 2},
			values:    []int{1, 9, 4},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "invalid difference of 4",
			cells:     []int{0, 1},
			values:    []int{1, 5},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "invalid difference of 3",
			cells:     []int{0, 1},
			values:    []int{3, 6},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "invalid difference of 1",
			cells:     []int{0, 1},
			values:    []int{5, 6},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "partial cells with empty",
			cells:     []int{0, 1, 2},
			values:    []int{1, 0, 7},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "valid first pair, invalid second",
			cells:     []int{0, 1, 2},
			values:    []int{1, 6, 8},
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "reverse difference valid",
			cells:     []int{0, 1},
			values:    []int{9, 4},
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "exactly 5 difference",
			cells:     []int{0, 1, 2, 3},
			values:    []int{1, 6, 1, 6},
			wantValid: true,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw, err := constraints.NewGermanWhispersConstraint(tt.cells)
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

			valid, err := gw.IsValid(board)
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

func TestGermanWhispersConstraintIsValidNilBoard(t *testing.T) {
	gw, err := constraints.NewGermanWhispersConstraint([]int{0, 1, 2})
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	valid, err := gw.IsValid(nil)
	if err == nil {
		t.Error("expected error for nil board, got none")
	}
	if valid {
		t.Error("expected invalid result for nil board")
	}
}
