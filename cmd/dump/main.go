package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hinshun/vt10x"
	"github.com/kr/pty"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ptm, pts, err := pty.Open()
	if err != nil {
		return err
	}
	defer pts.Close()
	defer ptm.Close()

	c := exec.Command(os.Getenv("SHELL"))
	c.Stdout = pts
	c.Stdin = pts
	c.Stderr = pts

	term := vt10x.New(vt10x.WithWriter(ptm))

	rows, cols := term.Size()
	vt10x.ResizePty(ptm, cols, rows)

	go func() {
		br := bufio.NewReader(ptm)
		for {
			err := term.Parse(br)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		}
	}()

	err = c.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	fmt.Println(term.String())
	return nil
}
