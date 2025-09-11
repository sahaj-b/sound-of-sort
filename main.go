package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sahaj-b/sound-of-sort/algos"
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

	fps        int
	imgMode    bool
	horizontal bool
	intArr     []int
	imgArr     []string
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
	app.imgMode = *imgFlag
	app.horizontal = *horizontalFlag

	app.delay.Store(int64(*initialDelay * float64(time.Millisecond)))
	app.volume.Store(math.Float64bits(*initialVolume))
	for i, s := range algos.Sorts {
		if s.Arg == *initialSort {
			app.currentSortIndex.Store(int32(i))
			break
		}
	}

	size := *initialSize
	if app.imgMode {
		img, err := readImageFromStdin()
		makeRectangular(img)
		if err != nil {
			log.Fatal("Error reading image from stdin:", err)
		}
		if len(img) == 0 {
			log.Fatal("No image data provided in stdin")
		}
		if app.horizontal {
			// horizontal: rows are the elements
			size = len(img)
			app.imgArr = img
		} else {
			transposed, err := transpose(img)
			if err != nil {
				log.Fatal("Error processing image:", err)
			}
			size = len(transposed)
			app.imgArr = transposed
		}
	}
	app.currentSize.Store(int32(size))
	app.intArr = getSequenceArr(0, size)
	setArrBounds(0, size-1)
	shuffleArr(app.intArr)

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
			if app.imgMode {
				app.intArr = getSequenceArr(0, int(currentSizeVal))
				setArrBounds(0, int(currentSizeVal)-1)
			} else {
				app.intArr = getSequenceArr(1, int(currentSizeVal))
				setArrBounds(1, int(currentSizeVal))
			}
			shuffleArr(app.intArr)
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
		var graph []string
		defer func() {
			if r := recover(); r != nil {
				// just skip this frame; likely transient length mismatch
			}
		}()
		if app.imgMode {
			if app.horizontal {
				graph = imgGraphHorizontal(state.Arr, app.imgArr, state.Colors)
			} else {
				graph = imgGraph(state.Arr, app.imgArr, state.Colors)
			}
		} else {
			graph = arrGraph(state.Arr, state.Colors)
		}
		if len(graph) == 0 {
			continue
		}
		app.render(
			graph,
			state.SortName,
		)
	}
}

func (app *App) runSortCycle() bool {
	fmt.Print(clear)

	arrToSort := make([]int, len(app.intArr))
	copy(arrToSort, app.intArr)

	currentIndex := app.currentSortIndex.Load()
	currentSortName := algos.Sorts[currentIndex].Name

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
		algos.Sorts[currentIndex].Fun(sortCtx, arr)
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second / time.Duration(app.fps))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				currentArr, currentColors := arr.GetState()
				app.stateChan <- VisState{
					Arr:      currentArr,
					Colors:   currentColors,
					SortName: currentSortName,
				}
				arr.ClearColors()
			case <-sortCtx.Done():
				return
			}
		}
	}()

inputLoop:
	for {
		select {
		case input := <-app.inputChan:
			if app.handleInput(input) {
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
	log.SetFlags(0)
	app := NewApp()
	app.Run()
}
