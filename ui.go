package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func initUI() (restore func()) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal("Error setting terminal to raw mode:", err)
	}

	fmt.Print(hideCursor)

	fmt.Println("Press any key to see its value (Ctrl+C to exit):")

	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print(showCursor)
	}
}

func getInput() (string, error) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return "", err
	}
	return string(buf[:n]), nil
}

func handleInput(input string) {
	switch input {
	case "q", "\x03": // Ctrl+C
		fmt.Println("\nExiting...")
		os.Exit(0)
	}
}
