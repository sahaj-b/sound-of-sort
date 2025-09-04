package main

import "math/rand"

func main() {
	defer initUI()()
	rawArr := getRandArr(1, 30, 10)
	arr := newArrObj(rawArr)
	colors := make([]string, arr.len())
	printGraph(arrGraph(rawArr, colors))
	bubbleSort(arr)
}

func getRandArr(min, max, length int) []int {
	if min > max || length < 1 {
		panic("bruh")
	}
	arr := make([]int, length)
	for i := range arr {
		arr[i] = min + rand.Intn(max-min+1)
	}
	return arr
}
