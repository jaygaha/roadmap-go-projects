package expense

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/expenses-tracker-cli/models"
)

// CSVFile is the path to the CSV file
var CSVFile = "expenses.csv"

// checks if the CSV file exists, if not, creates it
func Init() error {
	if _, err := os.Stat(CSVFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(CSVFile)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"ID", "Date", "Description", "Amount"}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}
	return nil
}

// reads the CSV file and returns the expenses
func LoadExpenses() ([]models.Expense, error) {
	file, err := os.Open(CSVFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var expenses []models.Expense

	for i, record := range records {
		// skip the header row
		if i == 0 {
			continue
		}
		// skip invalid rows
		if len(record) != 4 {
			continue
		}

		id, _ := strconv.Atoi(record[0])
		amount, _ := strconv.ParseFloat(record[3], 64)

		expense := models.Expense{
			ID:          id,
			Date:        record[1],
			Description: record[2],
			Amount:      amount,
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func SaveExpenses(expenses []models.Expense) error {
	file, err := os.Create(CSVFile)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Date", "Description", "Amount"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, expense := range expenses {
		record := []string{
			strconv.Itoa(expense.ID),
			expense.Date,
			expense.Description,
			strconv.FormatFloat(expense.Amount, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func AddExpense(e models.Expense) (models.Expense, error) {
	records, err := LoadExpenses()
	if err != nil {
		return models.Expense{}, err
	}

	newID := 1
	if len(records) > 0 {
		newID = records[len(records)-1].ID + 1
	}
	e.ID = newID

	// open the file in append mode
	file, err := os.OpenFile(CSVFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return models.Expense{}, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(e.ID),
		e.Date,
		e.Description,
		strconv.FormatFloat(e.Amount, 'f', 2, 64),
	}

	// write the new record to the CSV file
	if err := writer.Write(record); err != nil {
		return models.Expense{}, err
	}

	return e, nil
}

func GetExpenseByID(id int) (models.Expense, int, error) {
	expenses, err := LoadExpenses()
	if err != nil {
		return models.Expense{}, -1, err
	}
	for i, e := range expenses {
		if e.ID == id {
			return e, i, nil
		}
	}
	return models.Expense{}, -1, errors.New("expense not found")
}

func UpdateExpense(updatedExpense models.Expense) error {
	expenses, err := LoadExpenses()
	if err != nil {
		return err
	}

	// open the file in write mode
	file, err := os.OpenFile(CSVFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the header row
	headers := []string{"ID", "Date", "Description", "Amount"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// write the updated records to the CSV file
	for _, e := range expenses {
		record := []string{
			strconv.Itoa(e.ID),
			e.Date,
			e.Description,
			strconv.FormatFloat(e.Amount, 'f', 2, 64),
		}
		writer.Write(record)
	}

	return nil
}
