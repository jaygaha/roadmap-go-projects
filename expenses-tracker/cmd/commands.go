package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/handlers"
)

type Command struct {
	Name        string
	Description string
	Handler     func(args []string)
}

var commands = []Command{
	{Name: "add", Description: "Add a new expense", Handler: handlers.AddExpenseHandler},
	{Name: "list", Description: "List all expenses", Handler: handlers.ListExpensesHandler},
	{Name: "update", Description: "Update an expense", Handler: handlers.UpdateExpenseHandler},
	{Name: "delete", Description: "Delete an expense", Handler: handlers.DeleteExpenseHandler},
	{Name: "summary", Description: "Get total expenses", Handler: handlers.SummaryExpenseHandler},
}

func Run() {
	// init storage
	if err := expense.Init(); err != nil {
		log.Fatalf("Error initializing storage: %v", err)
		return
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	for _, cmd := range commands {
		if cmd.Name == cmdName {
			cmd.Handler(cmdArgs)
			return
		}
	}

	fmt.Printf("Unknown command: %s\n", cmdName)
	printUsage()
}

func printUsage() {
	fmt.Println("Usage: expense-tracker <command> [--flags]")
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("  %-10s %s\n", cmd.Name, cmd.Description)
	}
}
