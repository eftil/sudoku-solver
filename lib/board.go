package lib

import (
	"fmt"

	"github.com/eftil/sudoku-solver.git/lib/logger"
	"github.com/eftil/sudoku-solver.git/lib/observer"
	"github.com/eftil/sudoku-solver.git/lib/utils"
)

type Board struct {
	board       [81]*Cell
	constraints []Constraint
	observers   []observer.CellObserver
}

// BoardError represents errors from board operations
type BoardError struct {
	Message string
}

func (e *BoardError) Error() string {
	return e.Message
}

// NewBoard creates a new board with all cells initialized
func NewBoard() *Board {
	logger.Info("Creating new Sudoku board...")

	b := &Board{
		observers: make([]observer.CellObserver, 0),
	}

	// Initialize all cells
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			b.board[row*9+col] = NewCell(row, col, b)
		}
	}

	logger.Info("Board created successfully with 81 cells")
	return b
}

func (b *Board) Set(row, col, value int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		logger.Error("Invalid board position: row=%d, col=%d", row, col)
		return &BoardError{Message: fmt.Sprintf("invalid position: row=%d, col=%d", row, col)}
	}

	// Initialize cell if it doesn't exist
	if b.board[row*9+col] == nil {
		logger.Debug("Initializing missing cell at R%dC%d", row+1, col+1)
		b.board[row*9+col] = NewCell(row, col, b)
	}

	logger.Info("Setting cell R%dC%d to value %d", row+1, col+1, value)
	return b.board[row*9+col].SetValue(value)
}

func (b *Board) Get(row, col int) int {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return 0
	}
	if b.board[row*9+col] == nil {
		return 0
	}
	return b.board[row*9+col].GetValue()
}

// GetCell returns the cell at the given index (0-80)
func (b *Board) GetCell(index int) *Cell {
	if index < 0 || index > 80 {
		return nil
	}
	return b.board[index]
}

// GetCellAt returns the cell at the given row and column
func (b *Board) GetCellAt(row, col int) *Cell {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return nil
	}
	return b.board[row*9+col]
}

func (b *Board) Print() {
	for i := range 81 {
		if b.board[i] != nil {
			fmt.Print(b.board[i].GetValue())
		} else {
			fmt.Print("0")
		}
		if (i+1)%9 == 0 {
			fmt.Println()
		}
	}
}

func (b *Board) GetRow(row int) [9]int {
	rowData := [9]int{}
	for i := 0; i < 9; i++ {
		if b.board[row*9+i] != nil {
			rowData[i] = b.board[row*9+i].GetValue()
		}
	}
	return rowData
}

func (b *Board) GetColumn(col int) [9]int {
	colData := [9]int{}
	for row := 0; row < 9; row++ {
		if b.board[row*9+col] != nil {
			colData[row] = b.board[row*9+col].GetValue()
		}
	}
	return colData
}

func (b *Board) GetBox(box int) [9]int {
	boxData := [9]int{}
	boxRow := (box / 3) * 3
	boxCol := (box % 3) * 3

	idx := 0
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if b.board[(boxRow+r)*9+(boxCol+c)] != nil {
				boxData[idx] = b.board[(boxRow+r)*9+(boxCol+c)].GetValue()
			}
			idx++
		}
	}
	return boxData
}

// AddConstraint adds a constraint to the board and registers it as an observer of its cells
func (b *Board) AddConstraint(c Constraint) {
	logger.Info("Adding constraint: %s - %s", c.GetName(), c.GetDescription())

	b.constraints = append(b.constraints, c)

	// Set the board reference on the constraint (using type assertion to access SetBoard)
	if bc, ok := c.(interface{ SetBoard(*Board) }); ok {
		bc.SetBoard(b)
	}

	// Register the constraint as an observer of all its cells
	// This is the elegant observer pattern in action!
	affectedCount := 0
	for _, cellIndex := range c.GetCells() {
		if cellIndex >= 0 && cellIndex <= 80 && b.board[cellIndex] != nil {
			b.board[cellIndex].AddObserver(c) // Constraint observes the cell
			affectedCount++
		}
	}

	logger.Debug("Constraint '%s' observing %d cells", c.GetName(), affectedCount)
}

