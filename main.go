package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleHeight = 4
const PaddleSymbol = 0x2588

type Paddle struct {
	row, column, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle
var debugLog string

func main() {
	InitScreen()
	GameState()
	inputChannel := UserInput()

	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := ReadInput(inputChannel)
		HandleUserInput(key)

	}

}

func PrintString(row, column int, str string) {
	for _, c := range str {
		screen.SetContent(column, row, c, nil, tcell.StyleDefault)
		column += 1
	}
}
func Print(row, column, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(column+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	Print(player1.row, player1.column, player1.width, player1.height, PaddleSymbol)
	Print(player2.row, player2.column, player2.width, player2.height, PaddleSymbol)
	screen.Show()

}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1.row > 0 {
		player1.row--
	} else if key == "Rune[s]" && player1.row+player1.height < screenHeight {
		player1.row++
	} else if key == "Up" && player2.row > 0 {
		player2.row--
	} else if key == "Down" && player2.row+player2.height < screenHeight {
		player2.row++
	}
}

func UserInput() chan string {
	inputChannel := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChannel <- ev.Name()
			}
		}
	}()
	return inputChannel
}

func GameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1 = &Paddle{
		row: paddleStart, column: 0, width: 1, height: PaddleHeight,
	}
	player2 = &Paddle{
		row: paddleStart, column: width - 1, width: 1, height: PaddleHeight,
	}
}

func ReadInput(inputChannel chan string) string {
	var key string
	select {
	case key = <-inputChannel:
	default:
		key = ""
	}
	return key
}
