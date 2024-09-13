package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
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

func getDSN() string {
	dbType := cel.DB.DataType
	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		}
		return dsn
	}

	return "mysql://" + cel.BuildDSN()
}

func showHelp() {
	color.Yellow(`Avaliable commands:
	help                       - show the help commands
	version                    - print application version
	migrate                    - runs all up migrations that have not benn run previously
	migrate down               - runs all down migrations in reverse order, and then all up migrations
	migrate reset              - runs all down migrations in reverse order, and then all up migrations
	make migration <name>      - creates two new up and down migrations in the migrations folder
	make auth                  - creates and run migrations for authentication tables, and creates models and middles
	make handler <name>        - creates a stub handler in the handers directory
	make model <name>          - creates a new model in the models directory
	make session               - creates a table in the database as a session store
	`)
}
