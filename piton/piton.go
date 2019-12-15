/*
SERPENT - a simple program to play a famous game in text mode
Copyright 2019 Eugenio Menegatti
myindievg@gmail.com

	 This file is part of SERPENT.
	 The file COPYING describes the terms under which SERPENT is distributed.

   SERPENT is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   SERPENT is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with SERPENT.  If not, see <http://www.gnu.org/licenses/>.
*/

package piton

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	term "github.com/nsf/termbox-go"
)

/*Coord describes the Y, X coord on the board*/
type Coord struct {
	x int
	y int
}

/*GameStatus retains the status for all the things that matters in a game*/
type GameStatus struct {
	fruits []Coord
}

/*GameSequence is the type of the array that describes a game*/
type GameSequence []int

/*BoardType is the type of the board*/
type BoardType [][]int

const startingX = 5
const startingY = 5

const playableBoardW = 20
const playableBoardH = 10

const none = -1

/*Esc is the ESC key*/
const Esc = 0

/*Right direction*/
const Right = 1

/*Left direction*/
const Left = 2

/*Up direction*/
const Up = 3

/*Down direction*/
const Down = 4

/*Q1 is north east quadrant*/
const Q1 = 5

/*Q2 is north west quadrant*/
const Q2 = 6

/*Q3 is south west quadrant*/
const Q3 = 7

/*Q4 is south east quadrant*/
const Q4 = 8

const space = 9
const enter = 10

/*Wall value*/
const Wall = -1

/*Empty value*/
const Empty = 0

/*Snake value*/
const Snake = 1

/*Neck value*/
const Neck = 2

/*Fruit value*/
const Fruit = -2

/*Poison value*/
const Poison = -3

/*MaxGameSequenceLength is the maximum number of moves the player can do during the game*/
const MaxGameSequenceLength = 10000

/*MaxGameScore is the maximum possible score before the game ends*/
const MaxGameScore = 1000

/*CurrentBoard holds the current board */
var CurrentBoard BoardType
var startingBoard BoardType
var snakeXHead, snakeYHead, snakeXTail, snakeYTail int
var fruitX, fruitY int
var direction int
var currentFruitIndex int
var score int

/*GetSnakeDirection returns the direction the snake is moving to*/
func GetSnakeDirection() int {
	return direction
}

/*GetSnakeX returns the x coord of the head of the snake*/
func GetSnakeX() int {
	return snakeXHead
}

/*GetSnakeY returns the y coord of the head of the snake*/
func GetSnakeY() int {
	return snakeYHead
}

/*IsDanger returns true if the cell is dangerous*/
func IsDanger(cell int) bool {
	if cell == Wall || cell >= Snake || cell == Poison {
		return true
	}
	return false
}

/*IsFruit returns true if the cell is a fruit*/
func IsFruit(cell int) bool {
	if cell == Fruit {
		return true
	}
	return false
}

/*SnakeHeadNeighbor returns the cell next to the snake's in the given direction*/
func SnakeHeadNeighbor(where int) int {
	var cell int
	switch where {
	case Right:
		cell = CurrentBoard[snakeYHead][snakeXHead+1]
	case Left:
		cell = CurrentBoard[snakeYHead][snakeXHead-1]
	case Up:
		cell = CurrentBoard[snakeYHead+1][snakeXHead]
	case Down:
		cell = CurrentBoard[snakeYHead-1][snakeXHead]
	}
	return cell
}

/*FruitLocation returns the quadrant of the fruit relative to the pivot location*/
func FruitLocation(pivotX int, pivotY int) int {
	var where int
	if fruitY == pivotY && fruitX > pivotX {
		where = Right
	}
	if fruitY == pivotY && fruitX < pivotX {
		where = Left
	}
	if fruitX == pivotY && fruitY < pivotY {
		where = Up
	}
	if fruitX == pivotY && fruitY > pivotY {
		where = Down
	}

	if fruitY < pivotY && fruitX > pivotX {
		where = Q1
	}
	if fruitY < pivotY && fruitX < pivotX {
		where = Q2
	}
	if fruitY > pivotY && fruitX < pivotX {
		where = Q3
	}
	if fruitY > pivotY && fruitX > pivotX {
		where = Q4
	}
	return where
}

/*Init all in Severus*/
func Init() {
	rand.Seed(time.Now().UTC().UnixNano())
	initBoard()
}

/*NewGame init the game*/
func NewGame(game *GameStatus) {
	score = 0
	initBoard()
	initSnake()
	initFruit(game)
}

