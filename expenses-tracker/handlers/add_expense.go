package handlers

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/models"
)

func AddExpenseHandler(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	description := fs.String("description", "", "Description of the expense")
	amount := fs.Float64("amount", 0, "Amount of the expense")

	fs.Parse(args)

	if *description == "" || *amount == 0 {
		fmt.Println("Description and amount are required.")
		fs.Usage()
		return
	}

	records, err := expense.LoadExpenses()
	if err != nil {
		log.Fatalf("Failed to load expenses: %v", err)
	}

	newID := len(records) + 1

	e := models.Expense{
		ID:          newID,
		Date:        time.Now().Format("2006-01-02"),
		Description: *description,
		Amount:      *amount,
	}

	records = append(records, e)

	if err := expense.SaveExpenses(records); err != nil {
		log.Fatalf("Failed to save expenses: %v", err)
	}

	fmt.Printf("Expense added successfully (ID: %d)\n", newID)
}
