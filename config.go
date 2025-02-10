package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// userHome enables testable home directory resolution across OS environments.
type userHome interface {
	dir() (string, error)
}

// defaultUserHome implements userHome using actual OS home directory lookup.
type defaultUserHome struct{}

func (u *defaultUserHome) dir() (string, error) {
	return os.UserHomeDir()
}

// initialize loads config file and validates API key for TMDB access.
func initialize(userHome userHome, fileName string) error {
	home, err := userHome.dir()
	if err != nil {
		return fmt.Errorf("get user home directory: %w", err)
	}
	cfgPath := filepath.Join(home, ".go-tmdb-cli", fileName)
	byt, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("read the configuration file: %w ", err)
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(byt)); err != nil {
		return fmt.Errorf("parse the configuration file: %w", err)
	}
	return nil
}
