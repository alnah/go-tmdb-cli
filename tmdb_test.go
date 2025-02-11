package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestUnitDeduplicate(t *testing.T) {
	fakeMovies := movies{
		fakeMovieList[0], fakeMovieList[0],
		fakeMovieList[1], fakeMovieList[1],
		fakeMovieList[2], fakeMovieList[2],
	}
	expected := movies{fakeMovieList[0], fakeMovieList[1], fakeMovieList[2]}
	got := fakeMovies.deduplicate()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected movies %+v, got %+v", expected, got)
	}
}

func TestUnitSortByField(t *testing.T) {
	fakeMovies := movies{fakeMovieList[0], fakeMovieList[1], fakeMovieList[2]}

	testCases := []struct {
		name    string
		param   string
		want    movies
		wantErr bool
	}{
		{
			name:  "sort by release date field ascending order",
			param: "date,asc",
			want:  movies{fakeMovieList[0], fakeMovieList[1], fakeMovieList[2]},
		},
		{
			name:  "sort by release date field descending order",
			param: "date,desc",
			want:  movies{fakeMovies[2], fakeMovies[1], fakeMovies[0]},
		},
		{
			name:  "sort by original title field ascending order",
			param: "otitle,asc",
			want:  movies{fakeMovieList[0], fakeMovieList[2], fakeMovieList[1]},
		},
		{
			name:  "sort by original title field descending order",
			param: "otitle,desc",
			want:  movies{fakeMovieList[1], fakeMovieList[2], fakeMovieList[0]},
		},
		{
			name:  "sort by title field ascending order",
			param: "title,asc",
			want:  movies{fakeMovieList[2], fakeMovieList[0], fakeMovieList[1]},
		},
		{
			name:  "sort by title field descending order",
			param: "title,desc",
			want:  movies{fakeMovieList[1], fakeMovieList[0], fakeMovieList[2]},
		},
		{
			name:  "sort by vote average field ascending order",
			param: "average,asc",
			want:  movies{fakeMovieList[1], fakeMovieList[0], fakeMovieList[2]},
		},
		{
			name:  "sort by vote average field descending order",
			param: "average,desc",
			want:  movies{fakeMovieList[2], fakeMovieList[0], fakeMovieList[1]},
		},
		{
			name:  "sort by vote count field ascending order",
			param: "votes,asc",
			want:  movies{fakeMovieList[1], fakeMovieList[0], fakeMovieList[2]},
		},
		{
			name:  "sort by vote count field descending order",
			param: "votes,desc",
			want:  movies{fakeMovieList[2], fakeMovieList[0], fakeMovieList[1]},
		},
		{
			name:    "invalid field",
			param:   "invalid,asc", // It could be any valid order
			wantErr: true,
		},
		{
			name:    "invalid order",
			param:   "title,invalid", // It could be any valid field
			wantErr: true,
		},
		{
			name:    "length value too big",
			param:   "too,big,value",
			wantErr: true,
		},
		{
			name:    "length value too small",
			param:   "toosmallvalue",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := fakeMovies.sortByField(tc.param)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				if !reflect.DeepEqual(tc.want, got) {
					t.Errorf("expected movies %+v, but got %+v", tc.want, got)
				}
			}
		})
	}
}

func TestUnitList(t *testing.T) {
	testCases := []struct {
		name    string
		param   string
		want    string
		wantErr bool
	}{
		{
			name:  "now playing",
			param: "now_playing",
			want:  "https://api.themoviedb.org/3/movie/now_playing?",
		},
		{
			name:  "popular",
			param: "popular",
			want:  "https://api.themoviedb.org/3/movie/popular?",
		},
		{
			name:  "top rated",
			param: "top_rated",
			want:  "https://api.themoviedb.org/3/movie/top_rated?",
		},
		{
			name:  "upcoming",
			param: "upcoming",
			want:  "https://api.themoviedb.org/3/movie/upcoming?",
		},
		{
			name:    "invalid param",
			param:   "invalid",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			builder := newURLBuilder()
			// Act
			got, err := builder.list(tc.param)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				assertURL(t, tc.want, got)
			}
		})
	}
}

