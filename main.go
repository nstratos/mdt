package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/nstratos/mdt/draw"

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
	Mode      rune
	TotalTime int
	Offset    int
	BaseHz    float64
	StartHz   float64
	EndHz     float64
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
	'd': "Language thought",
	'e': "Language voice",
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

// Logs to .txt file in program's directory, named: S-E hz day date month time
// where S is start hz and E is end hz, e.g. '15-19 hz wed 27 dec 22.09.txt'
func logCaptures(c Config) error {
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
	captures = make([]Capture, 0)
	return nil
}

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

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	// Drawing title
	x, y := 0, 0
	_, y = printText(x, y, title)
	_, y = printText(x, y+2, "Press 'space' to start capturing keys, 'Esc' to quit")
	keepY := y
	_, y = printOptions(x, y+2, c)
	_, y = printKeyLabels(x+25, keepY+2)
	draw.Box()
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
				go timer(5, letter, endTimer, c)
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
			return
		case <-expired:
			end <- true
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
			printText(0, 1, fmt.Sprintf("Mouse clicked %d, %d", ev.MouseX, ev.MouseY))
			termbox.Flush()
		}
	}
}

// CurrentHz = Seconds * (EndHz - StartHz) / ((TotalTime - Offset) * 60)
func currentHz(seconds int, c Config) float64 {
	return (float64(seconds) * (c.EndHz - c.StartHz) / float64((c.TotalTime-c.Offset)*60)) + c.StartHz
}

func showResults(c Config) {
	printText(0, 0, "RESULTS")
	l := 1
	for _, capt := range captures {
		printText(0, l, fmt.Sprintf("%.2fhz @ %.0f base hz, on %v %v", capt.Hz, c.BaseHz, capt.Timestamp(), capt.Label()))
		l += 1
	}
	termbox.Flush()
}

func printOptions(x, y int, c Config) (finalX, finalY int) {
	x, y = printText(x, y, "Mode: "+strconv.QuoteRuneToASCII(c.Mode))
	x, y = printText(x, y+1, "TotalTime: "+strconv.Itoa(c.TotalTime)+" min")
	x, y = printText(x, y+1, "Offset: "+strconv.Itoa(c.Offset))
	x, y = printText(x, y+1, fmt.Sprintf("Base: %.2f hz", c.BaseHz))
	x, y = printText(x, y+1, fmt.Sprintf("Start: %.2f hz", c.StartHz))
	x, y = printText(x, y+1, fmt.Sprintf("End: %.2f hz", c.EndHz))
	return x, y
}

func printKeyLabels(x, y int) (finalX, finalY int) {
	x, y = printText(x, y, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('q'), labels['q']))
	x, y = printText(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('a'), labels['a']))
	x, y = printText(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('w'), labels['w']))
	x, y = printText(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('s'), labels['s']))
	x, y = printText(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('e'), labels['e']))
	x, y = printText(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('d'), labels['d']))
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
