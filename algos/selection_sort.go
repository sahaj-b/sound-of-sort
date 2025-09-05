package algos

import "context"

func selectionSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr.Get(ctx, j) < arr.Get(ctx, minIdx) {
				minIdx = j
			}
		}
		if minIdx != i {
			arr.Swap(ctx, i, minIdx)
		}
	}
}
