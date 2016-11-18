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
	newX := w.CursorX + x
	if newX < 0 {
		newX = 0
	} else if newX >= w.Width {
		newX = w.CursorX
	}
	newY := w.CursorY + y
	if newY < 0 {
		newY = 0
	} else if newY >= w.Height {
		newY = w.CursorY
	}

	w.CursorX, w.CursorY = newX, newY
}

func (w *Window) Resize(width, height int) {
	w.Width = width
	w.Height = height

	if w.CursorX > width {
		w.CursorX = width
	}
	if w.CursorY > height {
		w.CursorY = height
	}
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
