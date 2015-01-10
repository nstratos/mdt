package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"bitbucket.org/nstratos/mdt/draw"

	"github.com/nsf/termbox-go"
)

// Global holder of captured key presses.
var captures = make([]draw.Capture, 0)

// Logs to .txt file in program's directory, named: S-E hz day date month time
// where S is start hz and E is end hz, e.g. '15-19 hz wed 27 dec 22.09.txt'
func logCaptures(c draw.Config) error {
	if len(captures) == 0 {
		return nil
	}
	format := "Mon 02 Jan 15.04"
	filename := fmt.Sprintf("%v-%v hz %v", c.StartHz, c.EndHz, time.Now().Format(format))
	f, err := os.Create(filename + ".txt")
	//err := ioutil.WriteFile(, []byte(filename+"\n"), 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("%v\r\nMode: %v\r\n", filename, strconv.QuoteRune(c.Mode)))
	if err != nil {
		return err
	}
	for _, capt := range captures {
		_, err = f.WriteString(
			fmt.Sprintf("%.2fhz @ %.0f base hz, on %v %v\r\n",
				capt.Hz, c.BaseHz, capt.Timestamp(), capt.Label()))
		if err != nil {
			return err
		}
	}
	// Emptying capture holder.
	captures = nil
	captures = make([]draw.Capture, 0)
	return nil
}

func main() {
	// Loading configuration from config.json
	c := draw.Config{}
	c.Load()

	// Initializing termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	// Drawing title
	x, y := 0, 0
	_, y = draw.Title()
	mx, my := termbox.Size()
	draw.Text(0, 0, fmt.Sprintf("%vx%v", mx, my))
	//_, y = draw.Text(x, y+1, "Press 'space' to start capturing keys, 'Esc' to quit")
	keepY := y
	_, y = draw.DrawOptions(0, y+1, c)
	sb := draw.NewStatusBar(0, y+1, 54, "Press 'space' to start capturing keys, 'Esc' to quit.")
	sb.Draw()
	draw.DrawKeyLabels(x+25, keepY+1)

	termbox.Flush()

	letter := make(chan rune)
	start := make(chan bool)
	done := make(chan bool)
	endTimer := make(chan bool)
	defer close(letter)
	defer close(start)
	defer close(done)
	defer close(endTimer)
	go captureEvents(letter, start, done)
	capturing := false
	timerEnded := false
loop:
	for {
		select {
		case <-start:
			capturing = !capturing
			if capturing {
				go timer(5, letter, endTimer, c, *sb)
				timerEnded = false
			}
			if !capturing && !timerEnded {
				endTimer <- true
				//showResults(c)
				logCaptures(c)
			}
		case timerEnded, _ = <-endTimer:
			capturing = false
			//showResults(c)
			logCaptures(c)
		case l := <-letter:
			// If the timer is on, we keep resending the letter to the channel so
			// that it will be eventually captured by the timer. If the timer is
			// not on, we discard the letter. Without this case the channel would
			// block forever if it was sent a letter without timer to consume it.
			if capturing {
				letter <- l
			}
		case <-done:
			break loop
		}
	}
}

func timer(maxSeconds int, letter chan rune, end chan bool, c draw.Config, sb draw.StatusBar) {
	seconds := 0
	// min, sec := 0, 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	sb.UpdateTimer(seconds)
	for {
		select {
		case l := <-letter:
			capture := draw.Capture{Value: l, Seconds: seconds, Hz: draw.CurrentHz(seconds, c)}
			captures = append(captures, capture)
			//printText(0, 0, fmt.Sprintf("Recorded %v (%v)", strconv.QuoteRune(l), currentHz(seconds, c)))
			sb.UpdateText(fmt.Sprintf("Recorded %v (%.2fhz) \"%v\"", strconv.QuoteRune(l), draw.CurrentHz(seconds, c), draw.Labels[l]))
			//termbox.Flush()
		case <-end:
			return
		case <-expired:
			end <- true
			return
		case <-tick:
			seconds += 1
			sb.UpdateTimer(seconds)
		}
	}
}

func captureEvents(letter chan rune, start chan bool, done chan bool) {
	started := false
	for {
		ev := termbox.PollEvent()
		switch {
		case ev.Key == termbox.KeyEsc:
			done <- true
		case ev.Key == termbox.KeySpace:
			started = !started
			start <- started
		case supportedLabel(ev.Ch):
			letter <- ev.Ch
		case ev.Type == termbox.EventMouse:
			cell := draw.GetCell(ev.MouseX, ev.MouseY)
			draw.Text(2, 22, fmt.Sprintf("Mouse clicked (%d, %d) = %v", ev.MouseX, ev.MouseY, strconv.QuoteRuneToASCII(cell.Ch)))
			//input := draw.GetInput(ev.MouseX, ev.MouseX+3, ev.MouseY)
			//printText(2, 22, fmt.Sprintf("Mouse clicked (%d, %d) = %v", ev.MouseX, ev.MouseY, input))
			termbox.Flush()
			// case ev.Key == termbox.KeyEnter:
			// 	draw.Input()
		}
	}
}

// 'a' = 97		-> visual imagination
// 'd' = 100	-> language thought
// 'e' = 101	-> language voice
// 'q' = 113	-> visual memory
// 's' = 115	-> auditory imagination
// 'w' = 119	-> auditory memory
func supportedLabel(key rune) bool {
	if key == 'a' || key == 'd' || key == 'e' ||
		key == 'q' || key == 's' || key == 'w' {
		return true
	}
	return false
}
