package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	token, err := getTMDBToken(".env")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(token)
}

func getTMDBToken(filepath string) (string, error) {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	viper.SetConfigType("ENV")
	if err = viper.ReadConfig(bytes.NewBuffer(byt)); err != nil {
		return "", err
	}

	return viper.GetString("TOKEN"), nil
}
