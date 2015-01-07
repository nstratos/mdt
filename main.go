package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const title = `
      ___           ___           ___     
     /\__\         /\  \         /\  \    
    /::|  |       /::\  \        \:\  \   
   /:|:|  |      /:/\:\  \        \:\  \  
  /:/|:|__|__   /:/  \:\__\       /::\  \ 
 /:/ |::::\__\ /:/__/ \:|__|     /:/\:\__\
 \/__/~~/:/  / \:\  \ /:/  /    /:/  \/__/
       /:/  /   \:\  /:/  /    /:/  /     
      /:/  /     \:\/:/  /     \/__/      
     /:/  /       \::/__/                 
     \/__/         ~~                     
`

type Config struct {
	Mode     rune
	Duration int
	Offset   int
	BaseHz   float64
	StartHz  float64
	EndHz    float64
}

func (c Config) Save() error {
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", json, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Load() error {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}

var labels = map[rune]string{
	'a': "Visual imagination",
	'd': "?",
	'e': "?",
	'q': "Visual memory",
	's': "Auditory imagination",
	'w': "Auditory memory",
}

type Capture struct {
	Value   rune
	Seconds int
	Hz      float64
}

func (c *Capture) Label() string {
	return labels[c.Value]
}

func (c *Capture) Timestamp() string {
	min := c.Seconds / 60
	sec := c.Seconds % 60
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

// Global holder of captured key presses.
var captures = make([]Capture, 0)

// Current line to print text.
var line = 2

func main() {
	// Loading configuration from config.json
	c := Config{}
	c.Load()

	// Initializing termbox
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	// Drawing title
	x, y := 0, 0
	_, y = printText(x, y, title)
	_, y = printText(x, y+2, "Press 'space' to start capturing keys, 'Esc' to quit")
	_, y = printOptions(x, y+2, c)
	termbox.Flush()

	/*started := false
		//lock := false
		timerDone := make(chan bool)
		defer close(timerDone)
	captureloop:
		for {
			select {
			case <-timerDone:
				printText(0, 0, "--------------------------------------------------")
				termbox.Flush()
				started = false
				break captureloop
			default:
				//printText(0, 0, fmt.Sprintf("CAPTURE LOOP started=%v, lock=%v", started, lock))
				//termbox.Flush()
				switch ev := termbox.PollEvent(); ev.Type {
				case termbox.EventKey:
					switch {
					// Esc -> quit
					case ev.Key == termbox.KeyEsc:
						break captureloop
					// Capturing + space -> session end
					case started && (ev.Key == termbox.KeySpace):
						started = false
						_, y = printText(0, 0, "Session stopped manually")
						termbox.Flush()
						//if !lock {
						timerDone <- true
						//}
						break captureloop
					case !started && (ev.Key == termbox.KeySpace):
						started = true
						//lock = true
						_, y = printText(0, 0, "Session started")
						termbox.Flush()
						//startTimer(c.Duration * 60)
						go startTimer(5, timerDone)
						//startTimerAlt(5)
						//break captureloop
					case started && supportedLabel(ev.Ch):
						printText(0, 0, fmt.Sprintf("Got %v", strconv.QuoteRune(ev.Ch)))
						termbox.Flush()
					}
				case termbox.EventError:
					log.Println(ev.Err)
				}
			}

		}
	uiloop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch {
				case ev.Key == termbox.KeyEsc:
					break uiloop
				}
			case termbox.EventError:
				log.Println(ev.Err)
			}
		}*/

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
	//timerRunning := false
	timerEnded := false
loop:
	for {
		select {
		// case l := <-letter:
		// 	printText(0, 0, fmt.Sprintf("Got %v", strconv.QuoteRune(l)))
		// 	termbox.Flush()
		// 	if capturing {
		// 		record(l)
		// 	}
		case <-start:
			capturing = !capturing
			printTextNext(fmt.Sprintf("1. capturing: %v, timerEnded: %v", capturing, timerEnded))
			termbox.Flush()
			if capturing {
				//timerRunning = true
				//printText(0, 2, fmt.Sprintf("1. capturing: %v, timerRunning: %v", capturing, timerRunning))
				printTextNext("STARTING NEW TIMER")
				termbox.Flush()
				go timer(5, letter, endTimer, c)
				timerEnded = false
			}

			if !capturing && !timerEnded {
				printTextNext(fmt.Sprintf("2. capturing: %v, timerEnded: %v", capturing, timerEnded))
				termbox.Flush()
				//_, ok := <-endTimer
				//printText(0, 2, fmt.Sprintf("OK: %v", ok))
				endTimer <- true
				showResults(c)
				//break loop
			}
			if !capturing && timerEnded {
				printTextNext(fmt.Sprintf("3. capturing: %v, timerEnded: %v", capturing, timerEnded))
				termbox.Flush()
			}
		case timerEnded, _ = <-endTimer:
			capturing = false
			printTextNext(fmt.Sprintf("END TIMER capturing: %v, timerEnded: %v", capturing, timerEnded))
			termbox.Flush()
			showResults(c)
			//break loop
		case <-done:
			break loop
		}
	}

	//<-done
}

func timer(maxSeconds int, letter chan rune, end chan bool, c Config) {
	seconds := 0
	min, sec := 0, 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	for {
		select {
		case l := <-letter:
			capture := Capture{Value: l, Seconds: seconds, Hz: currentHz(seconds, c)}
			captures = append(captures, capture)
			printText(0, 0, fmt.Sprintf("Recorded %v (%v)", strconv.QuoteRune(l), currentHz(seconds, c)))
			termbox.Flush()
		case <-end:
			printText(0, 1, "timer stopped")
			termbox.Flush()
			return
		case <-expired:
			printText(0, 1, "Session ended")
			termbox.Flush()
			end <- true
			//close(end)
			return
		case <-tick:
			seconds = seconds + 1
			if sec == 59 {
				sec = -1
				min += 1
			}
			sec += 1
			printText(0, 1, fmt.Sprintf("%d min %d sec", min, sec))
			termbox.Flush()
		}
	}
}

func captureEvents(letter chan rune, start chan bool, done chan bool) {
	started := false
	//eventsloop:
	for {
		ev := termbox.PollEvent()
		switch {
		case ev.Key == termbox.KeyEsc:
			done <- true
			//return
		//break eventsloop
		case ev.Key == termbox.KeySpace:
			started = !started
			start <- started
		case supportedLabel(ev.Ch):
			letter <- ev.Ch
		}
	}
}

func currentHz(seconds int, c Config) float64 {
	return (float64(seconds) * (c.EndHz - c.StartHz) / float64((c.Duration-c.Offset)*60)) + c.StartHz
}

func showResults(c Config) {
	printText(0, 0, "RESULTS")
	l := 1
	for _, capt := range captures {
		printText(0, l, fmt.Sprintf("%.2fhz @ %.0f base hz, on %v %v", capt.Hz, c.BaseHz, capt.Timestamp(), capt.Label()))
		l += 1
	}
	termbox.Flush()
	// resultsloop:
	// 	for {
	// 		ev := termbox.PollEvent()
	// 		switch {
	// 		case ev.Key == termbox.KeyEsc:
	// 			break resultsloop
	// 		}
	// 	}
}

func captureKeysAlt(letter chan rune, done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			ev := termbox.PollEvent()
			switch {
			case ev.Key == termbox.KeyEsc:
				done <- true
				return
			case ev.Key == termbox.KeySpace:
				printText(0, 0, "Stopped fucking manually")
				termbox.Flush()
				done <- true
				return
			case ev.Ch < 128:
				letter <- ev.Ch
			default:
				return
			}
		}
	}
}

