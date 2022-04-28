package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/nstratos/mdt/ui"

	"github.com/nsf/termbox-go"
)

// Global holder of captured key presses.
var captures = make([]ui.Capture, 0)

var (
	version     = "devel"
	showVersion = flag.Bool("v", false, "print program version and exit")
)

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
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("%v\r\nMode: %v\r\n", filename, c.Mode))
	if err != nil {
		return err
	}
	for _, capt := range captures {
		_, err = f.WriteString(
			fmt.Sprintf("%.2fhz @ %.2f base hz, on %v %v\r\n",
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
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], version, runtime.Version())
		os.Exit(0)
	}

	if err := ui.Init(version); err != nil {
		log.Println("Could not initialize: ", err)
		if werr := ioutil.WriteFile("debug.txt", []byte(fmt.Sprintf("%s", err)), 0644); werr != nil {
			log.Fatalln(werr)
			os.Exit(1)
		}
		os.Exit(1)
	}
	ui.DrawAll()
	defer ui.Close()

	letter := make(chan rune)
	input := make(chan *ui.Entry)
	start := make(chan bool)
	done := make(chan bool)
	endTimer := make(chan bool)
	defer close(letter)
	defer close(input)
	defer close(start)
	defer close(done)
	defer close(endTimer)
	go captureEvents(letter, input, start, done)
	capturing := false
	timerEnded := false
loop:
	for {
		select {
		case <-start:
			ui.DeselectAllInputs()
			capturing = !capturing
			if capturing {
				c := ui.GetConfig()
				go timer(c.TotalTime*60, c.Offset*60, letter, endTimer)
				timerEnded = false
			}
			if !capturing && !timerEnded {
				endTimer <- true
				if err := logCaptures(); err != nil {
					ui.Debug(fmt.Sprintf("Error logging to txt file: %v", err))
				}
			}
		case timerEnded, _ = <-endTimer:
			capturing = false
			if err := logCaptures(); err != nil {
				ui.Debug(fmt.Sprintf("Error logging to txt file: %v", err))
			}
		case l := <-letter:
			// If the timer is on, we keep resending the letter to the channel so
			// that it will be eventually captured by the timer. If the timer is
			// not on, we discard the letter. Without this case the channel would
			// block forever if it was sent a letter without timer to consume it.
			if capturing {
				letter <- l
			}
		case in := <-input:
			if si := ui.SelectedInput(); si != nil {
				if in.Enter {
					if err := si.Valid(); err != nil {
						ui.UpdateText(fmt.Sprintf("Invalid value (%v)", err))
						continue
					}
					m, err := si.ValueMap()
					if err != nil {
						ui.UpdateText(fmt.Sprintf("%v", err))
						continue
					}
					c := ui.GetConfig()
					if err := c.Update(m); err != nil {
						ui.UpdateText(fmt.Sprintf("%v", err))
						continue
					}
					if err := c.Validate(); err != nil {
						ui.UpdateText(fmt.Sprintf("Invalid value (%v)", err))
						continue
					}
					if err := c.Save(); err != nil {
						ui.UpdateText(fmt.Sprintf("Could not save (%v)", err))
						continue
					}
					ui.UpdateConfig(c)
					ui.ReloadInputs(c)
					ui.UpdateText("Configuration changed successfully.")
					ui.DeselectAllInputs()
				} else {
					si.SetBuf(in)
				}
			}
		case <-done:
			if ui.SelectedInput() == nil {
				break loop
			}
			ui.DeselectAllInputs()
			ui.ResetText()
		}
	}
}

func timer(maxSeconds, offsetSeconds int, letter chan rune, end chan bool) {
	seconds := 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	ui.UpdateTimer(seconds)
	ui.UpdateText("New Session started, press 'space' to stop, 'Esc' to quit.")
	ui.Debug(fmt.Sprintf("Key Capturing starts in %v", ui.FormatTimer(offsetSeconds)))
	for {
		select {
		case l := <-letter:
			// If user has set an offset it means that we have to wait for that amount
			// of seconds. Thus unless it reaches 0 we ignore label keypresses.
			if offsetSeconds == 0 {
				capture := ui.Capture{Value: l, Seconds: seconds, Hz: ui.CurrentHz(seconds)}
				captures = append(captures, capture)
				ui.UpdateText(ui.RecordedKeyText(l, seconds))
			}
		case <-end:
			ui.UpdateText("Session stopped manually.")
			return
		case <-expired:
			end <- true
			ui.UpdateText("Session ended.")
			return
		case <-tick:
			seconds++
			ui.UpdateTimer(seconds)
			if offsetSeconds == 0 {
				ui.Debug("Key Capturing has started")
			} else {
				offsetSeconds--
				ui.Debug(fmt.Sprintf("Key Capturing starts in %v", ui.FormatTimer(offsetSeconds)))
			}

		}
	}
}

func captureEvents(letter chan rune, input chan *ui.Entry, start, done chan bool) {
	started := false
	for {
		ev := termbox.PollEvent()
		switch {
		case ev.Key == termbox.KeyEsc:
			done <- true
		case ev.Key == termbox.KeySpace:
			started = !started
			start <- started
		case ui.AllowedEntry(ev):
			input <- ui.NewEntry(ev)
		case supportedLabel(ev.Ch):
			letter <- ev.Ch
		case ev.Type == termbox.EventResize:
			ui.DrawAll()
		case ev.Type == termbox.EventMouse:
			cell := ui.GetCell(ev.MouseX, ev.MouseY)
			if cell.Input != nil {
				if cell.Input.Type == ui.InputSwitch {
					ui.DeselectAllInputs()
					if err := cell.Input.Switch(); err != nil {
						ui.UpdateText(fmt.Sprintf("%v", err))
						//ui.Debug(fmt.Sprintf("switch. %v", err))
						continue
					}
				} else {
					cell.Input.SetSelected(true)
				}
			} else {
				ui.DeselectAllInputs()
			}
		}
	}
}

// 'a' = 97	-> visual imagination
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
