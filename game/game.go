package game

import (
	"fmt"
	"io"
	"time"
)

type Point struct {
	X int
	Y int
}

type Config struct {
	Width    int
	Height   int
	TickRate time.Duration
}

func DefaultConfig() Config {
	return Config{Width: 40, Height: 20, TickRate: 120 * time.Millisecond}
}

type Snake struct {
	Head Point
	Body []Point
}

type GameState struct {
	Snake *Snake
}

func render_game(w io.Writer, game_state *GameState) {
	snake := game_state.Snake
	fmt.Fprintln(w, "Snake head is", snake.Head)
}

func Run(r io.Reader, w io.Writer, cfg Config) error {
	if cfg.Width <= 0 {
		return fmt.Errorf("invalid width: %d", cfg.Width)
	}

	input_chan := make(chan byte, 8)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := r.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				input_chan <- buf[0]
			}
		}
	}()

	snake := Snake{
		Head: Point{X: 5, Y: 5},
		Body: []Point{
			{X: 5, Y: 5},
			{X: 4, Y: 5},
			{X: 3, Y: 5},
		},
	}
	direction := Point{X: 0, Y: 0}
	game_state := GameState{Snake: &snake}
	ticker := time.NewTicker(cfg.TickRate).C

	fmt.Fprintln(w, "Hello, World! FROM game.go")
	// var current_time time.Time
	for {
		select {
		case <-ticker:
			// current_time = t
		case cur_input := <-input_chan:
			switch cur_input {
			case 'w':
				direction.X = 0
				direction.Y = 1
			case 'a':
				direction.X = -1
				direction.Y = 0
			case 's':
				direction.X = 0
				direction.Y = -1
			case 'd':
				direction.X = 1
				direction.Y = 0
			}

			if cur_input == 3 {
				return nil
			}

			fmt.Fprintln(w, "Input was: ", cur_input)
		}
		snake.Head.X = snake.Head.X + direction.X
		snake.Head.Y = snake.Head.Y + direction.Y
		render_game(w, &game_state)
	}

	return nil
}
