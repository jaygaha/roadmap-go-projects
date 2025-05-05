# Expenses Tracker CLI app using `Go`

This project is a command-line application that allows users to track and manage their expenses. Built with `Go`, it provides a simple interface for adding, listing, updating, and deleting expenses, as well as generating expense summaries.

## Features

- Add new expenses with description and amount
- List all recorded expenses
- Update existing expenses
- Delete expenses
- Generate expense summaries (total or filtered by month/year)
- Data persistence using **CSV** storage

## Installation

### Prerequisites

- Go 1.24 or higher

### Steps

1. Clone the repository
   ```bash
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   cd roadmap-go-projects/expenses-tracker
   ```

2. Build the application
   ```bash
   go build -o expenses-tracker
   ```

3. Run the application
   ```bash
   ./expenses-tracker
   ```

## Usage

### Available Commands

```
Usage: expense-tracker <command> [--flags]
Available commands:
  add        Add a new expense
  list       List all expenses
  update     Update an expense
  delete     Delete an expense
  summary    Get total expenses
```

### Adding an Expense

```bash
./expenses-tracker add --description "Groceries" --amount 45.67
```

### Listing All Expenses

```bash
./expenses-tracker list
```

Output example:
```
ID	Date		Description	Amount
1	2025-05-04	Groceries	45.67
2	2023-05-04	Gas		    35.50
```

### Updating an Expense

```bash
./expenses-tracker update --id 1 --description "Grocery shopping" --amount 50.25
```

### Deleting an Expense

```bash
./expenses-tracker delete --id 1
```

### Getting Expense Summary

```bash
# Get total of all expenses
./expenses-tracker summary

# Get total for a specific month
./expenses-tracker summary --month 5

# Get total for a specific year
./expenses-tracker summary --year 2025

# Get total for a specific month and year
./expenses-tracker summary --month 5 --year 2025
```

## Project Structure

```
expenses-tracker/
├── cmd/
│   └── commands.go       # Command definitions and CLI entry point
├── expense/
│   └── storage.go        # Data storage operations (CSV handling)
├── handlers/
│   ├── add_expense.go    # Handler for adding expenses
│   ├── delete_expense.go # Handler for deleting expenses
│   ├── list_expense.go   # Handler for listing expenses
│   ├── summary_expense.go# Handler for expense summaries
│   └── update_expense.go # Handler for updating expenses
├── models/
│   └── expense.go        # Expense data model
├── expenses.csv          # Data storage file
├── go.mod                # Go module definition
└── main.go               # Application entry point
```

## Data Storage

Expenses are stored in a CSV file (`expenses.csv`) with the following structure:

- ID: Unique identifier for each expense
- Date: Date of the expense (YYYY-MM-DD format)
- Description: Description of the expense
- Amount: Cost of the expense

## Project Link

- https://roadmap.sh/projects/expenses-tracker