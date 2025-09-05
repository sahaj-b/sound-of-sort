package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		close(ch)
		return
	}
	ch <- string(buf[:n])
}

func handleInput(input string) bool {
	switch input {
	case "q", "\x03": // Ctrl+C
		return true
	}
	return false
}

func arrGraph(arr []int, colors []string) []string {
	graphHeight := termHeight - 3
	if len(arr) != len(colors) {
		log.Fatal("arr and colors must have the same length")
	}
	graphChar := "█▊" // █ ▉ ▊ ▋ ▌
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

func render(graph []string) {
	fmt.Print(moveToTop)
	fmt.Println()
	for _, line := range graph {
		fmt.Println(line + "\r")
	}
}
