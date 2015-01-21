package ui

import (
	"fmt"
	"strconv"

	"github.com/nsf/termbox-go"
)

const bufSize = 5

type InputType uint8

const (
	InputNumericInt InputType = iota
	InputNumericFloat
	InputSwitch
)

type Input struct {
	X      int
	Y      int
	LabelW int    // label width
	LabelT string // label text
	W      int    // width
	T      string // text
	a      bool   // attached draws with connected borders
	s      bool   // selected
	//buf    *buf
	*b
	Type  InputType
	Field ConfigField
}

type b struct {
	buf []rune
	cur *cur
}

type cur struct {
	i int // 0 < index < bufSize
	x int
	y int
}

func (in Input) newBuf() *b {
	buf := make([]rune, 0, bufSize)
	cur := &cur{i: 0, x: in.TextStartX(), y: in.TextY()}
	return &b{buf: buf, cur: cur}
}

func (in *Input) ClearBuf() {
	in.b = nil
	in.b = in.newBuf()
}

func (in *Input) SetBuf(e *Entry) {
	Debug(e.String())
	if e.Ch != 0 {
		in.bufAppend(e.Ch)
	}
	if e.Backspace {
		in.bufBackspace()
	}
	in.bufShow()
	flush()
}

func (in *Input) Switch() error {
	c := GetConfig()
	if c.Mode == 'A' {
		c.Mode = 'B'
	} else {
		c.Mode = 'A'
	}
	if err := c.Save(); err != nil {
		return err
	}
	UpdateConfig(c)
	UpdateConfig(c)
	ReloadInputs(c)
	return nil
}

func (in *Input) ValueMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	var val interface{}
	var err error
	if in.Type == InputNumericInt {
		if val, err = strconv.Atoi(string(in.buf)); err != nil {
			return nil, err
		}
	}
	if in.Type == InputNumericFloat {
		if val, err = strconv.ParseFloat(string(in.buf), 64); err != nil {
			return nil, err
		}
	}
	m[in.Field.Val()] = val
	return m, nil
}

func (in *Input) Valid() error {
	if in.Type == InputNumericInt {
		if _, err := strconv.Atoi(string(in.buf)); err != nil {
			return err
		}
	}
	if in.Type == InputNumericFloat {
		if _, err := strconv.ParseFloat(string(in.buf), 64); err != nil {
			return err
		}
	}
	return nil
}

func (in *Input) bufParseFloat() (float64, error) {
	return strconv.ParseFloat(string(in.buf), 64)
}

func (in *Input) bufParseInt() (int, error) {
	return strconv.Atoi(string(in.buf))
}

func (in *Input) bufAppend(r rune) {
	if len(in.buf) < cap(in.buf) {
		in.buf = append(in.buf, r)
		in.cur.x += 1
		in.cur.i += 1
		setCursor(in.cur.x, in.cur.y)
	}
}

func (in *Input) bufBackspace() {
	if len(in.buf) > 0 {
		in.buf = in.buf[0 : len(in.buf)-1]
		in.cur.x -= 1
		in.cur.i -= 1
		setCursor(in.cur.x, in.cur.y)
	}
}

func (in Input) bufShow() {
	in.SetText(string(in.buf))
}

func NewInput(x, y, labelW int, labelT string, w int, t string, a bool, it InputType, cf ConfigField) *Input {
	in := &Input{x, y, labelW, labelT, w, t, a, false, nil, it, cf}
	in.b = in.newBuf()
	return in
}

func (in Input) TextStartX() int {
	return in.X + in.LabelW + 3
}

func (in Input) TextEndX() int {
	return in.TextStartX() + in.W
}

func (in Input) TextY() int {
	return in.Y + 1
}

func (in Input) MaxX() int {
	return in.X + in.LabelW + 3 + in.W
}

func (in Input) MaxY() int {
	return in.Y + 2
}

func (in Input) ClearText() {
	fill(in.TextStartX(), in.TextY(), in.W, 1, ' ')
}

func (in Input) SetText(s string) {
	in.ClearText()
	text(in.TextStartX(), in.TextY(), s)
}

func (in Input) ResetText() {
	text(in.TextStartX(), in.TextY(), in.T)
	flush()
}

