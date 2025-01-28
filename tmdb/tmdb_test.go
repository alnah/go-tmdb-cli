package tmdb_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alnah/go-tmdb-cli/tmdb"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	fakeToken = "test_token"
	perPage   = 20
)

var (
	fakeMovies = tmdb.MovieCollection{
		tmdb.Movie{
			Id:            1,
			Date:          "2023-01-01",
			OriginalTitle: "L'Aube de l'Aventure",
			Title:         "Epic Journey Begins",
			Average:       8.5,
			Votes:         100,
		},
		tmdb.Movie{
			Id:            2,
			Date:          "2023-02-01",
			OriginalTitle: "Rise of the Heroes",
			Title:         "Rise of the Heroes",
			Average:       7.0,
			Votes:         50,
		},
		tmdb.Movie{
			Id:            3,
			Date:          "2023-03-01",
			OriginalTitle: "O Confronto Final",
			Title:         "Clash of Titans",
			Average:       9.0,
			Votes:         200,
		},
		tmdb.Movie{
			Id:            4,
			Date:          "2023-04-01",
			OriginalTitle: "The Quest for Knowledge",
			Title:         "À la recherche du savoir",
			Average:       8.0,
			Votes:         75,
		},
		tmdb.Movie{
			Id:            5,
			Date:          "2023-05-01",
			OriginalTitle: "A Ascensão da Inovação",
			Title:         "Rise of Innovation",
			Average:       7.5,
			Votes:         120,
		},
		tmdb.Movie{
			Id:            6,
			Date:          "2023-06-01",
			OriginalTitle: "L'Aube d'une Nouvelle Ère",
			Title:         "A New Dawn",
			Average:       8.2,
			Votes:         90,
		},
		tmdb.Movie{
			Id:            7,
			Date:          "2023-07-01",
			OriginalTitle: "A Última Resistência",
			Title:         "The Final Stand",
			Average:       9.1,
			Votes:         150,
		},
		tmdb.Movie{
			Id:            8,
			Date:          "2023-08-01",
			OriginalTitle: "La Grande Aventure",
			Title:         "The Great Adventure",
			Average:       8.7,
			Votes:         200,
		},
		tmdb.Movie{
			Id:            9,
			Date:          "2023-09-01",
			OriginalTitle: "A Jornada Continua",
			Title:         "The Journey Continues",
			Average:       7.8,
			Votes:         110,
		},
		tmdb.Movie{
			Id:            10,
			Date:          "2023-10-01",
			OriginalTitle: "L'Héritage des Héros",
			Title:         "The Legacy of Heroes",
			Average:       9.3,
			Votes:         250,
		},
		tmdb.Movie{
			Id:            11,
			Date:          "2023-11-01",
			OriginalTitle: "O Poder da Unidade",
			Title:         "The Power of Unity",
			Average:       8.4,
			Votes:         130,
		},
		tmdb.Movie{
			Id:            12,
			Date:          "2023-12-01",
			OriginalTitle: "L'Appel à l'Aventure",
			Title:         "The Call of Adventure",
			Average:       8.9,
			Votes:         180,
		},
		tmdb.Movie{
			Id:            13,
			Date:          "2024-01-01",
			OriginalTitle: "A Ascensão da Fênix",
			Title:         "The Rise of the Phoenix",
			Average:       9.0,
			Votes:         220,
		},
		tmdb.Movie{
			Id:            14,
			Date:          "2024-02-01",
			OriginalTitle: "Le Voyage d'une Vie",
			Title:         "The Journey of a Lifetime",
			Average:       7.6,
			Votes:         95,
		},
		tmdb.Movie{
			Id:            15,
			Date:          "2024-03-01",
			OriginalTitle: "O Coração de um Campeão",
			Title:         "The Heart of a Champion",
			Average:       8.3,
			Votes:         140,
		},
		tmdb.Movie{
			Id:            16,
			Date:          "2024-04-01",
			OriginalTitle: "L'Esprit d'Aventure",
			Title:         "The Spirit of Adventure",
			Average:       8.1,
			Votes:         115,
		},
		tmdb.Movie{
			Id:            17,
			Date:          "2024-05-01",
			OriginalTitle: "O Legado dos Bravos",
			Title:         "The Legacy of the Brave",
			Average:       9.2,
			Votes:         260,
		},
		tmdb.Movie{
			Id:            18,
			Date:          "2024-06-01",
			OriginalTitle: "Le Pouvoir des Rêves",
			Title:         "The Power of Dreams",
			Average:       8.6,
			Votes:         170,
		},
		tmdb.Movie{
			Id:            19,
			Date:          "2024-07-01",
			OriginalTitle: "O Chamado do Destino",
			Title:         "The Call of Destiny",
			Average:       9.4,
			Votes:         300,
		},
		tmdb.Movie{
			Id:            20,
			Date:          "2024-08-01",
			OriginalTitle: "Le Voyage de l'Espoir",
			Title:         "The Journey of Hope",
			Average:       8.8,
			Votes:         190,
		},
		tmdb.Movie{
			Id:            21,
			Date:          "2024-09-01",
			OriginalTitle: "A Ascensão dos Titãs",
			Title:         "The Rise of the Titans",
			Average:       9.1,
			Votes:         240,
		},
		tmdb.Movie{
			Id:            22,
			Date:          "2024-10-01",
			OriginalTitle: "Le Cœur de l'Océan",
			Title:         "The Heart of the Ocean",
			Average:       8.5,
			Votes:         160,
		},
		tmdb.Movie{
			Id:            23,
			Date:          "2024-11-01",
			OriginalTitle: "O Espírito da Floresta",
			Title:         "The Spirit of the Forest",
			Average:       8.0,
			Votes:         130,
		},
		tmdb.Movie{
			Id:            24,
			Date:          "2024-12-01",
			OriginalTitle: "L'Appel de la Nature",
			Title:         "The Call of the Wild",
			Average:       9.0,
			Votes:         210,
		},
		tmdb.Movie{
			Id:            25,
			Date:          "2025-01-01",
			OriginalTitle: "A Jornada das Estrelas",
			Title:         "The Journey of the Stars",
			Average:       8.4,
			Votes:         175,
		},
		tmdb.Movie{
			Id:            26,
			Date:          "2025-02-01",
			OriginalTitle: "La Montée des Gardiens",
			Title:         "The Rise of the Guardians",
			Average:       9.3,
			Votes:         280,
		},
		tmdb.Movie{
			Id:            27,
			Date:          "2025-03-01",
			OriginalTitle: "O Coração do Guerreiro",
			Title:         "The Heart of the Warrior",
			Average:       8.7,
			Votes:         150,
		},
		tmdb.Movie{
			Id:            28,
			Date:          "2025-04-01",
			OriginalTitle: "L'Esprit du Ciel",
			Title:         "The Spirit of the Sky",
			Average:       8.9,
			Votes:         165,
		},
		tmdb.Movie{
			Id:            29,
			Date:          "2025-05-01",
			OriginalTitle: "A Lenda do Dragão",
			Title:         "The Legend of the Dragon",
			Average:       9.5,
			Votes:         300,
		},
		tmdb.Movie{
			Id:            30,
			Date:          "2025-06-01",
			OriginalTitle: "Le Dernier Voyage",
			Title:         "The Last Voyage",
			Average:       8.3,
			Votes:         140,
		},
		tmdb.Movie{
			Id:            31,
			Date:          "2025-07-01",
			OriginalTitle: "O Caminho da Esperança",
			Title:         "The Path of Hope",
			Average:       8.6,
			Votes:         175,
		},
		tmdb.Movie{
			Id:            32,
			Date:          "2025-08-01",
			OriginalTitle: "La Lumière de l'Aube",
			Title:         "The Light of Dawn",
			Average:       9.0,
			Votes:         200,
		},
		tmdb.Movie{
			Id:            33,
			Date:          "2025-09-01",
			OriginalTitle: "A Última Esperança",
			Title:         "The Last Hope",
			Average:       9.2,
			Votes:         220,
		},
		tmdb.Movie{
			Id:            34,
			Date:          "2025-10-01",
			OriginalTitle: "Le Voyage des Rêves",
			Title:         "The Journey of Dreams",
			Average:       8.8,
			Votes:         190,
		},
		tmdb.Movie{
			Id:            35,
			Date:          "2025-11-01",
			OriginalTitle: "O Chamado da Aventura",
			Title:         "The Call of Adventure",
			Average:       9.1,
			Votes:         210,
		},
		tmdb.Movie{
			Id:            36,
			Date:          "2025-12-01",
			OriginalTitle: "La Force de l'Amitié",
			Title:         "The Strength of Friendship",
			Average:       8.4,
			Votes:         160,
		},
		tmdb.Movie{
			Id:            37,
			Date:          "2026-01-01",
			OriginalTitle: "A Jornada do Guerreiro",
			Title:         "The Warrior's Journey",
			Average:       9.3,
			Votes:         250,
		},
		tmdb.Movie{
			Id:            38,
			Date:          "2026-02-01",
			OriginalTitle: "Le Dernier Héros",
			Title:         "The Last Hero",
			Average:       8.7,
			Votes:         180,
		},
		tmdb.Movie{
			Id:            39,
			Date:          "2026-03-01",
			OriginalTitle: "O Legado dos Heróis",
			Title:         "The Legacy of Heroes",
			Average:       9.0,
			Votes:         230,
		},
		tmdb.Movie{
			Id:            40,
			Date:          "2026-04-01",
			OriginalTitle: "La Quête de la Vérité",
			Title:         "The Quest for Truth",
			Average:       8.5,
			Votes:         200,
		},
	}

	fakeRes = tmdb.TMDBResponse{
		Results:     fakeMovies[:perPage],
		Page:        1,
		TotalPages:  2,
		TotalMovies: 40,
	}
)

