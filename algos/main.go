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
	{"Bubble Sort", "bubble", bubbleSort},
	{"Selection Sort", "selection", selectionSort},
	{"Insertion Sort", "insertion", insertionSort},
	{"Gnome Sort", "gnome", gnomeSort},
	{"Cocktail Shaker Sort", "cocktail", cocktailShakerSort},
	{"Pancake Sort", "pancake", pancakeSort},
	{"Shell Sort", "shell", shellSort},
	{"Merge Sort", "merge", mergeSort},
	{"Quick Sort", "quick", quickSort},
	{"Heap Sort", "heap", heapSort},
	{"Bitonic Sort", "bitonic", bitonicSort},
	{"Timsort", "timsort", timsort},
	{"Radix Sort (LSD)", "radix", radixSortLSD},
	{"Tournament Sort", "tournament", tournamentSort},
	{"Introsort", "introsort", introsort},
	{"Odd-Even Sort", "oddeven", oddEvenSort},
	{"Cycle Sort", "cycle", cycleSort},
	{"Strand Sort", "strand", strandSort},
	{"Bogo Sort", "bogo", bogoSort},
}
