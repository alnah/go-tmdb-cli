package fetch_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	d "github.com/alnah/go-tmdb-cli/internal/data"
	f "github.com/alnah/go-tmdb-cli/internal/fetch"
	tf "github.com/alnah/go-tmdb-cli/test/fake"
	"github.com/stretchr/testify/assert"
)

const (
	fakeToken     = "test_token"
	moviesPerPage = 20
)

var page1 = f.Response{
	Movies:      tf.FakeMovies[:moviesPerPage],
	Page:        1,
	TotalPages:  2,
	TotalMovies: 40,
}

func TestUnitFetchMovies(t *testing.T) {
	testCases := []struct {
		name     string
		query    *MockQueryParams
		maxItems int
		want     d.Movies
	}{
		{
			name:     "return first page results",
			query:    &MockQueryParams{t: t, fetchAll: false},
			maxItems: 50,
			want:     page1.Movies,
		},
		{
			name:     "return multi-page results",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 50,
			want:     tf.FakeMovies[:moviesPerPage*2],
		},
		{
			name:     "return results limited to max items",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 15,
			want:     tf.FakeMovies[:15],
		},
		{
			name:     "return partial final page",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 30,
			want:     tf.FakeMovies[:30],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			movies, err := f.FetchMovies(tc.query, fakeToken, tc.maxItems)

			assert.NoError(t, err)
			assert.Len(t, movies, len(tc.want))
			assert.Equal(t, tc.want, movies)
		})
	}
}

type MockQueryParams struct {
	t           *testing.T
	fetchAll    bool
	currentPage int
}

func (m *MockQueryParams) BuildQuery() (string, error) {
	m.currentPage++
	var res f.Response

	switch {
	case m.fetchAll && m.currentPage <= 2:
		res = f.Response{
			Movies:      tf.FakeMovies[(m.currentPage-1)*moviesPerPage : m.currentPage*moviesPerPage],
			Page:        m.currentPage,
			TotalPages:  2,
			TotalMovies: 40,
		}
	case !m.fetchAll:
		res = page1
		res.TotalPages = 1 // Force single page result
	default:
		m.t.Fatal("Unexpected extra page request")
	}

	server := newServer(m.t, res)
	return server.URL, nil
}

func newServer(t *testing.T, res f.Response) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer "+fakeToken, r.Header.Get("Authorization"))

			byt, err := json.Marshal(res)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(byt)
			assert.NoError(t, err)
		}))
}
