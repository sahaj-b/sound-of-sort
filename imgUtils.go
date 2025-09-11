package main

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
)

// parseCells groups any contiguous SGR (ANSI) sequences with the following rune as one visual cell.
// Trailing SGR sequences (like a final reset) get merged into the last cell so colors don't leak.
func parseCells(s string) []string {
	if s == "" {
		return nil
	}
	cells := []string{}
	pending := ""
	for i := 0; i < len(s); {
		if s[i] == '\x1b' { // start of ANSI
			j := i + 1
			for j < len(s) && s[j] != 'm' { // scan until 'm'
				j++
			}
			if j < len(s) { // include 'm'
				j++
				pending += s[i:j]
				i = j
				continue
			}
			// truncated sequence: stop
			break
		}
		r, sz := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && sz == 1 { // invalid byte, treat as raw
			cells = append(cells, pending+string(s[i]))
			pending = ""
			i++
			continue
		}
		cells = append(cells, pending+string(r))
		pending = ""
		i += sz
	}
	if pending != "" && len(cells) > 0 { // attach any trailing SGR (usually reset)
		cells[len(cells)-1] += pending
	}
	return cells
}

func transpose(img []string) ([]string, error) {
	if len(img) == 0 {
		return nil, nil
	}

	cellsPerLine := make([][]string, len(img))
	maxWidth := 0
	for i, line := range img {
		cells := parseCells(line)
		cellsPerLine[i] = cells
		if len(cells) > maxWidth {
			maxWidth = len(cells)
		}
	}
	if maxWidth == 0 {
		return nil, nil
	}

	transposed := make([]string, maxWidth)
	for col := 0; col < maxWidth; col++ {
		var b strings.Builder
		for row := 0; row < len(cellsPerLine); row++ {
			if col < len(cellsPerLine[row]) {
				b.WriteString(cellsPerLine[row][col])
			} else {
				b.WriteByte(' ')
			}
		}
		transposed[col] = b.String()
	}
	return transposed, nil
}

func makeRectangular(img []string) {
	if len(img) == 0 {
		return
	}
	maxWidth := getLineWidth(img[0])
	for _, line := range img {
		if w := getLineWidth(line); w > maxWidth {
			maxWidth = w
		}
	}
	for i, line := range img {
		w := getLineWidth(line)
		if w < maxWidth {
			img[i] = line + strings.Repeat(" ", maxWidth-w)
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
	return utf8.RuneCountInString(stripansi.Strip(str))
}
