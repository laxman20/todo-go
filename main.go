package main

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"os"
	"path/filepath"
)

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (t Todo) String() string {
	doneLabel := " "
	if t.Done {
		doneLabel = "x"
	}
	return fmt.Sprintf("[%s] %s", doneLabel, t.Text)
}

func (t *Todo) toggle() {
	t.Done = !t.Done
}

func makeTodo(text string) Todo {
	return Todo{Text: text, Done: false}
}

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

func loadTodos() tea.Msg {
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return fmt.Errorf("Could not get data file path: %w\n", err)
	}
	if _, err = os.Stat(dataFilePath); os.IsNotExist(err) {
		return todoLoadMsg([]Todo{})
	}
	dataFile, err := os.OpenFile(dataFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error opening file: %w\n", err)
	}
	defer dataFile.Close()
	data, err := io.ReadAll(dataFile)
	if err != nil {
		return fmt.Errorf("Could not load data: %w\n", err)
	}
	newTodos := []Todo{}
	err = json.Unmarshal(data, &newTodos)
	if err != nil {
		return fmt.Errorf("Invalid JSON read from file: %w\n", err)
	}
	return todoLoadMsg(newTodos)
}

func writeTodos(todos []Todo) tea.Cmd {
	return func() tea.Msg {
		bytes, err := json.Marshal(todos)
		if err != nil {
			return fmt.Errorf("Could not serialize todos: %w\n", err)
		}
		dataFilePath, err := getDataFilePath()
		if err != nil {
			return fmt.Errorf("Could not get data file path: %w\n", err)
		}
		err = os.WriteFile(dataFilePath, bytes, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Could not write to file: %w\n", err)
		}
		return todoSaveMsg{}
	}
}

type State int

const (
	NORMAL State = iota
	APPEND
	EDIT
)

type model struct {
	state     State
	todos     []Todo
	cursor    int
	textInput textinput.Model
	err       error
}

type todoLoadMsg []Todo
type todoSaveMsg struct{}

func (m *model) up() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) down() {
	if m.cursor < len(m.todos)-1 {
		m.cursor++
	}
}

func (m *model) gotoAdd() {
	m.state = APPEND
	m.textInput.Focus()
}

func (m *model) goToEdit() {
	m.state = EDIT
	m.textInput.SetValue(m.todos[m.cursor].Text)
	m.textInput.Focus()
}

func (m *model) gotoNormal() {
	m.state = NORMAL
	m.textInput.Blur()
	m.textInput.Reset()
}

func (m *model) toggle() {
	m.todos[m.cursor].toggle()
}

func (m *model) addTodo(text string) {
	if m.state == APPEND {
		m.todos = append(m.todos, makeTodo(text))
	} else if m.state == EDIT {
		m.todos[m.cursor].Text = text
	}
}

func (m *model) removeTodo() {
	if len(m.todos) == 0 {
		return
	}
	idx := m.cursor
	m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
}

func (m model) Init() tea.Cmd {
	return loadTodos
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case todoLoadMsg:
		m.todos = msg
		ti := textinput.NewModel()
		ti.CharLimit = 100
		m.textInput = ti

		return m, nil
	case todoSaveMsg:
		return m, tea.Quit
	case error:
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.state == NORMAL {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j":
				m.down()
				return m, nil
			case "k":
				m.up()
				return m, nil
			case " ":
				m.toggle()
				return m, nil
			case "A":
				m.gotoAdd()
				return m, nil
			case "i":
				m.goToEdit()
				return m, nil
			case "D":
				m.removeTodo()
				return m, nil
			case "q":
				return m, writeTodos(m.todos)
			}
		}
	}

	if m.state == APPEND || m.state == EDIT {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.gotoNormal()
				return m, nil
			case tea.KeyEnter:
				text := m.textInput.Value()
				if len(text) > 0 {
					m.addTodo(text)
					m.gotoNormal()
				}
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("An error occurred: %v\n", m.err)
	}
	s := "Todos:\n"
	if len(m.todos) == 0 {
		s += "  No todos!\n"
	}
	for i, todo := range m.todos {
		if m.state == EDIT && m.cursor == i {
			s += m.textInput.View() + "\n"
		} else {
			cursorTxt := " "
			if m.cursor == i {
				cursorTxt = "*"
			}
			s += fmt.Sprintf("  %s %s\n", cursorTxt, todo)
		}
	}
	if m.state == APPEND {
		s += m.textInput.View()
	}
	return s
}

func main() {
	if err := tea.NewProgram(model{}).Start(); err != nil {
		fmt.Printf("There was an error: %v\n", err)
		os.Exit(1)
	}
}
