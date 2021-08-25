package model

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/data"
	"github.com/laxman20/todo-go/todo"
)

type state int

const (
	NORMAL state = iota
	APPEND
	EDIT
)

type Model struct {
	state     state
	todos     []todo.Todo
	cursor    int
	textInput textinput.Model
	err       error
}

type todoLoadMsg []todo.Todo
type todoSaveMsg struct{}

func (m *Model) up() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) down() {
	if m.cursor < len(m.todos)-1 {
		m.cursor++
	}
}

func (m *Model) gotoAdd() {
	m.state = APPEND
	m.textInput.Focus()
}

func (m *Model) goToEdit() {
	m.state = EDIT
	m.textInput.SetValue(m.todos[m.cursor].Text)
	m.textInput.Focus()
}

func (m *Model) gotoNormal() {
	m.state = NORMAL
	m.textInput.Blur()
	m.textInput.Reset()
}

func (m *Model) toggle() {
	m.todos[m.cursor].Toggle()
}

func (m *Model) addTodo(text string) {
	if m.state == APPEND {
		m.todos = append(m.todos, todo.NewTodo(text))
	} else if m.state == EDIT {
		m.todos[m.cursor].Text = text
	}
}

func (m *Model) removeTodo() {
	if len(m.todos) == 0 {
		return
	}
	idx := m.cursor
	m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
}

func (m Model) Init() tea.Cmd {
	return loadTodos
}

func loadTodos() tea.Msg {
	todos, err := data.LoadJSON()
	if err != nil {
		return err
	}
	return todoLoadMsg(todos)
}

func writeTodos(model Model) tea.Cmd {
	return func() tea.Msg {
		err := data.WriteJSON(model.todos)
		if err != nil {
			return err
		}
		return todoSaveMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return m, writeTodos(m)
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

func (m Model) View() string {
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
