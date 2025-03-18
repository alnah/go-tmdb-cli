package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type contextKey string

const dependencies contextKey = "deps"

// Dependencies provides shared services for CLI commands to access TMDB API.
type Dependencies struct {
	URLBuilder *urlBuilder
	Client     *httpClient
}

// newRootCmd creates the root command to organize all subcommands and CLI setup.
func newRootCmd(fileName string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "go-tmdb-cli",
		Args:  cobra.NoArgs,
		Short: "A CLI for The Movie Database (TMDB)",
		Long: `A simple command-line interface (CLI) to fetch data from The
Movie Database (TMDB), and display it in the terminal.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := initialize(&defaultUserHome{}, fileName)
			if err != nil {
				return err
			}
			apiKey := viper.GetString("api_key")
			if apiKey == "" {
				return fmt.Errorf(`missing API key in ~/.go-tmdb-cli/%s,
please ensure you include your API key in the following format:
  api_key: YOUR_API_KEY`, fileName)
			}
			deps := &Dependencies{
				URLBuilder: newURLBuilder(),
				Client:     newHTTPClient(apiKey),
			}
			ctx := context.WithValue(cmd.Context(), dependencies, deps)
			cmd.SetContext(ctx)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.AddCommand(
		completionCommand(),
		newListCmd(),
		newDiscoverCmd(),
		newInfoCmd(),
	)
	return rootCmd
}

// newInfoCmd defines the command to show CLI version and authorship details.
func newInfoCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "info",
		Args:  cobra.NoArgs,
		Short: "Display version number, author and licence",
		Long:  "All software has a version, an author, and a license. These are the details for Go TMDB-CLI.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Go TMDB-CLI v1.0.0")
			cmd.Println("Copyright (c) 2025 Alexis Nahan <alexis.nahan@gmail.com>")
			cmd.Println("Licensed under the Apache License v2.0")
		},
	}
	return versionCmd
}

// newListCmd creates the command to display pre-defined movie categories.
func newListCmd() *cobra.Command {
	var isNowPlaying, isPopular, isTopRated, isUpcoming bool
	movieListCmd := &cobra.Command{
		Use:   "list",
		Short: "Display a ready-made movie list",
		Long: `Retrieve and display a curated list of movies from The Movie 
Database (TMDB), including categories such as now playing, popular, top rated, 
and upcoming, formatted as a user-friendly table.`,
		Example: `  go-tmdb-cli list -n
  go-tmdb-cli list -p
  go-tmdb-cli list -t
  go-tmdb-cli list -u`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				_ = cmd.Help()
				return nil
			}
			deps, err := getDependencies(cmd)
			if err != nil {
				return err
			}
			var url string
			switch {
			case isNowPlaying:
				url, _ = deps.URLBuilder.list("now_playing")
			case isPopular:
				url, _ = deps.URLBuilder.list("popular")
			case isTopRated:
				url, _ = deps.URLBuilder.list("top_rated")
			case isUpcoming:
				url, _ = deps.URLBuilder.list("upcoming")
			}
			tmdbRes, err := asyncFetchMovies(deps.Client, url, 20)
			if err != nil {
				return err
			}
			got := formatResults(tmdbRes)
			cmd.Println(got)
			return nil
		},
	}
	flags := map[string]struct {
		alias   string
		help    string
		enabled *bool
	}{
		"now": {"n", "now playing movies", &isNowPlaying},
		"pop": {"p", "popular movies", &isPopular},
		"top": {"t", "top rated movies", &isTopRated},
		"up":  {"u", "upcoming movies", &isUpcoming},
	}
	for name, flag := range flags {
		movieListCmd.Flags().BoolVarP(flag.enabled, name, flag.alias, false, flag.help)
	}
	return movieListCmd
}

// newDiscoverCmd builds the command for advanced movie searches with filters.
func newDiscoverCmd() *cobra.Command {
	discoverCmd := &cobra.Command{
		Use:   "discover",
		Short: "Discover movies based on various criteria",
		Long: `Discover enables users to explore a diverse selection of films 
that align with their interests and preferences, for more refined searches.`,
		Example: `  go-tmdb-cli discover  -l=en  -y=2000,2005  -g=comedy,action  -a=6.5,10   -v=100,50000  -m=100  -s=average,desc
  go-tmdb-cli discover  -l=fr  -y=1960,gte   -g=history        -a=7,gte    -v100,gte     -m=50   -s=title,asc
  go-tmdb-cli discover  -l=pt  -y=1960,lte   -w=comedy         -a=9.0,lte  -v=2000,lte   -m=10   -s=votes,asc
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().NFlag() == 0 {
				_ = cmd.Help()
				return nil
			}
			deps, err := getDependencies(cmd)
			if err != nil {
				return err
			}
			var url, sort, maxItems string
			q := queryParams{}
			flags := map[string]*string{
				"language":       &q.Language,
				"year":           &q.Year,
				"average":        &q.VoteAverage,
				"votes":          &q.VoteCount,
				"genres":         &q.WithGenres,
				"without-genres": &q.WithoutGenres,
				"sort":           &sort,
				"max-items":      &maxItems,
			}
			for name, value := range flags {
				if flagValue, _ := cmd.Flags().GetString(name); flagValue != "" {
					*value = flagValue
				}
			}
			url, err = deps.URLBuilder.discover(q)
			if err != nil {
				return err
			}
			var wantItems int
			if maxItems == "" {
				wantItems = 20
			} else {
				wantItems, err = strconv.Atoi(maxItems)
				if err != nil {
					return fmt.Errorf(`validation error: items must be an integer, e.g. "50"`)
				}
			}
			movies, err := asyncFetchMovies(deps.Client, url, wantItems)
			if err != nil {
				return err
			}
			if sort != "" {
				_, err = movies.sortByField(sort)
				if err != nil {
					return err
				}
			}
			output := formatResults(movies)
			cmd.Println(output)
			return nil
		},
	}
	flags := []struct {
		name  string
		alias string
		help  string
	}{
		{"language", "l", "original language (not the country!)"},
		{"year", "y", "primary release year or dates"},
		{"average", "a", "votes average"},
		{"votes", "v", "vote counts"},
		{"genres", "g", "with one or many genres"},
		{"without-genres", "w", "without one or many genres"},
		{"sort", "s", "sort by field and order"},
		{"max-items", "m", fmt.Sprintf("maximum number of movies, default 20, max %d", APIMaxItems)},
	}
	for _, flag := range flags {
		discoverCmd.Flags().StringP(flag.name, flag.alias, "", flag.help)
	}
	return discoverCmd
}

// completionCommand generates shell autocompletion scripts (hidden helper).
func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "completion",
		Short:  "Generate the autocompletion script for the specified shell",
		Hidden: true,
	}
}

// getDependencies retrieves API clients from context for command execution.
func getDependencies(cmd *cobra.Command) (*Dependencies, error) {
	deps, ok := cmd.Context().Value(dependencies).(*Dependencies)
	if !ok {
		return nil, fmt.Errorf("retrieve dependencies from context")
	}
	return deps, nil
}

// formatResults converts movie data into a formatted table for terminal output.
func formatResults(movies movies) string {
	if len(movies) == 0 {
		return "No results available. Please try another query."
	}
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{
		"#",
		"Original Title",
		"Release Date",
		"Title",
		"Average",
		"Votes",
	})
	table.SetRowLine(true)
	table.SetBorder(true)
	table.SetColumnSeparator("│")
	table.SetRowSeparator("⎯")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for i, r := range movies {
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			r.OriginalTitle,
			r.ReleaseDate,
			r.Title,
			fmt.Sprintf("%.1f", r.VoteAverage),
			fmt.Sprintf("%d", r.VoteCount),
		})
	}
	table.Render()
	return buf.String()
}
