package handlers

import (
	"flag"
	"fmt"
	"log"

	"slices"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
)

func DeleteExpenseHandler(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	id := fs.Int("id", 0, "ID of the expense to delete")

	fs.Parse(args)

	if *id == 0 {
		fmt.Println("ID is required.")
		fs.Usage()
		return
	}

	records, err := expense.LoadExpenses()
	if err != nil {
		log.Fatalf("Failed to load expenses: %v", err)
	}

	_, index, err := expense.GetExpenseByID(*id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Delete the expense from the slice
	records = slices.Delete(records, index, index+1)

	// save
	if err := expense.SaveExpenses(records); err != nil {
		log.Fatalf("Failed to save expenses: %v", err)
	}

	fmt.Println("Expense deleted successfully")
}