// ValidateAll checks if all constraints on the board are satisfied
func (b *Board) ValidateAll() (bool, error) {
	logger.Info("Validating all %d constraints...", len(b.constraints))

	for _, constraint := range b.constraints {
		valid, err := constraint.IsValid(b)
		if err != nil {
			logger.Error("Error validating constraint '%s': %v", constraint.GetName(), err)
			return false, fmt.Errorf("error validating %s: %w", constraint.GetName(), err)
		}
		if !valid {
			logger.Warn("Constraint validation failed: %s", constraint.GetName())
			return false, nil
		}
		logger.Debug("Constraint '%s' validation: PASSED", constraint.GetName())
	}

	logger.Info("All constraints validated successfully")
	return true, nil
}

// GetConstraints returns all constraints on the board
func (b *Board) GetConstraints() []Constraint {
	return b.constraints
}

// ApplyPencilMarkConstraints applies advanced solving techniques (naked/hidden pairs, etc.)
// to all constraints that enforce uniqueness. Returns true if any candidates were eliminated.
func (b *Board) ApplyPencilMarkConstraints() bool {
	logger.SolvingStep("Pencil Mark", "Applying pencil mark constraints (naked/hidden subsets)")

	changed := false
	constraintsApplied := 0

	for _, constraint := range b.constraints {
		if constraint.RequiresUniqueness() {
			constraintsApplied++
			if constraint.ApplyPencilMarkConstraints(b) {
				changed = true
				logger.Debug("Pencil mark technique found eliminations in: %s", constraint.GetName())
			}
		}
	}

	if changed {
		logger.Info("Pencil mark constraints eliminated candidates")
	} else {
		logger.Debug("Pencil mark constraints did not find any eliminations")
	}

	return changed
}

// ApplyPencilMarkConstraintsUntilStable repeatedly applies pencil mark constraints
// until no more changes occur. Returns the number of iterations performed.
func (b *Board) ApplyPencilMarkConstraintsUntilStable() int {
	logger.Info("Applying pencil mark constraints until stable...")

	iterations := 0
	for {
		changed := b.ApplyPencilMarkConstraints()
		iterations++
		if !changed {
			logger.Info("Pencil mark constraints stabilized after %d iteration(s)", iterations)
			break
		}
		logger.Debug("Pencil mark iteration %d: Changes detected, continuing...", iterations)
	}

	return iterations
}

// ApplyAdvancedTechniques applies advanced solving techniques like X-Wings, Swordfish, and XY-Wings
// Returns true if any candidates were eliminated
func (b *Board) ApplyAdvancedTechniques() bool {
	logger.SolvingStep("Advanced", "Trying advanced solving techniques...")

	changed := false

	// Try X-Wings (2x2 patterns)
	logger.Debug("Attempting X-Wing technique...")
	if b.applyXWings() {
		changed = true
		logger.Info("X-Wing technique found eliminations")
	}

	// Try Swordfish (3x3 patterns)
	logger.Debug("Attempting Swordfish technique...")
	if b.applySwordfish() {
		changed = true
		logger.Info("Swordfish technique found eliminations")
	}

	// Try XY-Wings
	logger.Debug("Attempting XY-Wing technique...")
	if b.applyXYWings() {
		changed = true
		logger.Info("XY-Wing technique found eliminations")
	}

	if !changed {
		logger.Debug("No advanced techniques found any eliminations")
	}

	return changed
}

// applyXWings implements the X-Wing technique
// When a candidate appears in exactly 2 cells in each of 2 rows, and those cells are in the same columns,
// that candidate can be eliminated from other cells in those columns (and vice versa for columns/rows)
func (b *Board) applyXWings() bool {
	changed := false

	// Try X-Wings in rows (eliminate from columns)
	logger.Debug("Checking for X-Wings in rows...")
	if b.applyXWingsInDirection(true) {
		changed = true
		logger.SolvingStep("X-Wing", "Found X-Wing pattern in rows")
	}

	// Try X-Wings in columns (eliminate from rows)
	logger.Debug("Checking for X-Wings in columns...")
	if b.applyXWingsInDirection(false) {
		changed = true
		logger.SolvingStep("X-Wing", "Found X-Wing pattern in columns")
	}

	return changed
}

