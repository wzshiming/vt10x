package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gdamore/tcell"
	"github.com/hinshun/vt10x"
	"github.com/kr/pty"
)

func main() {
	err := goterm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func goterm() error {
	var state vt10x.State
	cmd := exec.Command(os.Getenv("SHELL"), "-i")
	pty, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	term, err := vt10x.Create(&state, pty)
	if err != nil {
		return err
	}
	defer term.Close()

	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer s.Fini()

	os.Setenv("LINES", "24")
	os.Setenv("COLUMNS", "80")

	err = s.Init()
	if err != nil {
		return err
	}

	width, height := s.Size()
	vt10x.ResizePty(pty, width-2, height-2)
	term.Resize(width-2, height-2)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	endc := make(chan bool)
	updatec := make(chan struct{}, 1)
	go func() {
		defer close(endc)
		for {
			err := term.Parse()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
			select {
			case updatec <- struct{}{}:
			default:
			}
		}
	}()

	go func() {
		io.Copy(pty, os.Stdin)
	}()

	eventc := make(chan tcell.Event, 4)
	go func() {
		for {
			eventc <- s.PollEvent()
		}
	}()

	for {
		select {
		case event := <-eventc:
			switch ev := event.(type) {
			case *tcell.EventResize:
				width, height = ev.Size()
				vt10x.ResizePty(pty, width-2, height-2)
				term.Resize(width-2, height-2)
				s.Sync()
			}
		case <-endc:
			return nil
		case <-updatec:
			update(s, &state, width-2, height-2)
		}
	}
}

func update(s tcell.Screen, state *vt10x.State, w, h int) {
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	for i := 0; i < h+2; i++ {
		s.SetContent(0, i, '|', nil, tcell.StyleDefault)
		s.SetContent(w+1, i, '|', nil, tcell.StyleDefault)
	}
	for i := 0; i < w+2; i++ {
		s.SetContent(i, 0, '─', nil, tcell.StyleDefault)
		s.SetContent(i, h+1, '─', nil, tcell.StyleDefault)
	}

	s.SetContent(0, 0, '┌', nil, tcell.StyleDefault)
	s.SetContent(w+1, 0, '┐', nil, tcell.StyleDefault)
	s.SetContent(w+1, h+1, '┘', nil, tcell.StyleDefault)
	s.SetContent(0, h+1, '└', nil, tcell.StyleDefault)

	state.Lock()
	defer state.Unlock()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c, fg, bg := state.Cell(x, y)

			var style tcell.Style
			style.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
			s.SetContent(x+1, y+1, c, nil, style)
		}
	}
	if state.CursorVisible() {
		curx, cury := state.Cursor()
		curx += 1
		cury += 1

		s.ShowCursor(curx, cury)
	} else {
		s.HideCursor()
	}

	s.Show()
}
