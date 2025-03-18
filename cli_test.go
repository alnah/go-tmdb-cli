package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestIntegrationRootCmd(t *testing.T) {
	testCases := []struct {
		name              string
		missingConfigFile bool
		missingAPIKey     bool
		wantHelp          bool
		wantErr           bool
	}{
		{
			name: "init deps context on persistent pre run",
		},
		{
			name:     "display help on run",
			wantHelp: true,
		},
		{
			name:              "error when missing config file",
			wantErr:           true,
			missingConfigFile: true,
		},
		{
			name:          "error when missing api key in config file",
			wantErr:       true,
			missingAPIKey: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var root *cobra.Command
			home, _ := os.UserHomeDir()
			cfgPath := filepath.Join(home, ".go-tmdb-cli")
			file, _ := os.CreateTemp(cfgPath, "config_*.yaml")
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})
			if !tc.missingAPIKey {
				file.WriteString("api_key: valid_api_key")
			}
			if tc.missingConfigFile {
				root = newRootCmd(filepath.Base("missing_config.yaml"))
			} else {
				root = newRootCmd(filepath.Base(file.Name()))
			}
			// Act
			got, err := executeCommand(root)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
				_, ok := root.Context().Value(dependencies).(*Dependencies)
				if !ok {
					t.Error("retrieve dependencies from context")
				}
				if tc.wantHelp {
					assertContains(t, got, []string{"Usage", "Available Commands", "Flags"})
				}
			}
		})
	}
}

func TestIntegrationListCmd(t *testing.T) {
	testCases := []struct {
		name          string
		flag          string
		wantHelp      bool
		wantNoResults bool
		wantErr       bool
	}{
		{name: "now playing", flag: "--now"},
		{name: "popular", flag: "--pop"},
		{name: "top rated", flag: "--top"},
		{name: "upcoming", flag: "--up"},
		{name: "help", wantHelp: true},
		{name: "no results", flag: "--now", wantNoResults: true},
		{name: "error", flag: "--now", wantErr: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			home, _ := os.UserHomeDir()
			cfgPath := filepath.Join(home, ".go-tmdb-cli")
			file, _ := os.CreateTemp(cfgPath, "config_*.yaml")
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})
			file.WriteString("api_key: valid_api_key")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var byt []byte
				requireAPIKey(t, w, r)
				w.Header().Set("Content-Type", "application/json")
				if tc.wantErr {
					byt, _ = json.Marshal("Invalid JSON format")
				} else if tc.wantNoResults {
					byt, _ = json.Marshal(&fakeEmptyRes)
					t.Cleanup(func() { fakeResPage1.TotalResults = len(fakeResPage1.Results) })
				} else {
					byt, _ = json.Marshal(&fakeResPage1)
				}
				w.Write(byt)
			}))
			t.Cleanup(func() { ts.Close() })
			root := newRootCmd(filepath.Base(file.Name()))
			root.PersistentPreRunE = nil // Disable to prevent overriding mock
			mockCtx := context.WithValue(context.Background(), dependencies, &Dependencies{
				URLBuilder: &urlBuilder{
					BaseURL:  ts.URL,
					ListPath: "/movie/%s?",
				},
				Client: newHTTPClient("valid_api_key"),
			})
			root.SetContext(mockCtx)
			// Act
			got, err := executeCommand(root, "list", tc.flag)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else if tc.wantHelp {
				assertNoError(t, err)
				assertContains(t, got, []string{"Usage", "Examples", "Flags"})
			} else if tc.wantNoResults {
				assertNoError(t, err)
				assertPrintNoResults(t, got)
			} else {
				assertNoError(t, err)
				assertContains(t, got, []string{"ORIGINAL TITLE", "RELEASE DATE", "TITLE", "AVERAGE", "VOTES"})
			}
		})
	}
}

