package main

import "github.com/nsf/termbox-go"

type Window struct {
	Width, Height    int
	CursorX, CursorY int
}

func NewWindow(w, h int) *Window {
	return &Window{Width: w, Height: h}
}

func (w *Window) MoveCursor(x, y int) {
	w.CursorX, w.CursorY = w.CursorX+x, w.CursorY+y
	w.EnsureCursorWithinWindow()
}

func (w *Window) Resize(width, height int) {
	w.Width = width
	w.Height = height
	w.EnsureCursorWithinWindow()
}

func (w *Window) EnsureCursorWithinWindow() {
	w.CursorX = min(w.Width-1, max(0, w.CursorX))
	w.CursorY = min(w.Height-1, max(0, w.CursorY))
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

var (
	window *Window
	keyMap = map[rune]func(){
		'h': func() { window.MoveCursor(-1, 0) },
		'j': func() { window.MoveCursor(0, +1) },
		'k': func() { window.MoveCursor(0, -1) },
		'l': func() { window.MoveCursor(+1, 0) },
	}
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	window = NewWindow(termbox.Size())

	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
				return
			} else {
				handleKeyPress(e)
			}
		case termbox.EventResize:
			window.Resize(e.Width, e.Height)
		}
		termbox.SetCursor(window.CursorX, window.CursorY)
		termbox.Flush()
	}
}

func handleKeyPress(event termbox.Event) {
	handler, ok := keyMap[event.Ch]
	if ok {
		handler()
	}
}
