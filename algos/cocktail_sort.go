package algos

import "context"

func cocktailShakerSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	swapped := true
	start := 0
	end := n - 1

	for swapped {
		swapped = false
		for i := start; i < end; i++ {
			if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
				arr.Swap(ctx, i, i+1)
				swapped = true
			}
		}
		if !swapped {
			break
		}
		swapped = false
		end--
		for i := end - 1; i >= start; i-- {
			if arr.Get(ctx, i) > arr.Get(ctx, i+1) {
				arr.Swap(ctx, i, i+1)
				swapped = true
			}
		}
		start++
	}
}
