package lib_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
)

func TestBoardSetGet(t *testing.T) {
	board := lib.NewBoard()

	// Test setting and getting values
	tests := []struct {
		row   int
		col   int
		value int
	}{
		{0, 0, 5},
		{0, 8, 9},
		{8, 0, 1},
		{8, 8, 3},
		{4, 4, 7},
	}

	for _, tt := range tests {
		err := board.Set(tt.row, tt.col, tt.value)
		if err != nil {
			t.Errorf("Set(%d, %d, %d) failed: %v", tt.row, tt.col, tt.value, err)
		}
		got := board.Get(tt.row, tt.col)
		if got != tt.value {
			t.Errorf("Set(%d, %d, %d) then Get(%d, %d) = %d, want %d",
				tt.row, tt.col, tt.value, tt.row, tt.col, got, tt.value)
		}
	}
}

func TestBoardGetRow(t *testing.T) {
	board := lib.NewBoard()

	// Set a row
	expected := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for col := 0; col < 9; col++ {
		err := board.Set(0, col, expected[col])
		if err != nil {
			t.Fatalf("Set(0, %d, %d) failed: %v", col, expected[col], err)
		}
	}

	// Get the row
	got := board.GetRow(0)
	if got != expected {
		t.Errorf("GetRow(0) = %v, want %v", got, expected)
	}
}

func TestBoardGetColumn(t *testing.T) {
	board := lib.NewBoard()

	// Set a column
	expected := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for row := 0; row < 9; row++ {
		err := board.Set(row, 0, expected[row])
		if err != nil {
			t.Fatalf("Set(%d, 0, %d) failed: %v", row, expected[row], err)
		}
	}

	// Get the column
	got := board.GetColumn(0)
	if got != expected {
		t.Errorf("GetColumn(0) = %v, want %v", got, expected)
	}
}

func TestBoardGetBox(t *testing.T) {
	board := lib.NewBoard()

	// Set box 0 (top-left 3x3)
	expected := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	idx := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			err := board.Set(r, c, expected[idx])
			if err != nil {
				t.Fatalf("Set(%d, %d, %d) failed: %v", r, c, expected[idx], err)
			}
			idx++
		}
	}

	// Get the box
	got := board.GetBox(0)
	if got != expected {
		t.Errorf("GetBox(0) = %v, want %v", got, expected)
	}
}

func TestBoardGetBoxMiddle(t *testing.T) {
	board := lib.NewBoard()

	// Set box 4 (center 3x3)
	expected := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	idx := 0
	for r := 3; r < 6; r++ {
		for c := 3; c < 6; c++ {
			err := board.Set(r, c, expected[idx])
			if err != nil {
				t.Fatalf("Set(%d, %d, %d) failed: %v", r, c, expected[idx], err)
			}
			idx++
		}
	}

	// Get the box
	got := board.GetBox(4)
	if got != expected {
		t.Errorf("GetBox(4) = %v, want %v", got, expected)
	}
}

func TestBoardAddConstraint(t *testing.T) {
	board := lib.NewBoard()

	rc, err := constraints.NewRowConstraint(0)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}

	board.AddConstraint(rc)

	constraintList := board.GetConstraints()
	if len(constraintList) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(constraintList))
	}
}

func TestBoardValidateAll(t *testing.T) {
	board := lib.NewBoard()

	// Add row constraint
	rc, err := constraints.NewRowConstraint(0)
	if err != nil {
		t.Fatalf("failed to create constraint: %v", err)
	}
	board.AddConstraint(rc)

	// Empty board should be valid
	valid, err := board.ValidateAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("empty board should be valid")
	}

	// Add valid values
	err = board.Set(0, 0, 1)
	if err != nil {
		t.Fatalf("Set(0, 0, 1) failed: %v", err)
	}
	err = board.Set(0, 1, 2)
	if err != nil {
		t.Fatalf("Set(0, 1, 2) failed: %v", err)
	}
	valid, err = board.ValidateAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("board with unique values should be valid")
	}

	// Add duplicate value
	err = board.Set(0, 2, 1)
	if err != nil {
		t.Fatalf("Set(0, 2, 1) failed: %v", err)
	}
	valid, err = board.ValidateAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if valid {
		t.Error("board with duplicate values should be invalid")
	}
}

func TestBoardValidateAllMultipleConstraints(t *testing.T) {
	board := lib.NewBoard()

	// Add all standard sudoku constraints
	for i := 0; i < 9; i++ {
		rc, err := constraints.NewRowConstraint(i)
		if err != nil {
			t.Fatalf("failed to create row constraint: %v", err)
		}
		board.AddConstraint(rc)

		cc, err := constraints.NewColumnConstraint(i)
		if err != nil {
			t.Fatalf("failed to create column constraint: %v", err)
		}
		board.AddConstraint(cc)

		bc, err := constraints.NewBoxConstraint(i)
		if err != nil {
			t.Fatalf("failed to create box constraint: %v", err)
		}
		board.AddConstraint(bc)
	}

	// Set up a valid partial sudoku
	err := board.Set(0, 0, 5)
	if err != nil {
		t.Fatalf("Set(0, 0, 5) failed: %v", err)
	}
	err = board.Set(0, 1, 3)
	if err != nil {
		t.Fatalf("Set(0, 1, 3) failed: %v", err)
	}
	err = board.Set(1, 0, 6)
	if err != nil {
		t.Fatalf("Set(1, 0, 6) failed: %v", err)
	}

	valid, err := board.ValidateAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("valid partial sudoku should pass validation")
	}

	// Violate a constraint
	err = board.Set(0, 2, 5) // duplicate 5 in row 0
	if err != nil {
		t.Fatalf("Set(0, 2, 5) failed: %v", err)
	}
	valid, err = board.ValidateAll()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if valid {
		t.Error("invalid sudoku should fail validation")
	}
}

