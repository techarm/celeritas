package main

import (
	"os"

	"github.com/joho/godotenv"
)

func setup() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	cel.RootPath = path
	cel.DB.DataType = os.Getenv("DATABASE_TYPE")

	return nil
}
