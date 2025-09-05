package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type VisState struct {
	Arr      []int
	Colors   []string
	SortName string
}

func main() {
	restoreUI := initUI()
	defer restoreUI()
	initAudio()

	var currentSortIndex atomic.Int32
	var currentSize atomic.Int32
	var shuffleRequested atomic.Bool
	delay.Store(int64(2 * time.Millisecond))
	volume.Store(math.Float64bits(0.1))

	inputChan := make(chan string, 1)
	go inputReader(inputChan)

	const arrSize = 150
	currentSize.Store(int32(arrSize))
	originalArr := getSequenceArr(1, arrSize)
	setArrBounds(1, arrSize)
	shuffleArr(originalArr)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	stateChan := make(chan VisState, 1)

	go func() {
		// just render loop
		for state := range stateChan {
			render(
				arrGraph(state.Arr, state.Colors),
				state.SortName,
				time.Duration(delay.Load()),
				math.Float64frombits(volume.Load()),
				int(currentSize.Load()),
			)
		}
	}()

	previousSize := currentSize.Load()

	for {
		select {
		case <-appCtx.Done():
			return
		default:
		}

		currentSizeVal := currentSize.Load()
		if currentSizeVal != previousSize || shuffleRequested.Load() {
			originalArr = getSequenceArr(1, int(currentSizeVal))
			setArrBounds(1, int(currentSizeVal))
			shuffleArr(originalArr)
			previousSize = currentSizeVal
			shuffleRequested.Store(false)
		}

		fmt.Print(clear)

		arrToSort := make([]int, len(originalArr))
		copy(arrToSort, originalArr)

		currentIndex := currentSortIndex.Load()
		currentSortName := sorts[currentIndex].name

		arr := newVisualizer(arrToSort, stateChan, currentSortName)

		sortCtx, sortCancel := context.WithCancel(appCtx)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			// sort goroutine
			defer wg.Done()
			defer func() {
				// recover from context.Canceled panic
				if r := recover(); r != nil && r != context.Canceled {
					panic(r)
				}
			}()

			// initial state render
			stateChan <- VisState{Arr: arrToSort, Colors: make([]string, len(arrToSort)), SortName: currentSortName}
			sorts[currentIndex].fun(sortCtx, arr)

			// final state render (all green bby)
			finalColors := make([]string, len(arrToSort))
			for i := range finalColors {
				finalColors[i] = green
			}
			stateChan <- VisState{Arr: arrToSort, Colors: finalColors, SortName: currentSortName}
		}()

	inputLoop:
		for {
			select {
			case input := <-inputChan:
				if handleInput(input, &currentSortIndex, &currentSize, &shuffleRequested) {
					appCancel()
					sortCancel()
					break inputLoop
				}

				if currentIndex != currentSortIndex.Load() {
					sortCancel()
					break inputLoop
				}

				if currentSizeVal != currentSize.Load() {
					sortCancel()
					break inputLoop
				}

				if shuffleRequested.Load() {
					sortCancel()
					break inputLoop
				}

			case <-sortCtx.Done():
				break inputLoop

			case <-appCtx.Done():
				break inputLoop
			}
		}

		wg.Wait()
	}
}
