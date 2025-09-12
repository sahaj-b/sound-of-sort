package algos

import "context"

const introsortThreshold = 10
const maxDepth = 32

func introsort(ctx context.Context, arr ArrObj) {
	introsortRecurse(ctx, arr, 0, arr.Len()-1, maxDepth)
}

func introsortRecurse(ctx context.Context, arr ArrObj, low, high, depth int) {
	if high <= low {
		return
	}

	size := high - low + 1

	if size <= introsortThreshold {
		introsortInsertionSort(ctx, arr, low, high)
		return
	}

	if depth == 0 {
		introsortHeapSort(ctx, arr, low, high)
		return
	}

	pivot := introsortPartition(ctx, arr, low, high)
	introsortRecurse(ctx, arr, low, pivot-1, depth-1)
	introsortRecurse(ctx, arr, pivot+1, high, depth-1)
}

func introsortInsertionSort(ctx context.Context, arr ArrObj, low, high int) {
	for i := low + 1; i <= high; i++ {
		key := arr.Get(ctx, i)
		j := i - 1
		for j >= low && arr.Get(ctx, j) > key {
			arr.Set(ctx, j+1, arr.Get(ctx, j))
			j--
		}
		arr.Set(ctx, j+1, key)
	}
}

func introsortPartition(ctx context.Context, arr ArrObj, low, high int) int {
	medianOfThree(ctx, arr, low, high)
	pivot := arr.Get(ctx, high)
	i := low - 1

	for j := low; j < high; j++ {
		if arr.Get(ctx, j) <= pivot {
			i++
			arr.Swap(ctx, i, j)
		}
	}
	arr.Swap(ctx, i+1, high)
	return i + 1
}

func medianOfThree(ctx context.Context, arr ArrObj, low, high int) {
	mid := (low + high) / 2

	if arr.Get(ctx, mid) < arr.Get(ctx, low) {
		arr.Swap(ctx, low, mid)
	}
	if arr.Get(ctx, high) < arr.Get(ctx, low) {
		arr.Swap(ctx, low, high)
	}
	if arr.Get(ctx, high) < arr.Get(ctx, mid) {
		arr.Swap(ctx, mid, high)
	}
}

func introsortHeapSort(ctx context.Context, arr ArrObj, low, high int) {
	n := high - low + 1

	for i := n/2 - 1; i >= 0; i-- {
		introsortHeapify(ctx, arr, low, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr.Swap(ctx, low, low+i)
		introsortHeapify(ctx, arr, low, i, 0)
	}
}

func introsortHeapify(ctx context.Context, arr ArrObj, offset, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr.Get(ctx, offset+left) > arr.Get(ctx, offset+largest) {
		largest = left
	}

	if right < n && arr.Get(ctx, offset+right) > arr.Get(ctx, offset+largest) {
		largest = right
	}

	if largest != i {
		arr.Swap(ctx, offset+i, offset+largest)
		introsortHeapify(ctx, arr, offset, n, largest)
	}
}