func initBoard() {
	initBoard := BoardType{
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1},
		{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	}
	startingBoard = initBoard
	CurrentBoard = initBoard
}

func initFruit(game *GameStatus) {
	currentFruitIndex = 0
	respawnFruit(game)
}

func initSnake() {
	snakeXHead = startingX
	snakeYHead = startingY
	snakeXTail = startingX
	snakeYTail = startingY + 2

	direction = none
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	CurrentBoard[snakeYTail-1][snakeXTail] = Snake + 1
	CurrentBoard[snakeYTail][snakeXTail] = Snake + 2
}

/*Close releases everything*/
func Close() {
}

/*InitTerm initializes the terminal*/
func InitTerm() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
}

/*CloseTerm closes the terminal*/
func CloseTerm() {
	term.Close()
}

func getNewFruitCoords(board *BoardType) (int, int) {
	fruitX = -1
	fruitY = -1
watchDog:
	for count := 0; count < 10; count++ {
		x := rand.Intn(playableBoardW) + 1
		y := rand.Intn(playableBoardH) + 1
		if isEmpty(board, x, y) {
			fruitX = x
			fruitY = y
			break watchDog
		}
	}

	return fruitX, fruitY
}

func isInsideSnakeBody(x int, y int) bool {
	if CurrentBoard[y][x] >= Snake {
		return true
	}
	return false
}

func isEmpty(board *BoardType, x int, y int) bool {
	if (*board)[y][x] == Empty {
		return true
	}
	return false
}

/*ClearConsole clears the console*/
func ClearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

/*OutputBoard prints the board on screen*/
func OutputBoard(board BoardType) {
	//fmt.Println("board h = ", len(board))
	//fmt.Println("board w = ", len(board[0]))
	for _, row := range CurrentBoard {
		for _, checker := range row {
			if checker != Wall {
				if checker == Empty {
					fmt.Print(".")
				}
				if checker == Snake { // Snake's head
					fmt.Print("@")
				} else {
					if checker > Snake {
						fmt.Print("O")
					}
				}
				if checker == Fruit {
					fmt.Print("F")
				}
				if checker == Poison {
					fmt.Print("P")
				}
			}
		}
		fmt.Println()
	}
}

func bufferedReadKey() byte {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	char := []byte(input)[0]
	fmt.Print(string(char))
	return char
}

/*GenKeyboardEventQueue generates a channel to retrieve keybord events*/
func GenKeyboardEventQueue() chan term.Event {
	eventQueue := make(chan term.Event)
	go func() {
		for {
			eventQueue <- term.PollEvent()
		}
	}()
	return eventQueue
}

/*KeyPressed gets a key from the keyboard, delaying if waitFlag is true*/
func KeyPressed(keyboard chan term.Event, waitFlag bool) rune {
	// https://github.com/nsf/termbox-go/issues/7
	var key rune
	key = none
	var waitTime int
	if waitFlag {
		waitTime = 250
	} else {
		waitTime = 0
	}
	select {
	case ev := <-keyboard:
		if ev.Type == term.EventKey {
			switch ev.Key {
			case term.KeyEsc:
				term.Sync()
				fmt.Println("Esc pressed")
				key = Esc
			case term.KeyArrowUp:
				key = Up
				term.Sync()
				fmt.Println("Arrow Up pressed")
			case term.KeyArrowDown:
				key = Down
				term.Sync()
				fmt.Println("Arrow Down pressed")
			case term.KeyArrowLeft:
				key = Left
				term.Sync()
				fmt.Println("Arrow Left pressed")
			case term.KeyArrowRight:
				key = Right
				term.Sync()
				fmt.Println("Arrow Right pressed")
			case term.KeySpace:
				key = space
				term.Sync()
				fmt.Println("Space pressed")
			case term.KeyEnter:
				key = enter
				term.Sync()
				fmt.Println("Enter pressed")
			default:
				key = none
				term.Sync()
				fmt.Println("ASCII : ", ev.Ch)
			}
		}
	case <-time.After(time.Duration(waitTime) * time.Millisecond):
	}
	return key
}

func snakeSetRight() {
	direction = Right
}

func snakeSetLeft() {
	direction = Left
}

func snakeSetUp() {
	direction = Up
}

func snakeSetDown() {
	direction = Down
}

