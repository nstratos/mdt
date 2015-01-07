package draw

import "github.com/nsf/termbox-go"

const edit_box_width = 50

func Box() {
	const coldef = termbox.ColorDefault
	//termbox.Clear(coldef, coldef)
	_, h := termbox.Size()

	midy := h - 2
	//midx := (w - edit_box_width) / 2
	midx := 1

	// unicode box drawing chars around the edit box
	termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy, '│', coldef, coldef)
	termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy-1, '┐', coldef, coldef)
	termbox.SetCell(midx+edit_box_width, midy+1, '┘', coldef, coldef)
	fill(midx, midy-1, edit_box_width, 1, termbox.Cell{Ch: '─'})
	fill(midx, midy+1, edit_box_width, 1, termbox.Cell{Ch: '─'})

	//edit_box.Draw(midx, midy, edit_box_width, 1)
	//termbox.SetCursor(midx+edit_box.CursorX(), midy)
	termbox.SetCursor(midx, midy)

	tbprint(midx+6, midy+3, coldef, coldef, "Press ESC to quit")
	//termbox.Flush()
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
