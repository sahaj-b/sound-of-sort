package main

type IntArr []int

type arrObj interface {
	get(ind int) int
	set(ind, val int)
	swap(i, j int)
}

func (arr *IntArr) get(ind int) int {
	colors := make([]string, len(*arr))
	colors[ind] = green
	printGraph(arrGraph(*arr, colors))
	return (*arr)[ind]
}

func (arr *IntArr) set(ind, val int) {
	colors := make([]string, len(*arr))
	colors[ind] = red
	(*arr)[ind] = val
	printGraph(arrGraph(*arr, colors))
}

func (arr *IntArr) swap(i, j int) {
	colors := make([]string, len(*arr))
	colors[i] = red
	colors[j] = red
	(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
	printGraph(arrGraph(*arr, colors))
}
