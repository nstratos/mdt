package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type key struct {
	x  int
	y  int
	ch rune
}

var K_a = []key{{7, 8, 'a'}}
var K_A = []key{{7, 8, 'A'}}
var K_s = []key{{10, 8, 's'}}
var K_S = []key{{10, 8, 'S'}}
var K_d = []key{{13, 8, 'd'}}
var K_D = []key{{13, 8, 'D'}}
var K_f = []key{{16, 8, 'f'}}
var K_F = []key{{16, 8, 'F'}}
var K_g = []key{{19, 8, 'g'}}
var K_G = []key{{19, 8, 'G'}}
var K_h = []key{{22, 8, 'h'}}
var K_H = []key{{22, 8, 'H'}}
var K_j = []key{{25, 8, 'j'}}
var K_J = []key{{25, 8, 'J'}}
var K_k = []key{{28, 8, 'k'}}
var K_K = []key{{28, 8, 'K'}}
var K_l = []key{{31, 8, 'l'}}
var K_L = []key{{31, 8, 'L'}}
var K_SEMICOLON = []key{{34, 8, ';'}}
var K_PARENTHESIS = []key{{34, 8, ':'}}
var K_QUOTE = []key{{37, 8, '\''}}
var K_DOUBLEQUOTE = []key{{37, 8, '"'}}
var K_K_4 = []key{{65, 8, '4'}}
var K_K_5 = []key{{68, 8, '5'}}
var K_K_6 = []key{{71, 8, '6'}}
var K_LSHIFT = []key{{1, 10, 'S'}, {2, 10, 'H'}, {3, 10, 'I'}, {4, 10, 'F'}, {5, 10, 'T'}}
var K_z = []key{{9, 10, 'z'}}
var K_Z = []key{{9, 10, 'Z'}}
var K_x = []key{{12, 10, 'x'}}
var K_X = []key{{12, 10, 'X'}}
var K_c = []key{{15, 10, 'c'}}
var K_C = []key{{15, 10, 'C'}}
var K_v = []key{{18, 10, 'v'}}
var K_V = []key{{18, 10, 'V'}}
var K_b = []key{{21, 10, 'b'}}
var K_B = []key{{21, 10, 'B'}}
var K_n = []key{{24, 10, 'n'}}
var K_N = []key{{24, 10, 'N'}}
var K_m = []key{{27, 10, 'm'}}
var K_M = []key{{27, 10, 'M'}}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	// termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// draw_keyboard()
	// termbox.Flush()
	// inputmode := 0
	// ctrlxpressed := false
	fmt.Println("starting to capture...")
	ev := termbox.PollEvent()
	dispatch_press(&ev)

}

func dispatch_press(ev *termbox.Event) {
	fmt.Println("this should print the key")
	fmt.Printf("%v", ev.Ch)
	// if ev.Ch < 128 {
	// 	if ev.Ch == 0 && ev.Key < 128 {
	// 		k = &combos[ev.Key]
	// 	} else {
	// 		k = &combos[ev.Ch]
	// 	}
	// }
	/*if ev.Mod&termbox.ModAlt != 0 {
		draw_key(K_LALT, termbox.ColorWhite, termbox.ColorRed)
		draw_key(K_RALT, termbox.ColorWhite, termbox.ColorRed)
	}

	var k *combo
	if ev.Key >= termbox.KeyArrowRight {
		k = &func_combos[0xFFFF-ev.Key]
	} else if ev.Ch < 128 {
		if ev.Ch == 0 && ev.Key < 128 {
			k = &combos[ev.Key]
		} else {
			k = &combos[ev.Ch]
		}
	}
	if k == nil {
		return
	}

	keys := k.keys
	for _, k := range keys {
		draw_key(k, termbox.ColorWhite, termbox.ColorRed)
	}*/
}
