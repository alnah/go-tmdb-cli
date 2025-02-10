package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type mockUserHome struct{}

func (m *mockUserHome) dir() (string, error) {
	return "", fmt.Errorf("home dir not found")
}

func TestUnitInitialize(t *testing.T) {
	testCases := []struct {
		name             string
		fileContent      string
		fileName         string
		wantMockUserHome bool
		wantErr          bool
	}{
		{
			name:        "valid config",
			fileContent: "api_key: api_value",
			wantErr:     false,
		},
		{
			name:             "missing home dir",
			wantMockUserHome: true,
			wantErr:          true,
		},
		{
			name:     "missing config file",
			fileName: "missing_config.yaml",
			wantErr:  true,
		},
		{
			name:        "invalid yaml format",
			fileContent: "invalid:yaml:content",
			wantErr:     true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			var configFile string
			var err error
			home, _ := os.UserHomeDir()
			rootDir := filepath.Join(home, ".go-tmdb-cli")
			os.MkdirAll(rootDir, 0o755)
			file, _ := os.CreateTemp(rootDir, "config_*.yaml")
			t.Cleanup(func() {
				file.Close()
				os.Remove(file.Name())
			})
			file.WriteString(tc.fileContent)
			if tc.fileName == "missing_config.yaml" {
				configFile = tc.fileName
			} else {
				configFile = filepath.Base(file.Name())
			}
			// Act
			if tc.wantMockUserHome {
				err = initialize(&mockUserHome{}, configFile)
			} else {
				err = initialize(&defaultUserHome{}, configFile)
			}
			// Assert
			if tc.wantErr {
				assertNotNil(t, err)
			} else {
				assertNoError(t, err)
			}
		})
	}
}
