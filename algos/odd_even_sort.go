package algos

import "context"

func oddEvenSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	sorted := false

	for !sorted {
		sorted = true

		for i := 1; i < n-1; i += 2 {
			if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
				arr.Swap(ctx, i, i+1)
				sorted = false
			}
		}

		for i := 0; i < n-1; i += 2 {
			if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
				arr.Swap(ctx, i, i+1)
				sorted = false
			}
		}
	}
}
