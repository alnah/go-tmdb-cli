package main

import (
	"bytes"
	"fmt"

	c "github.com/alnah/go-tmdb-cli/internal/config"
	d "github.com/alnah/go-tmdb-cli/internal/data"
	f "github.com/alnah/go-tmdb-cli/internal/fetch"
	q "github.com/alnah/go-tmdb-cli/internal/query"

	"github.com/olekukonko/tablewriter"
)

func main() {
	token, err := c.GetTMDBToken(".env")
	if err != nil {
		fmt.Println(err.Error())
	}

	query := &q.QueryParams{
		MaxItems: 100,
		Language: "en",
		Genres:   &q.Genres{"comedy", "horror"},
		Date: &q.Date{
			StartDate: "1888-01-01", StartOption: "gte",
			EndDate: "2025-01-26", EndOption: "lte",
		},
		Average: &q.Average{
			StartAverage: 7, StartOption: "gte",
		},
		Vote: &q.Vote{
			StartVote: 100, StartOption: "gte",
		},
	}

	movies, err := f.FetchMovies(*query, token, 50)
	if err != nil {
		fmt.Println(err.Error())
	}

	sortedMovies, err := movies.Sort()
	if err != nil {
		fmt.Println(err.Error())
	}

	table := renderTable(sortedMovies)

	if len(movies) == 0 {
		fmt.Println("No results.")
	} else {
		fmt.Println(table)
	}
}

func renderTable(results d.Movies) string {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{
		"Title",
		"Original Title",
		"Date",
		"Average",
		"Votes",
	})

	for _, r := range results {
		table.Append([]string{
			r.Title,
			r.OriginalTitle,
			r.Date,
			fmt.Sprintf("%.1f", r.Average),
			fmt.Sprintf("%d", r.Votes),
		})
	}

	table.Render()
	return buffer.String()
}
