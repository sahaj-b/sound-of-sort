package main

import "fmt"

func main() {
	defer initUI()()
	initAudio()
	n := 50
	rawArr := getSequenceArr(1, n)
	setArrBounds(1, n)
	shuffleArr(rawArr)
	arr := newArrObj(rawArr)
	colors := make([]string, arr.len())
	render(arrGraph(rawArr, colors))

	done := make(chan struct{})

	// go selectionSort(arr, done)
	// go bubbleSort(arr, done)
	// go insertionSort(arr, done)
	go quickSort(arr, done)

	inputs := make(chan string)
	go inputReader(inputs)
	exit := false

	for !exit {
		select {
		case <-done:
			fmt.Println("Sorting complete!\r")
			return
		case input, ok := <-inputs:
			if !ok {
				exit = true
				break
			}
			exit = handleInput(input)
		}
	}
	fmt.Println("Exiting...\r")
}
