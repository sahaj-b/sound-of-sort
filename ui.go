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

	"github.com/sahaj-b/sound-of-sort/algos"
	"golang.org/x/term"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
	red        = "\x1b[31m"
	green      = "\x1b[32m"
	reset      = "\x1b[0m"
	clear      = "\x1b[2J\x1b[H"
	clearLine  = "\x1b[K"
	moveToTop  = "\x1b[H"
	bggray     = "\x1b[48;5;236m"
	cyan       = "\x1b[36m"
)

var (
	termWidth  int
	termHeight int
	graphChar  string
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

func handleInput(input string, currentSortIndex *atomic.Int32, currentSize *atomic.Int32, shuffleRequested *atomic.Bool, delay *atomic.Int64, volume *atomic.Uint64) bool {
	switch input {
	case "q", "\x03": // Ctrl+C
		return true

	case "w": // increase delay
		delay.Add(100 * int64(time.Microsecond))

	case "s": // decrease delay
		// use a CAS loop
		for {
			currentDelay := delay.Load()

			newDelay := max(50*int64(time.Microsecond), currentDelay-(100*int64(time.Microsecond)))

			// this returns false if the currentDelay got changed since we loaded it
			if delay.CompareAndSwap(currentDelay, newDelay) {
				break
			}
		}

	case "A", "\x1b[A": // Up arrow
		oldBits := volume.Load()
		oldFloat := math.Float64frombits(oldBits)

		newFloat := min(oldFloat+0.005, 1.0)

		newBits := math.Float64bits(newFloat)
		volume.Store(newBits)

	case "B", "\x1b[B": // Down arrow
		oldFloat := math.Float64frombits(volume.Load())
		newFloat := max(0.0, oldFloat-0.01)
		newBits := math.Float64bits(newFloat)
		volume.Store(newBits)

	case "p", "\x1b[D": // Left arrow
		newIndex := currentSortIndex.Add(1)
		currentSortIndex.Store(newIndex % int32(len(algos.Sorts)))
	case "n", "\x1b[C": // Right arrow
		newIndex := currentSortIndex.Add(-1)
		if newIndex < 0 {
			newIndex = int32(len(algos.Sorts) - 1)
		}
		currentSortIndex.Store(newIndex)

	case "a": // decrease array size
		newSize := max(0, currentSize.Add(-10))
		currentSize.Store(newSize)

	case "d": // increase array size
		newSize := currentSize.Add(10)
		currentSize.Store(newSize)

	case "r": // shuffle array
		shuffleRequested.Store(true)
	}
	return false
}

func arrGraph(arr []int, colors []string) []string {
	graphHeight := termHeight - 3
	if len(arr) != len(colors) {
		log.Fatal("arr and colors must have the same length")
	}
	graphChar = "█▊" // █ ▇ ▉ ▊ ▋ ▌
	if termWidth < 2*len(arr) {
		graphChar = "▊"
	}
	output := make([]string, graphHeight)
	denom := maxVal - minVal
	if denom == 0 {
		denom = 1
	}

	for i := range output {
		for j := range arr {
			val := graphHeight * (arr[j] - minVal) / denom
			if graphHeight-i <= val || i == graphHeight-1 {
				output[i] += colors[j] + graphChar + reset
			} else {
				output[i] += strings.Repeat(" ", utf8.RuneCountInString(graphChar))
			}
		}
	}
	return output
}

func render(graph []string, sortName string, currentDelay time.Duration, currentVol float64, currentSize int) {
	fmt.Print(moveToTop)
	sortStr := "←/→: " + sortName
	volumeStr := "↑/↓ Volume: " + fmt.Sprintf("%.1f", currentVol*100)
	delayStr := "w/s Delay: " + currentDelay.String()
	sizeStr := fmt.Sprintf("a/d Arr Size: %d", currentSize)
	quitStr := "q / Ctrl+C to quit"

	statusStr := bggray + cyan + " " + sortStr + " " + reset + " " + bggray + cyan + " " + volumeStr + " " + reset + " " + bggray + cyan + " " + delayStr + " " + reset + " " + bggray + cyan + " " + sizeStr + " " + reset + " " + bggray + red + " " + quitStr + " " + reset
	statusLen := len(volumeStr) + len(delayStr) + len(sizeStr) + len(sortStr) + len(quitStr) + 3*5
	statusPadding := max(0, (termWidth-statusLen)/2)
	fmt.Println(strings.Repeat(" ", statusPadding) + statusStr + "\r\n\r")
	graphWidth := utf8.RuneCountInString(graphChar) * currentSize
	graphPadding := max(0, (termWidth-graphWidth)/2)
	for _, line := range graph {
		fmt.Println(strings.Repeat(" ", graphPadding) + line + clearLine + "\r")
	}
}
