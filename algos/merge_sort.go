package algos

import "context"

func mergeSort(ctx context.Context, arr ArrObj) {
	mergeSortRecurse(ctx, arr, 0, arr.Len()-1)
}

func mergeSortRecurse(ctx context.Context, arr ArrObj, left, right int) {
	if left < right {
		mid := (left + right) / 2
		mergeSortRecurse(ctx, arr, left, mid)
		mergeSortRecurse(ctx, arr, mid+1, right)
		merge(ctx, arr, left, mid, right)
	}
}

func merge(ctx context.Context, arr ArrObj, left, mid, right int) {
	n1 := mid - left + 1
	n2 := right - mid

	leftArr := make([]int, n1)
	rightArr := make([]int, n2)

	for i := range n1 {
		leftArr[i] = arr.Get(ctx, left+i)
	}
	for j := range n2 {
		rightArr[j] = arr.Get(ctx, mid+1+j)
	}

	i, j, k := 0, 0, left
	for i < n1 && j < n2 {
		if leftArr[i] <= rightArr[j] {
			arr.Set(ctx, k, leftArr[i])
			i++
		} else {
			arr.Set(ctx, k, rightArr[j])
			j++
		}
		k++
	}

	for i < n1 {
		arr.Set(ctx, k, leftArr[i])
		i++
		k++
	}

	for j < n2 {
		arr.Set(ctx, k, rightArr[j])
		j++
		k++
	}
}
