package ui

import (
	"sync"

	"github.com/nsf/termbox-go"
)

// Labels is a map of each allowed key press (rune) to it's full string label.
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

var (
	mu        sync.Mutex
	inputs    []*Input
	keys      []*KeyLabel
	statusBar *StatusBar
	config    Config
)

const (
	inputLabelWidth      = 10
	inputWidth           = 10
	inputMinutesBufWidth = 3
	inputHzBufWidth      = 5
	keyLabelWidth        = 3
	keyWidth             = 21
	statusBarWidth       = 60
	statusBarDefaultText = "Press 'space' to start capturing keys, 'Esc' to quit."
)

// Init must be called before any other function. It initializes
// configuration and termbox.
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

// Close should be deferred after initialization. It finalizes termbox library.
func Close() {
	mu.Lock()
	termbox.Close()
	mu.Unlock()
}

// DrawAll draws the title, the inputs, the key labels and the status bar.
// It should also be called when a resize event is received.
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
	in1 := NewInput(x, y+0, lw, "Mode", w, 0, config.ModeS(), false, InputSwitch, configMode)
	in2 := NewInput(x, y+2, lw, "TotalTime", w, inputMinutesBufWidth, config.TotalTimeS(), true, InputNumericInt, configTotalTime)
	in3 := NewInput(x, y+4, lw, "Offset", w, inputMinutesBufWidth, config.OffsetS(), true, InputNumericInt, configOffset)
	in4 := NewInput(x, y+6, lw, "BaseHz", w, inputHzBufWidth, config.BaseHzS(), true, InputNumericFloat, configBaseHz)
	in5 := NewInput(x, y+8, lw, "StartHz", w, inputHzBufWidth, config.StartHzS(), true, InputNumericFloat, configStartHz)
	in6 := NewInput(x, y+10, lw, "EndHz", w, inputHzBufWidth, config.EndHzS(), true, InputNumericFloat, configEndHz)
	inputs = nil
	inputs = append(inputs, in1, in2, in3, in4, in5, in6)
	for _, in := range inputs {
		in.Draw()
	}
	return in6.MaxX(), in6.MaxY()
}

// ReloadInputs updates each input with the values of a new configuration.
// It should be called after receiving a valid value from en enabled input.
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

// Cell wraps a termbox cell. A single conceptual entity on the screen. A
// screen considered a 2d array of cells. Each cell also holds a reference
// to an input if it exists, thus by clicking a cell that contains an input
// we can easily access it.
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
			j++
			i = 0
		}
		cells[i][j].Ch = c.Ch
		cells[i][j].Fg = c.Fg
		cells[i][j].Bg = c.Bg
		i++
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

// Cells initializes all cells on screen and registers the inputs on the cells
// that should have them.
func Cells() [][]Cell {
	c := cells()
	return registerInputs(c)
}

// GetCell returns a cell based on the x, y points of the screen.
func GetCell(x, y int) Cell {
	c := Cells()
	return c[x][y]
}

// Scan scans the cells of the screen that are contained at a certain y,
// between startX and endX. It returns the content of the scanned cells as
// a string.
func Scan(startX, endX, y int) string {
	c := Cells()
	var runes []rune
	for x := startX; x <= endX; x++ {
		runes = append(runes, c[x][y].Ch)
	}
	return string(runes)
}

// KeyLabel holds a label with each allowed key press and right next, the
// coresponding text that describes the label.
type KeyLabel struct {
	X      int
	Y      int
	LabelW int    // label width
	LabelT string // label text
	W      int    // rest of width
	T      string // text describing the key
	a      bool
	S      bool
}

// NewKeyLabel creates a new KeyLabel.
func NewKeyLabel(x, y, labelW int, labelT string, w int, t string, a bool) *KeyLabel {
	return &KeyLabel{x, y, labelW, labelT, w, t, a, false}
}

// MaxX returns the maximum x that a KeyLabel reaches on the screen.
func (kl KeyLabel) MaxX() int {
	return kl.X + kl.LabelW + 3 + kl.W
}

// MaxY returns the maximum y that a KeyLabel reaches on the screen.
func (kl KeyLabel) MaxY() int {
	return kl.Y + 2
}

// Draw draws the KeyLabel.
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

// StatusBar holds the data for drawing a status bar with a specified width
// and text. It also has space to the left for the timer.
type StatusBar struct {
	X          int
	Y          int
	Width      int
	Text       string
	timerWidth int
}

// NewStatusBar returns a new StatusBar.
func NewStatusBar(x, y, width int, text string) *StatusBar {
	return &StatusBar{x, y, width, text, 6}
}

// MaxX returns the maximum x that a StatusBar can reach on the screen.
func (sb StatusBar) MaxX() int {
	return sb.X + sb.timerWidth + 2 + sb.Width
}

// MaxY returns the maximum y that a StatusBar can reach on the screen.
func (sb StatusBar) MaxY() int {
	return sb.Y + 2
}

// Draw draws the StatusBar.
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

	flush()
}

// UpdateTimer updates the timer of the status bar by a specified amount of
// seconds.
func (sb StatusBar) UpdateTimer(seconds int) {
	text(sb.X+1, sb.Y+1, FormatTimer(seconds))
	flush()
}

// UpdateText updates the text of the status bar.
func (sb StatusBar) UpdateText(t string) {
	//Text(sb.X+sb.timerWidth+2, sb.Y+1, text[0:sb.Width]) // text[0:sb.Width] panics for some reason
	fill(sb.X+sb.timerWidth+2, sb.Y+1, sb.Width, 1, ' ')
	text(sb.X+sb.timerWidth+2, sb.Y+1, t)
	flush()
}

// UpdateTimer is a helper function that updates the status bar's timer.
func UpdateTimer(seconds int) {
	if statusBar != nil {
		statusBar.UpdateTimer(seconds)
	}
}

// ResetTimer clears the status bar's timer.
func ResetTimer() {
	if statusBar != nil {
		fill(statusBar.X+1, statusBar.Y+1, statusBar.timerWidth, 1, ' ')
		flush()
	}
}

// UpdateText is a helper function that updates the status bar's text.
func UpdateText(text string) {
	if statusBar != nil {
		statusBar.UpdateText(text)
	}
}

// ResetText clears the status bar's text.
func ResetText() {
	if statusBar != nil {
		statusBar.UpdateText(statusBarDefaultText)
	}
}

// Capture represents a captured key press at a specific second of the timer
// along with the value of Hz that was recorded.
type Capture struct {
	Value   rune
	Seconds int
	Hz      float64
}

// Label returns the description of the captured key value.
func (c *Capture) Label() string {
	return Labels[c.Value]
}

// Timestamp returns the timestamp of when a capture happened.
func (c *Capture) Timestamp() string {
	return FormatTimer(c.Seconds)
}
