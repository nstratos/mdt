package ui

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nsf/termbox-go"
)

// InputType is the type of each input which can be either an input that
// accepts integers, an input that accepts float or an input that switches
// value on click.
type InputType uint8

const (
	// InputNumericInt is an input that accepts integers.
	InputNumericInt InputType = iota
	// InputNumericFloat is an input that accepts floats.
	InputNumericFloat
	// InputSwitch is an input that switches value on click.
	InputSwitch
)

// Input represents an input box on the screen.
type Input struct {
	X      int
	Y      int
	LabelW int    // label width
	LabelT string // label text
	W      int    // width
	bufW   int    // buffer width
	T      string // text
	a      bool   // attached draws with connected borders
	s      bool   // selected
	*b            // buffer
	Type   InputType
	Field  ConfigField
}

// buffer containing the runes of each cell and a cursor.
type b struct {
	buf []rune
	cur *cur
}

// cursor of an input.
type cur struct {
	i int
	x int
	y int
}

func (in Input) newBuf() *b {
	buf := make([]rune, 0, in.bufW)
	cur := &cur{i: 0, x: in.TextStartX(), y: in.TextY()}
	return &b{buf: buf, cur: cur}
}

// ClearBuf clears the inputs buffer.
func (in *Input) ClearBuf() {
	in.b = nil
	in.b = in.newBuf()
}

// SetBuf sets the input's buffer.
func (in *Input) SetBuf(e *Entry) {
	//Debug(e.String())
	if e.Ch != 0 {
		in.bufAppend(e.Ch)
	}
	if e.Backspace {
		in.bufBackspace()
	}
	in.bufShow()
	flush()
}

// Switch switches the input between "Binaural" and "Isochronic" values.
func (in *Input) Switch() error {
	c := GetConfig()
	if c.Mode == "Binaural" {
		c.Mode = "Isochronic"
	} else {
		c.Mode = "Binaural"
	}
	if err := c.Save(); err != nil {
		return err
	}
	UpdateConfig(c)
	UpdateConfig(c)
	ReloadInputs(c)
	return nil
}

// ValueMap returns a map containing the value of the input. Each input
// corresponds to one configuration field. The key of the map is the
// configuration field's key. It is used to easilly update the JSON
// config file.
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

// Valid checks the current value of an input and returns an error if it's
// not valid.
func (in *Input) Valid() error {
	if in.Type == InputNumericInt {
		if _, err := strconv.Atoi(string(in.buf)); err != nil {
			return errors.New("Expecting number of minutes e.g. 60")
		}
	}
	if in.Type == InputNumericFloat {
		if _, err := strconv.ParseFloat(string(in.buf), 64); err != nil {
			return errors.New("Expecting decimal e.g. 50.65")
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
		in.cur.x++
		in.cur.i++
		setCursor(in.cur.x, in.cur.y)
	}
}

func (in *Input) bufBackspace() {
	if len(in.buf) > 0 {
		in.buf = in.buf[0 : len(in.buf)-1]
		in.cur.x--
		in.cur.i--
		setCursor(in.cur.x, in.cur.y)
	}
}

func (in Input) bufShow() {
	in.SetText(string(in.buf))
}

// NewInput returns a new Input.
func NewInput(x, y, labelW int, labelT string, w, bufW int, t string, a bool, it InputType, cf ConfigField) *Input {
	in := &Input{x, y, labelW, labelT, w, bufW, t, a, false, nil, it, cf}
	in.b = in.newBuf()
	return in
}

// TextStartX returns the x where the input's text starts.
func (in Input) TextStartX() int {
	return in.X + in.LabelW + 3
}

// TextEndX returns the x where the input's text ends.
func (in Input) TextEndX() int {
	return in.TextStartX() + in.W
}

// TextY returns the y of the input's text.
func (in Input) TextY() int {
	return in.Y + 1
}

// MaxX returns the maximum x that the input reaches on the screen.
func (in Input) MaxX() int {
	return in.X + in.LabelW + 3 + in.W
}

// MaxY returns the maximum y that the input reaches on the screen.
func (in Input) MaxY() int {
	return in.Y + 2
}

// ClearText clears the inputs text.
func (in Input) ClearText() {
	fill(in.TextStartX(), in.TextY(), in.W, 1, ' ')
}

// SetText sets the input's text.
func (in Input) SetText(s string) {
	in.ClearText()
	text(in.TextStartX(), in.TextY(), s)
	flush()
}

// ResetText clears the input's text and then resets it to it's original value.
func (in Input) ResetText() {
	in.ClearText()
	text(in.TextStartX(), in.TextY(), in.T)
	flush()
}

// Draw draws the input on the screen.
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

// Selected returns true of the input is selected.
func (in Input) Selected() bool {
	return in.s
}

// SetSelected sets the input as selected after deselecting all other inputs.
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
		if in.Type == InputNumericInt {
			UpdateText(fmt.Sprintf("Enter minutes (Previous value: %s)", in.T))
		} else if in.Type == InputNumericFloat {
			UpdateText(fmt.Sprintf("Enter hz (Previous value: %s)", in.T))
		}
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

	flush()
}

// SelectedInput returns the input that is currently selected.
func SelectedInput() *Input {
	var si *Input
	for _, input := range inputs {
		if input.Selected() {
			si = input
		}
	}
	return si
}

// DeselectAllInputs sets all inputs selected as false.
func DeselectAllInputs() {
	for _, in := range inputs {
		in.SetSelected(false)
		in.ResetText()
	}
}

// Entry represents a key entry that an input accpets. These include a
// character, backspace, delete and enter.
type Entry struct {
	Ch        rune
	Backspace bool
	Delete    bool
	Enter     bool
}

// String returns a string representation of the entry.
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

// NewEntry returns a new Entry based on a termbox event received.
func NewEntry(te termbox.Event) *Entry {
	if te.Key == termbox.KeyEnter {
		return &Entry{Enter: true}
	}
	if te.Key == termbox.KeyDelete {
		return &Entry{Delete: true}
	}
	if te.Key == termbox.KeyBackspace {
		return &Entry{Backspace: true}
	}
	if te.Key == termbox.KeyBackspace2 {
		return &Entry{Backspace: true}
	}
	if AllowedEntry(te) {
		return &Entry{Ch: te.Ch}
	}
	return nil
}

// AllowedEntry returns true if the termbox event received is a valid entry
// for the input. These include 0-9, '.', delete, backspaces and enter.
func AllowedEntry(te termbox.Event) bool {
	if te.Key == termbox.KeyDelete || te.Key == termbox.KeyBackspace || te.Key == termbox.KeyBackspace2 || te.Key == termbox.KeyEnter {
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
