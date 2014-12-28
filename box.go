package main

import (
	"fmt"
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

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.OutputNormal)

	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	// draw_keyboard()
	//termbox.Flush()
	// inputmode := 0
	// ctrlxpressed := false

	// dispatch_press(&ev)

	_, y := printText(0, 0, title)
	printText(0, y+1, "Press 'space' to start capturing keys, 'Esc' to quit")
	//fmt.Printf("before flush: (%d, %d)", x, y)
	termbox.Flush()
	// x, y = termbox.Size()
	// fmt.Printf("after flush: (%d, %d)", x, y)

	letter := make(chan rune)
	done := make(chan bool)
	defer close(letter)
	defer close(done)
	go captureKeys(letter, done)
	for {
		select {
		case l := <-letter:
			printText(0, 0, "Got "+strconv.QuoteRune(l))
			termbox.Flush()
		case <-done:
			fmt.Println("Done")
			return
		}

	}

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
	return x, y

}

func captureKeys(letter chan rune, done chan bool) {
	started := false
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeySpace {
			started = !started
			if started {
				printText(0, 0, "Started")
			} else {
				printText(0, 0, "Stopped")
			}
		}
		if ev.Ch < 128 && started {
			letter <- ev.Ch
		}
		if ev.Key == termbox.KeyEsc {
			done <- true
			return
		}
	}

}