func (b *Board) applyXWingsInDirection(rowBased bool) bool {
	changed := false

	// For each candidate 1-9
	for candidate := 1; candidate <= 9; candidate++ {
		// Build a map of line -> positions where candidate appears
		linePositions := make(map[int][]int)

		for line := 0; line < 9; line++ {
			positions := make([]int, 0)

			for pos := 0; pos < 9; pos++ {
				var cell *Cell
				if rowBased {
					cell = b.GetCellAt(line, pos)
				} else {
					cell = b.GetCellAt(pos, line)
				}

				if cell != nil && !cell.IsSolved() && cell.HasCandidate(candidate) {
					positions = append(positions, pos)
				}
			}

			// Only interested in lines with exactly 2 positions
			if len(positions) == 2 {
				linePositions[line] = positions
			}
		}

		// Now find pairs of lines with the same positions
		lines := make([]int, 0)
		for line := range linePositions {
			lines = append(lines, line)
		}

		// Check all pairs of lines
		for i := 0; i < len(lines); i++ {
			for j := i + 1; j < len(lines); j++ {
				line1, line2 := lines[i], lines[j]
				pos1, pos2 := linePositions[line1], linePositions[line2]

				// Check if positions are the same
				if len(pos1) == 2 && len(pos2) == 2 && pos1[0] == pos2[0] && pos1[1] == pos2[1] {
					// X-Wing found! Eliminate candidate from other cells in these positions
					direction := "rows"
					if !rowBased {
						direction = "columns"
					}
					logger.SolvingStep("X-Wing", "Found X-Wing for candidate %d in %s %d and %d at positions %v",
						candidate, direction, line1+1, line2+1, pos1)

					eliminatedCount := 0
					for otherLine := 0; otherLine < 9; otherLine++ {
						if otherLine != line1 && otherLine != line2 {
							for _, pos := range pos1 {
								var cell *Cell
								if rowBased {
									// Eliminate from column
									cell = b.GetCellAt(otherLine, pos)
								} else {
									// Eliminate from row
									cell = b.GetCellAt(pos, otherLine)
								}

								if cell != nil && !cell.IsSolved() && cell.HasCandidate(candidate) {
									cell.RemoveCandidate(candidate)
									changed = true
									eliminatedCount++
								}
							}
						}
					}
					logger.Info("X-Wing eliminated candidate %d from %d cell(s)", candidate, eliminatedCount)
				}
			}
		}
	}

	return changed
}

// applySwordfish implements the Swordfish technique (3x3 version of X-Wing)
func (b *Board) applySwordfish() bool {
	changed := false

	// Try Swordfish in rows (eliminate from columns)
	logger.Debug("Checking for Swordfish in rows...")
	if b.applySwordfishInDirection(true) {
		changed = true
		logger.SolvingStep("Swordfish", "Found Swordfish pattern in rows")
	}

	// Try Swordfish in columns (eliminate from rows)
	logger.Debug("Checking for Swordfish in columns...")
	if b.applySwordfishInDirection(false) {
		changed = true
		logger.SolvingStep("Swordfish", "Found Swordfish pattern in columns")
	}

	return changed
}

