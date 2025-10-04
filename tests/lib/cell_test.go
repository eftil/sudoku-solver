package lib_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
)

func TestNewCell(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	if cell == nil {
		t.Fatal("NewCell returned nil")
	}

	if cell.GetRow() != 0 {
		t.Errorf("Expected row 0, got %d", cell.GetRow())
	}

	if cell.GetCol() != 0 {
		t.Errorf("Expected col 0, got %d", cell.GetCol())
	}

	if cell.GetValue() != 0 {
		t.Errorf("Expected initial value 0, got %d", cell.GetValue())
	}

	if cell.GetBoard() != board {
		t.Error("Cell's board reference is incorrect")
	}
}

func TestCellSetValue(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	tests := []struct {
		name      string
		value     int
		shouldErr bool
	}{
		{"valid value 1", 1, false},
		{"valid value 5", 5, false},
		{"valid value 9", 9, false},
		{"valid value 0 (clear)", 0, false},
		{"invalid value -1", -1, true},
		{"invalid value 10", 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cell.SetValue(tt.value)
			if (err != nil) != tt.shouldErr {
				t.Errorf("SetValue(%d) error = %v, shouldErr %v", tt.value, err, tt.shouldErr)
			}

			if !tt.shouldErr && cell.GetValue() != tt.value {
				t.Errorf("After SetValue(%d), GetValue() = %d", tt.value, cell.GetValue())
			}
		})
	}
}

func TestCellCandidates(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(2, 3, board)

	// Initially should have all candidates
	if cell.CandidateCount() != 9 {
		t.Errorf("Expected 9 initial candidates, got %d", cell.CandidateCount())
	}

	// Test removing candidates
	cell.RemoveCandidate(5)
	if cell.HasCandidate(5) {
		t.Error("Candidate 5 should have been removed")
	}

	if cell.CandidateCount() != 8 {
		t.Errorf("Expected 8 candidates after removal, got %d", cell.CandidateCount())
	}

	// Remove more candidates
	cell.RemoveCandidate(1)
	cell.RemoveCandidate(2)
	cell.RemoveCandidate(3)

	if cell.CandidateCount() != 5 {
		t.Errorf("Expected 5 candidates, got %d", cell.CandidateCount())
	}

	// Test adding candidate back
	cell.AddCandidate(5)
	if !cell.HasCandidate(5) {
		t.Error("Candidate 5 should have been added back")
	}

	if cell.CandidateCount() != 6 {
		t.Errorf("Expected 6 candidates after adding, got %d", cell.CandidateCount())
	}
}

func TestCellSetValueClearsCandidates(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	// Initially should have candidates
	if cell.CandidateCount() != 9 {
		t.Errorf("Expected 9 initial candidates, got %d", cell.CandidateCount())
	}

	// Set a value
	err := cell.SetValue(7)
	if err != nil {
		t.Fatalf("SetValue failed: %v", err)
	}

	// Candidates should be cleared
	if cell.CandidateCount() != 0 {
		t.Errorf("Expected 0 candidates after setting value, got %d", cell.CandidateCount())
	}

	if cell.HasCandidate(7) {
		t.Error("Cell should not have candidates after value is set")
	}
}

func TestCellIsSolved(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	if cell.IsSolved() {
		t.Error("New cell should not be solved")
	}

	err := cell.SetValue(5)
	if err != nil {
		t.Fatalf("SetValue failed: %v", err)
	}

	if !cell.IsSolved() {
		t.Error("Cell with value should be solved")
	}

	err = cell.SetValue(0)
	if err != nil {
		t.Fatalf("SetValue(0) failed: %v", err)
	}

	if cell.IsSolved() {
		t.Error("Cell with value 0 should not be solved")
	}
}

// MockObserver for testing
type MockObserver struct {
	singleCandidateCalls     []struct{ row, col, candidate int }
	cellSolvedCalls          []struct{ row, col, value int }
	candidateEliminatedCalls []struct{ row, col, candidate, remainingCount int }
}

