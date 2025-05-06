package cmd

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/tmdb-cli-tool/api"
)

var movieListMap = map[string]string{
	"playing":  "now_playing",
	"popular":  "popular",
	"top":      "top_rated",
	"upcoming": "upcoming",
}

func FetchMovieList(movieListType string, page int) {
	listType, ok := movieListMap[movieListType]
	if !ok {
		fmt.Println("Invalid movie list type. Please choose from playing, popular, top, upcoming")
		return
	}

	movieLists, err := api.FetchMovieList(listType, page)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n🎬 %s movies:\n", transformUnderscoreToSpaceAndCapitalize(listType))
	fmt.Printf("Page: %d\n", page)

	for i, movie := range movieLists {
		fmt.Printf("\n%d.\n", i+1)
		fmt.Printf("#️⃣  %d\n", movie.ID)
		fmt.Printf("🎥 %s (%s)\n", movie.Title, movie.ReleaseDate)
		fmt.Printf("📝 %s\n", movie.Overview)
		fmt.Printf("⭐ %.1f\n", movie.VoteAverage)
	}

}

// transform underscore to space and capitalize the first letter
func transformUnderscoreToSpaceAndCapitalize(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '_' {
			s = s[:i] + " " + s[i+1:]
			i++
		}
	}

	// capitalize the first letter
	s = string(s[0]-32) + s[1:]
	return s
}
