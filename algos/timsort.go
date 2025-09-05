package algos

import "context"

func timsort(ctx context.Context, arr ArrObj) {
	const RUN = 32
	n := arr.Len()

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

func insertionSortRange(ctx context.Context, arr ArrObj, left, right int) {
	for i := left + 1; i <= right; i++ {
		key := arr.Get(ctx, i)
		j := i - 1
		for j >= left && arr.Get(ctx, j) > key {
			arr.Set(ctx, j+1, arr.Get(ctx, j))
			j--
		}
		arr.Set(ctx, j+1, key)
	}
}
