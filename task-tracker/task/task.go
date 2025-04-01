package task

import (
	"errors"
	"fmt"
	"time"

	"slices"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

// Custome type for status (enum)
type Status string

const (
	TodoStatus       Status = "todo"
	InProgressStatus Status = "in-progress"
	DoneStatus       Status = "done"
)

// Task struct
// these are the fields that represent a todo data
type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      Status     `json:"status"` // enum: todo, in-progress, done
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"` //default:null
}

// Filter struct
// these are the fields that represent a filter
type FilterStatus struct {
	Status Status `json:"status"`
}

// Task slice
// this is the list of tasks
type Tasks []Task

// Add task
func (t *Tasks) AddTask(description string) (string, error) {
	// Check if description is empty
	if description == "" {
		err := errors.New("description cannot be blank")
		fmt.Println(err.Error())

		return "", err
	}

	// initialize id
	id := 1

	// if list is not empty, get the last id and add 1
	if len(*t) > 0 {
		id = (*t)[len(*t)-1].Id + 1
	}

	// Set the values according to the given fields
	newTask := Task{
		Id:          id,
		Description: description,
		Status:      TodoStatus, // default: todo
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}

	// Append the new task to the list
	// using pointer *t, as we are modifying the list
	*t = append(*t, newTask)

	return fmt.Sprintf("Task added successfully (ID: %d)\n", id), nil
}

// Get tasks by filter
// camelcase vs capitalize
// getTasksByFilter: GetTasksByStatus
// local variable: Global variable
func (t *Tasks) getTasksByFilter(status FilterStatus) (Tasks, error) {
	// validate status
	if status.Status != "" {
		err := t.validateTaskStatus(status.Status)

		if err != nil {
			return nil, err
		}
	}

	// filter tasks by status if given
	if status.Status != "" {
		filteredTasks := Tasks{} // initialize empty slice

		for _, task := range *t {
			// if task status is equal to the given status, add it to the filtered tasks
			if task.Status == status.Status {
				filteredTasks = append(filteredTasks, task)
			}
		}

		return filteredTasks, nil
	}

	return *t, nil
}

// Render tasks in a table
func (t *Tasks) RenderTasks(status string) error {
	tasks, err := t.getTasksByFilter(FilterStatus{Status: Status(status)})

	if err != nil {
		fmt.Sprintln(err.Error())
		return nil
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgCyan).SprintfFunc()

	table := table.New("Id", "Description", "Status", "Created", "Updated")
	table.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Loop through the tasks and add them to the table
	for _, task := range tasks {
		// Handle nil UpdatedAt
		updated := "-"
		if task.UpdatedAt != nil {
			updated = task.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		// Add the task to the table
		table.AddRow(task.Id, task.Description, task.Status, task.CreatedAt.Format("2006-01-02 15:04:05"), updated)
	}

	// add - if no record found
	if len(tasks) == 0 {
		table.AddRow("-", "-", "-", "-", "-")
	}

	table.Print()

	return nil
}

// Update task
func (t *Tasks) UpdateTask(id int, description string, status FilterStatus) error {
	// Check if id is valid
	err := t.validateTaskId(id)

	if err != nil {
		return err
	}

	// Check if description is empty and status is also empty
	if description == "" && status.Status == "" {
		err := errors.New("description cannot be blank")
		fmt.Println(err.Error())

		return err
	}

	// validate status
	if status.Status != "" {
		err := t.validateTaskStatus(status.Status)

		if err != nil {
			return err
		}
	}

	// Loop through the tasks and update the task
	for i, task := range *t {
		if task.Id == id {
			// Set the values according to the given fields
			// update description if given
			if description != "" {
				(*t)[i].Description = description
			}
			// update status if given
			if status.Status != "" {
				(*t)[i].Status = status.Status
			}
			// update updated_at if either description or status is given
			if description != "" || status.Status != "" {
				updatedAt := time.Now()
				(*t)[i].UpdatedAt = &updatedAt
			}
			return nil
		}
	}

	return nil
}

// Delete task
func (t *Tasks) DeleteTask(id int) error {
	// Check if id is valid
	err := t.validateTaskId(id)
	if err != nil {
		return err
	}

	// Loop through the tasks and delete the task
	for i, task := range *t {
		if task.Id == id {
			// delete the task
			*t = slices.Delete((*t), i, i+1)

			return nil
		}
	}
	return nil
}

// validate id
func (t *Tasks) validateTaskId(id int) error {
	for _, task := range *t {
		if task.Id == id {
			return nil // nil means the task is exists and valid
		}
	}

	err := errors.New("task not found")
	fmt.Println(err.Error())

	return err
}

// validate status
func (t *Tasks) validateTaskStatus(status Status) error {
	if status == TodoStatus || status == InProgressStatus || status == DoneStatus {
		return nil // nil means the task is exists and valid
	}

	err := errors.New("invalid status")
	fmt.Println(err.Error())

	return err
}
