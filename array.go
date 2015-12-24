package sudoku

import (
	"sort"
)

func dedupArr(arr []int) []int {
	localArr := make([]int, len(arr))
	copy(localArr, arr)
	sort.Ints(localArr)
	for i := 0; i < len(localArr)-1; i++ {
		if localArr[i] == localArr[i+1] {
			localArr = append(localArr[:i], localArr[i+1:]...)
		}
	}
	if localArr[len(localArr)-2] == localArr[len(localArr)-1] {
		localArr = localArr[:len(localArr)-1]
	}
	return localArr
}

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
			arr = append(arr[:j], append(add[i:i], arr[j:]...)...)
			j++
		}
	}
	arr = dedupArr(arr)
	return arr
}

func subArr(sub, arr []int) []int {
	sub = dedupArr(sub)
	arr = dedupArr(arr)
	var i, j int
	for i < len(sub) && j < len(arr)-1 {
		if sub[i] == arr[j] {
			arr = append(arr[:j], arr[:j+1]...)
		} else if sub[i] < arr[j] {
			i++
		} else {
			j++
		}
	}
	return arr
}

func andArr(a, b []int) []int {
	var output []int
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
