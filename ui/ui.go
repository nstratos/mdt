package ui

import "github.com/nsf/termbox-go"

var Labels = map[rune]string{
	'a': "Visual imagination",
	'd': "Language thought",
	'e': "Language voice",
	'q': "Visual memory",
	's': "Auditory imagination",
	'w': "Auditory memory",
}

const title = `           _ _   
  _ __  __| | |_ 
 | '  \/ _  |  _|
 |_|_|_\__,_|\__|
                 
`

//var cells [][]Cell
var inputs []*Input
var keys []*KeyLabel
var statusBar *StatusBar
var config Config

const inputLabelWidth = 10
const inputWidth = 9
const inputMinutesBufWidth = 3
const inputHzBufWidth = 5
const keyLabelWidth = 3
const keyWidth = 21
const statusBarWidth = 60
const statusBarDefaultText = "Press 'space' to start capturing keys, 'Esc' to quit."

// Must be called before any other function.
func Init() error {
	if err := initConfig(); err != nil {
		return err
	}
	if err := initTermbox(); err != nil {
		return err
	}
	return nil
}

// Loading configuration from config.json
func initConfig() error {
	config = Config{}
	err := config.Load()
	if err != nil {
		return err
	}
	return nil
}

// Initializing termbox
func initTermbox() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
	return nil
}

func Close() {
	termbox.Close()
}

func DrawAll() {
	inputs = nil
	keys = nil
	statusBar = nil
	_, keysY := drawTitle(0, 0)
	keysX, sbY := drawInputs(0, keysY+1)
	_, _ = drawKeys(keysX+3, keysY+1)
	_, _ = drawStatusBar(0, sbY+1)
}

func drawTitle(x, y int) (maxX, maxY int) {
	return text(x, y, title)
}

func drawStatusBar(x, y int) (maxX, maxY int) {
	sbw := statusBarWidth
	statusBar = NewStatusBar(x, y, sbw, statusBarDefaultText)
	statusBar.Draw()
	return statusBar.MaxX(), statusBar.MaxY()
}

func drawInputs(x, y int) (maxX, maxY int) {
	const lw = inputLabelWidth
	const w = inputWidth
	in1 := NewInput(x, y+0, lw, "Mode", w, 0, config.ModeS(), false, InputSwitch, ConfigMode)
	in2 := NewInput(x, y+2, lw, "TotalTime", w, inputMinutesBufWidth, config.TotalTimeS(), true, InputNumericInt, ConfigTotalTime)
	in3 := NewInput(x, y+4, lw, "Offset", w, inputMinutesBufWidth, config.OffsetS(), true, InputNumericInt, ConfigOffset)
	in4 := NewInput(x, y+6, lw, "BaseHz", w, inputHzBufWidth, config.BaseHzS(), true, InputNumericFloat, ConfigBaseHz)
	in5 := NewInput(x, y+8, lw, "StartHz", w, inputHzBufWidth, config.StartHzS(), true, InputNumericFloat, ConfigStartHz)
	in6 := NewInput(x, y+10, lw, "EndHz", w, inputHzBufWidth, config.EndHzS(), true, InputNumericFloat, ConfigEndHz)
	inputs = nil
	inputs = append(inputs, in1, in2, in3, in4, in5, in6)
	for _, in := range inputs {
		in.Draw()
	}
	return in6.MaxX(), in6.MaxY()
}

func ReloadInputs(c Config) {
	inputs[0].T = c.ModeS()
	inputs[1].T = c.TotalTimeS()
	inputs[2].T = c.OffsetS()
	inputs[3].T = c.BaseHzS()
	inputs[4].T = c.StartHzS()
	inputs[5].T = c.EndHzS()
	for _, in := range inputs {
		in.ClearBuf()
		in.ResetText()
	}
}

func drawKeys(x, y int) (maxX, maxY int) {
	const lw = keyLabelWidth
	const w = keyWidth
	k1 := NewKeyLabel(x, y+0, lw, rtoa('q'), w, Labels['q'], false)
	k2 := NewKeyLabel(x, y+2, lw, rtoa('a'), w, Labels['a'], true)
	k3 := NewKeyLabel(x, y+4, lw, rtoa('w'), w, Labels['w'], true)
	k4 := NewKeyLabel(x, y+6, lw, rtoa('s'), w, Labels['s'], true)
	k5 := NewKeyLabel(x, y+8, lw, rtoa('e'), w, Labels['e'], true)
	k6 := NewKeyLabel(x, y+10, lw, rtoa('d'), w, Labels['d'], true)
	keys = append(keys, k1, k2, k3, k4, k5, k6)
	for _, k := range keys {
		k.Draw()
	}
	return k6.MaxX(), k6.MaxY()
}

type Cell struct {
	Input *Input
	termbox.Cell
}

func cells() [][]Cell {
	mx, my := termbox.Size()
	cellBuffer := termbox.CellBuffer()
	cells := make([][]Cell, mx)
	for k := range cells {
		cells[k] = make([]Cell, my)
	}
	i, j := 0, 0
	for _, c := range cellBuffer {
		if i == mx {
			j += 1
			i = 0
		}
		cells[i][j].Ch = c.Ch
		cells[i][j].Fg = c.Fg
		cells[i][j].Bg = c.Bg
		i += 1
	}
	return cells
}