func TestBoardApplyPencilMarkConstraints(t *testing.T) {
	board := lib.NewBoard()

	// Add standard constraints
	for i := 0; i < 9; i++ {
		rc, _ := constraints.NewRowConstraint(i)
		board.AddConstraint(rc)
		cc, _ := constraints.NewColumnConstraint(i)
		board.AddConstraint(cc)
		bc, _ := constraints.NewBoxConstraint(i)
		board.AddConstraint(bc)
	}

	// Set up a board state that would benefit from pencil mark techniques
	board.Set(0, 0, 1)
	board.Set(0, 1, 2)
	board.Set(1, 0, 3)

	// Apply pencil mark constraints
	changed := board.ApplyPencilMarkConstraints()

	// Just verify it runs without error - exact behavior depends on board state
	// In this simple case, it might or might not find eliminations
	_ = changed
}

func TestBoardApplyPencilMarkConstraintsUntilStable(t *testing.T) {
	board := lib.NewBoard()

	// Add standard constraints
	for i := 0; i < 9; i++ {
		rc, _ := constraints.NewRowConstraint(i)
		board.AddConstraint(rc)
	}

	// Apply until stable
	iterations := board.ApplyPencilMarkConstraintsUntilStable()

	if iterations < 1 {
		t.Error("Should have at least 1 iteration")
	}
}

func TestBoardApplyAdvancedTechniques(t *testing.T) {
	board := lib.NewBoard()

	// Add standard constraints
	for i := 0; i < 9; i++ {
		rc, _ := constraints.NewRowConstraint(i)
		board.AddConstraint(rc)
		cc, _ := constraints.NewColumnConstraint(i)
		board.AddConstraint(cc)
		bc, _ := constraints.NewBoxConstraint(i)
		board.AddConstraint(bc)
	}

	// Set up a simple board
	board.Set(0, 0, 5)

	// Apply advanced techniques
	changed := board.ApplyAdvancedTechniques()

	// Just verify it runs without error
	_ = changed
}

func TestBoardGetCellAt(t *testing.T) {
	board := lib.NewBoard()

	cell := board.GetCellAt(3, 4)
	if cell == nil {
		t.Fatal("GetCellAt returned nil for valid coordinates")
	}

	if cell.GetRow() != 3 || cell.GetCol() != 4 {
		t.Errorf("GetCellAt(3, 4) returned cell with wrong coordinates: (%d, %d)",
			cell.GetRow(), cell.GetCol())
	}

	// Test invalid coordinates
	if board.GetCellAt(-1, 0) != nil {
		t.Error("GetCellAt with negative row should return nil")
	}

	if board.GetCellAt(0, -1) != nil {
		t.Error("GetCellAt with negative col should return nil")
	}

	if board.GetCellAt(9, 0) != nil {
		t.Error("GetCellAt with row >= 9 should return nil")
	}

	if board.GetCellAt(0, 9) != nil {
		t.Error("GetCellAt with col >= 9 should return nil")
	}
}

func TestBoardGetCell(t *testing.T) {
	board := lib.NewBoard()

	cell := board.GetCell(40) // Middle of board
	if cell == nil {
		t.Fatal("GetCell returned nil for valid index")
	}

	if cell.GetIndex() != 40 {
		t.Errorf("GetCell(40) returned cell with wrong index: %d", cell.GetIndex())
	}

	// Test invalid indices
	if board.GetCell(-1) != nil {
		t.Error("GetCell with negative index should return nil")
	}

	if board.GetCell(81) != nil {
		t.Error("GetCell with index >= 81 should return nil")
	}
}

func TestBoardInvalidPosition(t *testing.T) {
	board := lib.NewBoard()

	tests := []struct {
		name     string
		row, col int
	}{
		{"negative row", -1, 0},
		{"negative col", 0, -1},
		{"row too large", 9, 0},
		{"col too large", 0, 9},
		{"both too large", 9, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := board.Set(tt.row, tt.col, 5)
			if err == nil {
				t.Errorf("Set(%d, %d, 5) should return error", tt.row, tt.col)
			}

			value := board.Get(tt.row, tt.col)
			if value != 0 {
				t.Errorf("Get(%d, %d) should return 0 for invalid position", tt.row, tt.col)
			}
		})
	}
}

func TestBoardInvalidValue(t *testing.T) {
	board := lib.NewBoard()

	tests := []struct {
		name  string
		value int
	}{
		{"negative value", -1},
		{"value too large", 10},
		{"value too large 2", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := board.Set(0, 0, tt.value)
			if err == nil {
				t.Errorf("Set(0, 0, %d) should return error", tt.value)
			}
		})
	}
}

func TestBoardPrint(t *testing.T) {
	board := lib.NewBoard()

	// Just verify Print doesn't panic
	board.Set(0, 0, 5)
	board.Set(4, 4, 9)

	// This function prints to stdout, we just verify it doesn't crash
	board.Print()
}
