package main

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var (
	// MovieList provides realistic sample data spanning different vote ranges and dates to
	// test table formatting edge cases. The 20-item dataset verifies pagination handling
	// and output truncation behavior without external dependencies.
	fakeMovieList = movies{
		{
			ID:            1,
			OriginalTitle: "L'Aube de l'Aventure",
			ReleaseDate:   "2023-01-01",
			Title:         "Epic Journey Begins",
			VoteAverage:   8.5,
			VoteCount:     100,
		},
		{
			ID:            2,
			OriginalTitle: "Rise of the Heroes",
			ReleaseDate:   "2023-02-01",
			Title:         "Rise of the Heroes",
			VoteAverage:   7.0,
			VoteCount:     50,
		},
		{
			ID:            3,
			OriginalTitle: "O Confronto Final",
			ReleaseDate:   "2023-03-01",
			Title:         "Clash of Titans",
			VoteAverage:   9.0,
			VoteCount:     200,
		},
		{
			ID:            4,
			OriginalTitle: "The Quest for Knowledge",
			ReleaseDate:   "2023-04-01",
			Title:         "À la recherche du savoir",
			VoteAverage:   8.0,
			VoteCount:     75,
		},
		{
			ID:            5,
			OriginalTitle: "A Ascensão da Inovação",
			ReleaseDate:   "2023-05-01",
			Title:         "Rise of Innovation",
			VoteAverage:   7.5,
			VoteCount:     120,
		},
		{
			ID:            6,
			OriginalTitle: "L'Aube d'une Nouvelle Ère",
			ReleaseDate:   "2023-06-01",
			Title:         "A New Dawn",
			VoteAverage:   8.2,
			VoteCount:     90,
		},
		{
			ID:            7,
			OriginalTitle: "A Última Resistência",
			ReleaseDate:   "2023-07-01",
			Title:         "The Final Stand",
			VoteAverage:   9.1,
			VoteCount:     150,
		},
		{
			ID:            8,
			OriginalTitle: "La Grande Aventure",
			ReleaseDate:   "2023-08-01",
			Title:         "The Great Adventure",
			VoteAverage:   8.7,
			VoteCount:     200,
		},
		{
			ID:            9,
			OriginalTitle: "A Jornada Continua",
			ReleaseDate:   "2023-09-01",
			Title:         "The Journey Continues",
			VoteAverage:   7.8,
			VoteCount:     110,
		},
		{
			ID:            10,
			OriginalTitle: "L'Héritage des Héros",
			ReleaseDate:   "2023-10-01",
			Title:         "The Legacy of Heroes",
			VoteAverage:   9.3,
			VoteCount:     250,
		},
		{
			ID:          11,
			Title:       "The Power of Unity",
			ReleaseDate: "2023-11-01",
			VoteAverage: 8.4,
			VoteCount:   130,
		},
		{
			ID:            12,
			OriginalTitle: "L'Appel à l'Aventure",
			ReleaseDate:   "2023-12-01",
			Title:         "The Call of Adventure",
			VoteAverage:   8.9,
			VoteCount:     180,
		},
		{
			ID:            13,
			OriginalTitle: "A Ascensão da Fênix",
			ReleaseDate:   "2024-01-01",
			Title:         "The Rise of the Phoenix",
			VoteAverage:   9.0,
			VoteCount:     220,
		},
		{
			ID:            14,
			OriginalTitle: "Le Voyage d'une Vie",
			ReleaseDate:   "2024-02-01",
			Title:         "The Journey of a Lifetime",
			VoteAverage:   7.6,
			VoteCount:     95,
		},
		{
			ID:            15,
			OriginalTitle: "O Coração de um Campeão",
			ReleaseDate:   "2024-03-01",
			Title:         "The Heart of a Champion",
			VoteAverage:   8.3,
			VoteCount:     140,
		},
		{
			ID:            16,
			OriginalTitle: "L'Esprit d'Aventure",
			ReleaseDate:   "2024-04-01",
			Title:         "The Spirit of Adventure",
			VoteAverage:   8.1,
			VoteCount:     115,
		},
		{
			ID:            17,
			OriginalTitle: "O Legado dos Bravos",
			ReleaseDate:   "2024-05-01",
			Title:         "The Legacy of the Brave",
			VoteAverage:   9.2,
			VoteCount:     260,
		},
		{
			ID:            18,
			OriginalTitle: "Le Pouvoir des Rêves",
			ReleaseDate:   "2024-06-01",
			Title:         "The Power of Dreams",
			VoteAverage:   8.6,
			VoteCount:     170,
		},
		{
			ID:            19,
			OriginalTitle: "O Chamado do Destino",
			ReleaseDate:   "2024-07-01",
			Title:         "The Call of Destiny",
			VoteAverage:   9.4,
			VoteCount:     300,
		},
		{
			ID:            20,
			OriginalTitle: "Le Voyage de l'Espoir",
			ReleaseDate:   "2024-08-01",
			Title:         "The Journey of Hope",
			VoteAverage:   8.8,
			VoteCount:     190,
		},
		{
			ID:            21,
			OriginalTitle: "A Ascensão dos Titãs",
			ReleaseDate:   "2024-09-01",
			Title:         "The Rise of the Titans",
			VoteAverage:   9.1,
			VoteCount:     240,
		},
		{
			ID:            22,
			OriginalTitle: "Le Cœur de l'Océan",
			ReleaseDate:   "2024-10-01",
			Title:         "The Heart of the Ocean",
			VoteAverage:   8.5,
			VoteCount:     160,
		},
		{
			ID:            23,
			OriginalTitle: "O Espírito da Floresta",
			ReleaseDate:   "2024-11-01",
			Title:         "The Spirit of the Forest",
			VoteAverage:   8.0,
			VoteCount:     130,
		},
		{
			ID:            24,
			OriginalTitle: "L'Appel de la Nature",
			ReleaseDate:   "2024-12-01",
			Title:         "The Call of the Wild",
			VoteAverage:   9.0,
			VoteCount:     210,
		},
		{
			ID:            25,
			OriginalTitle: "A Jornada das Estrelas",
			ReleaseDate:   "2025-01-01",
			Title:         "The Journey of the Stars",
			VoteAverage:   8.4,
			VoteCount:     175,
		},
		{
			ID:            26,
			OriginalTitle: "La Montée des Gardiens",
			ReleaseDate:   "2025-02-01",
			Title:         "The Rise of the Guardians",
			VoteAverage:   9.3,
			VoteCount:     280,
		},
		{
			ID:            27,
			OriginalTitle: "O Coração do Guerreiro",
			ReleaseDate:   "2025-03-01",
			Title:         "The Heart of the Warrior",
			VoteAverage:   8.7,
			VoteCount:     150,
		},
		{
			ID:            28,
			OriginalTitle: "L'Esprit du Ciel",
			ReleaseDate:   "2025-04-01",
			Title:         "The Spirit of the Sky",
			VoteAverage:   8.9,
			VoteCount:     165,
		},
		{
			ID:            29,
			OriginalTitle: "A Lenda do Dragão",
			ReleaseDate:   "2025-05-01",
			Title:         "The Legend of the Dragon",
			VoteAverage:   9.5,
			VoteCount:     300,
		},
		{
			ID:            30,
			OriginalTitle: "Le Dernier Voyage",
			ReleaseDate:   "2025-06-01",
			Title:         "The Last Voyage",
			VoteAverage:   8.3,
			VoteCount:     140,
		},
		{
			ID:            31,
			OriginalTitle: "O Caminho da Esperança",
			ReleaseDate:   "2025-07-01",
			Title:         "The Path of Hope",
			VoteAverage:   8.6,
			VoteCount:     175,
		},
		{
			ID:            32,
			OriginalTitle: "La Lumière de l'Aube",
			ReleaseDate:   "2025-08-01",
			Title:         "The Light of Dawn",
			VoteAverage:   9.0,
			VoteCount:     200,
		},
		{
			ID:            33,
			OriginalTitle: "A Última Esperança",
			ReleaseDate:   "2025-09-01",
			Title:         "The Last Hope",
			VoteAverage:   9.2,
			VoteCount:     220,
		},
		{
			ID:            34,
			OriginalTitle: "Le Voyage des Rêves",
			ReleaseDate:   "2025-10-01",
			Title:         "The Journey of Dreams",
			VoteAverage:   8.8,
			VoteCount:     190,
		},
		{
			ID:            35,
			OriginalTitle: "O Chamado da Aventura",
			ReleaseDate:   "2025-11-01",
			Title:         "The Call of Adventure",
			VoteAverage:   9.1,
			VoteCount:     210,
		},
		{
			ID:            36,
			OriginalTitle: "La Force de l'Amitié",
			ReleaseDate:   "2025-12-01",
			Title:         "The Strength of Friendship",
			VoteAverage:   8.4,
			VoteCount:     160,
		},
		{
			ID:            37,
			OriginalTitle: "A Jornada do Guerreiro",
			ReleaseDate:   "2026-01-01",
			Title:         "The Warrior's Journey",
			VoteAverage:   9.3,
			VoteCount:     250,
		},
		{
			ID:            38,
			OriginalTitle: "Le Dernier Héros",
			ReleaseDate:   "2026-02-01",
			Title:         "The Last Hero",
			VoteAverage:   8.7,
			VoteCount:     180,
		},
		{
			ID:            39,
			OriginalTitle: "O Legado dos Heróis",
			ReleaseDate:   "2026-03-01",
			Title:         "The Legacy of Heroes",
			VoteAverage:   9.0,
			VoteCount:     230,
		},
		{
			ID:            40,
			OriginalTitle: "La Quête de la Vérité",
			ReleaseDate:   "2026-04-01",
			Title:         "The Quest for Truth",
			VoteAverage:   8.5,
			VoteCount:     200,
		},
	}
	fakeResPage1 = tmdbResponse{
		Page:         1,
		Results:      fakeMovieList[:20],
		TotalPages:   2,
		TotalResults: len(fakeMovieList),
	}
	fakeResPage2 = tmdbResponse{
		Page:         2,
		Results:      fakeMovieList[20:],
		TotalPages:   2,
		TotalResults: len(fakeMovieList),
	}
	fakeEmptyRes = tmdbResponse{
		Page:         1,
		Results:      movies{},
		TotalResults: 0,
		TotalPages:   1,
	}
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buffer := new(bytes.Buffer)
	root.SetOut(buffer)
	root.SetErr(buffer)
	root.SetArgs(args)
	c, err = root.ExecuteC()
	return c, buffer.String(), err
}