func (b *Board) applySwordfishInDirection(rowBased bool) bool {
	changed := false

	// For each candidate 1-9
	for candidate := 1; candidate <= 9; candidate++ {
		// Build a map of line -> positions where candidate appears
		linePositions := make(map[int][]int)

		for line := 0; line < 9; line++ {
			positions := make([]int, 0)

			for pos := 0; pos < 9; pos++ {
				var cell *Cell
				if rowBased {
					cell = b.GetCellAt(line, pos)
				} else {
					cell = b.GetCellAt(pos, line)
				}

				if cell != nil && !cell.IsSolved() && cell.HasCandidate(candidate) {
					positions = append(positions, pos)
				}
			}

			// Only interested in lines with 2 or 3 positions
			if len(positions) >= 2 && len(positions) <= 3 {
				linePositions[line] = positions
			}
		}

		// Now find triples of lines that cover exactly 3 positions
		lines := make([]int, 0)
		for line := range linePositions {
			lines = append(lines, line)
		}

		// Check all triples of lines
		for i := 0; i < len(lines); i++ {
			for j := i + 1; j < len(lines); j++ {
				for k := j + 1; k < len(lines); k++ {
					line1, line2, line3 := lines[i], lines[j], lines[k]

					// Get union of positions
					posUnion := make(map[int]bool)
					for _, pos := range linePositions[line1] {
						posUnion[pos] = true
					}
					for _, pos := range linePositions[line2] {
						posUnion[pos] = true
					}
					for _, pos := range linePositions[line3] {
						posUnion[pos] = true
					}

					// If exactly 3 positions, we have a Swordfish
					if len(posUnion) == 3 {
						// Eliminate candidate from other cells in these positions
						positions := make([]int, 0)
						for pos := range posUnion {
							positions = append(positions, pos)
						}

						direction := "rows"
						if !rowBased {
							direction = "columns"
						}
						logger.SolvingStep("Swordfish", "Found Swordfish for candidate %d in %s %d, %d, %d at positions %v",
							candidate, direction, line1+1, line2+1, line3+1, positions)

						eliminatedCount := 0
						for otherLine := 0; otherLine < 9; otherLine++ {
							if otherLine != line1 && otherLine != line2 && otherLine != line3 {
								for _, pos := range positions {
									var cell *Cell
									if rowBased {
										cell = b.GetCellAt(otherLine, pos)
									} else {
										cell = b.GetCellAt(pos, otherLine)
									}

									if cell != nil && !cell.IsSolved() && cell.HasCandidate(candidate) {
										cell.RemoveCandidate(candidate)
										changed = true
										eliminatedCount++
									}
								}
							}
						}
						logger.Info("Swordfish eliminated candidate %d from %d cell(s)", candidate, eliminatedCount)
					}
				}
			}
		}
	}

	return changed
}

