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
	if parseArgs() {
		return
	}

	restoreUI := initUI()
	defer restoreUI()
	initAudio()

	var currentSortIndex atomic.Int32
	var currentSize atomic.Int32
	var shuffleRequested atomic.Bool

	delay.Store(int64(time.Duration(*initialDelay) * time.Millisecond))
	volume.Store(math.Float64bits(*initialVolume))
	fps := *fpsFlag
	for i, s := range sorts {
		if s.arg == *initialSort {
			currentSortIndex.Store(int32(i))
			break
		}
	}

	inputChan := make(chan string, 1)
	go inputReader(inputChan)

	arrSize := *initialSize
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

		arr := newVisualizer(arrToSort)

		sortCtx, sortCancel := context.WithCancel(appCtx)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil && r != context.Canceled {
					panic(r)
				}
			}()
			sorts[currentIndex].fun(sortCtx, arr)
		}()

		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Second / time.Duration(fps))
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					currentArr, currentColors := arr.getState()
					stateChan <- VisState{
						Arr:      currentArr,
						Colors:   currentColors,
						SortName: currentSortName,
					}
				case <-sortCtx.Done():
					return
				}
			}
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

				if currentIndex != currentSortIndex.Load() ||
					currentSizeVal != currentSize.Load() ||
					shuffleRequested.Load() {
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

		finalArr, _ := arr.getState()
		finalColors := make([]string, len(finalArr))
		for i := range finalColors {
			finalColors[i] = green
		}
		stateChan <- VisState{Arr: finalArr, Colors: finalColors, SortName: currentSortName}
	}
}
