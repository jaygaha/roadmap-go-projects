package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jaygaha/roadmap-go-projects/task-tracker/task"
)

const tasksFile = "tasks.json"

func Init() (task.Tasks, error) {
	// Initialize storage
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		// Create empty file
		file, err := os.Create(tasksFile)
		if err != nil {
			fmt.Println("Error creating tasks file:", err)
			return nil, err
		}

		defer file.Close()

		// Write empty JSON array to file
		_, err = file.Write([]byte("[]"))
		if err != nil {
			return nil, err
		}
	}

	// load tasks
	tasks, err := loadTasks()

	if err != nil {
		fmt.Println("Error loading tasks:", err)
		return nil, err
	}

	return tasks, nil
}

func loadTasks() (task.Tasks, error) {
	fileData, err := os.ReadFile(tasksFile)
	if err != nil {
		fmt.Println("Error reading tasks file:", err)
		return nil, err
	}

	// unmarshal json
	tasks := task.Tasks{}

	err = json.Unmarshal(fileData, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func SaveTasks(tasks task.Tasks) error {
	// marshal tasks
	tasksData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	// write tasks to file
	err = os.WriteFile(tasksFile, tasksData, os.ModeAppend.Perm()) //os.WriteFile writes data to a file named by filename. If the file does not exist, WriteFile creates it with given permissions
	if err != nil {
		return err
	}

	return nil
}
