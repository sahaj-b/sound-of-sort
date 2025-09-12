package algos

import "context"

func tournamentSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	if n <= 1 {
		return
	}

	for i := 0; i < n-1; i++ {
		winner := findWinner(ctx, arr, i, n-1)
		if winner != i {
			arr.Swap(ctx, i, winner)
		}
	}
}

func findWinner(ctx context.Context, arr ArrObj, start, end int) int {
	if start == end {
		return start
	}

	if start == end-1 {
		if arr.Get(ctx, start) < arr.Get(ctx, end) {
			return start
		}
		return end
	}

	mid := (start + end) / 2
	leftWinner := findWinner(ctx, arr, start, mid)
	rightWinner := findWinner(ctx, arr, mid+1, end)

	if arr.Get(ctx, leftWinner) < arr.Get(ctx, rightWinner) {
		return leftWinner
	}
	return rightWinner
}