func TestUnitValidationError(t *testing.T) {
	t.Run("return error string", func(t *testing.T) {
		err := &tmdb.ValidationError{Filter: "Test", Message: "description"}
		assert.Equal(t, "Test description.", err.Error())
	})
}

func TestUnitSortOrderError(t *testing.T) {
	t.Run("return error string", func(t *testing.T) {
		err := &tmdb.SortOrderError{}
		assert.Equal(t, `Sorting order must be "asc" or "desc".`, err.Error())
	})
}

func TestUnitBuildQuery(t *testing.T) {
	t.Run("return query", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *tmdb.QueryParams
			want  string
		}{
			{
				desc: "page year date average",
				query: &tmdb.QueryParams{
					Year:    2000,
					Dates:   &tmdb.Dates{StartDate: "2000-06-01", StartOption: "gte"},
					Average: &tmdb.Average{StartAverage: 8.0, StartOption: "gte"},
					Votes:   &tmdb.Votes{StartVote: 1000, StartOption: "gte"},
				},
				want: tmdb.APIBaseURL + tmdb.DiscoverPathURL +
					"primary_release_year=2000" +
					"&primary_release_date.gte=2000-06-01" +
					"&vote_average.gte=8.0&vote_count.gte=1000",
			},
			{
				desc: "popular movies list",
				query: &tmdb.QueryParams{
					MovieListPath: "popular",
				},
				want: fmt.Sprintf(tmdb.APIBaseURL+tmdb.MoviesListPathURL, "popular"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				got, err := tc.query.BuildQuery()
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("return error when error occured", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *tmdb.QueryParams
		}{
			{
				desc: "on discover",
				query: &tmdb.QueryParams{
					Dates: &tmdb.Dates{
						StartDate:   "1800-01-01", // movies don't exist at that time
						StartOption: "gte",
					},
				},
			},
			{
				desc: "on movies list",
				query: &tmdb.QueryParams{
					MovieListPath: "invalid",
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				var want *tmdb.ValidationError
				got, err := tc.query.BuildQuery()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			})
		}
	})
}

func TestUnitWithMoviesList(t *testing.T) {
	testCases := []struct {
		desc    string
		path    string
		wantErr bool
	}{
		{
			desc: "return now playing",
			path: "now_playing",
		},
		{
			desc: "return popular",
			path: "popular",
		},
		{
			desc: "return top rated",
			path: "top_rated",
		},
		{
			desc: "return upcoming",
			path: "upcoming",
		},
		{
			desc:    "return error when doesn't match wanted paths",
			path:    "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			switch {
			case tc.wantErr:
				query := tmdb.QueryParams{MovieListPath: tc.path}
				var want *tmdb.ValidationError
				got, err := query.WithMoviesList()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				query := tmdb.QueryParams{MovieListPath: tc.path}
				want := fmt.Sprintf(tmdb.MoviesListPathURL, tc.path)
				got, err := query.WithMoviesList()
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}

	t.Run("return empty string when no movies list", func(t *testing.T) {
		query := &tmdb.QueryParams{MovieListPath: ""}
		got, err := query.WithMoviesList()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitWithLanguage(t *testing.T) {
	testCases := []struct {
		desc    string
		iso     string
		want    string
		wantErr bool
	}{
		{
			desc: "return english",
			iso:  "en",
			want: "with_original_language=en",
		},
		{
			desc: "return french",
			iso:  "fr",
			want: "with_original_language=fr",
		},
		{
			desc: "return portuguese",
			iso:  "pt",
			want: "with_original_language=pt",
		},
		{
			desc: "return japanese",
			iso:  "jp",
			want: "with_original_language=jp",
		},
		{
			desc: "return mandarin",
			iso:  "zh",
			want: "with_original_language=zh",
		},
		{
			desc: "return guarani",
			iso:  "gn",
			want: "with_original_language=gn",
		},
		{
			desc:    "return error when not ISO 639-1 language code",
			iso:     "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &tmdb.QueryParams{Language: tc.iso}
			got, err := query.WithLanguage()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no language", func(t *testing.T) {
		query := &tmdb.QueryParams{Language: ""}
		got, err := query.WithLanguage()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitWithYear(t *testing.T) {
	testCases := []struct {
		desc    string
		value   int
		want    string
		wantErr bool
	}{
		{
			desc:  "return year",
			value: 2001,
			want:  "primary_release_year=2001",
		},
		{
			desc:    "return empty string when no year",
			wantErr: false,
		},
		{
			desc:    "return error when year out of range",
			value:   2147483648, // out of range
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &tmdb.QueryParams{Year: tc.value}
			got, err := query.WithYear()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitWithDates(t *testing.T) {
	testCases := []struct {
		desc    string
		date    tmdb.Dates
		want    string
		wantErr bool
	}{
		{
			desc: "return start date gte",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "gte",
			},
			want: "primary_release_date.gte=2001-01-01",
		},
		{
			desc: "return start date lte",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
			},
			want: "primary_release_date.lte=2001-01-01",
		},
		{
			desc: "return start date gte and end date lte",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "gte",
				EndDate:     "2002-02-02",
				EndOption:   "lte",
			},
			want: "primary_release_date.gte=2001-01-01&" +
				"primary_release_date.lte=2002-02-02",
		},
		{
			desc: "return start date lte and end date gte",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "gte",
			},
			want: "primary_release_date.lte=2001-01-01&" +
				"primary_release_date.gte=2002-02-02",
		},
		{
			desc: "return error when invalid start date format",
			date: tmdb.Dates{
				StartDate:   "01-01-2001", // invalid date format
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid start date range",
			date: tmdb.Dates{
				StartDate:   "1800-01-01", // movies don't exist at that time
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid start date option",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid end date format",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "01-01-2001", // invalid date format
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid end date range",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "1800-01-01", // movies don't exist at that time
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid end date option",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when same options",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return error when same dates",
			date: tmdb.Dates{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2001-01-01",
				EndOption:   "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &tmdb.QueryParams{Dates: &tc.date}
			got, err := query.WithDates()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &tmdb.QueryParams{Dates: nil}
		got, err := query.WithDates()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitWithAverage(t *testing.T) {
	testCases := []struct {
		desc    string
		average tmdb.Average
		want    string
		wantErr bool
	}{
		{
			desc: "return start average gte",
			average: tmdb.Average{
				StartAverage: 8.0,
				StartOption:  "gte",
			},
			want: "vote_average.gte=8.0",
		},
		{
			desc: "return start average lte",
			average: tmdb.Average{
				StartAverage: 8.0,
				StartOption:  "lte",
			},
			want: "vote_average.lte=8.0",
		},
		{
			desc: "return start average gte and end average lte",
			average: tmdb.Average{
				StartAverage: 7.0,
				StartOption:  "gte",
				EndAverage:   8.0,
				EndOption:    "lte",
			},
			want: "vote_average.gte=7.0&" +
				"vote_average.lte=8.0",
		},
		{
			desc: "return start average lte and end average gte",
			average: tmdb.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "gte",
			},
			want: "vote_average.lte=7.0&" +
				"vote_average.gte=8.0",
		},
		{
			desc: "return error when invalid start option",
			average: tmdb.Average{
				StartAverage: 8.0,
				StartOption:  "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when start average out of range",
			average: tmdb.Average{
				StartAverage: 100.0,
				StartOption:  "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when end average out of range",
			average: tmdb.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   100.0, // invalid average
				EndOption:    "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid end option",
			average: tmdb.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when same options",
			average: tmdb.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "lte",
			},
			wantErr: true,
		},

		{
			desc: "return error when same averages",
			average: tmdb.Average{
				StartAverage: 8.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &tmdb.QueryParams{Average: &tc.average}
			got, err := query.WithAverage()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no average", func(t *testing.T) {
		query := &tmdb.QueryParams{Average: nil}
		got, err := query.WithAverage()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitWithVotes(t *testing.T) {
	testCases := []struct {
		desc    string
		vote    tmdb.Votes
		want    string
		wantErr bool
	}{
		{
			desc: "return start vote gte",
			vote: tmdb.Votes{
				StartVote:   1000,
				StartOption: "gte",
			},
			want: "vote_count.gte=1000",
		},
		{
			desc: "return start vote lte",
			vote: tmdb.Votes{
				StartVote:   1000,
				StartOption: "lte",
			},
			want: "vote_count.lte=1000",
		},
		{
			desc: "return start vote gte and end vote lte",
			vote: tmdb.Votes{
				StartVote:   500,
				StartOption: "gte",
				EndVote:     1000,
				EndOption:   "lte",
			},
			want: "vote_count.gte=500&" +
				"vote_count.lte=1000",
		},
		{
			desc: "return start vote lte and end vote gte",
			vote: tmdb.Votes{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "gte",
			},
			want: "vote_count.lte=500&" +
				"vote_count.gte=1000",
		},
		{
			desc: "return error when invalid start option",
			vote: tmdb.Votes{
				StartVote:   1000,
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when start vote out of range",
			vote: tmdb.Votes{
				StartVote:   2147483648, // invalid vote
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when end vote out of range",
			vote: tmdb.Votes{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     2147483648, // invalid vote
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return error when invalid end option",
			vote: tmdb.Votes{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return error when same options",
			vote: tmdb.Votes{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return error when same vote",
			vote: tmdb.Votes{
				StartVote:   1000,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &tmdb.QueryParams{Votes: &tc.vote}
			got, err := query.WithVotes()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &tmdb.QueryParams{Average: nil}
		got, err := query.WithVotes()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetGenres(t *testing.T) {
	testCases := []struct {
		desc    string
		genres  tmdb.Genres
		want    string
		wantErr bool
	}{
		{
			desc:   "return one genre ID",
			genres: tmdb.Genres{"horror"},
			want:   "with_genres=27",
		},
		{
			desc:   "return many genre IDs",
			genres: tmdb.Genres{"horror", "comedy"},
			want:   "with_genres=27,35",
		},
		{
			desc:    "return error when not allowed genre",
			genres:  tmdb.Genres{"invalid"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := tmdb.QueryParams{Genres: &tc.genres}
			got, err := query.WithGenres()
			switch {
			case tc.wantErr:
				var want *tmdb.ValidationError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitDefaultSort(t *testing.T) {
	t.Run("return sorted movies", func(t *testing.T) {
		movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}
		got, err := movies.DefaultSort()
		assert.NoError(t, err)
		assert.Equal(t, movies, got)
	})
}

func TestUnitSortByDate(t *testing.T) {
	movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}

	testCases := []struct {
		description string
		order       string
		want        tmdb.MovieCollection
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        tmdb.MovieCollection{fakeMovies[2], fakeMovies[1], fakeMovies[0]},
		},
		{
			description: "error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := movies.SortByDate(tc.order)
			switch {
			case tc.wantErr:
				var want *tmdb.SortOrderError
				assert.Equal(t, movies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSortByOriginalTitle(t *testing.T) {
	movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}

	testCases := []struct {
		description string
		order       string
		want        tmdb.MovieCollection
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        tmdb.MovieCollection{fakeMovies[0], fakeMovies[2], fakeMovies[1]},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        tmdb.MovieCollection{fakeMovies[1], fakeMovies[2], fakeMovies[0]},
		},
		{
			description: "error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := movies.SortByOriginalTitle(tc.order)
			switch {
			case tc.wantErr:
				var want *tmdb.SortOrderError
				assert.Equal(t, movies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSortByTitle(t *testing.T) {
	movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}

	testCases := []struct {
		description string
		order       string
		want        tmdb.MovieCollection
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        tmdb.MovieCollection{fakeMovies[2], fakeMovies[0], fakeMovies[1]},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        tmdb.MovieCollection{fakeMovies[1], fakeMovies[0], fakeMovies[2]},
		},
		{
			description: "error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := movies.SortByTitle(tc.order)
			switch {
			case tc.wantErr:
				var want *tmdb.SortOrderError
				assert.Equal(t, movies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSortByAverage(t *testing.T) {
	movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}

	testCases := []struct {
		description string
		order       string
		want        tmdb.MovieCollection
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        tmdb.MovieCollection{fakeMovies[1], fakeMovies[0], fakeMovies[2]},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        tmdb.MovieCollection{fakeMovies[2], fakeMovies[0], fakeMovies[1]},
		},
		{
			description: "error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := movies.SortByAverage(tc.order)
			switch {
			case tc.wantErr:
				var want *tmdb.SortOrderError
				assert.Equal(t, movies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSortByVotes(t *testing.T) {
	movies := tmdb.MovieCollection{fakeMovies[0], fakeMovies[1], fakeMovies[2]}

	testCases := []struct {
		description string
		order       string
		want        tmdb.MovieCollection
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        tmdb.MovieCollection{fakeMovies[1], fakeMovies[0], fakeMovies[2]},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        tmdb.MovieCollection{fakeMovies[2], fakeMovies[0], fakeMovies[1]},
		},
		{
			description: "error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := movies.SortByVotes(tc.order)
			switch {
			case tc.wantErr:
				var want *tmdb.SortOrderError
				assert.Equal(t, movies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitFetchMovieCollection(t *testing.T) {
	testCases := []struct {
		name     string
		query    *MockQueryParams
		maxItems int
		want     tmdb.MovieCollection
	}{
		{
			name:     "return first page results",
			query:    &MockQueryParams{t: t, fetchAll: false},
			maxItems: 50,
			want:     fakeRes.Results,
		},
		{
			name:     "return multi-page results",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 50,
			want:     fakeMovies[:perPage*2],
		},
		{
			name:     "return results limited to max items",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 15,
			want:     fakeMovies[:15],
		},
		{
			name:     "return partial final page",
			query:    &MockQueryParams{t: t, fetchAll: true},
			maxItems: 30,
			want:     fakeMovies[:30],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			movies, err := tmdb.FetchMovieCollection(tc.query, fakeToken, tc.maxItems)

			assert.NoError(t, err)
			assert.Len(t, movies, len(tc.want))
			assert.Equal(t, tc.want, movies)
		})
	}
}

type MockQueryParams struct {
	t        *testing.T
	fetchAll bool
	page     int
}

func (m *MockQueryParams) BuildQuery() (string, error) {
	m.page++
	var res tmdb.TMDBResponse

	switch {
	case m.fetchAll && m.page <= 2:
		res = tmdb.TMDBResponse{
			Results:     fakeMovies[(m.page-1)*perPage : m.page*perPage],
			Page:        m.page,
			TotalPages:  2,
			TotalMovies: 40,
		}
	case !m.fetchAll:
		res = fakeRes
		res.TotalPages = 1 // Force single page result
	default:
		m.t.Fatal("Unexpected extra page request")
	}

	server := newServer(m.t, res)
	return server.URL, nil
}

func newServer(t *testing.T, res tmdb.TMDBResponse) *httptest.Server {
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

func TestUnitGetTMDBToken(t *testing.T) {
	t.Run("return TMDB API Token from env file", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`TOKEN="test"`)
		assert.NoError(t, err)

		got, err := tmdb.GetTMDBToken(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test", got)
	})

	t.Run("return error when env file doesn't exist", func(t *testing.T) {
		var want *fs.PathError
		_, err := tmdb.GetTMDBToken("test.env")
		assert.ErrorAs(t, err, &want)
	})

	t.Run("return config parse error when invalid data", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`invalid_data`)
		assert.NoError(t, err)

		var want viper.ConfigParseError
		_, err = tmdb.GetTMDBToken(file.Name())
		assert.ErrorAs(t, err, &want)
	})
}

func createTempFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", "test.env")
	assert.NoError(t, err)
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}
