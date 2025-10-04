package lib_test

import (
	"reflect"
	"testing"

	"github.com/eftil/sudoku-solver.git/lib/utils"
)

func TestHasUniqueNonZeros(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected bool
	}{
		{
			name:     "all zeros",
			values:   []int{0, 0, 0, 0, 0},
			expected: true,
		},
		{
			name:     "unique values",
			values:   []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "unique with zeros",
			values:   []int{1, 0, 2, 0, 3, 0, 4},
			expected: true,
		},
		{
			name:     "duplicate non-zero",
			values:   []int{1, 2, 3, 2, 5},
			expected: false,
		},
		{
			name:     "duplicate non-zero with zeros",
			values:   []int{1, 0, 2, 0, 2, 0, 3},
			expected: false,
		},
		{
			name:     "invalid value (negative)",
			values:   []int{1, 2, -1, 4, 5},
			expected: false,
		},
		{
			name:     "invalid value (>9)",
			values:   []int{1, 2, 10, 4, 5},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.HasUniqueNonZeros(tt.values)
			if result != tt.expected {
				t.Errorf("HasUniqueNonZeros(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func TestGenerateCombinations(t *testing.T) {
	tests := []struct {
		name     string
		n, k     int
		expected [][]int
	}{
		{
			name:     "choose 2 from 3",
			n:        3,
			k:        2,
			expected: [][]int{{0, 1}, {0, 2}, {1, 2}},
		},
		{
			name:     "choose 1 from 3",
			n:        3,
			k:        1,
			expected: [][]int{{0}, {1}, {2}},
		},
		{
			name:     "choose 3 from 3",
			n:        3,
			k:        3,
			expected: [][]int{{0, 1, 2}},
		},
		{
			name:     "choose 0 from 3",
			n:        3,
			k:        0,
			expected: [][]int{{}},
		},
		{
			name:     "k > n",
			n:        2,
			k:        3,
			expected: [][]int{},
		},
		{
			name:     "negative k",
			n:        3,
			k:        -1,
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GenerateCombinations(tt.n, tt.k)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GenerateCombinations(%d, %d) = %v, want %v", tt.n, tt.k, result, tt.expected)
			}
		})
	}
}

func TestContainsInt(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		target   int
		expected bool
	}{
		{"found in middle", []int{1, 2, 3, 4, 5}, 3, true},
		{"found at start", []int{1, 2, 3, 4, 5}, 1, true},
		{"found at end", []int{1, 2, 3, 4, 5}, 5, true},
		{"not found", []int{1, 2, 3, 4, 5}, 6, false},
		{"empty slice", []int{}, 1, false},
		{"single element found", []int{5}, 5, true},
		{"single element not found", []int{5}, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ContainsInt(tt.slice, tt.target)
			if result != tt.expected {
				t.Errorf("ContainsInt(%v, %d) = %v, want %v", tt.slice, tt.target, result, tt.expected)
			}
		})
	}
}

func TestIndexToRowCol(t *testing.T) {
	tests := []struct {
		index       int
		expectedRow int
		expectedCol int
	}{
		{0, 0, 0},
		{1, 0, 1},
		{8, 0, 8},
		{9, 1, 0},
		{10, 1, 1},
		{40, 4, 4},
		{80, 8, 8},
	}

	for _, tt := range tests {
		row, col := utils.IndexToRowCol(tt.index)
		if row != tt.expectedRow || col != tt.expectedCol {
			t.Errorf("IndexToRowCol(%d) = (%d, %d), want (%d, %d)",
				tt.index, row, col, tt.expectedRow, tt.expectedCol)
		}
	}

	// Test invalid indices
	row, col := utils.IndexToRowCol(-1)
	if row != -1 || col != -1 {
		t.Errorf("IndexToRowCol(-1) should return (-1, -1), got (%d, %d)", row, col)
	}

	row, col = utils.IndexToRowCol(81)
	if row != -1 || col != -1 {
		t.Errorf("IndexToRowCol(81) should return (-1, -1), got (%d, %d)", row, col)
	}
}

