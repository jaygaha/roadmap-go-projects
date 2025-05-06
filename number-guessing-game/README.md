# Number Guessing Game in `Go`

A simple command-line number guessing game written in `Go`. The game generates a random number between 1 and 100, and the player has to guess it within a limited number of attempts based on the selected difficulty level.

## Features

- **Multiple Difficulty Levels**:
  - Easy: 10 attempts
  - Medium: 5 attempts
  - Hard: 3 attempts

- **Helpful Hints**:
  - The game tells you if your guess is too high or too low
  - Special notification when you're getting close (within ±5 of the target number)
  - Warning when you're on your last attempt

- **Performance Tracking**:
  - Tracks the number of attempts used
  - Measures the time taken to guess the correct number

- **Replay Option**:
  - Option to play again after completing or failing a game

## Installation

### Prerequisites

- Go 1.24.0 or higher

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   ```

2. Navigate to the project directory:
   ```bash
   cd roadmap-go-projects/number-guessing-game
   ```

3. Build the game:
   ```bash
   go build
   ```

4. Run the game:
   ```bash
   ./number-guessing-game
   ```

## How to Play

1. Start the game
2. Select a difficulty level (1-3)
3. Enter your guess when prompted
4. Follow the hints to adjust your next guess
5. Try to guess the number within the allowed attempts
6. Choose to play again or quit when the game ends

## Project Structure

```
number-guessing-game/
├── config/
│   └── game.go         # Game configuration (difficulty levels)
├── games/
│   └── guess_no.go     # Main game logic
├── utils/
│   └── utils.go        # Utility functions (random number generation)
└── main.go             # Entry point
```

## Game Mechanics

- The game generates a random number between 1 and 100
- Players select a difficulty level that determines the number of allowed attempts
- After each guess, the game provides feedback:
  - If the guess is correct, the player wins
  - If the guess is too high or too low, the game provides a hint
  - If the guess is within ±5 of the target, the game indicates the player is getting close
- The game tracks and displays the time taken to guess correctly
- After winning or losing, players can choose to play again

## Project Link

- [Number Guessing Game](https://roadmap.sh/projects/number-guessing-game)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)