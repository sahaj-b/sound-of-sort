package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
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
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// No TTY available, run in non-interactive mode
		return func() {}
	}

	oldState, err := term.MakeRaw(int(tty.Fd()))
	if err != nil {
		log.Fatal("Error setting terminal to raw mode:", err)
	}

	width, height, err := term.GetSize(int(tty.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
	}

	termWidth, termHeight = width, height
	fmt.Print(hideCursor)
	fmt.Print(clear)
	return func() {
		term.Restore(int(tty.Fd()), oldState)
		fmt.Print(showCursor)
		tty.Close()
	}
}

func inputReader(ch chan string) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
	if err != nil {
		// if no tty is available
		close(ch)
		return
	}
	defer tty.Close()
	defer close(ch)
	buf := make([]byte, 3)
	for {
		n, err := tty.Read(buf)
		if err != nil || n == 0 {
			return
		}
		ch <- string(buf[:n])
	}
}

func (app *App) handleInput(input string) bool {
	switch input {
	case "q", "\x03": // Ctrl+C
		return true

	case "w": // increase delay
		app.delay.Add(100 * int64(time.Microsecond))

	case "s": // decrease delay
		// use a CAS loop
		for {
			currentDelay := app.delay.Load()

			newDelay := max(1*int64(time.Microsecond), currentDelay-(100*int64(time.Microsecond)))

			// this returns false if the currentDelay got changed since we loaded it
			if app.delay.CompareAndSwap(currentDelay, newDelay) {
				break
			}
		}

	case "\x1b[A": // Up arrow
		oldBits := app.volume.Load()
		oldFloat := math.Float64frombits(oldBits)

		newFloat := min(oldFloat+0.005, 1.0)

		newBits := math.Float64bits(newFloat)
		app.volume.Store(newBits)

	case "\x1b[B": // Down arrow
		oldFloat := math.Float64frombits(app.volume.Load())
		newFloat := max(0.0, oldFloat-0.01)
		newBits := math.Float64bits(newFloat)
		app.volume.Store(newBits)

	case "p", "\x1b[D": // Left arrow
		newIndex := app.currentSortIndex.Add(1)
		app.currentSortIndex.Store(newIndex % int32(len(algos.Sorts)))
	case "n", "\x1b[C": // Right arrow
		newIndex := app.currentSortIndex.Add(-1)
		if newIndex < 0 {
			newIndex = int32(len(algos.Sorts) - 1)
		}
		app.currentSortIndex.Store(newIndex)

	case "a": // decrease array size
		if app.imgMode {
			return false
		}
		newSize := max(0, app.currentSize.Add(-10))
		app.currentSize.Store(newSize)

	case "d": // increase array size
		if app.imgMode {
			return false
		}
		newSize := app.currentSize.Add(10)
		app.currentSize.Store(newSize)

	case "r": // shuffle array
		app.shuffleRequested.Store(true)
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

func imgGraph(arr []int, img []string, colors []string) []string {
	if len(img) == 0 {
		return img
	}

	firstCol := parseCells(img[0])
	height := len(firstCol)
	output := make([]string, height)
	for i, val := range arr {
		if val < 0 || val >= len(img) {
			for r := 0; r < height; r++ {
				output[r] += " "
			}
			continue
		}
		col := parseCells(img[val])
		for r := 0; r < height; r++ {
			if r < len(col) && col[r] != "" {
				if colors[i] != "" {
					stripped := stripansi.Strip(col[r])
					if len(stripped) > 0 {
						output[r] += colors[i] + string([]rune(stripped)[0]) + reset
					} else {
						output[r] += " "
					}
				} else {
					output[r] += col[r] + reset
				}
			} else {
				output[r] += " "
			}
		}
	}
	return output
}

func imgGraphHorizontal(arr []int, img []string, colors []string) []string {
	if len(img) == 0 {
		return img
	}
	// guard: sometimes a race elsewhere could hand us fewer colors; pad so we don't blow up
	if len(colors) < len(arr) {
		tmp := make([]string, len(arr))
		copy(tmp, colors)
		colors = tmp
	}
	output := make([]string, len(arr))
	for i, val := range arr {
		if val < 0 || val >= len(img) {
			output[i] = ""
			continue
		}
		line := img[val]
		c := ""
		if i < len(colors) {
			c = colors[i]
		}
		if c != "" {
			output[i] = c + line + reset
		} else {
			output[i] = line
		}
	}
	return output
}

func (app *App) render(graph []string, sortName string) {
	fmt.Print(moveToTop)
	arrSizeStr := fmt.Sprintf("a/d Arr Size: %d", app.currentSize.Load())
	if app.imgMode {
		arrSizeStr = fmt.Sprintf("Image Size: %d", app.currentSize.Load())
	}
	statusStrs := []string{
		"←/→: " + sortName,
		"↑/↓ Volume: " + fmt.Sprintf("%.1f", math.Float64frombits(app.volume.Load())*100),
		"w/s Delay: " + time.Duration(app.delay.Load()).String(),
		arrSizeStr,
		"r to reshuffle",
		"q to quit",
	}

	n := len(statusStrs)
	statusStr := bggray + cyan + " " + strings.Join(statusStrs[:n-1], " "+reset+" "+bggray+" "+cyan) + " " + reset + " " + bggray + red + " " + statusStrs[n-1] + " " + reset
	statusLen := 0
	for _, s := range statusStrs {
		statusLen += len(s)
	}
	statusLen += 3 * (len(statusStrs) - 1)
	statusPadding := max(0, (termWidth-statusLen)/2)
	fmt.Println(strings.Repeat(" ", statusPadding) + statusStr + "\r\n\r")

	w := getLineWidth(graph[0])
	graphPadding := max(0, (termWidth-w)/2)
	for _, line := range graph {
		fmt.Println(reset + strings.Repeat(" ", graphPadding) + line + reset + clearLine + "\r")
	}
}
