package terminal

import (
	"github.com/maxzender/jv/jsontree"
	"github.com/nsf/termbox-go"
)

type Terminal struct {
	Width, Height    int
	CursorX, CursorY int
}

func New() (*Terminal, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	w, h := termbox.Size()
	return &Terminal{Width: w, Height: h}, nil
}

func (t *Terminal) MoveCursor(x, y int) {
	t.CursorX, t.CursorY = t.CursorX+x, t.CursorY+y
	t.EnsureCursorWithinWindow()
}

func (t *Terminal) Resize(width, height int) {
	t.Width = width
	t.Height = height
	t.EnsureCursorWithinWindow()
}

func (t *Terminal) EnsureCursorWithinWindow() {
	t.CursorX = min(t.Width-1, max(0, t.CursorX))
	t.CursorY = min(t.Height-1, max(0, t.CursorY))
}

func (t *Terminal) Render(tree *jsontree.JsonTree) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	for y := 0; y < t.Height; y++ {
		if line := tree.Line(y); line != nil {
			lineLen := len(line)
			for x := 0; x < t.Width && x < lineLen; x++ {
				c := line[x]
				termbox.SetCell(x, y, c.Val, c.Color, termbox.ColorDefault)
			}
		}
	}

	termbox.SetCursor(t.CursorX, t.CursorY)
	termbox.Flush()
}

func (t *Terminal) Poll() termbox.Event {
	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			return e
		case termbox.EventResize:
			t.Resize(e.Width, e.Height)
		}
	}
}

func (t *Terminal) Close() {
	termbox.Close()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
