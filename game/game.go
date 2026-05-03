package game

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
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
	Head   Point
	Length int
}

type GameState struct {
	Snake     *Snake
	Direction *Point
	Food      *Point
}

func new_food(cfg Config) Point {
	return Point{
		X: rand.IntN(cfg.Width-cfg.CellWidth) + cfg.CellWidth - 1,
		Y: rand.IntN(cfg.Height-cfg.CellWidth) + cfg.CellWidth - 1,
	}
}

func render_game(w io.Writer, game_state *GameState, cfg *Config) {
	var buf bytes.Buffer
	buf.WriteString("\x1b[H")

	snake := game_state.Snake
	food := game_state.Food

	grid := make([][]rune, cfg.Height)
	for i := range grid {
		grid[i] = make([]rune, cfg.Width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}
	grid[food.Y][food.X] = '❤'
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
	fmt.Fprintf(w, "Snake is %d \n food is %d", snake, food)
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
		Head:   Point{X: 5, Y: 5},
		Length: 0,
	}
	direction := Point{X: 0, Y: 0}
	food := new_food(cfg)
	game_state := GameState{Snake: &snake, Direction: &direction,
		Food: &food}
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
		step_game(&game_state, &cfg)
		render_game(w, &game_state, &cfg)
	}

	return nil
}

func step_game(game_state *GameState, cfg *Config) {
	snake := game_state.Snake
	direction := game_state.Direction
	food := game_state.Food

	snake.Head.X = snake.Head.X + direction.X
	snake.Head.Y = snake.Head.Y + direction.Y
	if snake.Head.X == 0 {
		snake.Head.X = cfg.Width - cfg.CellWidth
	} else if snake.Head.X == cfg.Width-1 {
		snake.Head.X = 1
	}

	if snake.Head.Y == 0 {
		snake.Head.Y = cfg.Height - cfg.CellWidth
	} else if snake.Head.Y == cfg.Height-1 {
		snake.Head.Y = 1
	}
	if snake.Head == *food {
		snake.Length += 1
		old_food := *food
		for old_food == *food {
			*food = new_food(*cfg)
		}

	}

}