func getNextFruit(game *GameStatus) (int, int) {
	x := game.fruits[currentFruitIndex].x
	y := game.fruits[currentFruitIndex].y
	currentFruitIndex++
	return x, y
}

func respawnFruit(game *GameStatus) {
	if game == nil {
		fruitX, fruitY = getNewFruitCoords(&CurrentBoard)
		if fruitX != -1 && fruitY != -1 {
			CurrentBoard[fruitY][fruitX] = Fruit
		}
	} else {
		fruitX, fruitY = getNextFruit(game)
	}
	CurrentBoard[fruitY][fruitX] = Fruit

	score++
}

func snakeProceedGivenMatch(status *GameStatus) bool {
	return snakeProceedImpl(status)
}

func snakeProceed() bool {
	return snakeProceedImpl(nil)
}

func snakeProceedImpl(status *GameStatus) bool {
	var gameOver bool = false
	switch direction {
	case Right:
		if canGo, what := snakeCanGoRight(); canGo {
			if what == Fruit {
				snakeGrowRight()
				respawnFruit(status)
			} else {
				snakeMoveRight()
			}
		} else {
			if what == Neck {
				// do nothing
			} else {
				snakeDies()
				gameOver = true
			}
		}
	case Left:
		if canGo, what := snakeCanGoLeft(); canGo {
			if what == Fruit {
				snakeGrowLeft()
				respawnFruit(status)
			} else {
				snakeMoveLeft()
			}
		} else {
			if what == Neck {
				// do nothing
			} else {
				snakeDies()
				gameOver = true
			}
		}
	case Up:
		if canGo, what := snakeCanGoUp(); canGo {
			if what == Fruit {
				snakeGrowUp()
				respawnFruit(status)
			} else {
				snakeMoveUp()
			}
		} else {
			if what == Neck {
				// do nothing
			} else {
				snakeDies()
				gameOver = true
			}
		}
	case Down:
		if canGo, what := snakeCanGoDown(); canGo {
			if what == Fruit {
				snakeGrowDown()
				respawnFruit(status)
			} else {
				snakeMoveDown()
			}
		} else {
			if what == Neck {
				// do nothing
			} else {
				snakeDies()
				gameOver = true
			}
		}
	}

	return gameOver
}

func snakeCanGoRight() (bool, int) {
	if CurrentBoard[snakeYHead][snakeXHead+1] == Empty {
		return true, Empty
	}
	if CurrentBoard[snakeYHead][snakeXHead+1] >= Snake {
		if CurrentBoard[snakeYHead][snakeXHead+1] == Neck {
			return false, Neck
		}
		return false, Snake
	}
	if CurrentBoard[snakeYHead][snakeXHead+1] == Fruit {
		return true, Fruit
	}
	if CurrentBoard[snakeYHead][snakeXHead+1] == Poison {
		return true, Poison
	}
	return false, Wall
}

func snakeCanGoLeft() (bool, int) {
	if CurrentBoard[snakeYHead][snakeXHead-1] == Empty {
		return true, Empty
	}
	if CurrentBoard[snakeYHead][snakeXHead-1] >= Snake {
		if CurrentBoard[snakeYHead][snakeXHead-1] == Neck {
			return false, Neck
		}
		return false, Snake
	}
	if CurrentBoard[snakeYHead][snakeXHead-1] == Fruit {
		return true, Fruit
	}
	if CurrentBoard[snakeYHead][snakeXHead-1] == Poison {
		return true, Poison
	}
	return false, Wall

}

func snakeCanGoUp() (bool, int) {
	if CurrentBoard[snakeYHead-1][snakeXHead] == Empty {
		return true, Empty
	}
	if CurrentBoard[snakeYHead-1][snakeXHead] >= Snake {
		if CurrentBoard[snakeYHead-1][snakeXHead] == Neck {
			return false, Neck
		}
		return false, Snake
	}
	if CurrentBoard[snakeYHead-1][snakeXHead] == Fruit {
		return true, Fruit
	}
	if CurrentBoard[snakeYHead-1][snakeXHead] == Poison {
		return true, Poison
	}
	return false, Wall
}

func snakeCanGoDown() (bool, int) {
	if CurrentBoard[snakeYHead+1][snakeXHead] == Empty {
		return true, Empty
	}
	if CurrentBoard[snakeYHead+1][snakeXHead] >= Snake {
		if CurrentBoard[snakeYHead+1][snakeXHead] == Neck {
			return false, Neck
		}
		return false, Snake
	}
	if CurrentBoard[snakeYHead+1][snakeXHead] == Fruit {
		return true, Fruit
	}
	if CurrentBoard[snakeYHead+1][snakeXHead] == Poison {
		return true, Poison
	}
	return false, Wall
}

