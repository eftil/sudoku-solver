package lib_test

import (
	"testing"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/observer"
)

func TestCellNotifier(t *testing.T) {
	notifier := observer.NewCellNotifier()

	if notifier == nil {
		t.Fatal("NewCellNotifier returned nil")
	}

	if notifier.HasObservers() {
		t.Error("New notifier should not have observers")
	}
}

func TestCellNotifierAddRemoveObserver(t *testing.T) {
	notifier := observer.NewCellNotifier()
	mock := &MockObserver{}

	notifier.AddObserver(mock)
	if !notifier.HasObservers() {
		t.Error("Notifier should have observers after adding one")
	}

	notifier.RemoveObserver(mock)
	if notifier.HasObservers() {
		t.Error("Notifier should not have observers after removing the only one")
	}
}

func TestCellNotifierNotifications(t *testing.T) {
	notifier := observer.NewCellNotifier()
	mock := &MockObserver{}
	notifier.AddObserver(mock)

	// Test OnSingleCandidate
	notifier.NotifySingleCandidate(1, 2, 5)
	if len(mock.singleCandidateCalls) != 1 {
		t.Errorf("Expected 1 single candidate notification, got %d", len(mock.singleCandidateCalls))
	}

	if len(mock.singleCandidateCalls) > 0 {
		call := mock.singleCandidateCalls[0]
		if call.row != 1 || call.col != 2 || call.candidate != 5 {
			t.Errorf("Wrong parameters in OnSingleCandidate: row=%d, col=%d, candidate=%d",
				call.row, call.col, call.candidate)
		}
	}

	// Test OnCellSolved
	notifier.NotifyCellSolved(3, 4, 7)
	if len(mock.cellSolvedCalls) != 1 {
		t.Errorf("Expected 1 cell solved notification, got %d", len(mock.cellSolvedCalls))
	}

	if len(mock.cellSolvedCalls) > 0 {
		call := mock.cellSolvedCalls[0]
		if call.row != 3 || call.col != 4 || call.value != 7 {
			t.Errorf("Wrong parameters in OnCellSolved: row=%d, col=%d, value=%d",
				call.row, call.col, call.value)
		}
	}

	// Test OnCandidateEliminated
	notifier.NotifyCandidateEliminated(5, 6, 3, 4)
	if len(mock.candidateEliminatedCalls) != 1 {
		t.Errorf("Expected 1 candidate eliminated notification, got %d", len(mock.candidateEliminatedCalls))
	}

	if len(mock.candidateEliminatedCalls) > 0 {
		call := mock.candidateEliminatedCalls[0]
		if call.row != 5 || call.col != 6 || call.candidate != 3 || call.remainingCount != 4 {
			t.Errorf("Wrong parameters in OnCandidateEliminated: row=%d, col=%d, candidate=%d, remaining=%d",
				call.row, call.col, call.candidate, call.remainingCount)
		}
	}
}

func TestCellNotifierClearObservers(t *testing.T) {
	notifier := observer.NewCellNotifier()
	mock1 := &MockObserver{}
	mock2 := &MockObserver{}

	notifier.AddObserver(mock1)
	notifier.AddObserver(mock2)

	if !notifier.HasObservers() {
		t.Error("Should have observers")
	}

	notifier.ClearObservers()
	if notifier.HasObservers() {
		t.Error("Should not have observers after clearing")
	}
}

func TestAutoSolverObserver(t *testing.T) {
	autoSolver := observer.NewAutoSolverObserver()

	if autoSolver == nil {
		t.Fatal("NewAutoSolverObserver returned nil")
	}

	if !autoSolver.IsEnabled() {
		t.Error("New auto solver should be enabled")
	}

	if autoSolver.GetSolutionCount() != 0 {
		t.Error("New auto solver should have 0 solutions")
	}

	if len(autoSolver.GetCellsToSolve()) != 0 {
		t.Error("New auto solver should have 0 cells to solve")
	}
}

func TestAutoSolverObserverSingleCandidate(t *testing.T) {
	autoSolver := observer.NewAutoSolverObserver()

	// Trigger single candidate notification
	autoSolver.OnSingleCandidate(2, 3, 7)

	cellsToSolve := autoSolver.GetCellsToSolve()
	if len(cellsToSolve) != 1 {
		t.Errorf("Expected 1 cell to solve, got %d", len(cellsToSolve))
	}

	if cellsToSolve["2,3"] != 7 {
		t.Errorf("Expected cell (2,3) to have value 7, got %d", cellsToSolve["2,3"])
	}
}

func TestAutoSolverObserverCellSolved(t *testing.T) {
	autoSolver := observer.NewAutoSolverObserver()

	if autoSolver.GetSolutionCount() != 0 {
		t.Error("Initial solution count should be 0")
	}

	autoSolver.OnCellSolved(1, 2, 5)
	if autoSolver.GetSolutionCount() != 1 {
		t.Errorf("Expected solution count 1, got %d", autoSolver.GetSolutionCount())
	}

	autoSolver.OnCellSolved(3, 4, 7)
	if autoSolver.GetSolutionCount() != 2 {
		t.Errorf("Expected solution count 2, got %d", autoSolver.GetSolutionCount())
	}
}

func TestAutoSolverObserverEnableDisable(t *testing.T) {
	autoSolver := observer.NewAutoSolverObserver()

	autoSolver.Disable()
	if autoSolver.IsEnabled() {
		t.Error("Auto solver should be disabled")
	}

	// Notifications should not be recorded when disabled
	autoSolver.OnSingleCandidate(1, 1, 5)
	if len(autoSolver.GetCellsToSolve()) != 0 {
		t.Error("Disabled auto solver should not record cells to solve")
	}

	autoSolver.Enable()
	if !autoSolver.IsEnabled() {
		t.Error("Auto solver should be enabled")
	}

	// Notifications should work again
	autoSolver.OnSingleCandidate(1, 1, 5)
	if len(autoSolver.GetCellsToSolve()) != 1 {
		t.Error("Enabled auto solver should record cells to solve")
	}
}

func TestBoardObserver(t *testing.T) {
	board := lib.NewBoard()
	mock := &MockObserver{}

	board.AddObserver(mock)

	// Set a value, should trigger notification
	err := board.Set(2, 3, 7)
	if err != nil {
		t.Fatalf("Board.Set failed: %v", err)
	}

	if len(mock.cellSolvedCalls) != 1 {
		t.Errorf("Expected 1 cell solved notification from board, got %d", len(mock.cellSolvedCalls))
	}

	// Remove observer
	board.RemoveObserver(mock)

	// Create another mock to verify removal worked
	mock2 := &MockObserver{}
	board.AddObserver(mock2)

	err = board.Set(4, 5, 9)
	if err != nil {
		t.Fatalf("Board.Set failed: %v", err)
	}

	// Original mock should not receive new notifications
	if len(mock.cellSolvedCalls) != 1 {
		t.Error("Removed observer should not receive new notifications")
	}

	// New mock should receive notifications
	if len(mock2.cellSolvedCalls) != 1 {
		t.Error("New observer should receive notifications")
	}
}
