package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	tick = time.Second
)

var (
	ticker         *time.Ticker
	queues         chan termbox.Event
	isStarted      bool
	startX, startY int
)

func draw() {
	w, h := termbox.Size()
	clear()

	str := formatTime(time.Now())
	text := toText(str)

	if !isStarted {
		isStarted = true
		startX, startY = w/2-text.width()/2, h/2-text.height()/2
	}

	x, y := startX, startY
	for _, s := range text {
		echo(s, x, y)
		x += s.width()
	}

	flush()
}

func start() {
	ticker = time.NewTicker(tick)
}

func formatTime(t time.Time) string {
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func countdown() {
	var exitCode int

	start()

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
				exitCode = 1
				break loop
			}
		case <-ticker.C:
			draw()
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	queues = make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	draw()
	countdown()
}
