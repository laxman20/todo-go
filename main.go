package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/model"
	"os"
)

func main() {
	if err := tea.NewProgram(model.Model{}).Start(); err != nil {
		fmt.Printf("There was an error: %v\n", err)
		os.Exit(1)
	}
}
