package algos

import "context"

func bubbleSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr.Get(ctx, j) > arr.Get(ctx, j+1) {
				arr.Swap(ctx, j, j+1)
			}
		}
	}
}
