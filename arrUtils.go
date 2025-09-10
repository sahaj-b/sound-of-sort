package main

import (
	"math/rand"
)

func getSequenceArr(start, length int) []int {
	if length < 1 {
		panic("bruh")
	}
	arr := make([]int, length)
	for i := range arr {
		arr[i] = start + i
	}
	return arr
}

func shuffleArr(arr []int) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}
