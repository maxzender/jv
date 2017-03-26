package main

import (
	"io/ioutil"
	"os"
	"strings"

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
	expandedLines = map[int]struct{}{0: struct{}{}}
	segments      = make(map[int]int)
)

func main() {
	term, err := terminal.New()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	byteContent, err := ioutil.ReadAll(os.Stdin)
	content := string(byteContent)
	segments = parseSegments(content)

	for {
		term.Draw(string(content))
		e := term.Poll()
		if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
			return
		} else {
			handleKeypress(term, e)
		}
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

func parseSegments(content string) map[int]int {
	lines := strings.Split(content, "\n")
	resultSegments := make(map[int]int)
	bracketBalances := make(map[int]int)
	var bal int
	for num, line := range lines {
		for _, c := range line {
			switch c {
			case '{', '[':
				bal += 1
				bracketBalances[bal] = num
			case '}', ']':
				resultSegments[bracketBalances[bal]] = num
				bal -= 1
			}
		}
	}

	return resultSegments
}

func toggleLine(term *terminal.Terminal) {
}
