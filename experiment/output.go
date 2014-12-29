package main

import "github.com/nsf/termbox-go"

var output_mode = termbox.OutputNormal

func print_wide(x, y int, s string) {
	red := false
	for _, r := range s {
		c := termbox.ColorDefault
		if red {
			c = termbox.ColorRed
		}
		termbox.SetCell(x, y, r, termbox.ColorDefault, c)
		x += 2
		red = !red
	}
}

const hello_world = "こんにちは世界"

func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	print_wide(0, 0, hello_world)

	termbox.Flush()
}

var available_modes = []termbox.OutputMode{
	termbox.OutputNormal,
	termbox.OutputGrayscale,
	termbox.Output216,
	termbox.Output256,
}

var output_mode_index = 0

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	draw_all()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowUp, termbox.KeyArrowRight:
				draw_all()
			case termbox.KeyArrowDown, termbox.KeyArrowLeft:
				draw_all()
			}
		case termbox.EventResize:
			draw_all()
		}
	}
}
