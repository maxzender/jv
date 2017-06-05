package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/maxzender/jsonexplorer/colorwriter"
	"github.com/maxzender/jsonexplorer/jsonfmt"
	"github.com/maxzender/jsonexplorer/jsontree"
	"github.com/maxzender/jsonexplorer/terminal"
	termbox "github.com/nsf/termbox-go"
)

var (
	colorMap = map[jsonfmt.TokenType]termbox.Attribute{
		jsonfmt.DelimiterType: termbox.ColorDefault,
		jsonfmt.BoolType:      termbox.ColorBlue,
		jsonfmt.StringType:    termbox.ColorRed,
		jsonfmt.NumberType:    termbox.ColorYellow,
		jsonfmt.NullType:      termbox.ColorCyan,
	}
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
	writer := colorwriter.New(colorMap, termbox.ColorDefault)
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

	tree := jsontree.New(formattedJson)

	for {
		term.Render(tree)
		e := term.Poll()
		if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
			return 0
		}
		handleKeypress(term, tree, e)
	}
	return 0
}

func handleKeypress(t *terminal.Terminal, j *jsontree.JsonTree, e terminal.Event) {
	if e.Ch == 0 {
		switch e.Key {
		case termbox.KeyArrowUp:
			t.MoveCursor(0, -1)
		case termbox.KeyArrowDown:
			t.MoveCursor(0, +1)
		case termbox.KeyArrowLeft:
			t.MoveCursor(-1, 0)
		case termbox.KeyArrowRight:
			t.MoveCursor(+1, 0)
		case termbox.KeyEnter:
			j.ToggleLine(t.CursorY)
		case termbox.KeySpace:
			j.ToggleLine(t.CursorY)
		}
	} else {
		switch e.Ch {
		case 'h':
			t.MoveCursor(-1, 0)
		case 'j':
			t.MoveCursor(0, +1)
		case 'k':
			t.MoveCursor(0, -1)
		case 'l':
			t.MoveCursor(+1, 0)
		}
	}
}
