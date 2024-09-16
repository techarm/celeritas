package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1 string) error {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
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
	}

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
	make mail <name>           - creates two starter mail templates in the mail directory
	`)
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// check if current file is directory
	if fi.IsDir() {
		return nil
	}

	// onlcy check go files
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return nil
	}

	// we have a matching file
	if matched {
		// read file contents
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		// write the changed file
		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateSource() error {
	// walk entire project folder, including subfolders
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		return err
	}
	return nil
}
