package config_test

import (
	"io/fs"
	"os"
	"testing"

	c "github.com/alnah/go-tmdb-cli/internal/config"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetTMDBToken(t *testing.T) {
	t.Run("returns TMDB API Token from env file", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`TOKEN="test"`)
		assert.NoError(t, err)

		got, err := c.GetTMDBToken(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test", got)
	})

	t.Run("returns path error when env file doesn't exist", func(t *testing.T) {
		var want *fs.PathError
		_, err := c.GetTMDBToken("test.env")
		assert.ErrorAs(t, err, &want)
	})

	t.Run("returns config parse error when invalid data", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`invalid_data`)
		assert.NoError(t, err)

		var want viper.ConfigParseError
		_, err = c.GetTMDBToken(file.Name())
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
