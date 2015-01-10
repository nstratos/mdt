package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"bitbucket.org/nstratos/mdt/ui"

	"github.com/nsf/termbox-go"
)

// Global holder of captured key presses.
var captures = make([]ui.Capture, 0)

// Logs to .txt file in program's directory, named: S-E hz day date month time
// where S is start hz and E is end hz, e.g. '15-19 hz wed 27 dec 22.09.txt'
func logCaptures() error {
	c := ui.GetConfig()
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
	captures = make([]ui.Capture, 0)
	return nil
}

func main() {

	if err := ui.Init(); err != nil {
		log.Println("Could not initialize: ", err)
		if err := ioutil.WriteFile("debug.txt", []byte(fmt.Sprintf("%s", err)), 0644); err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		os.Exit(1)
	}
	ui.DrawAll()
	defer ui.Close()

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
				go timer(5, letter, endTimer)
				timerEnded = false
			}
			if !capturing && !timerEnded {
				endTimer <- true
				logCaptures()
			}
		case timerEnded, _ = <-endTimer:
			capturing = false
			logCaptures()
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

func timer(maxSeconds int, letter chan rune, end chan bool) {
	seconds := 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	ui.UpdateTimer(seconds)
	for {
		select {
		case l := <-letter:
			capture := ui.Capture{Value: l, Seconds: seconds, Hz: ui.CurrentHz(seconds)}
			captures = append(captures, capture)
			ui.UpdateText(ui.RecordedKeyText(l, seconds))
		case <-end:
			return
		case <-expired:
			end <- true
			return
		case <-tick:
			seconds += 1
			ui.UpdateTimer(seconds)
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
			cell := ui.GetCell(ev.MouseX, ev.MouseY)
			ui.Text(2, 22, fmt.Sprintf("Mouse clicked (%d, %d) = %v", ev.MouseX, ev.MouseY, strconv.QuoteRuneToASCII(cell.Ch)))
			//input := ui.GetInput(ev.MouseX, ev.MouseX+3, ev.MouseY)
			//printText(2, 22, fmt.Sprintf("Mouse clicked (%d, %d) = %v", ev.MouseX, ev.MouseY, input))
			termbox.Flush()
			// case ev.Key == termbox.KeyEnter:
			// 	ui.Input()
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
