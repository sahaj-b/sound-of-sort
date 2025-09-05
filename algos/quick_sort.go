package algos

import "context"

func quickSort(ctx context.Context, arr ArrObj) {
	quickSortRecurse(ctx, arr, 0, arr.Len()-1)
}

func quickSortRecurse(ctx context.Context, arr ArrObj, low, high int) {
	if low < high {
		pi := partition(ctx, arr, low, high)
		quickSortRecurse(ctx, arr, low, pi-1)
		quickSortRecurse(ctx, arr, pi+1, high)
	}
}

func partition(ctx context.Context, arr ArrObj, low, high int) int {
	pivot := arr.Get(ctx, high)
	i := low - 1
	for j := low; j < high; j++ {
		if arr.Get(ctx, j) < pivot {
			i++
			arr.Swap(ctx, i, j)
		}
	}
	arr.Swap(ctx, i+1, high)
	return i + 1
}
