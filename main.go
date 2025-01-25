package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	i "github.com/alnah/go-tmdb-cli/internal"
	"github.com/olekukonko/tablewriter"
)

type Result struct {
	Id            int     `json:"id"`
	Date          string  `json:"release_date"`
	OriginalTitle string  `json:"original_title"`
	Title         string  `json:"title"`
	Average       float64 `json:"vote_average"`
	Votes         int     `json:"vote_count"`
}

type Response struct {
	Page         int      `json:"page"`
	Results      []Result `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

func main() {
	token, err := i.GetTMDBToken(".env")
	if err != nil {
		fmt.Println(err.Error())
	}

	query := &i.QueryParams{
		Page:    toPtr(1),
		Year:    toPtr(2000),
		Average: &i.Average{Value: 7.0, Option: "gte"},
		Vote:    &i.Vote{Value: 1000, Option: "gte"},
	}

	url, err := query.BuildURL()
	if err != nil {
		fmt.Println(err.Error())
	}

	response := fetch(url, token)
	table := renderTable(response)

	if len(response.Results) == 0 {
		fmt.Println("No results.")
	} else {
		fmt.Println(table)
	}
}

func renderTable(response Response) string {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{
		"Title",
		"Original Title",
		"Date",
		"Average",
		"Votes",
	})

	for _, r := range response.Results {
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

func fetch(url string, token string) Response {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var data Response
	if err = json.Unmarshal(body, &data); err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func toPtr[T any](value T) *T {
	return &value
}
