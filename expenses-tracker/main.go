package main

import (
	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/cmd"
)

func main() {
	cmd.Run()
	// init storage
	// if err := expense.Init(); err != nil {
	// 	log.Fatalf("Error initializing storage: %v", err)
	// 	return
	// }

	// // load expenses
	// expenses, err := expense.LoadExpenses()
	// if err != nil {
	// 	log.Fatalf("Error loading expenses: %v", err)
	// 	return
	// }

	// log.Printf("Loaded %d expenses", len(expenses))
}