// applyXYWings implements the XY-Wing technique
// Finds a pivot cell with 2 candidates (XY), and two wing cells (XZ and YZ)
// If both wings share a candidate (Z), it can be eliminated from cells that see both wings
func (b *Board) applyXYWings() bool {
	changed := false

	// Find all cells with exactly 2 candidates (potential pivots and wings)
	cells2Cands := make([]*Cell, 0)
	for idx := 0; idx < 81; idx++ {
		cell := b.GetCell(idx)
		if cell != nil && !cell.IsSolved() && cell.CandidateCount() == 2 {
			cells2Cands = append(cells2Cands, cell)
		}
	}

	logger.Debug("Found %d cells with exactly 2 candidates for XY-Wing analysis", len(cells2Cands))

	// Try each cell as a pivot
	for _, pivot := range cells2Cands {
		pivotCands := utils.GetCandidatesAsSlice(pivot.GetCandidates())
		if len(pivotCands) != 2 {
			continue
		}
		X, Y := pivotCands[0], pivotCands[1]

		// Find cells that share a constraint with pivot and have exactly 2 candidates
		visibleCells := b.getVisibleCells(pivot)

		// Find two wings: one with {X, Z} and one with {Y, Z}
		for _, wing1 := range visibleCells {
			if wing1.CandidateCount() != 2 {
				continue
			}

			wing1Cands := utils.GetCandidatesAsSlice(wing1.GetCandidates())
			if len(wing1Cands) != 2 {
				continue
			}

			// Check if wing1 shares exactly one candidate with pivot
			sharedWithPivot1 := 0
			var Z1 int
			if wing1Cands[0] == X || wing1Cands[0] == Y {
				sharedWithPivot1++
			} else {
				Z1 = wing1Cands[0]
			}
			if wing1Cands[1] == X || wing1Cands[1] == Y {
				sharedWithPivot1++
			} else {
				Z1 = wing1Cands[1]
			}

			if sharedWithPivot1 != 1 {
				continue
			}

			// Now find wing2
			for _, wing2 := range visibleCells {
				if wing2 == wing1 || wing2.CandidateCount() != 2 {
					continue
				}

				wing2Cands := utils.GetCandidatesAsSlice(wing2.GetCandidates())
				if len(wing2Cands) != 2 {
					continue
				}

				// Check if wing2 shares exactly one candidate with pivot (different from wing1)
				sharedWithPivot2 := 0
				var Z2 int
				if wing2Cands[0] == X || wing2Cands[0] == Y {
					sharedWithPivot2++
				} else {
					Z2 = wing2Cands[0]
				}
				if wing2Cands[1] == X || wing2Cands[1] == Y {
					sharedWithPivot2++
				} else {
					Z2 = wing2Cands[1]
				}

				if sharedWithPivot2 != 1 {
					continue
				}

				// Check if both wings share the same Z candidate
				if Z1 != Z2 {
					continue
				}
				Z := Z1

				logger.SolvingStep("XY-Wing", "Found XY-Wing: Pivot R%dC%d {%d,%d}, Wing1 R%dC%d, Wing2 R%dC%d, eliminating %d",
					pivot.GetRow()+1, pivot.GetCol()+1, X, Y,
					wing1.GetRow()+1, wing1.GetCol()+1,
					wing2.GetRow()+1, wing2.GetCol()+1, Z)

				// XY-Wing found! Eliminate Z from cells that see both wings
				wing1Visible := b.getVisibleCells(wing1)
				wing2Visible := b.getVisibleCells(wing2)

				eliminatedCount := 0

				for _, cell := range wing1Visible {
					if cell == pivot || cell == wing1 || cell == wing2 {
						continue
					}

					// Check if this cell also sees wing2
					seesWing2 := false
					for _, w2v := range wing2Visible {
						if w2v == cell {
							seesWing2 = true
							break
						}
					}

					if seesWing2 && cell.HasCandidate(Z) {
						cell.RemoveCandidate(Z)
						changed = true
						eliminatedCount++
					}
				}

				if eliminatedCount > 0 {
					logger.Info("XY-Wing eliminated candidate %d from %d cell(s)", Z, eliminatedCount)
				}
			}
		}
	}

	return changed
}

// getVisibleCells returns all cells that share at least one constraint with the given cell
func (b *Board) getVisibleCells(cell *Cell) []*Cell {
	visibleMap := make(map[*Cell]bool)

	// Find all constraints that include this cell
	for _, constraint := range b.constraints {
		cellIncluded := false
		for _, idx := range constraint.GetCells() {
			if idx == cell.GetIndex() {
				cellIncluded = true
				break
			}
		}

		if cellIncluded {
			// Add all other cells in this constraint
			for _, idx := range constraint.GetCells() {
				otherCell := b.GetCell(idx)
				if otherCell != nil && otherCell != cell {
					visibleMap[otherCell] = true
				}
			}
		}
	}

	// Convert map to slice
	visible := make([]*Cell, 0, len(visibleMap))
	for c := range visibleMap {
		visible = append(visible, c)
	}

	return visible
}

// AddObserver adds an observer to all cells on the board
func (b *Board) AddObserver(obs observer.CellObserver) {
	if obs == nil {
		return
	}

	b.observers = append(b.observers, obs)

	// Add observer to all cells
	for i := 0; i < 81; i++ {
		if b.board[i] != nil {
			b.board[i].AddObserver(obs)
		}
	}

	logger.Debug("Added observer to all board cells")
}

// RemoveObserver removes an observer from all cells
func (b *Board) RemoveObserver(obs observer.CellObserver) {
	// Remove from board's observer list
	for i, o := range b.observers {
		if o == obs {
			b.observers = append(b.observers[:i], b.observers[i+1:]...)
			break
		}
	}

	// Remove from all cells
	for i := 0; i < 81; i++ {
		if b.board[i] != nil && b.board[i].GetNotifier() != nil {
			b.board[i].GetNotifier().RemoveObserver(obs)
		}
	}

	logger.Debug("Removed observer from all board cells")
}
