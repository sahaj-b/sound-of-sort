package main

// type arrObj interface {
// 	get(ind int) int
// 	set(ind, val int)
// 	swap(i, j int)
// 	len() int
// }

func bubbleSort(arr arrObj, done chan struct{}) {
	defer close(done)
	n := arr.len()
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr.get(j) > arr.get(j+1) {
				arr.swap(j, j+1)
			}
		}
	}
}

func selectionSort(arr arrObj, done chan struct{}) {
	defer close(done)
	n := arr.len()
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr.get(j) < arr.get(minIdx) {
				minIdx = j
			}
		}
		if minIdx != i {
			arr.swap(i, minIdx)
		}
	}
}

func insertionSort(arr arrObj, done chan struct{}) {
	defer close(done)
	n := arr.len()
	for i := 1; i < n; i++ {
		key := arr.get(i)
		j := i - 1
		for j >= 0 && arr.get(j) > key {
			arr.set(j+1, arr.get(j))
			j--
		}
		arr.set(j+1, key)
	}
}

func quickSort(arr arrObj, done chan struct{}) {
	defer close(done)
	quickSortRecurse(arr, 0, arr.len()-1)
}

func quickSortRecurse(arr arrObj, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSortRecurse(arr, low, pi-1)
		quickSortRecurse(arr, pi+1, high)
	}
}

func partition(arr arrObj, low, high int) int {
	pivot := arr.get(high)
	i := low - 1
	for j := low; j < high; j++ {
		if arr.get(j) < pivot {
			i++
			arr.swap(i, j)
		}
	}
	arr.swap(i+1, high)
	return i + 1
}
