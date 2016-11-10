package main

import "github.com/nsf/termbox-go"

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventCh := make(chan termbox.Event)
	draw := drawer()

	go pollEvents(eventCh)

	for event := range eventCh {
		draw(event.Ch)
	}
}

func pollEvents(eventCh chan termbox.Event) {
	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
				close(eventCh)
			} else {
				eventCh <- e
			}
		}
	}
}

func drawer() func(rune) {
	i := 0
	return func(r rune) {
		termbox.SetCell(i, 0, r, termbox.ColorWhite, termbox.ColorDefault)
		i++
		termbox.Flush()
	}
}
