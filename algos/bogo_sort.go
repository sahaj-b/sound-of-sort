package algos

import (
	"context"
	"math/rand"
)

func bogoSort(ctx context.Context, arr ArrObj) {
	isSorted := func() bool {
		for i := 0; i < arr.Len()-1; i++ {
			if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
				return false
			}
		}
		return true
	}

	for !isSorted() {
		for i := arr.Len() - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			arr.Swap(ctx, i, j)
		}
	}
}
