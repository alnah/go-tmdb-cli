// cmd/usage.go
package main

import "fmt"

// showUsage displays the main help message.
func showUsage() {
	fmt.Print(`TMDB CLI - Movie & TV Database Command Line Interface

USAGE:
    tmdb <command> [options]

MOVIE COMMANDS:
    popular [count]              Show popular movies
    top-rated [count]            Show top-rated movies
    now-playing [count]          Show movies in theaters
    upcoming [count]             Show upcoming releases
    search <query>               Search for movies
    discover [options]           Discover movies with filters

TV SHOW COMMANDS:
    tv popular [count]           Show popular TV shows
    tv top-rated [count]         Show top-rated TV shows
    tv on-the-air [count]        Show TV shows currently airing
    tv search <query>            Search for TV shows
    tv discover [options]        Discover TV shows with filters

OTHER COMMANDS:
    version                      Show version information
    config                       Show configuration help
    help                         Show this help

SHORTCUTS:
    tmdb "movie title"           Auto-search for movie
    tmdb pop 50                  50 popular movies
    tmdb top                     Top-rated movies
    tmdb now                     Now playing movies
    tmdb soon                    Upcoming movies
    tmdb tv                      Show TV usage
    tmdb tv pop                  Popular TV shows
    tmdb tv top                  Top-rated TV shows
    tmdb tv air                  Currently airing TV shows

OPTIONS:
    --format table|json|csv      Output format (default: table)
    --max-items N                Maximum items to show (default: 20)
    --original-title             Show original language titles
    --no-header                  Don't show table headers
    --verbose                    Enable verbose logging
    --debug                      Enable debug logging

SEARCH OPTIONS:
    --language LANG              Language filter (e.g., en, fr, es)
    --year YEAR                  Year filter (for TV: first air year)
    --min-rating N               Minimum rating (0-10)
    --max-rating N               Maximum rating (0-10)

DISCOVER OPTIONS:
    --genres action,comedy       Include genres
    --exclude horror,thriller    Exclude genres
    --sort popularity|rating|votes|release_date|title
                                 Sort by field
    --order asc|desc             Sort order
    --min-votes N                Minimum vote count
    --max-votes N                Maximum vote count

EXAMPLES:
    tmdb popular                 # 20 popular movies
    tmdb "The Matrix"            # Search for The Matrix
    tmdb search "action 2023"    # Search action movies from 2023
    tmdb top-rated --format json # Top movies as JSON
    tmdb discover --genres action --year 2023 --min-rating 7
    tmdb discover --sort votes --order desc --min-votes 1000
    tmdb now-playing --original-title

    tmdb tv popular              # 20 popular TV shows
    tmdb tv search "Breaking"    # Search for TV shows with "Breaking"
    tmdb tv on-the-air 10        # 10 currently airing TV shows
    tmdb tv discover --genres drama --min-rating 8 --min-votes 500

CONFIGURATION:
    Set TMDB_API_KEY environment variable or create config file.
    Run 'tmdb config' for detailed configuration help.

Get your API key from: https://www.themoviedb.org/settings/api
`)
}

// showTVUsage displays TV-specific help message.
func showTVUsage() {
	fmt.Print(`TV Show Commands:

USAGE:
    tmdb tv <subcommand> [options]

SUBCOMMANDS:
    popular [count]              Show popular TV shows
    top-rated [count]            Show top-rated TV shows
    on-the-air [count]           Show TV shows currently airing
    search <query>               Search for TV shows
    discover [options]           Discover TV shows with filters

OPTIONS:
    --format table|json|csv      Output format (default: table)
    --max-items N                Maximum items to show (default: 20)
    --original-title             Show original language names
    --no-header                  Don't show table headers

DISCOVER OPTIONS:
    --language LANG              Language filter
    --year YEAR                  First air year
    --min-rating N               Minimum rating (0-10)
    --max-rating N               Maximum rating (0-10)
    --min-votes N                Minimum vote count
    --max-votes N                Maximum vote count
    --genres comedy,drama        Include genres
    --exclude reality            Exclude genres
    --sort popularity|rating|votes|first_air_date|name
                                 Sort by field
    --order asc|desc             Sort order

EXAMPLES:
    tmdb tv popular              # 20 popular TV shows
    tmdb tv search "The Office"  # Search for The Office
    tmdb tv on-the-air --format json
    tmdb tv discover --genres comedy --min-rating 8
    tmdb tv discover --sort votes --min-votes 5000
`)
}
