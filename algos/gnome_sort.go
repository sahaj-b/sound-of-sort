package algos

import "context"

func gnomeSort(ctx context.Context, arr ArrObj) {
	index := 0
	n := arr.Len()
	for index < n {
		if index == 0 {
			index++
		}
		if arr.Get(ctx, index) >= arr.Get(ctx, index-1) {
			index++
		} else {
			arr.Swap(ctx, index, index-1)
			index--
		}
	}
}
