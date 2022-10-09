package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const PaddleHeight = 4
const InitalBallVelocityRow = 1
const InitalBallVelocityCol = 2

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var debugLog string

var gameObjects []*GameObject

func main() {
	InitScreen()
	GameState()
	inputChannel := UserInput()

	for !IsGameOver() {
		HandleUserInput(ReadInput(inputChannel))
		UpdateState()
		DrawState()
		time.Sleep(50 * time.Millisecond)

	}

	screenWidth, screenHeigh := screen.Size()
	winner := GetWinner()
	PrintStringCentered(screenHeigh/2-1, screenWidth/2, "Game Over!")
	PrintStringCentered(screenHeigh/2, screenWidth/2, fmt.Sprintf("%s Wins!", winner))
	screen.Show()

	time.Sleep(3 * time.Second)
	screen.Fini()

}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func GetWinner() string {
	screenWidth, _ := screen.Size()
	if ball.col < 0 {
		return "Player 2"
	} else if ball.col >= screenWidth {
		return "Player 1"
	} else {
		return " "
	}
}

func IsGameOver() bool {
	return GetWinner() != " "
}

func UpdateState() {

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if CollideWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if CollideWithPaddle(ball, player1Paddle) || CollideWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}

}

func CollideWithPaddle(obj *GameObject, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}
	return collidesOnColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}
func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()

}

func CollideWithWall(obj *GameObject) bool {
	_, screenHeigh := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeigh

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
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
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

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}
	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: InitalBallVelocityRow, velCol: InitalBallVelocityCol,
		symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
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
