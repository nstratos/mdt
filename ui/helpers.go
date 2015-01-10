package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

// Given the seconds passed, it returns the current Gz based on the formula:
// CurrentHz = Seconds * (EndHz - StartHz) / ((TotalTime - Offset) * 60)
func CurrentHz(seconds int) float64 {
	return (float64(seconds) * (config.EndHz - config.StartHz) / float64((config.TotalTime-config.Offset)*60)) + config.StartHz
}

func RecordedKeyText(key rune, seconds int) string {
	return fmt.Sprintf("Recorded %v (%.2fhz) \"%v\"", strconv.QuoteRune(key), CurrentHz(seconds), Labels[key])
}

func text(x, y int, s string) (maxX, maxY int) {
	mx := 0
	tempx := x
	text := []string{s}
	if strings.Contains(s, "\n") {
		text = strings.Split(s, "\n")
	}
	for _, t := range text {
		for _, r := range t {
			termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
			x += 1
		}
		if x > mx {
			mx = x
		}
		x = tempx
		y += 1
	}
	// Because we always icrease it one more.
	y -= 1
	return mx, y
}

// Draws text on the screen. When it encounters a new line
// it continues to draw from the next line.
func Text(x, y int, s string) (maxX, maxY int) {
	mx, my := text(x, y, s)
	termbox.Flush()
	return mx, my
}

func tbfill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func fill(x, y, w, h int, r rune) {
	tbfill(x, y, w, h, termbox.Cell{Ch: r})
}

func Fill(x, y, w, h int, r rune) {
	fill(x, y, w, h, r)
	termbox.Flush()
}

func formatTimer(seconds int) string {
	min := seconds / 60
	sec := seconds % 60
	m := fmt.Sprintf("%v", min)
	if min < 10 {
		m = fmt.Sprintf("0%v", min)
	}
	s := fmt.Sprintf("%v", sec)
	if sec < 10 {
		s = fmt.Sprintf("0%v", sec)
	}
	return fmt.Sprintf("%v:%v", m, s)
}

func rtoa(r rune) string {
	return strconv.QuoteRuneToASCII(r)
}
