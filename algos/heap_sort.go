package algos

import "context"

func heapSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()

	for i := n/2 - 1; i >= 0; i-- {
		heapify(ctx, arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr.Swap(ctx, 0, i)
		heapify(ctx, arr, i, 0)
	}
}

func heapify(ctx context.Context, arr ArrObj, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr.Get(ctx, left) > arr.Get(ctx, largest) {
		largest = left
	}

	if right < n && arr.Get(ctx, right) > arr.Get(ctx, largest) {
		largest = right
	}

	if largest != i {
		arr.Swap(ctx, i, largest)
		heapify(ctx, arr, n, largest)
	}
}