func TestIntegrationDiscoverCmd(t *testing.T) {
	testCases := []struct {
		name          string
		flag          string
		wantHelp      bool
		wantNoResults bool
		wantFetchErr  bool
		wantErr       bool
	}{
		{name: "language", flag: `--language=fr`},
		{name: "valid year", flag: `--year=2000`},
		{name: "valid year gte", flag: `--year=2000,gte`},
		{name: "valid year lte", flag: `--year=2000,lte`},
		{name: "valid years", flag: `--year=2000,2000`},
		{name: "valid average", flag: "--average=8.0,9.0"},
		{name: "valid average gte", flag: "--average=8.0,gte"},
		{name: "valid average lte", flag: "--average=8.0,lte"},
		{name: "valid votes", flag: "--votes=1000,2000"},
		{name: "valid votes gte", flag: "--votes=1000,gte"},
		{name: "valid votes lte", flag: "--votes=1000,lte"},
		{name: "valid one genre", flag: "--genres=drama"},
		{name: "valid many genres", flag: "--genres=comedy,horror,science-fiction"},
		{name: "valid one genre", flag: "--without-genres=drama"},
		{name: "valid many genres", flag: "--without-genres=comedy,horror,science-fiction"},
		{name: "valid sort", flag: "--sort=average,desc"},
		{name: "help", wantHelp: true},
		{name: "year error", flag: "--year=1", wantErr: true},                           // Parsing error
		{name: "average error", flag: "--average=11", wantErr: true},                    // Above max average
		{name: "votes error", flag: "--votes=-1", wantErr: true},                        // Below min average
		{name: "genres error", flag: "--genres=invalid", wantErr: true},                 // Below min average
		{name: "without genres error", flag: "--without-genres=invalid", wantErr: true}, // Below min average
		{name: "sort error", flag: "--sort=invalid,desc", wantErr: true},                // Invalid field fort sorting
		{name: "fetch error", flag: "--language=pt", wantFetchErr: true, wantErr: true},
		{name: "max items error", flag: "--max-items=abc", wantErr: true},
		{name: "no results", flag: `--language=fr`, wantNoResults: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var byt []byte
			var url string
			home, _ := os.UserHomeDir()
			cfgPath := filepath.Join(home, ".go-tmdb-cli")
			file, _ := os.CreateTemp(cfgPath, "config_*.yaml")
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})
			file.WriteString("api_key: valid_api_key")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requireAPIKey(t, w, r)
				if tc.wantNoResults {
					byt, _ = json.Marshal(fakeEmptyRes)
					t.Cleanup(func() { fakeResPage1.TotalResults = len(fakeResPage1.Results) })
				} else {
					byt, _ = json.Marshal(fakeResPage1)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(byt)
			}))
			t.Cleanup(func() { ts.Close() })
			root := newRootCmd(filepath.Base(file.Name()))
			root.PersistentPreRunE = nil // Disable to prevent overriding mock
			if tc.wantFetchErr {
				url = "https://not_found"
			} else {
				url = ts.URL
			}
			mockCtx := context.WithValue(context.Background(), dependencies, &Dependencies{
				URLBuilder: &urlBuilder{
					BaseURL:      url,
					DiscoverPath: "/discover/movie?",
				},
				Client: newHTTPClient("valid_api_key"),
			})
			root.SetContext(mockCtx)
			// Act
			got, err := executeCommand(root, "discover", tc.flag)
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else if tc.wantHelp {
				assertNoError(t, err)
				assertContains(t, got, []string{"Usage", "Examples", "Flags"})
			} else if tc.wantNoResults {
				assertNoError(t, err)
				assertPrintNoResults(t, got)
			} else {
				assertNoError(t, err)
				assertContains(t, got, []string{"ORIGINAL TITLE", "RELEASE DATE", "TITLE", "AVERAGE", "VOTES"})
			}
		})
	}
}

func TestIntegrationInfoCmd(t *testing.T) {
	// Arrange
	home, _ := os.UserHomeDir()
	cfgPath := filepath.Join(home, ".go-tmdb-cli")
	file, _ := os.CreateTemp(cfgPath, "config_*.yaml")
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})
	file.WriteString("api_key: valid_api_key")
	root := newRootCmd(filepath.Base(file.Name()))
	// Act
	got, err := executeCommand(root, "info")
	// Assert
	assertNoError(t, err)
	assertContains(t, got, []string{"v", "Alexis Nahan", "Apache"})
}
