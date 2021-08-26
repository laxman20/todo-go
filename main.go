package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/model"
)

func main() {
	if err := tea.NewProgram(model.Model{}, tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("There was an error: %v\n", err)
		os.Exit(1)
	}
}
