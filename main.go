package main

import (
	"os"
)

func main() {
	rootCmd := newRootCmd("config.yaml")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
