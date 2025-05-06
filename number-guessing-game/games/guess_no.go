package games

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jaygaha/roadmap-go-projects/number-guessing-game/config"
	"github.com/jaygaha/roadmap-go-projects/number-guessing-game/utils"
)

// level holds the difficuly level of the game
var difficulyLevel = map[int]config.Difficulty{
	1: {Name: "Easy", MaxAttempts: 10},
	2: {Name: "Medium", MaxAttempts: 5},
	3: {Name: "Hard", MaxAttempts: 3},
}

func StartGuessingGame() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
	fmt.Println("You have 5 chances to guess the correct number.")

	// bufio.NewReader reads input from the console
	reader := bufio.NewReader(os.Stdin)

	// infinite loop
	for {
		fmt.Println()
		fmt.Println("Please select the difficulty level:")
		fmt.Println("1. Easy (10 chances)")
		fmt.Println("2. Medium (5 chances)")
		fmt.Println("3. Hard (3 chances)")
		fmt.Println()

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		// validate user input
		difficulty, ok := difficulyLevel[choice]
		if !ok {
			fmt.Println("Invalid choice. Please try again.")
			return
		}

		fmt.Println()
		fmt.Printf("Great! You selected %s difficulty.\n", difficulty.Name)
		fmt.Println("Let's start the game!")
		fmt.Println()

		targetNumber := utils.GenerateRandomNumber(1, 100)
		attempts := 0
		maxAttempts := difficulty.MaxAttempts
		startTime := time.Now()

		for attempts < maxAttempts {
			var guess int
			fmt.Print("Enter your guess: ")
			fmt.Scanln(&guess)
			attempts++

			if guess == targetNumber {
				fmt.Println("Congratulations! You guessed the correct number in", attempts, "attempts.")

				elapsedTime := time.Since(startTime)
				fmt.Printf("You took %.2f seconds to complete the game.\n", elapsedTime.Seconds())

				break
			} else if guess < targetNumber {
				fmt.Printf("Incorrect! The number is greater than %d.\n", guess)
			} else {
				fmt.Printf("Incorrect! The number is less than %d.\n", guess)
			}

			// hints the user if they are close to the correct number; ±5 of target number
			if utils.CloseGuessNumber(targetNumber-guess) <= 5 {
				fmt.Println("You're getting close!")
			}

			// alert user that last attempt
			if attempts == maxAttempts-1 {
				fmt.Println("This is your last attempt. Be careful!")
			}

			if attempts == maxAttempts {
				fmt.Println("Sorry, you've reached the maximum number of attempts.")
				fmt.Println("The correct number was", targetNumber)
			}
		}

		// ask user if they want to play again
		fmt.Println()
		fmt.Print("Do you want to play again? (yes/q): ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		if input != "yes" {
			fmt.Println("Thank you for playing! See you!!")
			break
		}
	}
}
