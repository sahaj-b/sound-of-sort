package main

import (
	"context"
	"sync/atomic"
	"time"
)

type arrObj interface {
	get(ctx context.Context, ind int) int
	set(ctx context.Context, ind, val int)
	swap(ctx context.Context, i, j int)
	len() int
	check(ctx context.Context)
}

type visArr struct {
	arr       []int
	stateChan chan<- VisState // rendering happens when this channel is sent to(state updates)
	sortName  string
}

const (
	readClr  = red
	writeClr = green
)

var delay atomic.Int64

func newVisualizer(arr []int, stateChan chan<- VisState, sortName string) arrObj {
	return &visArr{
		arr:       arr,
		stateChan: stateChan,
		sortName:  sortName,
	}
}

func (v *visArr) check(ctx context.Context) {
	select {
	case <-ctx.Done():
		panic(ctx.Err())
	default:
	}
	time.Sleep(time.Duration(delay.Load()))
}

func (v *visArr) len() int {
	return len(v.arr)
}

func (v *visArr) get(ctx context.Context, ind int) int {
	v.check(ctx)
	colors := make([]string, len(v.arr))
	colors[ind] = readClr

	v.stateChan <- VisState{
		Arr:      v.arr,
		Colors:   colors,
		SortName: v.sortName,
	}

	playBeepArr(v.arr[ind])
	return v.arr[ind]
}

func (v *visArr) set(ctx context.Context, ind, val int) {
	v.check(ctx)
	colors := make([]string, len(v.arr))
	colors[ind] = writeClr
	v.arr[ind] = val

	v.stateChan <- VisState{
		Arr:      v.arr,
		Colors:   colors,
		SortName: v.sortName,
	}

	playBeepArr(v.arr[ind])
}

func (v *visArr) swap(ctx context.Context, i, j int) {
	v.check(ctx)
	colors := make([]string, len(v.arr))
	colors[i] = writeClr
	colors[j] = writeClr
	v.arr[i], v.arr[j] = v.arr[j], v.arr[i]

	v.stateChan <- VisState{
		Arr:      v.arr,
		Colors:   colors,
		SortName: v.sortName,
	}

	playBeepArr(v.arr[i])
	playBeepArr(v.arr[j])
}
