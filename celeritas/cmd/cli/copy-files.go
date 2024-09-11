package main

import (
	"embed"
	"os"
)

//go:embed templates/*
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return err
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		return err
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
