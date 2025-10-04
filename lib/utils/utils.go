package utils

// HasUniqueNonZeros checks if all non-zero values in a slice are unique
func HasUniqueNonZeros(values []int) bool {
	seen := make(map[int]bool)
	for _, v := range values {
		if v == 0 {
			continue // 0 means empty cell, skip it
		}
		if v < 1 || v > 9 {
			return false // invalid value
		}
		if seen[v] {
			return false // duplicate found
		}
		seen[v] = true
	}
	return true
}

// GenerateCombinations generates all combinations of n items taken k at a time
func GenerateCombinations(n, k int) [][]int {
	if k > n || k < 0 {
		return [][]int{}
	}

	result := make([][]int, 0)
	combination := make([]int, k)

	var generate func(start, depth int)
	generate = func(start, depth int) {
		if depth == k {
			temp := make([]int, k)
			copy(temp, combination)
			result = append(result, temp)
			return
		}

		for i := start; i < n; i++ {
			combination[depth] = i
			generate(i+1, depth+1)
		}
	}

	generate(0, 0)
	return result
}

// ContainsInt checks if an int is in a slice of ints
func ContainsInt(slice []int, target int) bool {
	for _, val := range slice {
		if val == target {
			return true
		}
	}
	return false
}

// IndexToRowCol converts a cell index (0-80) to row and column coordinates
func IndexToRowCol(index int) (row, col int) {
	if index < 0 || index > 80 {
		return -1, -1
	}
	return index / 9, index % 9
}

// RowColToIndex converts row and column coordinates to a cell index (0-80)
func RowColToIndex(row, col int) int {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return -1
	}
	return row*9 + col
}

// GetBoxNumber returns the box number (0-8) for a given row and column
func GetBoxNumber(row, col int) int {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return -1
	}
	return (row/3)*3 + (col / 3)
}

// GetBoxCoordinates returns the starting row and column for a given box number
func GetBoxCoordinates(box int) (startRow, startCol int) {
	if box < 0 || box > 8 {
		return -1, -1
	}
	return (box / 3) * 3, (box % 3) * 3
}

// GetCandidatesAsSlice converts a candidate map to a sorted slice
func GetCandidatesAsSlice(candidates map[int]bool) []int {
	result := make([]int, 0, len(candidates))
	for i := 1; i <= 9; i++ {
		if candidates[i] {
			result = append(result, i)
		}
	}
	return result
}

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
