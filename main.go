package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eftil/sudoku-solver.git/lib"
	"github.com/eftil/sudoku-solver.git/lib/constraints"
	"github.com/eftil/sudoku-solver.git/lib/logger"
	"github.com/eftil/sudoku-solver.git/lib/observer"
)

func main() {
	// Configure logger
	// Change to logger.DEBUG to see detailed solving steps
	logger.SetLevel(logger.INFO)
	logger.SetOutput(os.Stdout)

	fmt.Println("=== Sudoku Solver - Comprehensive Demo ===")

	// Create a new board
	board := lib.NewBoard()

	// Create and add an auto-solver observer
	autoSolver := observer.NewAutoSolverObserver()
	board.AddObserver(autoSolver)

	fmt.Println("\n✓ Observer system enabled - will detect cells with single candidates")

	// Add standard sudoku constraints (rows, columns, boxes)
	logger.Info("Adding standard Sudoku constraints...")
	for i := 0; i < 9; i++ {
		rowConstraint, err := constraints.NewRowConstraint(i)
		if err != nil {
			log.Fatalf("Failed to create row constraint: %v", err)
		}
		board.AddConstraint(rowConstraint)

		colConstraint, err := constraints.NewColumnConstraint(i)
		if err != nil {
			log.Fatalf("Failed to create column constraint: %v", err)
		}
		board.AddConstraint(colConstraint)

		boxConstraint, err := constraints.NewBoxConstraint(i)
		if err != nil {
			log.Fatalf("Failed to create box constraint: %v", err)
		}
		board.AddConstraint(boxConstraint)
	}

	fmt.Println("\n=== Example 1: Standard Sudoku ===")

	// Set up a simple pattern that demonstrates candidate elimination
	logger.Info("\nSetting up initial values...")
	board.Set(0, 0, 5)
	board.Set(0, 1, 3)
	board.Set(0, 2, 4)
	board.Set(0, 3, 6)
	board.Set(0, 4, 7)
	board.Set(0, 5, 8)
	board.Set(0, 6, 9)
	board.Set(0, 7, 1)
	// Cell R1C9 now has only candidate 2 remaining!

	// Print the board
	fmt.Println("\nCurrent board state:")
	board.Print()

	// Validate all constraints
	fmt.Println("\n=== Validation ===")
	valid, err := board.ValidateAll()
	if err != nil {
		log.Fatalf("Error during validation: %v", err)
	}

	if valid {
		fmt.Println("✓ All constraints are satisfied!")
	} else {
		fmt.Println("✗ Some constraints are violated")
	}

	// Print all constraints
	fmt.Println("\n=== Active Constraints ===")
	for i, constraint := range board.GetConstraints() {
		fmt.Printf("%d. %s\n", i+1, constraint.GetName())
	}

	fmt.Println("\n=== Example 2: Advanced Sudoku with Special Constraints ===")

	// Create a new board for variant sudoku
	variantBoard := lib.NewBoard()
	autoSolver2 := observer.NewAutoSolverObserver()
	variantBoard.AddObserver(autoSolver2)

	// Add standard constraints
	for i := 0; i < 9; i++ {
		rc, _ := constraints.NewRowConstraint(i)
		variantBoard.AddConstraint(rc)
		cc, _ := constraints.NewColumnConstraint(i)
		variantBoard.AddConstraint(cc)
		bc, _ := constraints.NewBoxConstraint(i)
		variantBoard.AddConstraint(bc)
	}

	// Add a killer cage constraint (cells in top-left that sum to 15)
	killerCells := []int{0, 1, 9} // row 0 col 0, row 0 col 1, row 1 col 0
	killerConstraint, err := constraints.NewKillerCageConstraint(killerCells, 15)
	if err != nil {
		log.Fatalf("Failed to create killer cage: %v", err)
	}
	variantBoard.AddConstraint(killerConstraint)
	fmt.Printf("\n✓ Added Killer Cage: %s\n", killerConstraint.GetDescription())

	// Add a German whispers line
	whisperCells := []int{4, 13, 22} // diagonal line from top center
	whisperConstraint, err := constraints.NewGermanWhispersConstraint(whisperCells)
	if err != nil {
		log.Fatalf("Failed to create German whispers: %v", err)
	}
	variantBoard.AddConstraint(whisperConstraint)
	fmt.Printf("✓ Added German Whispers: %s\n", whisperConstraint.GetDescription())

	// Add a Renban line
	renbanCells := []int{36, 37, 38} // horizontal line in middle
	renbanConstraint, err := constraints.NewRenbanConstraint(renbanCells)
	if err != nil {
		log.Fatalf("Failed to create Renban: %v", err)
	}
	variantBoard.AddConstraint(renbanConstraint)
	fmt.Printf("✓ Added Renban: %s\n", renbanConstraint.GetDescription())

	// Set some values that satisfy the killer cage
	variantBoard.Set(0, 0, 5)
	variantBoard.Set(0, 1, 6)
	variantBoard.Set(1, 0, 4)

	fmt.Println("\nVariant board state:")
	variantBoard.Print()

	// Validate the variant board
	fmt.Println("\n=== Variant Board Validation ===")
	valid, err = variantBoard.ValidateAll()
	if err != nil {
		log.Fatalf("Error during validation: %v", err)
	}

	if valid {
		fmt.Println("✓ All constraints (including special rules) are satisfied!")
	} else {
		fmt.Println("✗ Some constraints are violated")
	}

	// Demonstrate pencil mark techniques
	fmt.Println("\n=== Demonstrating Advanced Solving Techniques ===")
	logger.Info("\nApplying pencil mark constraints...")

	iterations := board.ApplyPencilMarkConstraintsUntilStable()
	fmt.Printf("Pencil mark techniques converged after %d iteration(s)\n", iterations)

	// Try advanced techniques
	logger.Info("\nTrying advanced techniques (X-Wing, Swordfish, XY-Wing)...")
	if board.ApplyAdvancedTechniques() {
		fmt.Println("✓ Advanced techniques found eliminations")
	} else {
		fmt.Println("• Advanced techniques did not find additional eliminations")
	}

	// Show observer results
	fmt.Println("\n=== Observer Statistics ===")
	fmt.Printf("Total cells auto-solved: %d\n", autoSolver.GetSolutionCount())
	if len(autoSolver.GetCellsToSolve()) > 0 {
		fmt.Printf("Cells identified for auto-solving: %d\n", len(autoSolver.GetCellsToSolve()))
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nTo see detailed solving steps, set logger.SetLevel(logger.DEBUG) at the top of main()")
	fmt.Println("The observer pattern allows automatic detection and solving of cells with single candidates")
	fmt.Println("All solving decisions are logged with explanations for transparency")
}
