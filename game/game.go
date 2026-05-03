package game

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type Point struct {
	X int
	Y int
}

type Config struct {
	Width     int
	Height    int
	CellWidth int
	TickRate  time.Duration
}

func DefaultConfig() Config {
	return Config{Width: 15, Height: 15, CellWidth: 2, TickRate: 120 * time.Millisecond}
}

type Snake struct {
	Head Point
	Body []Point
}

type GameState struct {
	Snake *Snake
}

func render_game(w io.Writer, game_state *GameState, cfg *Config) {
	var buf bytes.Buffer
	buf.WriteString("\x1b[H")

	snake := game_state.Snake
	grid := make([][]rune, cfg.Height)
	for i := range grid {
		grid[i] = make([]rune, cfg.Width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	grid[snake.Head.Y][snake.Head.X] = '█'

	// border: top and bottom rows
	for x := 0; x < cfg.Width; x++ {
		grid[0][x] = '-'
		grid[cfg.Height-1][x] = '-'
	}
	// border: left and right columns
	for y := 0; y < cfg.Height; y++ {
		grid[y][0] = '|'
		grid[y][cfg.Width-1] = '|'
	}
	// corners
	grid[0][0] = '+'
	grid[0][cfg.Width-1] = '+'
	grid[cfg.Height-1][0] = '+'
	grid[cfg.Height-1][cfg.Width-1] = '+'

	for _, row := range grid {
		for _, cell := range row {
			buf.WriteRune(cell)
			buf.WriteRune(cell)
		}
		buf.WriteString("\x1b[K\r\n")
	}
	w.Write(buf.Bytes()) // ONE syscall, the whole frame
	fmt.Fprintf(w, "Snake head is %", snake.Head)
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
				direction.Y = -1
			case 'a':
				direction.X = -1
				direction.Y = 0
			case 's':
				direction.X = 0
				direction.Y = 1
			case 'd':
				direction.X = 1
				direction.Y = 0
			}

			if cur_input == 3 {
				return nil
			}

			// fmt.Fprintln(w, "Input was: ", cur_input)
		}
		snake.Head.X = snake.Head.X + direction.X
		snake.Head.Y = snake.Head.Y + direction.Y
		if snake.Head.X == 0 {
			snake.Head.X = cfg.Width - 1
		}
		if snake.Head.Y == 0 {
			snake.Head.Y = cfg.Height - 1
		}

		if snake.Head.X == cfg.Width {
			snake.Head.X = 0
		}
		if snake.Head.Y == cfg.Height {
			snake.Head.Y = 0
		}

		render_game(w, &game_state, &cfg)
	}

	return nil
}
