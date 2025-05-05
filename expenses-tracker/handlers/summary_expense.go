package handlers

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/expense"
)

func SummaryExpenseHandler(args []string) {
	fs := flag.NewFlagSet("summary", flag.ExitOnError)

	month := fs.Int("month", 0, "Month for which to calculate the summary")
	year := fs.Int("year", 0, "Year for which to calculate the summary")

	fs.Parse(args)

	expenses, err := expense.LoadExpenses()
	if err != nil {
		log.Fatalf("Failed to load expenses: %v", err)
	}

	var total float64
	for _, e := range expenses {
		if *month != 0 && strings.Split(e.Date, "-")[1] != fmt.Sprintf("%02d", *month) {
			continue
		}
		if *year != 0 && strings.Split(e.Date, "-")[0] != fmt.Sprintf("%d", *year) {
			continue
		}

		total += e.Amount
	}

	fmt.Printf("Total expenses: $%.2f\n", total)
}
