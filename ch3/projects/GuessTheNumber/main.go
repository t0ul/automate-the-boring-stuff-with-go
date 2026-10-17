package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	fmt.Println("Game: Guess a number between 0 and 9")
	fmt.Println("You have three (3) tries")

	// math/rand/v2 is seeded for you; no clock tricks needed
	secretNumber := rand.IntN(10)

	var guess int
	for try := 1; try <= 3; try++ {
		fmt.Printf("Attempt %d. Enter your number: ", try)
		if _, err := fmt.Scan(&guess); err != nil {
			fmt.Println("I don't recognize that input. Try a number like '1'")
			continue
		}
		switch {
		case guess == secretNumber:
			fmt.Println("You win!")
			return
		case guess < secretNumber:
			fmt.Println("Sorry, wrong guess; the number is bigger")
		case guess > secretNumber:
			fmt.Println("Sorry, wrong guess; the number is smaller")
		}
	}
	fmt.Printf("Game over! The correct number was %d\n", secretNumber)
}
