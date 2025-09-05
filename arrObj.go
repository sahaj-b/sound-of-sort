package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type arrObj interface {
	get(ctx context.Context, ind int) int
	set(ctx context.Context, ind, val int)
	swap(ctx context.Context, i, j int)
	len() int
	getState() ([]int, []string)
	clearColors()
	check(ctx context.Context)
}

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

func newVisualizer(arr []int, delay *atomic.Int64, volume *atomic.Uint64) arrObj {
	return &visArr{
		arr:    arr,
		colors: make([]string, len(arr)),
		delay:  delay,
		volume: volume,
	}
}

func (v *visArr) check(ctx context.Context) {
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

func (v *visArr) len() int {
	return len(v.arr)
}

func (v *visArr) clearColors() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.colors = make([]string, len(v.arr))
}

func (v *visArr) getState() ([]int, []string) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	arrCopy := make([]int, len(v.arr))
	copy(arrCopy, v.arr)
	colorsCopy := make([]string, len(v.colors))
	copy(colorsCopy, v.colors)

	return arrCopy, colorsCopy
}

func (v *visArr) get(ctx context.Context, ind int) int {
	v.check(ctx)
	v.mu.RLock()
	defer v.mu.RUnlock()

	v.colors[ind] = readClr
	playBeepArr(v.arr[ind], v.volume)

	return v.arr[ind]
}

func (v *visArr) set(ctx context.Context, ind, val int) {
	v.check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[ind] = val
	v.colors[ind] = writeClr

	playBeepArr(v.arr[ind], v.volume)
}

func (v *visArr) swap(ctx context.Context, i, j int) {
	v.check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[i], v.arr[j] = v.arr[j], v.arr[i]
	v.colors[i] = writeClr
	v.colors[j] = writeClr

	playBeepArr(v.arr[i], v.volume)
	playBeepArr(v.arr[j], v.volume)
}
