package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sahaj-b/sound-of-sort/algos"
)

type visArr struct {
	mu                  sync.RWMutex
	arr                 []int
	colors              []string
	delay               *atomic.Int64
	skipReadWriteColors bool
	app                 *App
}

const (
	readClr  = red
	writeClr = green
)

func newVisualizer(arr []int, delay *atomic.Int64, skipReadWriteColors bool, app *App) algos.ArrObj {
	return &visArr{
		arr:                 arr,
		colors:              make([]string, len(arr)),
		delay:               delay,
		skipReadWriteColors: skipReadWriteColors,
		app:                 app,
	}
}

func (v *visArr) Check(ctx context.Context) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
	}
	d := v.delay.Load()
	if d > 0 {
		time.Sleep(time.Duration(d))
	}
}

func (v *visArr) Len() int {
	return len(v.arr)
}

func (v *visArr) ClearColors() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.colors = make([]string, len(v.arr))
}

func (v *visArr) GetState() ([]int, []string) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	arrCopy := make([]int, len(v.arr))
	copy(arrCopy, v.arr)
	colorsCopy := make([]string, len(v.colors))
	copy(colorsCopy, v.colors)

	return arrCopy, colorsCopy
}

func (v *visArr) Get(ctx context.Context, ind int) int {
	v.Check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	val := v.arr[ind]
	if !v.skipReadWriteColors {
		v.colors[ind] = readClr
	}
	v.app.playBeepArr(val)
	return val
}

func (v *visArr) Set(ctx context.Context, ind, val int) {
	v.Check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[ind] = val
	if !v.skipReadWriteColors {
		v.colors[ind] = writeClr
	}

	v.app.playBeepArr(v.arr[ind])
}

func (v *visArr) Swap(ctx context.Context, i, j int) {
	v.Check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[i], v.arr[j] = v.arr[j], v.arr[i]
	if !v.skipReadWriteColors {
		v.colors[i] = writeClr
		v.colors[j] = writeClr
	}

	v.app.playBeepArr(v.arr[i])
	v.app.playBeepArr(v.arr[j])
}