func TestRowColToIndex(t *testing.T) {
	tests := []struct {
		row, col      int
		expectedIndex int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{0, 8, 8},
		{1, 0, 9},
		{1, 1, 10},
		{4, 4, 40},
		{8, 8, 80},
	}

	for _, tt := range tests {
		index := utils.RowColToIndex(tt.row, tt.col)
		if index != tt.expectedIndex {
			t.Errorf("RowColToIndex(%d, %d) = %d, want %d",
				tt.row, tt.col, index, tt.expectedIndex)
		}
	}

	// Test invalid coordinates
	if utils.RowColToIndex(-1, 0) != -1 {
		t.Error("RowColToIndex with negative row should return -1")
	}

	if utils.RowColToIndex(0, -1) != -1 {
		t.Error("RowColToIndex with negative col should return -1")
	}

	if utils.RowColToIndex(9, 0) != -1 {
		t.Error("RowColToIndex with row >= 9 should return -1")
	}

	if utils.RowColToIndex(0, 9) != -1 {
		t.Error("RowColToIndex with col >= 9 should return -1")
	}
}

func TestGetBoxNumber(t *testing.T) {
	tests := []struct {
		row, col    int
		expectedBox int
	}{
		{0, 0, 0},
		{0, 1, 0},
		{0, 2, 0},
		{0, 3, 1},
		{0, 6, 2},
		{3, 0, 3},
		{4, 4, 4},
		{6, 6, 8},
		{8, 8, 8},
	}

	for _, tt := range tests {
		box := utils.GetBoxNumber(tt.row, tt.col)
		if box != tt.expectedBox {
			t.Errorf("GetBoxNumber(%d, %d) = %d, want %d",
				tt.row, tt.col, box, tt.expectedBox)
		}
	}

	// Test invalid coordinates
	if utils.GetBoxNumber(-1, 0) != -1 {
		t.Error("GetBoxNumber with negative row should return -1")
	}

	if utils.GetBoxNumber(0, 9) != -1 {
		t.Error("GetBoxNumber with col >= 9 should return -1")
	}
}

func TestGetBoxCoordinates(t *testing.T) {
	tests := []struct {
		box                      int
		expectedRow, expectedCol int
	}{
		{0, 0, 0},
		{1, 0, 3},
		{2, 0, 6},
		{3, 3, 0},
		{4, 3, 3},
		{5, 3, 6},
		{6, 6, 0},
		{7, 6, 3},
		{8, 6, 6},
	}

	for _, tt := range tests {
		row, col := utils.GetBoxCoordinates(tt.box)
		if row != tt.expectedRow || col != tt.expectedCol {
			t.Errorf("GetBoxCoordinates(%d) = (%d, %d), want (%d, %d)",
				tt.box, row, col, tt.expectedRow, tt.expectedCol)
		}
	}

	// Test invalid box numbers
	row, col := utils.GetBoxCoordinates(-1)
	if row != -1 || col != -1 {
		t.Error("GetBoxCoordinates(-1) should return (-1, -1)")
	}

	row, col = utils.GetBoxCoordinates(9)
	if row != -1 || col != -1 {
		t.Error("GetBoxCoordinates(9) should return (-1, -1)")
	}
}

func TestGetCandidatesAsSlice(t *testing.T) {
	tests := []struct {
		name       string
		candidates map[int]bool
		expected   []int
	}{
		{
			name:       "all candidates",
			candidates: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true},
			expected:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:       "subset of candidates",
			candidates: map[int]bool{2: true, 5: true, 7: true},
			expected:   []int{2, 5, 7},
		},
		{
			name:       "single candidate",
			candidates: map[int]bool{9: true},
			expected:   []int{9},
		},
		{
			name:       "no candidates",
			candidates: map[int]bool{},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetCandidatesAsSlice(tt.candidates)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetCandidatesAsSlice(%v) = %v, want %v", tt.candidates, result, tt.expected)
			}
		})
	}
}

func TestMaxMin(t *testing.T) {
	if utils.Max(5, 10) != 10 {
		t.Error("Max(5, 10) should be 10")
	}

	if utils.Max(10, 5) != 10 {
		t.Error("Max(10, 5) should be 10")
	}

	if utils.Max(5, 5) != 5 {
		t.Error("Max(5, 5) should be 5")
	}

	if utils.Min(5, 10) != 5 {
		t.Error("Min(5, 10) should be 5")
	}

	if utils.Min(10, 5) != 5 {
		t.Error("Min(10, 5) should be 5")
	}

	if utils.Min(5, 5) != 5 {
		t.Error("Min(5, 5) should be 5")
	}
}
