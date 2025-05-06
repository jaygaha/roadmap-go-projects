package main

import (
	"flag"
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/tmdb-cli-tool/cmd"
	"github.com/jaygaha/roadmap-go-projects/tmdb-cli-tool/config"
)

func main() {
	// load the config
	if err := config.LoadEnv(); err != nil {
		fmt.Println("Error loading config: ", err)
		return
	}

	movieLists := flag.String("type", "", "Type of movie lists: playing, popular, top, upcoming")
	page := flag.Int("page", 1, "Page number")

	flag.Parse()

	if *movieLists == "" {
		flag.Usage()
		return
	}

	// get the movie lists
	cmd.FetchMovieList(*movieLists, *page)
}