func snakeCanGoDirection(desiredDirection int) (bool, int) {
	if desiredDirection == Right {
		return snakeCanGoRight()
	}
	if desiredDirection == Left {
		return snakeCanGoLeft()
	}
	if desiredDirection == Up {
		return snakeCanGoUp()
	}
	if desiredDirection == Down {
		return snakeCanGoDown()
	}
	return false, -1
}

func snakeDies() {

}

func moveSnake(x int, y int, snake int) {
	if CurrentBoard[y][x-1] == snake {
		if y == snakeYTail && x-1 == snakeXTail {
			CurrentBoard[y][x-1] = Empty
			snakeXTail = x
			return
		}
		CurrentBoard[y][x-1] = snake + 1
		moveSnake(x-1, y, snake+1)
		return
	}
	if CurrentBoard[y][x+1] == snake {
		if y == snakeYTail && x+1 == snakeXTail {
			CurrentBoard[y][x+1] = Empty
			snakeXTail = x
			return
		}
		CurrentBoard[y][x+1] = snake + 1
		moveSnake(x+1, y, snake+1)
		return
	}
	if CurrentBoard[y-1][x] == snake {
		if y-1 == snakeYTail && x == snakeXTail {
			CurrentBoard[y-1][x] = Empty
			snakeYTail = y
			return
		}
		CurrentBoard[y-1][x] = snake + 1
		moveSnake(x, y-1, snake+1)
		return
	}
	if CurrentBoard[y+1][x] == snake {
		if y+1 == snakeYTail && x == snakeXTail {
			CurrentBoard[y+1][x] = Empty
			snakeYTail = y
			return
		}
		CurrentBoard[y+1][x] = snake + 1
		moveSnake(x, y+1, snake+1)
		return
	}
}

func growSnake(x int, y int, snake int) {
	if CurrentBoard[y][x-1] == snake {
		if y == snakeYTail && x-1 == snakeXTail {
			CurrentBoard[y][x-1] = snake + 1
			//snakeXTail = x
			return
		}
		CurrentBoard[y][x-1] = snake + 1
		growSnake(x-1, y, snake+1)
		return
	}
	if CurrentBoard[y][x+1] == snake {
		if y == snakeYTail && x+1 == snakeXTail {
			CurrentBoard[y][x+1] = snake + 1
			//snakeXTail = x
			return
		}
		CurrentBoard[y][x+1] = snake + 1
		growSnake(x+1, y, snake+1)
		return
	}
	if CurrentBoard[y-1][x] == snake {
		if y-1 == snakeYTail && x == snakeXTail {
			CurrentBoard[y-1][x] = snake + 1
			//snakeYTail = y
			return
		}
		CurrentBoard[y-1][x] = snake + 1
		growSnake(x, y-1, snake+1)
		return
	}
	if CurrentBoard[y+1][x] == snake {
		if y+1 == snakeYTail && x == snakeXTail {
			CurrentBoard[y+1][x] = snake + 1
			//snakeYTail = y
			return
		}
		CurrentBoard[y+1][x] = snake + 1
		growSnake(x, y+1, snake+1)
		return
	}
}

func snakeMoveRight() {
	snakeXHead = snakeXHead + 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	moveSnake(snakeXHead, snakeYHead, Snake)
}

func snakeMoveLeft() {
	snakeXHead = snakeXHead - 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	moveSnake(snakeXHead, snakeYHead, Snake)
}

func snakeMoveUp() {
	snakeYHead = snakeYHead - 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	moveSnake(snakeXHead, snakeYHead, Snake)
}

func snakeMoveDown() {
	snakeYHead = snakeYHead + 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	moveSnake(snakeXHead, snakeYHead, Snake)
}

func snakeGrowRight() {
	snakeXHead = snakeXHead + 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	growSnake(snakeXHead, snakeYHead, Snake)
}

func snakeGrowLeft() {
	snakeXHead = snakeXHead - 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	growSnake(snakeXHead, snakeYHead, Snake)
}

func snakeGrowUp() {
	snakeYHead = snakeYHead - 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	growSnake(snakeXHead, snakeYHead, Snake)
}

