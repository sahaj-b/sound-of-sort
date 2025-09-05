package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
	red        = "\x1b[31m"
	green      = "\x1b[32m"
	reset      = "\x1b[0m"
	clear      = "\x1b[2J\x1b[H"
	moveToTop  = "\x1b[H"
)

var (
	termWidth  int
	termHeight int
)

func initUI() (restore func()) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Error setting terminal to raw mode:", err)
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
	}

	termWidth, termHeight = width, height
	fmt.Print(hideCursor)
	fmt.Print(clear)
	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print(showCursor)
	}
}

func inputReader(ch chan string) {
	defer close(ch)
	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			close(ch)
			return
		}
		ch <- string(buf[:n])
	}
}

func handleInput(input string, currentSortIndex *atomic.Int32) bool {
	switch input {
	case "q", "\x03": // Ctrl+C
		return true

	case "w": // increase delay
		delay.Add(1 * int64(time.Millisecond))

	case "s": // decrease delay
		// Use a CAS loop to prevent it from going below zero.
		for {
			currentDelay := delay.Load()

			newDelay := max(0, currentDelay-(2*int64(time.Millisecond)))

			// this returns false if the currentDelay got changed since we loaded it
			if delay.CompareAndSwap(currentDelay, newDelay) {
				break
			}
		}

	case "A", "\x1b[A": // Up arrow
		// CAS loop for float volume increase.
		for {
			oldBits := volume.Load()
			oldFloat := math.Float64frombits(oldBits)

			newFloat := min(oldFloat+0.005, 1.0)

			newBits := math.Float64bits(newFloat)
			if volume.CompareAndSwap(oldBits, newBits) {
				break
			}
		}

	case "B", "\x1b[B": // Down arrow
		// CAS loop for float volume decrease.
		for {
			oldBits := volume.Load()
			oldFloat := math.Float64frombits(oldBits)

			newFloat := max(0.0, oldFloat-0.01)

			newBits := math.Float64bits(newFloat)
			if volume.CompareAndSwap(oldBits, newBits) {
				break
			}
		}

	case "p", "\x1b[D": // Left arrow
		newIndex := currentSortIndex.Add(1)
		currentSortIndex.Store(newIndex % int32(len(sorts)))
	case "n", "\x1b[C": // Right arrow
		newIndex := currentSortIndex.Add(-1)
		if newIndex < 0 {
			newIndex = int32(len(sorts) - 1)
		}
		currentSortIndex.Store(newIndex)
	}
	return false
}

func arrGraph(arr []int, colors []string) []string {
	graphHeight := termHeight - 3
	if len(arr) != len(colors) {
		log.Fatal("arr and colors must have the same length")
	}
	graphChar := "█▊" // █ ▇ ▉ ▊ ▋ ▌
	if termWidth < 2*len(arr) {
		graphChar = "▊"
	}
	maxVal := 0
	minVal := 0
	for _, v := range arr {
		maxVal = max(maxVal, v)
		minVal = min(minVal, v)
	}
	output := make([]string, graphHeight)
	for i := range output {
		for j := range arr {
			val := graphHeight * (arr[j] - minVal) / (maxVal - minVal)
			if graphHeight-i <= val {
				output[i] += colors[j] + graphChar + reset
			} else {
				output[i] += strings.Repeat(" ", utf8.RuneCountInString(graphChar))
			}
		}
	}
	return output
}

func render(graph []string, sortName string, currentDelay time.Duration, currentVol float64) {
	fmt.Print(moveToTop)
	sortStr := "←/→: " + sortName
	volumeStr := "↑/↓ Volume: " + fmt.Sprintf("%.1f", currentVol*100)
	delayStr := "w/s Delay: " + currentDelay.String()
	quitStr := "q / Ctrl+C to quit"

	statusStr := sortStr + " | " + volumeStr + " | " + delayStr + " | " + quitStr
	paddingNeeded := max(0, (termWidth-len(statusStr))/2)
	statusPadding := strings.Repeat(" ", paddingNeeded)
	fmt.Println(statusPadding + statusStr + "\r")
	for _, line := range graph {
		fmt.Println(line + "\r")
	}
}
