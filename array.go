package sudoku

import (
	"sort"
)

// Holds some functions needed for working with arrays

func inArr(val int, arr []int) bool {
	for _, each := range arr {
		if each == val {
			return true
		}
	}
	return false
}

func anyInArr(tgt, arr []int) bool {
	for _, t := range tgt {
		for _, a := range arr {
			if t == a {
				return true
			}
		}
	}
	return false
}

func addArr(add, arr []int) []int {
	add = dedupArr(add)
	arr = dedupArr(arr)
	var i, j int
	for i < len(add) && j < len(arr)-1 {
		if add[i] == arr[j] {
			j++
		} else if add[i] < arr[j] {
			i++
		} else if add[i] > arr[i] {
			arr = append(arr[:j], add[i], arr[j:]...)
			j++
		}
	}
	arr = defupArr(arr)
	return arr
}

func subArr(sub, arr []int) []int {
	sub = dedupArr(sub)
	arr = dedupArr(arr)
	var i, j int
	for i < len(sub) && j < len(arr)-1 {
		if sub[i] == arr[j] {
			arr = append(arr[:j], arr[j+1]...)
		} else if sub[i] < arr[j] {
			i++
		} else {
			j++
		}
	}
	return arr
}

func dedupArr(arr []int) []int {
	sort.Ints(arr)
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] == arr[i+1] {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func andArr(a, b []int) []int {
	var output []int
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {

		}
	}
}
