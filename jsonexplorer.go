package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/maxzender/jsonexplorer/contentview"
	"github.com/maxzender/jsonexplorer/terminal"
	termbox "github.com/nsf/termbox-go"
)

var (
	keyMap = map[rune]func(*terminal.Terminal){
		'h': func(t *terminal.Terminal) { t.MoveCursor(-1, 0) },
		'j': func(t *terminal.Terminal) { t.MoveCursor(0, +1) },
		'k': func(t *terminal.Terminal) { t.MoveCursor(0, -1) },
		'l': func(t *terminal.Terminal) { t.MoveCursor(+1, 0) },
	}
	specialKeyMap = map[termbox.Key]func(*terminal.Terminal){
		termbox.KeyEnter: func(t *terminal.Terminal) { toggleLine(t) },
	}
	view *contentview.ContentView
)

func main() {
	term, err := terminal.New()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	byteContent, err := ioutil.ReadAll(os.Stdin)
	content := strings.Split(string(byteContent), "\n")
	view = contentview.New(content)

	for {
		term.Draw(view.Content())
		e := term.Poll()
		if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
			return
		}
		handleKeypress(term, e)
	}
}

func handleKeypress(term *terminal.Terminal, event terminal.Event) {
	var handler func(*terminal.Terminal)
	var ok bool
	if event.Ch == 0 {
		handler, ok = specialKeyMap[event.Key]
	} else {
		handler, ok = keyMap[event.Ch]
	}

	if ok {
		handler(term)
	}
}

func toggleLine(term *terminal.Terminal) {
	view.ToggleLine(term.CursorY)
}
