package main

import (
	"bufio"
	"os"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
)

func transpose(img []string) ([]string, error) {
	if len(img) == 0 {
		return nil, nil
	}

	// Strip ANSI codes to get the actual character positions
	stripped := make([]string, len(img))
	for i, line := range img {
		stripped[i] = stripansi.Strip(line)
	}

	if len(stripped[0]) == 0 {
		return nil, nil
	}

	// Create transposed output based on stripped character positions
	transposed := make([]string, len(stripped[0]))

	for col := 0; col < len(stripped[0]); col++ {
		var newRow string
		for row := 0; row < len(stripped); row++ {
			if col < len(stripped[row]) {
				// Extract the character at this position from the original (with ANSI codes)
				char := getCharAtPosition(img[row], col)
				newRow += char
			} else {
				newRow += " "
			}
		}
		transposed[col] = newRow
	}

	return transposed, nil
}

// Extract character at visual position, preserving ANSI codes
func getCharAtPosition(line string, pos int) string {
	if pos < 0 {
		return " "
	}

	visualPos := 0
	i := 0

	for i < len(line) {
		if line[i] == '\x1b' {
			// Find the end of the ANSI sequence
			j := i
			for j < len(line) && line[j] != 'm' {
				j++
			}
			if j < len(line) {
				j++ // include the 'm'
			}
			if visualPos == pos {
				// Return the ANSI sequence + the character that follows
				ansiCode := line[i:j]
				if j < len(line) {
					r, _ := utf8.DecodeRuneInString(line[j:])
					return ansiCode + string(r)
				}
				return ansiCode + " "
			}
			i = j
		} else {
			if visualPos == pos {
				r, _ := utf8.DecodeRuneInString(line[i:])
				return string(r)
			}
			_, size := utf8.DecodeRuneInString(line[i:])
			i += size
			visualPos++
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
