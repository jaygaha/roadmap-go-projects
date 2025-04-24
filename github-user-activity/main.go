package main

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/github-user-activity/cmd"
)

func main() {
	// username, filter, err := cmd.ExecuteCommand()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(username)
	// fmt.Println(filter)
	if err := cmd.ExecuteCommand(); err != nil {
		fmt.Println(err)
	}

}
