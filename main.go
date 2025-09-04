package main

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

func main() {
	defer initUI()
}
