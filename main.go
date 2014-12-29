package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

const title = `
      ___           ___           ___     
     /\__\         /\  \         /\  \    
    /::|  |        \:\  \       /::\  \   
   /:|:|  |         \:\  \     /:/\:\  \  
  /:/|:|__|__       /::\  \   /:/  \:\__\ 
 /:/ |::::\__\     /:/\:\__\ /:/__/ \:|__|
 \/__/~~/:/  /    /:/  \/__/ \:\  \ /:/  /
       /:/  /    /:/  /       \:\  /:/  / 
      /:/  /     \/__/         \:\/:/  /  
     /:/  /                     \::/__/   
     \/__/                       ~~       
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

	letter := make(chan rune)
	done := make(chan bool)
	defer close(letter)
	defer close(done)
	go captureKeys(letter, done)
	for {
		select {
		case l := <-letter:
			// printText(0, 0, "Got "+strconv.QuoteRune(l))
			printText(0, 0, fmt.Sprintf("Got %v", strconv.QuoteRune(l)))
			termbox.Flush()
		case <-done:
			fmt.Println("Done")
			return
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

func captureKeys(letter chan rune, done chan bool) {
	started := false
	for {
		ev := termbox.PollEvent()
		switch {
		case ev.Key == termbox.KeyEsc:
			done <- true
			return
		case ev.Key == termbox.KeySpace:
			started = !started
			if started {
				printText(0, 0, "Started")
				termbox.Flush()
			} else {
				printText(0, 0, "Stopped")
				termbox.Flush()
			}
		case ev.Ch < 128 && started:
			letter <- ev.Ch
		}
	}

}