func snakeGrowDown() {
	snakeYHead = snakeYHead + 1
	CurrentBoard[snakeYHead][snakeXHead] = Snake
	growSnake(snakeXHead, snakeYHead, Snake)
}

/*HumanPlay lets you play the game*/
func HumanPlay() int {
	InitTerm()
	defer CloseTerm()
	keyboard := GenKeyboardEventQueue()
	var gameOver bool = false
mainLoop:
	for !gameOver {
		ClearConsole()
		OutputBoard(CurrentBoard)
		key := KeyPressed(keyboard, true)
		switch key {
		case Right:
			snakeSetRight()
		case Left:
			snakeSetLeft()
		case Up:
			snakeSetUp()
		case Down:
			snakeSetDown()
		case Esc:
			break mainLoop
		}
		gameOver = snakeProceed()
	}
	score--
	return score
}

func generateFruits(game *GameStatus) {
	for i := 0; i < MaxGameScore; i++ {
		game.fruits[i].x, game.fruits[i].y = getNewFruitCoords(&startingBoard)
	}
}

/*GenerateGameParams generates a configuration fot the game*/
func GenerateGameParams() GameStatus {
	var status GameStatus
	status.fruits = make([]Coord, MaxGameScore)
	generateFruits(&status)
	return status
}

/*PlayAlone lets the computer play 1 game and returns the game sequence*/
func PlayAlone(verboseFlag bool, game *GameStatus) GameSequence {
	gameOver := false
	//var currentGameSequence GameSequence
	//currentGameSequence = make([]int, MaxGameSequenceLength)
	currentGameSequence := []int{}
	i := 0
	for !gameOver && i < MaxGameSequenceLength {
		if verboseFlag {
			ClearConsole()
			OutputBoard(CurrentBoard)
		}
		direction = getRandomValidMove()
		switch direction {
		case Up:
			fmt.Print("U")
		case Down:
			fmt.Print("D")
		case Right:
			fmt.Print("R")
		case Left:
			fmt.Print("L")
		}
		gameOver = snakeProceedGivenMatch(game)
		//currentGameSequence[i] = direction
		currentGameSequence = append(currentGameSequence, direction)
		i++
	}
	score--
	/*fmt.Println()
	fmt.Println("Final score: ", score)
	fmt.Println()*/
	return currentGameSequence
}

/*getRandomValidMove returns a valid move. Here Valid doesn't mean the snake won't die*/
func getRandomValidMove() int {
	var randDirection int
findValidDir:
	for {
		randDirection = rand.Intn(4) + 1
		if canGo, what := snakeCanGoDirection(randDirection); canGo {
			break findValidDir
		} else {
			if what != Neck {
				break findValidDir
			}
		}
	}
	return randDirection
}

/*GetRandomSolution generates a solution by playing a game by choosing random moves
until the snake dies
*/
func GetRandomSolution(game *GameStatus) (GameSequence, int) {
	gameSequence := PlayAlone(false, game)
	return gameSequence, score
}

/*GetContinuingSolution continue the game starting from the given sequence*/
func GetContinuingSolution(game *GameStatus, gameSequence *GameSequence, prevScore int) (GameSequence, int) {
	score = prevScore
	newGameSequence := ReplayGame(false, game, gameSequence)
	return newGameSequence, score
}

/*ReplayGame lets the computer play 1 game and returns the game sequence*/
func ReplayGame(verboseFlag bool, game *GameStatus, inputGameSequence *GameSequence) GameSequence {
	gameOver := false
	//var currentGameSequence GameSequence
	//currentGameSequence = make([]int, MaxGameSequenceLength)
	currentGameSequence := []int{}
	currentGameSequence = append(currentGameSequence, *inputGameSequence...)
	i := 0
	for !gameOver && i < MaxGameSequenceLength {
		if verboseFlag {
			ClearConsole()
			OutputBoard(CurrentBoard)
		}

		if i < len(*inputGameSequence) {
			direction = (*inputGameSequence)[i]
		} else {
			direction = getRandomValidMove()
		}

		switch direction {
		case Up:
			fmt.Print("U")
		case Down:
			fmt.Print("D")
		case Right:
			fmt.Print("R")
		case Left:
			fmt.Print("L")
		}
		//gameOver = snakeProceedGivenMatch(game)
		gameOver = snakeProceed()
		currentGameSequence = append(currentGameSequence, direction)
		i++
	}
	score--
	/*fmt.Println()
	fmt.Println("Final score: ", score)
	fmt.Println()*/
	return currentGameSequence
}
