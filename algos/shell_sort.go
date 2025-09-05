package algos

import "context"

func shellSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	for gap := n / 2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			temp := arr.Get(ctx, i)
			j := i
			for ; j >= gap && arr.Get(ctx, j-gap) > temp; j -= gap {
				arr.Set(ctx, j, arr.Get(ctx, j-gap))
			}
			arr.Set(ctx, j, temp)
		}
	}
}
