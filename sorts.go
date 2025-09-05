package main

import "context"

type sortFunc func(ctx context.Context, arr arrObj)

var sorts = []struct {
	name string
	arg  string
	fun  sortFunc
}{
	{"Quick Sort", "quick", quickSort},
	{"Bubble Sort", "bubble", bubbleSort},
	{"Selection Sort", "selection", selectionSort},
	{"Insertion Sort", "insertion", insertionSort},
}

func bubbleSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr.get(ctx, j) > arr.get(ctx, j+1) {
				arr.swap(ctx, j, j+1)
			}
		}
	}
}

func selectionSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr.get(ctx, j) < arr.get(ctx, minIdx) {
				minIdx = j
			}
		}
		if minIdx != i {
			arr.swap(ctx, i, minIdx)
		}
	}
}

func insertionSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	for i := 1; i < n; i++ {
		key := arr.get(ctx, i)
		j := i - 1
		for j >= 0 && arr.get(ctx, j) > key {
			arr.set(ctx, j+1, arr.get(ctx, j))
			j--
		}
		arr.set(ctx, j+1, key)
	}
}

func quickSort(ctx context.Context, arr arrObj) {
	quickSortRecurse(ctx, arr, 0, arr.len()-1)
}

func quickSortRecurse(ctx context.Context, arr arrObj, low, high int) {
	if low < high {
		pi := partition(ctx, arr, low, high)
		quickSortRecurse(ctx, arr, low, pi-1)
		quickSortRecurse(ctx, arr, pi+1, high)
	}
}

func partition(ctx context.Context, arr arrObj, low, high int) int {
	pivot := arr.get(ctx, high)
	i := low - 1
	for j := low; j < high; j++ {
		if arr.get(ctx, j) < pivot {
			i++
			arr.swap(ctx, i, j)
		}
	}
	arr.swap(ctx, i+1, high)
	return i + 1
}
