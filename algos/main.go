package algos

import "context"

type ArrObj interface {
	Get(ctx context.Context, ind int) int
	Set(ctx context.Context, ind, val int)
	Swap(ctx context.Context, i, j int)
	Len() int
	GetState() ([]int, []string)
	ClearColors()
	Check(ctx context.Context)
}

type sortFunc func(ctx context.Context, arr ArrObj)

var Sorts = []struct {
	Name string
	Arg  string
	Fun  sortFunc
}{
	{"Quick Sort", "quick", quickSort},
	{"Bubble Sort", "bubble", bubbleSort},
	{"Selection Sort", "selection", selectionSort},
	{"Insertion Sort", "insertion", insertionSort},
	{"Merge Sort", "merge", mergeSort},
	{"Heap Sort", "heap", heapSort},
	{"Shell Sort", "shell", shellSort},
	{"Cocktail Shaker Sort", "cocktail", cocktailShakerSort},
	{"Gnome Sort", "gnome", gnomeSort},
	{"Pancake Sort", "pancake", pancakeSort},
	{"Radix Sort (LSD)", "radix", radixSortLSD},
	{"Timsort", "timsort", timsort},
	{"Bitonic Sort", "bitonic", bitonicSort},
	{"Bogo Sort", "bogo", bogoSort},
}
