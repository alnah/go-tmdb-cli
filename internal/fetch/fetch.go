package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	d "github.com/alnah/go-tmdb-cli/internal/data"
	q "github.com/alnah/go-tmdb-cli/internal/query"
)

type Response struct {
	Movies      d.Movies `json:"results"`
	Page        int      `json:"page"`
	TotalPages  int      `json:"total_pages"`
	TotalMovies int      `json:"total_results"`
}

func FetchMovies(query q.Query, token string, maxItems int) (d.Movies, error) {
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
	var movies d.Movies
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

		movies = append(movies, data.Movies...)
		if len(movies) > maxItems || data.Page >= data.TotalPages {
			break
		}

		data.Page = data.Page + 1
		url, err = query.BuildQuery()
		if err != nil {
			return d.Movies{}, err
		}

		req, err = newRequest("GET", url, token)
		if err != nil {
			return d.Movies{}, err
		}
	}

	if len(movies) > maxItems {
		movies = movies[:maxItems]
	}

	return movies, nil
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
