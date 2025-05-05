package handlers

import (
	"flag"
	"fmt"
	"log"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
)

func UpdateExpenseHandler(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	id := fs.Int("id", 0, "ID of the expense to update")
	description := fs.String("description", "", "Description of the expense")
	amount := fs.Float64("amount", 0, "Amount of the expense")

	fs.Parse(args)

	// validate id
	if *id <= 0 {
		fmt.Println("Invalid ID. ID must be a positive integer.")
		fs.Usage()
		return
	}

	expenses, err := expense.LoadExpenses()
	if err != nil {
		log.Fatalf("Failed to load expenses: %v", err)
	}

	foundExpense, index, err := expense.GetExpenseByID(*id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update fields
	if *description != "" {
		foundExpense.Description = *description
	}
	if *amount != -1 && *amount != 0 {
		foundExpense.Amount = *amount
	}

	expenses[index] = foundExpense

	if err := expense.SaveExpenses(expenses); err != nil {
		log.Fatalf("Failed to save expenses: %v", err)
	}

	fmt.Println("Expense updated successfully!")
}
