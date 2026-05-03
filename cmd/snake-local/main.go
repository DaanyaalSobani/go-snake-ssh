package main

import (
	"log"
	"os"

	"github.com/DaanyaalSobani/go-snake-ssh/game"
	"golang.org/x/term"
)

func main() {
	term.MakeRaw(int(os.Stdin.Fd()))

	if err := game.Run(os.Stdin, os.Stdout, game.DefaultConfig()); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(game.Point{X: 1, Y: 1})

}
