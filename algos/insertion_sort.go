package algos

import "context"

func insertionSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	for i := 1; i < n; i++ {
		key := arr.Get(ctx, i)
		j := i - 1
		for j >= 0 && arr.Get(ctx, j) > key {
			arr.Set(ctx, j+1, arr.Get(ctx, j))
			j--
		}
		arr.Set(ctx, j+1, key)
	}
}
