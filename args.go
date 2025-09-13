package main

import (
	"flag"
	"os"

	"github.com/sahaj-b/sound-of-sort/algos"
)

var (
	initialVolume         = flag.Float64("volume", 10, "Initial volume level (0 to 100)")
	initialDelay          = flag.Float64("delay", 5, "Initial delay in milliseconds between operations")
	initialSize           = flag.Int("size", 100, "Initial array size")
	initialSort           = flag.String("sort", "quick", "Initial sorting algorithm")
	fpsFlag               = flag.Int("fps", 60, "Frames per second for rendering")
	listFlag              = flag.Bool("list", false, "List available sorting algorithms")
	imgFlag               = flag.Bool("img", false, "Enable image mode (Get ascii image from stdin)")
	horizontalFlag        = flag.Bool("horiz", false, "Horizontal image mode (rows instead of columns)")
	noColorsFlag          = flag.Bool("no-colors", false, "Strip ALL ANSI colors from output")
	noReadWriteColorsFlag = flag.Bool("no-rw-colors", false, "Disable read/write highlighting colors only")
	helpFlag              = flag.Bool("help", false, "Show help message")
)

func parseArgs() bool {
	flag.Parse()

	if os.Getenv("NO_COLOR") != "" {
		*noColorsFlag = true
		*noReadWriteColorsFlag = true
	}

	if *initialVolume < 0.0 || *initialVolume > 1.0 {
		*initialVolume = 0.1
	}
	if *initialDelay < 0.0 {
		*initialDelay = 5.0
	}
	if *fpsFlag <= 0 {
		*fpsFlag = 60
	}
	if *listFlag {
		printAvailableSorts()
		return true
	}
	if *helpFlag {
		flag.Usage()
		return true
	}
	return false
}

func printAvailableSorts() {
	println("Available sorting algorithms:")
	for _, s := range algos.Sorts {
		println("-", s.Name, "(arg:", s.Arg+")")
	}
}