func (mo *MockObserver) OnSingleCandidate(row, col, candidate int) {
	mo.singleCandidateCalls = append(mo.singleCandidateCalls, struct{ row, col, candidate int }{row, col, candidate})
}

func (mo *MockObserver) OnCellSolved(row, col, value int) {
	mo.cellSolvedCalls = append(mo.cellSolvedCalls, struct{ row, col, value int }{row, col, value})
}

func (mo *MockObserver) OnCandidateEliminated(row, col, candidate, remainingCount int) {
	mo.candidateEliminatedCalls = append(mo.candidateEliminatedCalls,
		struct{ row, col, candidate, remainingCount int }{row, col, candidate, remainingCount})
}

func TestCellObserver(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(3, 4, board)

	mock := &MockObserver{}
	cell.AddObserver(mock)

	// Set value should trigger OnCellSolved
	err := cell.SetValue(7)
	if err != nil {
		t.Fatalf("SetValue failed: %v", err)
	}

	if len(mock.cellSolvedCalls) != 1 {
		t.Errorf("Expected 1 OnCellSolved call, got %d", len(mock.cellSolvedCalls))
	}

	if len(mock.cellSolvedCalls) > 0 {
		call := mock.cellSolvedCalls[0]
		if call.row != 3 || call.col != 4 || call.value != 7 {
			t.Errorf("OnCellSolved called with wrong parameters: row=%d, col=%d, value=%d",
				call.row, call.col, call.value)
		}
	}
}

func TestCellObserverSingleCandidate(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(5, 6, board)

	mock := &MockObserver{}
	cell.AddObserver(mock)

	// Remove candidates until only one remains
	for i := 1; i <= 8; i++ {
		cell.RemoveCandidate(i)
	}

	// Should have triggered OnSingleCandidate
	if len(mock.singleCandidateCalls) != 1 {
		t.Errorf("Expected 1 OnSingleCandidate call, got %d", len(mock.singleCandidateCalls))
	}

	if len(mock.singleCandidateCalls) > 0 {
		call := mock.singleCandidateCalls[0]
		if call.row != 5 || call.col != 6 || call.candidate != 9 {
			t.Errorf("OnSingleCandidate called with wrong parameters: row=%d, col=%d, candidate=%d",
				call.row, call.col, call.candidate)
		}
	}
}

func TestCellGetNotifier(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	notifier := cell.GetNotifier()
	if notifier == nil {
		t.Error("Cell notifier should not be nil")
	}

	if !notifier.HasObservers() {
		// This is expected for a new cell
	}

	// Add an observer
	mock := &MockObserver{}
	cell.AddObserver(mock)

	if !notifier.HasObservers() {
		t.Error("Notifier should have observers after adding one")
	}
}

func TestCellGetCandidates(t *testing.T) {
	board := lib.NewBoard()
	cell := lib.NewCell(0, 0, board)

	candidates := cell.GetCandidates()
	if len(candidates) != 9 {
		t.Errorf("Expected 9 candidates, got %d", len(candidates))
	}

	// Verify all candidates 1-9 are present
	for i := 1; i <= 9; i++ {
		if !candidates[i] {
			t.Errorf("Candidate %d should be present", i)
		}
	}

	// After setting value, candidates should be empty
	cell.SetValue(5)
	candidates = cell.GetCandidates()
	if len(candidates) != 0 {
		t.Errorf("Expected 0 candidates after setting value, got %d", len(candidates))
	}
}

func TestCellIndex(t *testing.T) {
	board := lib.NewBoard()

	tests := []struct {
		row, col      int
		expectedIndex int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{0, 8, 8},
		{1, 0, 9},
		{1, 1, 10},
		{8, 8, 80},
		{4, 4, 40},
	}

	for _, tt := range tests {
		cell := lib.NewCell(tt.row, tt.col, board)
		if cell.GetIndex() != tt.expectedIndex {
			t.Errorf("Cell(%d,%d): expected index %d, got %d",
				tt.row, tt.col, tt.expectedIndex, cell.GetIndex())
		}
	}
}
