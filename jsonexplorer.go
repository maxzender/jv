package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

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
	expandedLines = make(map[int]struct{})
	segments      = make(map[int]int)
)

// finds pairs of line numbers that describe the section between two matching brackets
func parseSegments(lines []string) {
	bracketBalances := make(map[int]int)
	var bal int
	for num, line := range lines {
		for _, c := range line {
			switch c {
			case '{', '[':
				bal += 1
				bracketBalances[bal] = num
			case '}', ']':
				segments[bracketBalances[bal]] = num
				bal -= 1
			}
		}
	}
}

func print(content string, w *Window) {
	lines := strings.Split(content, "\n")
	lineCount := len(lines)
	parseSegments(lines)
	skipTill := 0
	lineCursor := 0

	for y := 0; y < w.Height && y < lineCount; y++ {
		if y >= skipTill {
			_, expanded := expandedLines[y]
			if !expanded {
				skipTill = segments[y]
			}

			lineLen := len(lines[y])
			for x := 0; x < w.Width && x < lineLen; x++ {
				current := rune(lines[y][x])
				termbox.SetCell(x, lineCursor, current, termbox.ColorWhite, termbox.ColorDefault)
			}
			lineCursor += 1
		}
	}
}

func toggleLine(num int) {
	_, expanded := expandedLines[num]
	if expanded {
		delete(expandedLines, num)
	} else {
		expandedLines[num] = struct{}{}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	window = NewWindow(termbox.Size())

	expandedLines[0] = struct{}{}
	content, err := ioutil.ReadAll(os.Stdin)
	print(string(content), window)

	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
				return
			} else if e.Key == termbox.KeyEnter {
				toggleLine(window.CursorY)
			} else {
				handleKeyPress(e)
			}
		case termbox.EventResize:
			window.Resize(e.Width, e.Height)
		}
		termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)
		print(string(content), window)
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
