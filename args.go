package main

import (
	"flag"

	"github.com/sahaj-b/sound-of-sort/algos"
)

var (
	initialVolume = flag.Float64("volume", 10, "Initial volume level (0 to 100)")
	initialDelay  = flag.Int64("delay", 5, "Initial delay in milliseconds between operations")
	initialSize   = flag.Int("size", 100, "Initial array size")
	initialSort   = flag.String("sort", "quick", "Initial sorting algorithm")
	fpsFlag       = flag.Int("fps", 60, "Frames per second for rendering")
	listFlag      = flag.Bool("list", false, "List available sorting algorithms")
	imgFlag       = flag.Bool("img", false, "Enable image mode (Get ascii image from stdin)")
	helpFlag      = flag.Bool("help", false, "Show help message")
)

func parseArgs() bool {
	flag.Parse()
	if *initialVolume < 0.0 || *initialVolume > 1.0 {
		*initialVolume = 0.1
	}
	if *initialDelay < 0 {
		*initialDelay = 5
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
