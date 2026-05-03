package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DaanyaalSobani/go-snake-ssh/game"
	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatalf("enter raw mode: %v", err)
	}
	defer term.Restore(fd, oldState)

	// query the terminal size via stdout's fd — on Windows / git-bash, stdin may be
	// a pipe wrapper rather than a real console handle, but stdout is reliable.
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Printf("get terminal size: %v (using defaults)", err)
		w, h = 60, 20 // sane fallback
	}

	// size the game to the terminal:
	//   - logical cells render 2 chars wide (CellWidth = 2), so width is half
	//   - subtract 1 row for the HUD line
	cfg := game.DefaultConfig()
	cfg.Width = w / 2
	cfg.Height = h - 1

	if cfg.Width < 12 || cfg.Height < 8 {
		fmt.Fprintf(os.Stdout,
			"Terminal too small (%dx%d). Resize to at least 24x9 and rerun.\r\n", w, h)
		return
	}

	if err := game.Run(os.Stdin, os.Stdout, cfg); err != nil {
		log.Print(err) // not log.Fatal — that calls os.Exit, which skips defers
	}
}