func (in Input) Draw() {
	x := in.X
	y := in.Y
	lw := in.LabelW
	lt := in.LabelT
	w := in.W
	t := in.T

	fill(x, y+0, 1, 1, '┌')
	fill(x, y+1, 1, 1, '│')
	fill(x, y+2, 1, 1, '└')
	fill(x+1, y+0, lw, 1, '─')
	text(x+1, y+1, lt)
	fill(x+1, y+2, lw, 1, '─')
	fill(x+lw+1, y+0, 1, 1, '┐')
	fill(x+lw+1, y+1, 1, 1, '│')
	fill(x+lw+1, y+2, 1, 1, '┘')
	fill(x+lw+2, y+0, 1, 1, ' ')
	fill(x+lw+2, y+1, 1, 1, ' ')
	fill(x+lw+2, y+2, 1, 1, ' ')
	fill(x+lw+3, y+0, w, 1, ' ')
	text(x+lw+3, y+1, t)
	fill(x+lw+3, y+2, w, 1, ' ')
	fill(x+lw+3+w, y+0, 1, 1, ' ')
	fill(x+lw+3+w, y+1, 1, 1, ' ')
	fill(x+lw+3+w, y+2, 1, 1, ' ')
	if in.a {
		fill(x, y+0, 1, 1, '├')
		fill(x+lw+1, y+0, 1, 1, '┤')
	}
}

func (in Input) Selected() bool {
	return in.s
}

func (in *Input) SetSelected(selected bool) {
	x := in.X
	y := in.Y
	lw := in.LabelW
	w := in.W

	if selected {
		DeselectAllInputs()
		in.s = true
		fill(x+lw+2, y+0, 1, 1, '┌')
		fill(x+lw+2, y+1, 1, 1, '│')
		fill(x+lw+2, y+2, 1, 1, '└')
		fill(x+lw+3, y+0, w, 1, '─')
		fill(x+lw+3, y+2, w, 1, '─')
		fill(x+lw+3+w, y+0, 1, 1, '┐')
		fill(x+lw+3+w, y+1, 1, 1, '│')
		fill(x+lw+3+w, y+2, 1, 1, '┘')
		in.ClearText()
		in.bufShow()
		setCursor(in.cur.x, in.cur.y)
		//termbox.SetCursor(in.TextStartX(), in.TextY())
	} else {
		in.s = false
		fill(x+lw+2, y+0, 1, 1, ' ')
		fill(x+lw+2, y+1, 1, 1, ' ')
		fill(x+lw+2, y+2, 1, 1, ' ')
		fill(x+lw+3, y+0, w, 1, ' ')
		fill(x+lw+3, y+2, w, 1, ' ')
		fill(x+lw+3+w, y+0, 1, 1, ' ')
		fill(x+lw+3+w, y+1, 1, 1, ' ')
		fill(x+lw+3+w, y+2, 1, 1, ' ')
		in.ResetText()
		in.ClearBuf()
		termbox.HideCursor()
	}

	termbox.Flush()
}

func SelectedInput() *Input {
	var si *Input
	for _, input := range inputs {
		if input.Selected() {
			si = input
		}
	}
	return si
}

func DeselectAllInputs() {
	//termbox.HideCursor()
	for _, in := range inputs {
		in.SetSelected(false)
		in.ResetText()
	}
}

type Entry struct {
	Ch        rune
	Backspace bool
	Delete    bool
	Enter     bool
}

func (e Entry) String() string {
	if e.Backspace {
		return fmt.Sprintf("Entry = 'Backspace'")
	}
	if e.Delete {
		return fmt.Sprintf("Entry = 'Delete'")
	}
	if e.Enter {
		return fmt.Sprintf("Entry = 'Enter'")
	}
	return fmt.Sprintf("Entry = %v", rtoa(e.Ch))

}

func NewEntry(te termbox.Event) *Entry {
	if te.Key == termbox.KeyEnter {
		return &Entry{Enter: true}
	}
	if te.Key == termbox.KeyDelete {
		return &Entry{Delete: true}
	}
	if te.Key == termbox.KeyBackspace2 {
		return &Entry{Backspace: true}
	}
	if AllowedEntry(te) {
		return &Entry{Ch: te.Ch}
	}
	return nil
}

func AllowedEntry(te termbox.Event) bool {
	if te.Key == termbox.KeyDelete || te.Key == termbox.KeyBackspace2 || te.Key == termbox.KeyEnter {
		return true
	}
	key := te.Ch
	if key == '0' || key == '1' || key == '2' || key == '3' || key == '4' ||
		key == '5' || key == '6' || key == '7' || key == '8' || key == '9' ||
		key == '.' {
		return true
	}
	return false
}
