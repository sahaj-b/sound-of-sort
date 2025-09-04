package main

import "time"

type IntArr []int

type arrObj interface {
	get(ind int) int
	set(ind, val int)
	swap(i, j int)
	len() int
}

const (
	delay = 5 * time.Millisecond
)

func (arr *IntArr) len() int {
	return len(*arr)
}

func (arr *IntArr) get(ind int) int {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[ind] = green
	render(arrGraph(*arr, colors))
	return (*arr)[ind]
}

func (arr *IntArr) set(ind, val int) {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[ind] = red
	(*arr)[ind] = val
	render(arrGraph(*arr, colors))
}

func (arr *IntArr) swap(i, j int) {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[i] = red
	colors[j] = red
	(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
	render(arrGraph(*arr, colors))
}

func newArrObj(arr []int) arrObj {
	a := IntArr(arr)
	return &a
}
