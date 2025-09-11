package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sahaj-b/sound-of-sort/algos"
)

type visArr struct {
	mu     sync.RWMutex
	arr    []int
	colors []string
	delay  *atomic.Int64
	volume *atomic.Uint64
}

const (
	readClr  = red
	writeClr = green
)

func newVisualizer(arr []int, delay *atomic.Int64, volume *atomic.Uint64) algos.ArrObj {
	return &visArr{
		arr:    arr,
		colors: make([]string, len(arr)),
		delay:  delay,
		volume: volume,
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
	v.colors[ind] = readClr
	playBeepArr(val, v.volume)
	return val
}

func (v *visArr) Set(ctx context.Context, ind, val int) {
	v.Check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[ind] = val
	v.colors[ind] = writeClr

	playBeepArr(v.arr[ind], v.volume)
}

func (v *visArr) Swap(ctx context.Context, i, j int) {
	v.Check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[i], v.arr[j] = v.arr[j], v.arr[i]
	v.colors[i] = writeClr
	v.colors[j] = writeClr

	playBeepArr(v.arr[i], v.volume)
	playBeepArr(v.arr[j], v.volume)
}
