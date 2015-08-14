package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

// CurrentHz given the seconds passed, it returns the current Gz based on the formula:
// CurrentHz = ( Seconds * Abs((EndHz - StartHz)) / ((TotalTime - Offset) * 60) ) + StartHz
// for Offset < TotalTime
// func CurrentHz(seconds int) float64 {
// 	offset := config.Offset
// 	if offset >= config.TotalTime {
// 		offset = 0
// 	}
// 	return (float64(seconds) * (math.Abs(config.EndHz - config.StartHz)) / float64((config.TotalTime-offset)*60)) + config.StartHz
// }
func CurrentHz(currentSecs int) float64 {

	hzPerSecond := float64((config.EndHz - config.StartHz) / (float64((config.TotalTime - config.Offset) * 60)))
	secondsSinceOffset := float64(currentSecs - (config.Offset * 60))

	currentHz := float64(hzPerSecond*secondsSinceOffset + config.StartHz)

	return currentHz
}

// RecordedKeyText returns a message indicating the key pressed, it's hz value and a timestamp of when it was received.
func RecordedKeyText(key rune, seconds int) string {
	return fmt.Sprintf("Recorded %v (%.2fhz) on %v \"%v\"", strconv.QuoteRune(key), CurrentHz(seconds), FormatTimer(seconds), Labels[key])
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
			x++
		}
		if x > mx {
			mx = x
		}
		x = tempx
		y++
	}
	// Because we always icrease it one more.
	y--
	return mx, y
}

// Text draws text on the screen. When it encounters a new line
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

// Fill fills cells of the screen with a rune r, starting from x, y and
// reaching a certain width w and height h.
func Fill(x, y, w, h int, r rune) {
	fill(x, y, w, h, r)
	termbox.Flush()
}

// FormatTimer accepts seconds and returns them in a timer format.
func FormatTimer(seconds int) string {
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

// Debug prints a debug message on the screen.
func Debug(s string) {
	x, y := termbox.Size()
	fill(2, y-2, x-1, 1, ' ')
	Text(2, y-2, s)
}

func flush() {
	termbox.Flush()
}

func setCursor(x, y int) {
	termbox.SetCursor(x, y)
}
