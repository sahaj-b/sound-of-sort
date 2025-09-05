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
	delay    = 0 * time.Millisecond
	readClr  = red
	writeClr = green
)

func (arr *IntArr) len() int {
	return len(*arr)
}

func (arr *IntArr) get(ind int) int {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[ind] = readClr
	render(arrGraph(*arr, colors))
	playBeepArr((*arr)[ind])
	return (*arr)[ind]
}

func (arr *IntArr) set(ind, val int) {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[ind] = writeClr
	(*arr)[ind] = val
	render(arrGraph(*arr, colors))
	playBeepArr((*arr)[ind])
}

func (arr *IntArr) swap(i, j int) {
	time.Sleep(delay)
	colors := make([]string, len(*arr))
	colors[i] = writeClr
	colors[j] = writeClr
	(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
	render(arrGraph(*arr, colors))
	playBeepArr((*arr)[i])
	playBeepArr((*arr)[j])
}

func newArrObj(arr []int) arrObj {
	a := IntArr(arr)
	return &a
}
