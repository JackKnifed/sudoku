package sudoku

import (
	"sort"
)

// Holds some functions needed for working with arrays

func dedupArr(arr []int) []int {
	localArr := make([]int, len(arr))
	copy(localArr, arr)
	sort.Ints(localArr)
	for i := 0; i < len(localArr)-1; i++ {
		if localArr[i] == localArr[i+1] {
			localArr = append(localArr[:i], localArr[i+1:]...)
		}
	}
	if len(localArr) > 1 && localArr[len(localArr)-2] == localArr[len(localArr)-1] {
		localArr = localArr[:len(localArr)-1]
	}
	return localArr
}

func inArr(arr []int, val int) bool {
	for _, each := range arr {
		if each == val {
			return true
		}
	}
	return false
}

func anyInArr(arr, tgt []int) bool {
	for _, t := range tgt {
		for _, a := range arr {
			if t == a {
				return true
			}
		}
	}
	return false
}

// preforms a union of a and b and removes anhy duplicates
func addArr(a, b []int) []int {
	return dedupArr(append(a, b...))
	// a = dedupArr(a)
	// b = dedupArr(b)
	// var i, j int
	// for i < len(a) && j < len(b)-1 {
	// 	if a[i] == b[j] {
	// 		j++
	// 	} else if a[i] < b[j] {
	// 		i++
	// 	} else if a[i] > b[i] {
	// 		b = append(b[:j], append(a[i:i], b[j:]...)...)
	// 		j++
	// 	}
	// }
	// b = dedupArr(b)
	// return b
}

// Subtracts b from a - removing any intersections from a and returning
func subArr(a, b []int) []int {
	b = dedupArr(b)
	a = dedupArr(a)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else if i+1 < len(a) {
			a = append(a[:i], a[i+1:]...)
		} else {
			a = a[:i]
		}
	}
	return a
}

func andArr(a, b []int) []int {
	output := []int{}
	a = dedupArr(a)
	b = dedupArr(b)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			output = append(output, a[i])
			i++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}
	return output
}
