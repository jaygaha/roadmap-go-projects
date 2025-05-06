# TMDB CLI Tool in `Go`

A command-line interface tool that allows you to fetch and display movie information from The Movie Database (TMDB) API. This tool provides a simple way to browse different categories of movies directly from your terminal.

## Features

- Fetch movies by different categories (now playing, popular, top rated, upcoming)
- Display detailed movie information including title, release date, overview, and rating
- Support for pagination to browse through multiple pages of results
- Emoji-rich output for better readability

## Prerequisites

- Go 1.24 or higher
- TMDB API key (get it from [https://www.themoviedb.org/settings/api](https://www.themoviedb.org/settings/api))

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/tmdb-cli-tool.git
cd tmdb-cli-tool
```

2. Create a `.env` file in the root directory based on the provided `.env.example`:

```bash
cp .env.example .env
```

3. Edit the `.env` file and add your TMDB API key:

```
TMDB_API_KEY=your_api_key_here
TMDB_API_URL=https://api.themoviedb.org/3
```

4. Build the application:

```bash
go build -o tmdb-app
```

## Usage

The basic command structure is:

```bash
./tmdb-app -type [movie_list_type] -page [page_number]
```

### Available Movie List Types

- `playing` - Now playing movies in theaters
- `popular` - Currently popular movies
- `top` - Top rated movies of all time
- `upcoming` - Upcoming movie releases

### Examples

Fetch currently playing movies (first page):
```bash
./tmdb-app -type playing
```

Fetch popular movies (second page):
```bash
./tmdb-app -type popular -page 2
```

Fetch top rated movies (third page):
```bash
./tmdb-app -type top -page 3
```

Fetch upcoming movies:
```bash
./tmdb-app -type upcoming
```

## Project Structure

- **main.go**: Entry point of the application, handles command-line flags and initiates the movie fetch process
- **api/client.go**: Contains functions to interact with the TMDB API
- **cmd/commands.go**: Implements the command logic for fetching and displaying movie information
- **config/config.go**: Handles loading environment variables from the .env file
- **models/movie.go**: Defines data structures for movie information and API responses

## Dependencies

This project uses standard Go libraries and does not require any external dependencies beyond the Go standard library.

## Environment Variables

| Variable | Description |
|----------|-------------|
| TMDB_API_KEY | Your TMDB API key |
| TMDB_API_URL | The base URL for TMDB API (default: https://api.themoviedb.org/3) |


## Acknowledgments

- [The Movie Database (TMDB)](https://www.themoviedb.org/) for providing the API
- [roadmap.sh](https://roadmap.sh/projects/tmdb-cli) for the project inspiration