func registerInputs(cells [][]Cell) [][]Cell {
	for _, in := range inputs {
		for x := in.TextStartX(); x <= in.TextEndX(); x++ {
			cells[x][in.TextY()-1].Input = in
			cells[x][in.TextY()+0].Input = in
			cells[x][in.TextY()+1].Input = in
		}
	}
	return cells
}

func Cells() [][]Cell {
	c := cells()
	return registerInputs(c)
}

func GetCell(x, y int) Cell {
	c := Cells()
	return c[x][y]
}

// Scans the cells of a line and returns it's contents as a string.
func Scan(startX, endX, y int) string {
	c := Cells()
	runes := make([]rune, 0)
	for x := startX; x <= endX; x++ {
		runes = append(runes, c[x][y].Ch)
	}
	return string(runes)
}

type KeyLabel struct {
	X      int
	Y      int
	LabelW int
	LabelT string
	W      int
	T      string
	a      bool
	S      bool
}

func NewKeyLabel(x, y, labelW int, labelT string, w int, t string, a bool) *KeyLabel {
	return &KeyLabel{x, y, labelW, labelT, w, t, a, false}
}

func (kl KeyLabel) Start() int {
	return kl.X + kl.LabelW + 3
}

func (kl KeyLabel) End() int {
	return kl.Start() + kl.W
}
func (kl KeyLabel) MaxX() int {
	return kl.X + kl.LabelW + 3 + kl.W
}

func (kl KeyLabel) MaxY() int {
	return kl.Y + 2
}

func (kl KeyLabel) Draw() {
	x := kl.X
	y := kl.Y
	lw := kl.LabelW
	lt := kl.LabelT
	w := kl.W
	t := kl.T

	fill(x, y+0, 1, 1, '┌')
	fill(x, y+1, 1, 1, '│')
	fill(x, y+2, 1, 1, '└')
	fill(x+1, y+0, lw, 1, '─')
	text(x+1, y+1, lt)
	fill(x+1, y+2, lw, 1, '─')
	fill(x+lw+1, y+0, 1, 1, '─')
	fill(x+lw+1, y+1, 1, 1, ' ')
	fill(x+lw+1, y+2, 1, 1, '─')
	fill(x+lw+2, y+0, w, 1, '─')
	text(x+lw+2, y+1, t)
	fill(x+lw+2, y+2, w, 1, '─')
	fill(x+lw+2+w, y+0, 1, 1, '┐')
	fill(x+lw+2+w, y+1, 1, 1, '│')
	fill(x+lw+2+w, y+2, 1, 1, '┘')
	if kl.a {
		fill(x, y+0, 1, 1, '├')
		fill(x+lw+2+w, y+0, 1, 1, '┤')
	}
}

type StatusBar struct {
	X          int
	Y          int
	Width      int
	Text       string
	timerWidth int
}

func NewStatusBar(x, y, width int, text string) *StatusBar {
	return &StatusBar{x, y, width, text, 6}
}

func (sb StatusBar) MaxX() int {
	return sb.X + sb.timerWidth + 2 + sb.Width
}

func (sb StatusBar) MaxY() int {
	return sb.Y + 2
}

func (sb StatusBar) Draw() {
	x := sb.X
	y := sb.Y
	w := sb.Width
	t := sb.Text
	tw := sb.timerWidth

	// unicode box drawing chars around the edit box
	fill(x, y+0, 1, 1, '╔')
	fill(x, y+1, 1, 1, '║')
	fill(x, y+2, 1, 1, '╚')
	fill(x+1, y+0, tw, 1, '═')
	fill(x+1, y+2, tw, 1, '═')
	fill(x+tw+1, y+0, 1, 1, '╤')
	fill(x+tw+1, y+1, 1, 1, '│')
	fill(x+tw+1, y+2, 1, 1, '╧')
	fill(x+tw+2, y+0, w, 1, '═')
	fill(x+tw+2, y+2, w, 1, '═')
	fill(x+tw+2+w, y+0, 1, 1, '╗')
	fill(x+tw+2+w, y+1, 1, 1, '║')
	fill(x+tw+2+w, y+2, 1, 1, '╝')
	text(x+tw+2, y+1, t)

	termbox.Flush()
}

func (sb StatusBar) UpdateTimer(seconds int) {
	text(sb.X+1, sb.Y+1, FormatTimer(seconds))
	termbox.Flush()
}

func (sb StatusBar) UpdateText(t string) {
	//Text(sb.X+sb.timerWidth+2, sb.Y+1, text[0:sb.Width]) // text[0:sb.Width] panics for some reason
	fill(sb.X+sb.timerWidth+2, sb.Y+1, sb.Width, 1, ' ')
	text(sb.X+sb.timerWidth+2, sb.Y+1, t)
	termbox.Flush()
}

func UpdateTimer(seconds int) {
	if statusBar != nil {
		statusBar.UpdateTimer(seconds)
	}
}

func ResetTimer() {
	if statusBar != nil {
		fill(statusBar.X+1, statusBar.Y+1, statusBar.timerWidth, 1, ' ')
		flush()
	}
}

func UpdateText(text string) {
	if statusBar != nil {
		statusBar.UpdateText(text)
	}
}

func ResetText() {
	if statusBar != nil {
		statusBar.UpdateText(statusBarDefaultText)
	}
}

type Capture struct {
	Value   rune
	Seconds int
	Hz      float64
}

func (c *Capture) Label() string {
	return Labels[c.Value]
}

func (c *Capture) Timestamp() string {
	return FormatTimer(c.Seconds)
}
