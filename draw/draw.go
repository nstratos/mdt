package draw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

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

func Title() (maxX, maxY int) {
	return Text(0, 0, title)
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
	return formatTimer(c.Seconds)
}

func formatTimer(seconds int) string {
	min := seconds / 60
	sec := seconds % 60
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

func cells() [][]termbox.Cell {
	mx, my := termbox.Size()
	cells := termbox.CellBuffer()
	cells2d := make([][]termbox.Cell, mx)
	for k := range cells2d {
		cells2d[k] = make([]termbox.Cell, my)
	}
	i, j := 0, 0
	for _, c := range cells {
		if i == mx {
			j += 1
			i = 0
		}
		cells2d[i][j] = c
		i += 1
	}
	return cells2d
}

func GetCell(x, y int) termbox.Cell {
	c := cells()
	return c[x][y]
}

func GetInput(startX, endX, y int) string {
	c := cells()
	//l := endY - startY
	runes := make([]rune, 0)
	//i := 0
	for x := startX; x <= endX; x++ {
		//runes[i] = c[x][y].Ch
		//i += 1
		runes = append(runes, c[x][y].Ch)
	}
	return string(runes)
}

// func Box() {
// 	const coldef = termbox.ColorDefault
// 	//termbox.Clear(coldef, coldef)
// 	_, h := termbox.Size()

// 	midy := h - 2
// 	//midx := (w - edit_box_width) / 2
// 	midx := 1

// 	// unicode box drawing chars around the edit box
// 	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
// 	termbox.SetCell(midx+edit_box_width, midy, '│', coldef, coldef)
// 	termbox.SetCell(midx-1, midy-1, '╭', coldef, coldef)
// 	termbox.SetCell(midx-1, midy+1, '╰', coldef, coldef)
// 	termbox.SetCell(midx+edit_box_width, midy-1, '╮', coldef, coldef)
// 	termbox.SetCell(midx+edit_box_width, midy+1, '╯', coldef, coldef)
// 	fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
// 	fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})

// 	//edit_box.Draw(midx, midy, edit_box_width, 1)
// 	//termbox.SetCursor(midx+edit_box.CursorX(), midy)
// 	termbox.SetCursor(midx, midy)

// 	tbprint(midx+6, midy+3, coldef, coldef, "Press ESC to quit")
// 	//termbox.Flush()
// }

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

func (sb StatusBar) Draw() {
	const coldef = termbox.ColorDefault
	x := sb.X
	y := sb.Y
	width := sb.Width
	text := sb.Text
	timerWidth := sb.timerWidth

	// unicode box drawing chars around the edit box
	termbox.SetCell(x, y+0, '╔', coldef, coldef)
	termbox.SetCell(x, y+1, '║', coldef, coldef)
	termbox.SetCell(x, y+2, '╚', coldef, coldef)
	fill(x+1, y+0, timerWidth, 1, termbox.Cell{Ch: '═'})
	fill(x+1, y+2, timerWidth, 1, termbox.Cell{Ch: '═'})
	termbox.SetCell(x+timerWidth+1, y+0, '╤', coldef, coldef)
	termbox.SetCell(x+timerWidth+1, y+1, '│', coldef, coldef)
	termbox.SetCell(x+timerWidth+1, y+2, '╧', coldef, coldef)
	fill(x+timerWidth+2, y+0, width, 1, termbox.Cell{Ch: '═'})
	fill(x+timerWidth+2, y+2, width, 1, termbox.Cell{Ch: '═'})
	termbox.SetCell(x+timerWidth+2+width, y+0, '╗', coldef, coldef)
	termbox.SetCell(x+timerWidth+2+width, y+1, '║', coldef, coldef)
	termbox.SetCell(x+timerWidth+2+width, y+2, '╝', coldef, coldef)
	Text(x+timerWidth+2, y+1, text)

	termbox.Flush()
}

func (sb StatusBar) UpdateTimer(seconds int) {
	Text(sb.X+1, sb.Y+1, formatTimer(seconds))
	termbox.Flush()
}

func (sb StatusBar) UpdateText(text string) {
	//Text(sb.X+sb.timerWidth+2, sb.Y+1, text[0:sb.Width]) // panics for some reason
	fill(sb.X+sb.timerWidth+2, sb.Y+1, sb.Width, 1, termbox.Cell{Ch: ' '})
	Text(sb.X+sb.timerWidth+2, sb.Y+1, text)
	termbox.Flush()
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func Text(x, y int, s string) (maxX, maxY int) {
	mx := 0
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
		if x > mx {
			mx = x
		}
		x = tempx
		y += 1
	}
	// Because we always icrease it one more.
	y -= 1
	return mx, y
}

func DrawOptions(x, y int, c Config) (finalX, finalY int) {

	_, y = Text(x, y, "───────────────")
	_, y = Text(x, y+1, "Mode: "+strconv.QuoteRuneToASCII(c.Mode))
	_, y = Text(x, y+1, "───────────────")
	_, y = Text(x, y+1, "TotalTime: "+strconv.Itoa(c.TotalTime)+" min")
	_, y = Text(x, y+1, "───────────────")
	_, y = Text(x, y+1, "Offset: "+strconv.Itoa(c.Offset))
	_, y = Text(x, y+1, "───────────────")
	_, y = Text(x, y+1, fmt.Sprintf("Base: %.2f hz", c.BaseHz))
	_, y = Text(x, y+1, "───────────────")
	_, y = Text(x, y+1, fmt.Sprintf("Start: %.2f hz", c.StartHz))
	_, y = Text(x, y+1, "───────────────")
	_, y = Text(x, y+1, fmt.Sprintf("End: %.2f hz", c.EndHz))
	_, y = Text(x, y+1, "───────────────")
	return x, y
}

func DrawKeyLabels(x, y int) (finalX, finalY int) {
	_, y = Text(x, y, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('q'), Labels['q']))
	_, y = Text(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('a'), Labels['a']))
	_, y = Text(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('w'), Labels['w']))
	_, y = Text(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('s'), Labels['s']))
	_, y = Text(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('e'), Labels['e']))
	_, y = Text(x, y+1, fmt.Sprintf("%v = %v", strconv.QuoteRuneToASCII('d'), Labels['d']))
	return x, y
}

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

// CurrentHz = Seconds * (EndHz - StartHz) / ((TotalTime - Offset) * 60)
func CurrentHz(seconds int, c Config) float64 {
	return (float64(seconds) * (c.EndHz - c.StartHz) / float64((c.TotalTime-c.Offset)*60)) + c.StartHz
}
