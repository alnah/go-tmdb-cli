// Package internal provides core functionality for the TMDB CLI application.
package internal

import "fmt"

// ShowVersion displays version information.
func ShowVersion(version, buildDate, gitCommit string) {
	fmt.Printf("TMDB CLI %s\n", version)
	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}

// ShowUsage displays the main CLI usage information.
func ShowUsage(version string) {
	fmt.Printf("TMDB CLI %s\n\n", version)
	fmt.Println("A command-line interface for The Movie Database (TMDB)")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tmdb-cli [command] [arguments] [flags]")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println("  popular, pop              Get popular movies")
	fmt.Println("  top-rated, top, rated     Get top-rated movies")
	fmt.Println("  now-playing, playing      Get now playing movies")
	fmt.Println("  upcoming                  Get upcoming movies")
	fmt.Println("  search [query]            Search for movies")
	fmt.Println("  discover                  Discover movies with filters")
	fmt.Println("  tv [subcommand]           TV show commands")
	fmt.Println("  config                    Show configuration help")
	fmt.Println("  version, --version, -v    Show version information")
	fmt.Println("  help, --help, -h          Show this help message")
	fmt.Println("")
	fmt.Println("For more information, visit: https://github.com/alnah/tmdb-cli")
}

// ShowTVUsage displays TV command usage information.
func ShowTVUsage() {
	fmt.Println("TV Show Commands:")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tmdb-cli tv [subcommand] [arguments] [flags]")
	fmt.Println("")
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("  popular, pop              Get popular TV shows")
	fmt.Println("  top-rated, top, rated     Get top-rated TV shows")
	fmt.Println("  on-the-air, air, airing   Get TV shows currently on the air")
	fmt.Println("  search [query]            Search for TV shows")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  tmdb-cli tv popular 10")
	fmt.Println("  tmdb-cli tv search \"Breaking Bad\"")
}

// ShowSearchHelp displays search help suggestions.
func ShowSearchHelp() {
	fmt.Println("Try:")
	fmt.Printf("  - Check spelling and try again\n")
	fmt.Printf("  - Use fewer, more common words\n")
	fmt.Printf("  - Try the original language title\n")
}
