package algos

import "context"

func strandSort(ctx context.Context, arr ArrObj) {
	n := arr.Len()
	if n <= 1 {
		return
	}

	output := make([]int, 0, n)
	input := make([]int, n)
	for i := 0; i < n; i++ {
		input[i] = arr.Get(ctx, i)
	}

	for len(input) > 0 {
		strand := []int{input[0]}
		remaining := make([]int, 0, len(input)-1)

		for i := 1; i < len(input); i++ {
			if input[i] >= strand[len(strand)-1] {
				strand = append(strand, input[i])
			} else {
				remaining = append(remaining, input[i])
			}
		}

		output = strandMerge(output, strand)
		input = remaining

		for i := 0; i < len(output); i++ {
			arr.Set(ctx, i, output[i])
		}
	}
}

func strandMerge(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	for i < len(a) {
		result = append(result, a[i])
		i++
	}

	for j < len(b) {
		result = append(result, b[j])
		j++
	}

	return result
}
