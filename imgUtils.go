package main

import (
	"bufio"
	"os"

	"github.com/acarl005/stripansi"
)

func transpose(img []string) ([]string, error) {
	if len(img) == 0 {
		return nil, nil
	}

	// Strip ANSI codes to get visual dimensions
	stripped := make([]string, len(img))
	for i, line := range img {
		stripped[i] = stripansi.Strip(line)
	}

	if len(stripped) == 0 || len(stripped[0]) == 0 {
		return nil, nil
	}

	// Find max width
	maxWidth := 0
	for _, line := range stripped {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Create transposed output - each column becomes a row
	transposed := make([]string, maxWidth)

	for col := 0; col < maxWidth; col++ {
		var newRow string
		for row := 0; row < len(stripped); row++ {
			if col < len(stripped[row]) {
				// Get the character at this visual position from original line with colors
				char := getCharWithColorAtPosition(img[row], col)
				newRow += char
			} else {
				newRow += " "
			}
		}
		transposed[col] = newRow
	}

	return transposed, nil
}

// Extract character with ANSI colors at visual position
func getCharWithColorAtPosition(line string, visualPos int) string {
	if visualPos < 0 {
		return " "
	}

	currentVisualPos := 0
	i := 0

	for i < len(line) {
		if line[i] == '\x1b' {
			// Find end of ANSI sequence
			j := i
			for j < len(line) && line[j] != 'm' {
				j++
			}
			if j < len(line) {
				j++ // include 'm'
			}

			if currentVisualPos == visualPos {
				// Return ANSI sequence + following character
				ansiCode := line[i:j]
				if j < len(line) {
					return ansiCode + string(line[j])
				}
				return ansiCode + " "
			}
			i = j
		} else {
			if currentVisualPos == visualPos {
				return string(line[i])
			}
			i++
			currentVisualPos++
		}
	}

	return " "
}

func makeRectangular(img []string) {
	maxWidth := getLineWidth(img[0])
	for _, line := range img {
		maxWidth = max(maxWidth, len(stripansi.Strip(line)))
	}
	for i, line := range img {
		lineWidth := len(stripansi.Strip(line))
		if lineWidth < maxWidth {
			padding := make([]byte, maxWidth-lineWidth)
			for j := range padding {
				padding[j] = ' '
			}
			img[i] = line + string(padding)
		}
	}
}

func readImageFromStdin() ([]string, error) {
	var img []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		img = append(img, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return img, nil
}

func getLineWidth(str string) int {
	return len(stripansi.Strip(str))
}