func requireAPIKey(t testing.TB, w http.ResponseWriter, r *http.Request) {
	t.Helper()
	apiKey := r.Header.Get("Authorization")
	if apiKey != "Bearer valid_api_key" {
		w.WriteHeader(401)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertNotNil(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func assertURL(t testing.TB, got, want string) {
	if want != got {
		t.Errorf("expected URL to be %s, but got %s", want, got)
	}
}

func assertResponse(t testing.TB, want, got tmdbResponse) {
	t.Helper()
	if want.Page != got.Page {
		t.Errorf("expected Page to be %d, but got %d", want.Page, got.Page)
	}
	if want.TotalPages != got.TotalPages {
		t.Errorf("expected TotalPages to be %d, but got %d", want.TotalPages, got.TotalPages)
	}
	if want.TotalResults != got.TotalResults {
		t.Errorf("expected TotalResults to be %d, but got %d", want.TotalResults, got.TotalResults)
	}
	assertMovies(t, want.Results, want.Results)
}

func assertMovies(t testing.TB, want, got movies) {
	t.Helper()
	expectedMap := make(map[int]movie)
	for _, movie := range want {
		expectedMap[movie.ID] = movie
	}
	for _, movie := range got {
		if _, exists := expectedMap[movie.ID]; !exists {
			t.Errorf("unexpected movie in response: %+v", movie)
		}
	}
}

func assertPrintNoResults(t testing.TB, got string) {
	want := "No results available. Please try another query.\n"
	if want != got {
		t.Errorf("expected printed output to be %q, but got %q", want, got)
	}
}

func assertContains(t testing.TB, s string, sl []string) {
	t.Helper()
	for _, e := range sl {
		if !strings.Contains(s, e) {
			t.Errorf("expected output to contain %q", e)
		}
	}
}
