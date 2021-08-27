package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/data"
	"github.com/laxman20/todo-go/todo"
	"github.com/muesli/reflow/wordwrap"
)

type state int

const (
	NORMAL state = iota
	ADD
	EDIT
	NOTES
)

type Model struct {
	state     state
	todos     []todo.Todo
	cursor    int
	insertPos int
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

func (m *Model) swapAbove() {
	if m.cursor > 0 {
		m.todos[m.cursor], m.todos[m.cursor-1] = m.todos[m.cursor-1], m.todos[m.cursor]
		m.up()
	}
}

func (m *Model) swapBelow() {
	if m.cursor < len(m.todos)-1 {
		m.todos[m.cursor], m.todos[m.cursor+1] = m.todos[m.cursor+1], m.todos[m.cursor]
		m.down()
	}
}

func (m *Model) gotoAdd() {
	m.state = ADD
	m.textInput.Focus()
}

func (m *Model) goToEdit() {
	m.state = EDIT
	m.textInput.SetValue(m.todos[m.cursor].Text)
	m.textInput.Focus()
}

func (m *Model) gotoNotes() {
	m.state = NOTES
	m.textInput.SetValue(wordwrap.String(m.todos[m.cursor].Notes, 80))
	m.textInput.Prompt = ""
	m.textInput.Focus()
}

func (m *Model) gotoNormal() {
	m.state = NORMAL
	m.textInput.Blur()
	m.textInput.Prompt = "> "
	m.textInput.Reset()
}

func (m *Model) toggle() {
	m.todos[m.cursor].Toggle()
}

func (m *Model) addTodo(text string) {
	todos := m.todos
	if m.state == ADD {
		m.todos = append(m.todos, todo.Todo{})
		copy(m.todos[m.insertPos+1:], m.todos[m.insertPos:])
		m.todos[m.insertPos] = todo.NewTodo(text)
	} else if m.state == EDIT {
		todos[m.cursor].Text = text
	}
}

func (m *Model) removeTodo() {
	if len(m.todos) == 0 {
		return
	}
	idx := m.cursor
	m.todos = append(m.todos[:idx], m.todos[idx+1:]...)
	if m.cursor > len(m.todos)-1 {
		m.cursor = len(m.todos) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
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
