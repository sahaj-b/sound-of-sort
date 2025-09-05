package algos

import "context"

func bitonicSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	bitonicSortRecurse(ctx, arr, 0, n, true)
}

func bitonicSortRecurse(ctx context.Context, arr ArrObj, low, count int, dir bool) {
	if count > 1 {
		k := count / 2
		bitonicSortRecurse(ctx, arr, low, k, true)
		bitonicSortRecurse(ctx, arr, low+k, k, false)
		bitonicMerge(ctx, arr, low, count, dir)
	}
}

func bitonicMerge(ctx context.Context, arr ArrObj, low, count int, dir bool) {
	if count > 1 {
		k := count / 2
		for i := low; i < low+k; i++ {
			compareAndSwap(ctx, arr, i, i+k, dir)
		}
		bitonicMerge(ctx, arr, low, k, dir)
		bitonicMerge(ctx, arr, low+k, k, dir)
	}
}

func compareAndSwap(ctx context.Context, arr ArrObj, i, j int, dir bool) {
	isGreater := arr.Get(ctx, i) > arr.Get(ctx, j)
	if dir == isGreater {
		arr.Swap(ctx, i, j)
	}
}
