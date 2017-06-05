package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/maxzender/jsonexplorer/jsonfmt"
	"github.com/maxzender/jsonexplorer/terminal"
	"github.com/maxzender/jsonexplorer/treemodel"
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
	colorMap = map[jsonfmt.TokenType]termbox.Attribute{
		jsonfmt.DelimiterType: termbox.ColorWhite,
		jsonfmt.BoolType:      termbox.ColorBlue,
		jsonfmt.StringType:    termbox.ColorRed,
		jsonfmt.NumberType:    termbox.ColorYellow,
		jsonfmt.NullType:      termbox.ColorCyan,
	}
	tree *treemodel.TreeModel
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "print usage")
	flag.BoolVar(&showHelp, "help", false, "print usage")

	flag.Usage = usage
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	reader := os.Stdin
	var err error
	if len(os.Args) > 1 {
		reader, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer reader.Close()
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	os.Exit(run(content))
}

func run(content []byte) int {
	writer := NewColorWriter(colorMap, termbox.ColorDefault)
	formatter := jsonfmt.New(content, writer)
	if err := formatter.Format(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	formattedJson := writer.Lines

	term, err := terminal.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	defer term.Close()

	tree = treemodel.New(formattedJson)

	for {
		term.Draw(tree)
		e := term.Poll()
		if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
			return 0
		}
		handleKeypress(term, e)
	}
	return 0
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
	tree.ToggleLine(term.CursorY)
}