func startTimerAlt(maxSeconds int) {
	letter := make(chan rune)
	done := make(chan bool)
	defer close(letter)
	defer close(done)
	go captureKeysAlt(letter, done)

	min, sec := 0, 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	for {
		select {
		case <-expired:
			printText(0, 0, "Session ended")
			termbox.Flush()
			//done <- true
			//<-letter
			return
		case <-done:
			printText(0, 0, "manual stop.................")
			termbox.Flush()
			return
		case <-tick:
			if sec == 59 {
				sec = -1
				min += 1
			}
			sec += 1
			printText(0, 0, fmt.Sprintf("%d min %d sec", min, sec))
			termbox.Flush()
		case l := <-letter:
			printText(0, 0, fmt.Sprintf("Got %v", strconv.QuoteRune(l)))
			termbox.Flush()
		}
	}
}

func startTimer(maxSeconds int, end chan bool) {
	min, sec := 0, 0
	expired := time.NewTimer(time.Second * time.Duration(maxSeconds)).C
	tick := time.NewTicker(time.Second).C
	for {
		select {
		case <-end:
			return
		case <-expired:
			printText(0, 0, "Session ended")
			termbox.Flush()
			end <- true
			return
		case <-tick:
			if sec == 59 {
				sec = -1
				min += 1
			}
			sec += 1
			printText(0, 0, fmt.Sprintf("%d min %d sec", min, sec))
			termbox.Flush()
		}
	}
}

func printOptions(x, y int, c Config) (finalX, finalY int) {
	x, y = printText(x, y, "Mode: "+strconv.QuoteRuneToASCII(c.Mode))
	x, y = printText(x, y+1, "Duration: "+strconv.Itoa(c.Duration)+" min")
	x, y = printText(x, y+1, "Offset: "+strconv.Itoa(c.Offset))
	x, y = printText(x, y+1, fmt.Sprintf("Base: %.2f hz", c.BaseHz))
	x, y = printText(x, y+1, fmt.Sprintf("Start: %.2f hz", c.StartHz))
	x, y = printText(x, y+1, fmt.Sprintf("End: %.2f hz", c.EndHz))
	return x, y
}

func printText(x, y int, s string) (finalX, finalY int) {
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
		x = tempx
		y += 1
	}
	// Because we always icrease it one more.
	y -= 1
	return x, y
}
func printTextNext(s string) (finalX, finalY int) {
	line = line + 1
	return printText(0, line, s)
}

// 'a' = 97		-> visual imagination
// 'd' = 100	-> ?
// 'e' = 101	-> ?
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
