package main

import (
	"context"
	"math/rand"
)

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
	{"Merge Sort", "merge", mergeSort},
	{"Heap Sort", "heap", heapSort},
	{"Shell Sort", "shell", shellSort},
	{"Cocktail Shaker Sort", "cocktail", cocktailShakerSort},
	{"Gnome Sort", "gnome", gnomeSort},
	{"Pancake Sort", "pancake", pancakeSort},
	{"Radix Sort (LSD)", "radix", radixSortLSD},
	{"Timsort", "timsort", timsort},
	{"Bitonic Sort", "bitonic", bitonicSort},
	{"Bogo Sort", "bogo", bogoSort},
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

func mergeSort(ctx context.Context, arr arrObj) {
	mergeSortRecurse(ctx, arr, 0, arr.len()-1)
}

func mergeSortRecurse(ctx context.Context, arr arrObj, left, right int) {
	if left < right {
		mid := (left + right) / 2
		mergeSortRecurse(ctx, arr, left, mid)
		mergeSortRecurse(ctx, arr, mid+1, right)
		merge(ctx, arr, left, mid, right)
	}
}

func merge(ctx context.Context, arr arrObj, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	leftArr := make([]int, n1)
	rightArr := make([]int, n2)

	for i := range n1 {
		leftArr[i] = arr.get(ctx, left+i)
	}
	for j := range n2 {
		rightArr[j] = arr.get(ctx, mid+1+j)
	}

	i, j, k := 0, 0, left
	for i < n1 && j < n2 {
		if leftArr[i] <= rightArr[j] {
			arr.set(ctx, k, leftArr[i])
			i++
		} else {
			arr.set(ctx, k, rightArr[j])
			j++
		}
		k++
	}

	for i < n1 {
		arr.set(ctx, k, leftArr[i])
		i++
		k++
	}

	for j < n2 {
		arr.set(ctx, k, rightArr[j])
		j++
		k++
	}
}

func heapSort(ctx context.Context, arr arrObj) {
	n := arr.len()

	for i := n/2 - 1; i >= 0; i-- {
		heapify(ctx, arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr.swap(ctx, 0, i)
		heapify(ctx, arr, i, 0)
	}
}

func heapify(ctx context.Context, arr arrObj, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr.get(ctx, left) > arr.get(ctx, largest) {
		largest = left
	}

	if right < n && arr.get(ctx, right) > arr.get(ctx, largest) {
		largest = right
	}

	if largest != i {
		arr.swap(ctx, i, largest)
		heapify(ctx, arr, n, largest)
	}
}

func shellSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	for gap := n / 2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			temp := arr.get(ctx, i)
			j := i
			for ; j >= gap && arr.get(ctx, j-gap) > temp; j -= gap {
				arr.set(ctx, j, arr.get(ctx, j-gap))
			}
			arr.set(ctx, j, temp)
		}
	}
}

func cocktailShakerSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	swapped := true
	start := 0
	end := n - 1

	for swapped {
		swapped = false
		for i := start; i < end; i++ {
			if arr.get(ctx, i) > arr.get(ctx, i+1) {
				arr.swap(ctx, i, i+1)
				swapped = true
			}
		}
		if !swapped {
			break
		}
		swapped = false
		end--
		for i := end - 1; i >= start; i-- {
			if arr.get(ctx, i) > arr.get(ctx, i+1) {
				arr.swap(ctx, i, i+1)
				swapped = true
			}
		}
		start++
	}
}

func gnomeSort(ctx context.Context, arr arrObj) {
	index := 0
	n := arr.len()
	for index < n {
		if index == 0 {
			index++
		}
		if arr.get(ctx, index) >= arr.get(ctx, index-1) {
			index++
		} else {
			arr.swap(ctx, index, index-1)
			index--
		}
	}
}

func pancakeSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	for currSize := n; currSize > 1; currSize-- {
		mi := 0
		for i := 0; i < currSize; i++ {
			if arr.get(ctx, i) > arr.get(ctx, mi) {
				mi = i
			}
		}

		if mi != currSize-1 {
			pancakeFlip(ctx, arr, mi+1)
			pancakeFlip(ctx, arr, currSize)
		}
	}
}

func pancakeFlip(ctx context.Context, arr arrObj, k int) {
	start := 0
	for start < k-1 {
		arr.swap(ctx, start, k-1)
		start++
		k--
	}
}

