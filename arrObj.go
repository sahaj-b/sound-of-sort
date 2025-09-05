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
	check(ctx context.Context)
}

type visArr struct {
	mu     sync.RWMutex
	arr    []int
	colors []string
}

const (
	readClr  = red
	writeClr = green
)

var delay atomic.Int64

func newVisualizer(arr []int) arrObj {
	return &visArr{
		arr:    arr,
		colors: make([]string, len(arr)),
	}
}

func (v *visArr) check(ctx context.Context) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
	}
	d := delay.Load()
	if d > 0 {
		time.Sleep(time.Duration(d))
	}
}

func (v *visArr) len() int {
	return len(v.arr)
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

	v.colors = make([]string, len(v.arr))
	v.colors[ind] = readClr
	playBeepArr(v.arr[ind])

	return v.arr[ind]
}

func (v *visArr) set(ctx context.Context, ind, val int) {
	v.check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[ind] = val
	v.colors = make([]string, len(v.arr))
	v.colors[ind] = writeClr

	playBeepArr(v.arr[ind])
}

func (v *visArr) swap(ctx context.Context, i, j int) {
	v.check(ctx)
	v.mu.Lock()
	defer v.mu.Unlock()

	v.arr[i], v.arr[j] = v.arr[j], v.arr[i]
	v.colors = make([]string, len(v.arr))
	v.colors[i] = writeClr
	v.colors[j] = writeClr

	playBeepArr(v.arr[i])
	playBeepArr(v.arr[j])
}
