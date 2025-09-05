package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type VisState struct {
	Arr      []int
	Colors   []string
	SortName string
}

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	stateChan chan VisState
	inputChan chan string

	currentSortIndex atomic.Int32
	currentSize      atomic.Int32
	shuffleRequested atomic.Bool
	delay            atomic.Int64
	volume           atomic.Uint64

	fps         int
	originalArr []int
}

func NewApp() *App {
	if parseArgs() {
		os.Exit(0)
	}

	app := &App{}

	app.ctx, app.cancel = context.WithCancel(context.Background())
	app.stateChan = make(chan VisState, 1)
	app.inputChan = make(chan string, 1)
	app.fps = *fpsFlag

	app.delay.Store(int64(time.Duration(*initialDelay) * time.Millisecond))
	app.volume.Store(math.Float64bits(*initialVolume))
	for i, s := range sorts {
		if s.arg == *initialSort {
			app.currentSortIndex.Store(int32(i))
			break
		}
	}
	app.currentSize.Store(int32(*initialSize))

	app.originalArr = getSequenceArr(1, *initialSize)
	setArrBounds(1, *initialSize)
	shuffleArr(app.originalArr)

	return app
}

func (app *App) Run() {
	restoreUI := initUI()
	defer restoreUI()
	initAudio()

	go inputReader(app.inputChan)
	go app.renderLoop()

	previousSize := app.currentSize.Load()

	for {
		select {
		case <-app.ctx.Done():
			return
		default:
		}

		currentSizeVal := app.currentSize.Load()
		if currentSizeVal != previousSize || app.shuffleRequested.Load() {
			app.originalArr = getSequenceArr(1, int(currentSizeVal))
			setArrBounds(1, int(currentSizeVal))
			shuffleArr(app.originalArr)
			previousSize = currentSizeVal
			app.shuffleRequested.Store(false)
		}

		if !app.runSortCycle() {
			return
		}
	}
}

func (app *App) renderLoop() {
	for state := range app.stateChan {
		render(
			arrGraph(state.Arr, state.Colors),
			state.SortName,
			time.Duration(app.delay.Load()),
			math.Float64frombits(app.volume.Load()),
			int(app.currentSize.Load()),
		)
	}
}

func (app *App) runSortCycle() bool {
	fmt.Print(clear)

	arrToSort := make([]int, len(app.originalArr))
	copy(arrToSort, app.originalArr)

	currentIndex := app.currentSortIndex.Load()
	currentSortName := sorts[currentIndex].name

	arr := newVisualizer(arrToSort, &app.delay, &app.volume)
	sortCtx, sortCancel := context.WithCancel(app.ctx)
	defer sortCancel()

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
		ticker := time.NewTicker(time.Second / time.Duration(app.fps))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentArr, currentColors := arr.getState()
				app.stateChan <- VisState{
					Arr:      currentArr,
					Colors:   currentColors,
					SortName: currentSortName,
				}
				arr.clearColors()
			case <-sortCtx.Done():
				return
			}
		}
	}()

inputLoop:
	for {
		select {
		case input := <-app.inputChan:
			if handleInput(input, &app.currentSortIndex, &app.currentSize, &app.shuffleRequested, &app.delay, &app.volume) {
				app.cancel()
				return false
			}

			if currentIndex != app.currentSortIndex.Load() ||
				app.currentSize.Load() != int32(len(arrToSort)) ||
				app.shuffleRequested.Load() {
				sortCancel()
				break inputLoop
			}
		case <-sortCtx.Done():
			break inputLoop
		case <-app.ctx.Done():
			return false
		}
	}

	wg.Wait()

	return true
}

func main() {
	app := NewApp()
	app.Run()
}
