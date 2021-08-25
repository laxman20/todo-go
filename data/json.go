package data

import (
	"encoding/json"
	"fmt"
	"github.com/laxman20/todo-go/todo"
	"io"
	"os"
	"path/filepath"
)

func getDataFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("Could not load app data directory: %w\n", err)
	}
	if err = os.MkdirAll(filepath.Join(cacheDir, "todo-go"), os.ModePerm); err != nil {
		return "", fmt.Errorf("Could not create todo-go directory: %w\n", err)
	}
	dataFilePath := filepath.Join(cacheDir, "todo-go", "data.json")
	return dataFilePath, nil
}

func LoadJSON() ([]todo.Todo, error) {
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return nil, fmt.Errorf("Could not get data file path: %w\n", err)
	}
	if _, err = os.Stat(dataFilePath); os.IsNotExist(err) {
		return []todo.Todo{}, nil
	}
	dataFile, err := os.OpenFile(dataFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %w\n", err)
	}
	defer dataFile.Close()
	data, err := io.ReadAll(dataFile)
	if err != nil {
		return nil, fmt.Errorf("Could not load data: %w\n", err)
	}
	newTodos := []todo.Todo{}
	err = json.Unmarshal(data, &newTodos)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON read from file: %w\n", err)
	}
	return newTodos, nil
}

func WriteJSON(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Could not serialize data: %w\n", err)
	}
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return fmt.Errorf("Could not get data file path: %w\n", err)
	}
	err = os.WriteFile(dataFilePath, bytes, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not write to file: %w\n", err)
	}
	return nil
}
