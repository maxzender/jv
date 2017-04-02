package terminal

import "github.com/nsf/termbox-go"

type Event termbox.Event

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

func (t *Terminal) Draw(content []string) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	lineCount := len(content)
	for y := 0; y < t.Height && y < lineCount; y++ {
		lineLen := len(content[y])
		for x := 0; x < t.Width && x < lineLen; x++ {
			current := rune(content[y][x])
			termbox.SetCell(x, y, current, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.SetCursor(t.CursorX, t.CursorY)
	termbox.Flush()
}

func (t *Terminal) Poll() Event {
	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			return Event(e)
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
