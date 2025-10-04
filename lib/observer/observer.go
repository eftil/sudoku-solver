package observer

// CellObserver is an interface for objects that want to be notified of cell events
type CellObserver interface {
	// OnSingleCandidate is called when a cell has only one candidate remaining
	OnSingleCandidate(row, col, candidate int)

	// OnCellSolved is called when a cell's value is set
	OnCellSolved(row, col, value int)

	// OnCandidateEliminated is called when a candidate is removed from a cell
	OnCandidateEliminated(row, col, candidate int, remainingCount int)
}

// CellNotifier manages observers for cell events
type CellNotifier struct {
	observers []CellObserver
}

// NewCellNotifier creates a new cell notifier
func NewCellNotifier() *CellNotifier {
	return &CellNotifier{
		observers: make([]CellObserver, 0),
	}
}

// AddObserver adds an observer to the notifier
func (cn *CellNotifier) AddObserver(observer CellObserver) {
	if observer == nil {
		return
	}
	cn.observers = append(cn.observers, observer)
}

// RemoveObserver removes an observer from the notifier
func (cn *CellNotifier) RemoveObserver(observer CellObserver) {
	if observer == nil {
		return
	}

	for i, obs := range cn.observers {
		if obs == observer {
			cn.observers = append(cn.observers[:i], cn.observers[i+1:]...)
			return
		}
	}
}

// NotifySingleCandidate notifies all observers that a cell has a single candidate
func (cn *CellNotifier) NotifySingleCandidate(row, col, candidate int) {
	for _, observer := range cn.observers {
		observer.OnSingleCandidate(row, col, candidate)
	}
}

// NotifyCellSolved notifies all observers that a cell has been solved
func (cn *CellNotifier) NotifyCellSolved(row, col, value int) {
	for _, observer := range cn.observers {
		observer.OnCellSolved(row, col, value)
	}
}

// NotifyCandidateEliminated notifies all observers that a candidate was eliminated
func (cn *CellNotifier) NotifyCandidateEliminated(row, col, candidate, remainingCount int) {
	for _, observer := range cn.observers {
		observer.OnCandidateEliminated(row, col, candidate, remainingCount)
	}
}

// HasObservers returns true if there are any observers registered
func (cn *CellNotifier) HasObservers() bool {
	return len(cn.observers) > 0
}

// ClearObservers removes all observers
func (cn *CellNotifier) ClearObservers() {
	cn.observers = make([]CellObserver, 0)
}
