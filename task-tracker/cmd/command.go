package cmd

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/task-tracker/task"
)

const (
	CmdAdd            = "add"
	CmdList           = "list"
	CmdUpdate         = "update"
	CmdDelete         = "delete"
	CmdMarkInProgress = "in-progress"
	CmdMarkDone       = "done"
)

type CommandFlgs struct {
	Add    string
	List   string
	Update string
	Delete int
	Mark   string
	Args   []string
}

func ParseCommandFlgs() *CommandFlgs {
	cmdFlg := CommandFlgs{}

	// parse the command line flags
	flag.StringVar(&cmdFlg.List, "list", "", "List all tasks. Use filter to refine the tasks. Format: list <filter>\nAvailable filters: todo, in-progress, done")
	flag.StringVar(&cmdFlg.Add, "add", "", "Add a new task. Format: add 'new description'")
	flag.StringVar(&cmdFlg.Update, "update", "", "Edit a task by ID and new text. Format: Id 'new_text'")
	flag.IntVar(&cmdFlg.Delete, "delete", -1, "Delete a task by ID. Format: delete <id>")
	flag.StringVar(&cmdFlg.Mark, "mark", "", "Mark a task as done or in-progress. Format: mark-<status> <id>")

	flag.Parse() // Parse the command-line flags

	// parse the command line arguments
	cmdFlg.Args = flag.Args()

	return &cmdFlg
}

func (cmdFlg *CommandFlgs) Run(tasks *task.Tasks) {
	// Handle positional arguments
	if len(cmdFlg.Args) > 0 {
		// lower case the command name
		commandName := strings.ToLower(cmdFlg.Args[0])

		// check mark-<status> command
		commandName = strings.TrimPrefix(commandName, "mark-")

		switch commandName {
		case CmdList:
			cmdFlg.handleList(tasks)
		case CmdAdd:
			cmdFlg.handleAdd(tasks)
		case CmdUpdate:
			cmdFlg.handleUpdate(tasks)
		case CmdDelete:
			cmdFlg.handleDelete(tasks)
		case CmdMarkDone, CmdMarkInProgress:
			status := strings.TrimPrefix(commandName, "mark-")
			cmdFlg.handleToggleStatus(tasks, status)
		default:
			fmt.Println("unknown command")
			fmt.Println("\nUsage: task-tracker <command> [option]")
		}
	} else {
		fmt.Println("command is required")
		fmt.Println("\nUsage: task-tracker <command> [option]")
	}
}

func (c *CommandFlgs) handleAdd(tasks *task.Tasks) {
	if len(c.Args) > 1 {
		// add the task
		res, err := tasks.AddTask(c.Args[1])

		if err != nil {
			fmt.Println("Error adding task:", err)
			return
		}

		fmt.Println(res)
	} else {
		fmt.Println("description is required")
		fmt.Println("\nUsage: task-tracker add <text>")
	}
}

func (c *CommandFlgs) handleList(tasks *task.Tasks) {
	filter := ""
	// check if the filter is provided
	if len(c.Args) > 1 {
		filter = c.Args[1]
	}

	// list the tasks
	err := tasks.RenderTasks(filter)

	if err != nil {
		fmt.Println("Error rendering tasks:", err)
		return
	}
}

func (c *CommandFlgs) handleUpdate(tasks *task.Tasks) {
	if len(c.Args) < 3 {
		fmt.Println("id and new text are required")
		fmt.Println("\nUsage: task-tracker update <id> <new_text>")
		return
	}

	// convert the id to int
	id, err := strconv.Atoi(c.Args[1])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	err = tasks.UpdateTask(id, c.Args[2], task.FilterStatus{})
	if err != nil {
		return
	}
}

func (c *CommandFlgs) handleDelete(tasks *task.Tasks) {
	if len(c.Args) < 2 {
		fmt.Println("id is required")
		fmt.Println("\nUsage: task-tracker delete <id>")
		return
	}

	// convert the id to int
	id, err := strconv.Atoi(c.Args[1])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	err = tasks.DeleteTask(id)
	if err != nil {
		return
	}
}

func (c *CommandFlgs) handleToggleStatus(tasks *task.Tasks, status string) {
	if len(c.Args) < 2 {
		fmt.Println("id is required")
		fmt.Println("\nUsage: task-tracker mark-<status> <id>")
		return
	}

	filterStatus := task.FilterStatus{
		Status: task.Status(status),
	}

	// convert the id to int
	id, err := strconv.Atoi(c.Args[1])
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}

	err = tasks.UpdateTask(id, "", filterStatus)
	if err != nil {
		return
	}
}