func TestUnitDiscover(t *testing.T) {
	testCases := []struct {
		name    string
		query   queryParams
		want    string
		wantErr bool
	}{
		// Language
		{
			name: "valid iso language code",
			query: queryParams{
				Language: "fr", // it could be any ISO 639-1 code
			},
			want: "https://api.themoviedb.org/3/discover/movie?with_original_language=fr",
		},
		{
			name: "invalid iso language code too long",
			query: queryParams{
				Language: "aaa", // Not a two-letter ISO 639-1 code
			},
			wantErr: true,
		},
		{
			name: "invalid iso language code too short",
			query: queryParams{
				Language: "a", // Not a two-letter ISO 639-1 code
			},
			wantErr: true,
		},
		// Year(s)
		{
			name: "valid primary release year",
			query: queryParams{
				Year: "2000",
			},
			want: "https://api.themoviedb.org/3/discover/movie?primary_release_year=2000",
		},
		{
			name: "valid primary release dates",
			query: queryParams{
				Year: "2000,2010",
			},
			want: "https://api.themoviedb.org/3/discover/movie?primary_release_date.gte=2000-01-01&primary_release_date.lte=2010-12-31",
		},
		{
			name: "valid primary release date gte",
			query: queryParams{
				Year: "2000,gte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?primary_release_date.gte=2000-01-01",
		},
		{
			name: "valid primary release date lte",
			query: queryParams{
				Year: "2000,lte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?primary_release_date.lte=2000-01-01",
		},
		{
			name: "invalid non numeric primary release year",
			query: queryParams{
				Year: "abcd",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release year format",
			query: queryParams{
				Year: "1", // Can't be parsed as YYYY
			},
			wantErr: true,
		},
		{
			name: "invalid primary release year only comma",
			query: queryParams{
				Year: ",",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release year below min",
			query: queryParams{
				Year: "1887", // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid primary release year above current year",
			query: queryParams{
				Year: fmt.Sprintf("%d", time.Now().Year()+1), // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates first year not in valid yyyy format",
			query: queryParams{
				Year: "1,2000",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates second year not in yyyy format",
			query: queryParams{
				Year: "2000,1",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates first year below min",
			query: queryParams{
				Year: "1000,2000", // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates second year below min",
			query: queryParams{
				Year: "2000,1000", // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates first year value",
			query: queryParams{
				Year: "abcd,2000",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates second year value",
			query: queryParams{
				Year: "2000,abcd",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates both years value",
			query: queryParams{
				Year: "/",
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates first year value above current year",
			query: queryParams{
				Year: fmt.Sprintf("%d,2000", time.Now().Year()+1), // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid primary release dates second year value above current",
			query: queryParams{
				Year: fmt.Sprintf("2000,%d", time.Now().Year()+1), // No existing movies for this year
			},
			wantErr: true,
		},
		{
			name: "invalid year length value too big",
			query: queryParams{
				Year: "too,big,value",
			},
			wantErr: true,
		},
		// Vote Average
		{
			name: "valid vote average gte",
			query: queryParams{
				VoteAverage: "8.0,gte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_average.gte=8.0",
		},
		{
			name: "valid vote average lte",
			query: queryParams{
				VoteAverage: "8.0,lte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_average.lte=8.0",
		},
		{
			name: "valid vote average range",
			query: queryParams{
				VoteAverage: "7.0,8.0",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_average.gte=7.0&vote_average.lte=8.0",
		},
		{
			name: "invalid vote average only comma",
			query: queryParams{
				VoteAverage: ",",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average alone",
			query: queryParams{
				VoteAverage: "8.0",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average with comma",
			query: queryParams{
				VoteAverage: "8.0,",
			},
			wantErr: true,
		},
		{
			name: "invalid non numeric vote average",
			query: queryParams{
				VoteAverage: "abcd,gte",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average below min",
			query: queryParams{
				VoteAverage: "-0.1,gte", // Min is 0.0
			},
			wantErr: true,
		},
		{
			name: "invalid vote average above max",
			query: queryParams{
				VoteAverage: "10.1,gte", // Max is 10.1
			},
			wantErr: true,
		},
		{
			name: "invalid vote average first value",
			query: queryParams{
				VoteAverage: "xyz,8.0",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average second value",
			query: queryParams{
				VoteAverage: "1.0,xyz",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average length value too big",
			query: queryParams{
				VoteCount: "too,big,value",
			},
			wantErr: true,
		},
		{
			name: "invalid vote average length value too small",
			query: queryParams{
				VoteCount: "toosmallvalue",
			},
			wantErr: true,
		},
		// Vote Count
		{
			name: "valid vote count gte",
			query: queryParams{
				VoteCount: "1000,gte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.gte=1000",
		},
		{
			name: "valid vote count lte",
			query: queryParams{
				VoteCount: "1000,lte",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.lte=1000",
		},
		{
			name: "valid vote count range",
			query: queryParams{
				VoteCount: "500,1000",
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.gte=500&vote_count.lte=1000",
		},
		{
			name: "valid vote count gte with trailin quotes",
			query: queryParams{
				VoteCount: `"1000,gte"`,
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.gte=1000",
		},
		{
			name: "valid vote count lte with trailin quotes",
			query: queryParams{
				VoteCount: `"1000,lte"`,
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.lte=1000",
		},
		{
			name: "valid vote count range with trailin quotes",
			query: queryParams{
				VoteCount: `"500,1000"`,
			},
			want: "https://api.themoviedb.org/3/discover/movie?vote_count.gte=500&vote_count.lte=1000",
		},
		{
			name: "invalid vote count only comma",
			query: queryParams{
				VoteCount: ",",
			},
			wantErr: true,
		},
		{
			name: "invalid vote count alone",
			query: queryParams{
				VoteCount: "1000",
			},
			wantErr: true,
		},
		{
			name: "invalid vote count with comma",
			query: queryParams{
				VoteCount: "1000,",
			},
			wantErr: true,
		},
		{
			name: "invalid non numeric vote count",
			query: queryParams{
				VoteCount: "abcd,gte",
			},
			wantErr: true,
		},
		{
			name: "invalid vote count below min",
			query: queryParams{
				VoteCount: "-1,gte", // Min is 0
			},
			wantErr: true,
		},
		{
			name: "invalid vote count first value",
			query: queryParams{
				VoteCount: "xyz,1000",
			},
			wantErr: true,
		},
		{
			name: "invalid vote count second value",
			query: queryParams{
				VoteCount: "1000,xyz",
			},
			wantErr: true,
		},
		{
			name: "invalid vote count length value too big",
			query: queryParams{
				VoteCount: "too,big,value",
			},
			wantErr: true,
		},
		// With Genres
		{
			name: "one valid with genre",
			query: queryParams{
				WithGenres: "drama",
			},
			want: "https://api.themoviedb.org/3/discover/movie?with_genres=18",
		},
		{
			name: "many valid with genre",
			query: queryParams{
				WithGenres: "drama,history",
			},
			want: "https://api.themoviedb.org/3/discover/movie?with_genres=18,36",
		},
		{
			name: "one invalid with genre",
			query: queryParams{
				WithGenres: "invalid",
			},
			wantErr: true,
		},
		{
			name: "one invalid with genre in many with genres",
			query: queryParams{
				WithGenres: "drama,invalid",
			},
			wantErr: true,
		},
		// Without Genres
		{
			name: "one valid without genre",
			query: queryParams{
				WithoutGenres: "drama",
			},
			want: "https://api.themoviedb.org/3/discover/movie?without_genres=18",
		},
		{
			name: "many valid without genre",
			query: queryParams{
				WithoutGenres: "drama,history",
			},
			want: "https://api.themoviedb.org/3/discover/movie?without_genres=18,36",
		},
		{
			name: "one invalid without genre",
			query: queryParams{
				WithoutGenres: "invalid",
			},
			wantErr: true,
		},
		{
			name: "one invalid without genre in many without genres",
			query: queryParams{
				WithoutGenres: "drama,invalid",
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			urlBuilder := newURLBuilder()
			// Act
			got, err := urlBuilder.discover(tc.query)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				assertURL(t, tc.want, got)
			}
		})
	}
}

func TestUniFetchTMDBResponse(t *testing.T) {
	testCases := []struct {
		name           string
		apiKey         string
		wantErr        bool
		wantClientEr   bool
		wantServerErr  bool
		wantNetworkErr bool
		wantRequestErr bool
		wantJSONErr    bool
	}{
		{name: "get results"},
		{name: "client error", apiKey: "invalid_api_key", wantErr: true, wantClientEr: true},
		{name: "server error", wantErr: true, wantServerErr: true},
		{name: "network error", wantErr: true, wantNetworkErr: true},
		{name: "request error", wantErr: true, wantRequestErr: true},
		{name: "JSON format error", wantErr: true, wantJSONErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var (
				hc      *httpClient
				tmdbRes tmdbResponse
				err     error
			)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.wantClientEr {
					if r.Header.Get("Authorization") != "Bearer valid_api_key" {
						w.WriteHeader(401)
						return
					}
				}
				if tc.wantServerErr {
					w.WriteHeader(501)
					w.Write([]byte(`{"error": "Invalid service"}`))
					return
				}
				if tc.wantJSONErr {
					w.Write([]byte("Invalid JSON format"))
					return
				}
				byt, _ := json.Marshal(fakeResPage1)
				w.Write(byt)
			}))
			t.Cleanup(ts.Close)
			hc = newHTTPClient(tc.apiKey)
			// Act
			if tc.wantRequestErr {
				tmdbRes, err = fetchTMDBResponse(hc, ":invalid_url")
			} else if tc.wantNetworkErr {
				tmdbRes, err = fetchTMDBResponse(hc, "http://0.0.0.0:9999") // Non-routable IP
			} else {
				tmdbRes, err = fetchTMDBResponse(hc, ts.URL)
			}
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				assertResponse(t, fakeResPage1, tmdbRes)
			}
		})
	}
}

func TestUnitTestUniTFetchTMDBResponse_Retry(t *testing.T) {
	// Arrange
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requireAPIKey(t, w, r)
		attempts++
		if attempts < 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(429)
			return
		}
		byt, _ := json.Marshal(fakeResPage1)
		w.Write(byt)
	}))
	t.Cleanup(ts.Close)
	hc := newHTTPClient("valid_api_key")
	// Act
	tmdbRes, err := fetchTMDBResponse(hc, ts.URL)
	// Assert
	assertNoError(t, err)
	assertResponse(t, fakeResPage1, tmdbRes)
}

func TestUnitAsyncFetchMovies(t *testing.T) {
	testCases := []struct {
		name     string
		maxItems int
		want     movies
		wantErr  bool
	}{
		{
			name:     "first page",
			maxItems: 20,
			want:     fakeMovieList[:20],
		},
		{
			name:     "all pages",
			maxItems: 40,
			want:     fakeMovieList,
		},
		{
			name:     "first 10 results",
			maxItems: 10,
			want:     fakeMovieList[:10],
		},
		{
			name:     "first 30 results",
			maxItems: 30,
			want:     fakeMovieList[:30],
		},
		{
			name:     "max items exceeds API max items",
			maxItems: APIMaxItems + 1,
			wantErr:  true,
		},
		{
			name:     "propagate error",
			maxItems: 40,
			wantErr:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requireAPIKey(t, w, r)
				page := r.URL.Query().Get("page")
				switch {
				case tc.wantErr && page == "2":
					byt, _ := json.Marshal([]byte("invalid"))
					w.Write(byt)
					return
				case page == "1":
					byt, _ := json.Marshal(fakeResPage1)
					w.Write(byt)
					return
				case page == "2":
					byt, _ := json.Marshal(fakeResPage2)
					w.Write(byt)
					return
				default:
					http.NotFound(w, r)
				}
			}))
			t.Cleanup(func() { ts.Close() })
			hc := newHTTPClient("valid_api_key")
			// Act
			got, err := asyncFetchMovies(hc, ts.URL+"?", tc.maxItems)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				assertMovies(t, tc.want, got)
			}
		})
	}
}

func BenchmarkAsyncFetchMovies(b *testing.B) {
	testCases := []struct {
		maxItems int
	}{
		{maxItems: 20},
		{maxItems: 40},
		{maxItems: 10},
		{maxItems: 30},
	}
	hc := newHTTPClient("valid_api_key")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byt, _ := json.Marshal(fakeResPage1)
		w.Write(byt)
	}))
	defer ts.Close()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_, err := asyncFetchMovies(hc, ts.URL+"?", tc.maxItems)
			if err != nil {
				b.Fatalf("failed to fetch movies: %v", err)
			}
		}
	}
}
