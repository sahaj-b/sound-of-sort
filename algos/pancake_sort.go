package algos

import "context"

func pancakeSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	for currSize := n; currSize > 1; currSize-- {
		mi := 0
		for i := 0; i < currSize; i++ {
			if arr.Get(ctx, i) > arr.Get(ctx, mi) {
				mi = i
			}
		}

		if mi != currSize-1 {
			pancakeFlip(ctx, arr, mi+1)
			pancakeFlip(ctx, arr, currSize)
		}
	}
}

func pancakeFlip(ctx context.Context, arr ArrObj, k int) {
	start := 0
	for start < k-1 {
		arr.Swap(ctx, start, k-1)
		start++
		k--
	}
}
