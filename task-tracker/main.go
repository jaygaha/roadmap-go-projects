package main

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/task-tracker/cmd"
	"github.com/jaygaha/roadmap-go-projects/task-tracker/storage"
)

func main() {
	tasks, err := storage.Init()

	if err != nil {
		fmt.Println("Error initializing storage:", err)
		return
	}

	// handle the command
	cmd := cmd.ParseCommandFlgs()
	cmd.Run(&tasks)

	// finally save to the storage
	err = storage.SaveTasks(tasks)

	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
}
