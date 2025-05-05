package handlers

import (
	"fmt"
	"log"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
)

func ListExpensesHandler(args []string) {
	expenses, err := expense.LoadExpenses()
	if err != nil {
		log.Fatalf("Failed to load expenses: %v", err)
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses found.")
	} else {
		fmt.Println("ID\tDate\t\tDescription\tAmount")
		for _, e := range expenses {
			fmt.Printf("%d\t%s\t%s\t\t%.2f\n", e.ID, e.Date, e.Description, e.Amount)
		}
	}
}
