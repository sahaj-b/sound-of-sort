package algos

import "context"

func radixSortLSD(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	if n == 0 {
		return
	}

	maxVal := arr.Get(ctx, 0)
	for i := 1; i < n; i++ {
		if arr.Get(ctx, i) > maxVal {
			maxVal = arr.Get(ctx, i)
		}
	}

	for exp := 1; maxVal/exp > 0; exp *= 10 {
		buckets := make([][]int, 10)

		for i := range n {
			val := arr.Get(ctx, i)
			bucketIndex := (val / exp) % 10
			buckets[bucketIndex] = append(buckets[bucketIndex], val)
		}

		i := 0
		for _, bucket := range buckets {
			for _, val := range bucket {
				arr.Set(ctx, i, val)
				i++
			}
		}
	}
}
