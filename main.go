package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	c "github.com/alnah/go-tmdb-cli/internal/config"
	d "github.com/alnah/go-tmdb-cli/internal/data"
	q "github.com/alnah/go-tmdb-cli/internal/query"

	"github.com/olekukonko/tablewriter"
)

type Response struct {
	Movies      d.Movies `json:"results"`
	Page        int      `json:"page"`
	TotalPages  int      `json:"total_pages"`
	TotalMovies int      `json:"total_results"`
}

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

	movies, err := fetch(*query, token)
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

func fetch(query q.QueryParams, token string) (d.Movies, error) {
	if query.Page == 0 {
		query.Page = 1
	}

	if query.MaxItems == 0 {
		query.MaxItems = 50
	}

	var url string
	url, err := query.BuildQuery()
	if err != nil {
		return d.Movies{}, err
	}

	req, err := newRequest("GET", url, token)
	if err != nil {
		return d.Movies{}, err
	}

	var data Response
	var allResults d.Movies
	for {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return d.Movies{}, err
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		if err = json.Unmarshal(body, &data); err != nil {
			return d.Movies{}, err
		}

		allResults = append(allResults, data.Movies...)
		if len(allResults) > query.MaxItems || query.Page > data.TotalPages {
			break
		}

		query.Page++
		url, _ = query.BuildQuery()
		req, err = newRequest("GET", url, token)
		if err != nil {
			return d.Movies{}, err
		}
	}

	if len(allResults) > query.MaxItems {
		allResults = allResults[:query.MaxItems]
	}

	return allResults, nil
}

func newRequest(method, url, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}
