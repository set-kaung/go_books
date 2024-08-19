package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrWrongFileFormat = errors.New("the file must be in .json format")

func createLocalSources() (map[string]string, error) {
	appDir, err := creatAppDataDir()
	if err != nil {
		return nil, err
	}
	entries := map[string]string{
		"cache.json": filepath.Join(appDir, "cache.json"),
		"books.db":   filepath.Join(appDir, "books.db"),
	}

	for name, path := range entries {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			// Create the file if it doesn't exist
			file, err := os.Create(path)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s: %w", name, err)
			}
			file.Close()
		}
	}

	return entries, nil
}

func creatAppDataDir() (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Define the path to the Application Support directory
	appSupportPath := filepath.Join(homeDir, "Library", "Application Support")
	appDir := filepath.Join(appSupportPath, "go_books")
	if _, err := os.Stat(appDir); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("The default directory '%s' does not exist.\n", appDir)
		fmt.Print("Would you like to use the default path? (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response == "no" || response == "n" {
			fmt.Print("Please enter a different path: ")
			customPath, err := reader.ReadString('\n')
			if err != nil {
				return "", fmt.Errorf("failed to read path: %w", err)
			}

			customPath = strings.TrimSpace(customPath)
			if customPath != "" {
				appDir = customPath
			}

		}
		if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create application directory: %w", err)
		}
		fmt.Println("Directory created at:", appDir)

	} else {
		fmt.Println("Using existing directory:", appDir)
	}

	return appDir, nil
}
