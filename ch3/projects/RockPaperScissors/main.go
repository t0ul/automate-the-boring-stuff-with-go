package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// Move is one of the three choices. Giving it a named type keeps the rest of the
// program from mixing it up with any other int.
type Move int

const (
	Rock Move = iota
	Paper
	Scissors
)

var names = [...]string{Rock: "Rock", Paper: "Paper", Scissors: "Scissors"}

var validInput = map[string]Move{
	"r": Rock, "rock": Rock,
	"p": Paper, "paper": Paper,
	"s": Scissors, "scissors": Scissors,
}

// outcome is indexed by (computer - player + 3) % 3.
var outcome = [...]string{"You tied.", "You lost.", "You win!"}

func getInput() string {
	fmt.Print("Pick [r]ock, [p]aper, or [s]cissors (ctrl+c to quit): ")
	var input string
	fmt.Scanln(&input)
	// lower-casing input is a good habit: users type however they like
	return strings.ToLower(input)
}

func playerMove() Move {
	// loop until the input is one we understand
	for {
		if move, ok := validInput[getInput()]; ok {
			fmt.Println("You chose", names[move])
			return move
		}
		fmt.Println("I didn't understand that, please retry")
	}
}

func main() {
	// an endless loop: the game runs until the player presses ctrl+c
	for {
		player := playerMove()
		computer := Move(rand.IntN(3))
		fmt.Printf("Computer chose %s\n", names[computer])
		fmt.Printf("Result: %s\n\n", outcome[(computer-player+3)%3])
	}
}