func radixCountingSort(arr []int, exp int) {
	n := len(arr)
	output := make([]int, n)
	count := make([]int, 10)

	for i := range n {
		count[(arr[i]/exp)%10]++
	}

	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	for i := n - 1; i >= 0; i-- {
		index := (arr[i] / exp) % 10
		output[count[index]-1] = arr[i]
		count[index]--
	}

	for i := range n {
		arr[i] = output[i]
	}
}

func radixSortLSD(ctx context.Context, arr arrObj) {
	n := arr.len()
	if n == 0 {
		return
	}

	// Find the maximum number to know the number of digits
	maxVal := arr.get(ctx, 0)
	for i := 1; i < n; i++ {
		if arr.get(ctx, i) > maxVal {
			maxVal = arr.get(ctx, i)
		}
	}

	// Do counting sort for every digit.
	for exp := 1; maxVal/exp > 0; exp *= 10 {
		// Use buckets to store numbers based on the current digit
		buckets := make([][]int, 10)

		// Distribute array elements into buckets
		for i := range n {
			val := arr.get(ctx, i)
			bucketIndex := (val / exp) % 10
			buckets[bucketIndex] = append(buckets[bucketIndex], val)
		}

		// Gather elements from buckets back into the array
		i := 0
		for _, bucket := range buckets {
			for _, val := range bucket {
				arr.set(ctx, i, val)
				i++
			}
		}
	}
}

// bogoSort (Permutation Sort) is a highly inefficient sorting algorithm based on the
// generate and test paradigm. It successively generates permutations of its input
// until it finds one that is sorted.
func bogoSort(ctx context.Context, arr arrObj) {
	isSorted := func() bool {
		for i := 0; i < arr.len()-1; i++ {
			if arr.get(ctx, i) > arr.get(ctx, i+1) {
				return false
			}
		}
		return true
	}

	for !isSorted() {
		// Fisher-Yates shuffle
		for i := arr.len() - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			arr.swap(ctx, i, j)
		}
	}
}

func timsort(ctx context.Context, arr arrObj) {
	const RUN = 32
	n := arr.len()

	for i := 0; i < n; i += RUN {
		end := i + RUN - 1
		if end >= n {
			end = n - 1
		}
		insertionSortRange(ctx, arr, i, end)
	}

	for size := RUN; size < n; size = 2 * size {
		for left := 0; left < n; left += 2 * size {
			mid := left + size - 1
			right := left + 2*size - 1
			if right >= n {
				right = n - 1
			}

			if mid < right {
				merge(ctx, arr, left, mid, right)
			}
		}
	}
}

func insertionSortRange(ctx context.Context, arr arrObj, left, right int) {
	for i := left + 1; i <= right; i++ {
		key := arr.get(ctx, i)
		j := i - 1
		for j >= left && arr.get(ctx, j) > key {
			arr.set(ctx, j+1, arr.get(ctx, j))
			j--
		}
		arr.set(ctx, j+1, key)
	}
}

func bitonicSort(ctx context.Context, arr arrObj) {
	n := arr.len()
	// NOTE: This algorithm requires the input size to be a power of 2.
	// We'll proceed anyway, but it may not sort correctly for other sizes
	// without padding, which is complex with the given interface.
	bitonicSortRecurse(ctx, arr, 0, n, true)
}

func bitonicSortRecurse(ctx context.Context, arr arrObj, low, count int, dir bool) {
	if count > 1 {
		k := count / 2
		bitonicSortRecurse(ctx, arr, low, k, true)
		bitonicSortRecurse(ctx, arr, low+k, k, false)
		bitonicMerge(ctx, arr, low, count, dir)
	}
}

func bitonicMerge(ctx context.Context, arr arrObj, low, count int, dir bool) {
	if count > 1 {
		k := count / 2
		for i := low; i < low+k; i++ {
			compareAndSwap(ctx, arr, i, i+k, dir)
		}
		bitonicMerge(ctx, arr, low, k, dir)
		bitonicMerge(ctx, arr, low+k, k, dir)
	}
}

func compareAndSwap(ctx context.Context, arr arrObj, i, j int, dir bool) {
	isGreater := arr.get(ctx, i) > arr.get(ctx, j)
	if dir == isGreater {
		arr.swap(ctx, i, j)
	}
}
