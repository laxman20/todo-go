package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/data"
	"github.com/laxman20/todo-go/todo"
	"github.com/muesli/reflow/wordwrap"
	"github.com/pkg/browser"
)

var encodedEnter = "¬"

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
type urlOpened struct{}

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

func openTicket(todo *todo.Todo) tea.Cmd {
	return func() tea.Msg {
		if len(TicketUrl) == 0 || len(TicketPrefix) == 0 {
			return nil
		}
		prefixes := strings.SplitN(TicketPrefix, ",", -1)
		tickets := []string{}
		for _, prefix := range prefixes {
			r, _ := regexp.Compile(prefix + "-[0-9]+")
			tickets = append(tickets, r.FindAllString(todo.Text, -1)...)
			tickets = append(tickets, r.FindAllString(todo.Notes, -1)...)
		}
		if len(tickets) > 0 {
			url := strings.Replace(TicketUrl, "{}", tickets[0], 1)
			browser.OpenURL(url)
		}
		return urlOpened{}
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case todoLoadMsg:
		m.todos = msg
		ti := textinput.NewModel()
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
			case "J":
				m.swapBelow()
				return m, nil
			case "k":
				m.up()
				return m, nil
			case "K":
				m.swapAbove()
				return m, nil
			case " ":
				m.toggle()
				return m, nil
			case "o":
				m.insertPos = min(m.cursor+1, len(m.todos))
				m.gotoAdd()
				return m, nil
			case "O":
				m.insertPos = m.cursor
				m.gotoAdd()
				return m, nil
			case "A":
				m.insertPos = len(m.todos)
				m.gotoAdd()
				return m, nil
			case "i":
				m.goToEdit()
				return m, nil
			case "D":
				m.removeTodo()
				return m, nil
			case "enter":
				if len(m.todos) > 0 {
					m.gotoNotes()
				}
				return m, nil
			case "x":
				if len(m.todos) > 0 {
					return m, openTicket(&m.todos[m.cursor])
				}
				return m, nil
			case "q":
				return m, writeTodos(m)
			}
		}
	}

	if m.state == ADD || m.state == EDIT {
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
	}

	if m.state == NOTES {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.todos[m.cursor].Notes = m.textInput.Value()
				m.gotoNormal()
				return m, nil
			case tea.KeyEnter:
				o := m.textInput.Value()
				m.textInput.Reset()
				m.textInput.SetValue(o + encodedEnter)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func textInputView(m *Model) string {
	decoded := strings.Replace(m.textInput.View(), encodedEnter, "\n", -1)
	return wordwrap.String(decoded+"\n", 80)
}

func mainView(m *Model) string {
	s := "Todos:\n"
	if len(m.todos) == 0 {
		s += "  No todos!\n"
	}
	editView := textInputView(m)
	for i, todo := range m.todos {
		cursorTxt := " "
		if m.cursor == i {
			cursorTxt = "*"
		}
		todoView := fmt.Sprintf("  %s %s\n", cursorTxt, todo)
		if m.state == ADD && m.insertPos == i {
			s += editView + todoView
		} else if m.state == EDIT && m.cursor == i {
			s += editView
		} else {
			s += todoView
		}
	}
	if m.state == ADD && m.insertPos == len(m.todos) {
		s += editView
	}
	return s
}

func notesView(m *Model) string {
	todo := m.todos[m.cursor]
	status := "(PENDING)"
	if todo.Done {
		status = "(DONE)"
	}
	br := strings.Repeat("=", 80)
	header := fmt.Sprintf("%s %s\n%s", todo.Text, status, br)
	return fmt.Sprintf("%s\nNotes:\n%s", header, textInputView(m))
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("An error occurred: %v\n", m.err)
	}
	switch m.state {
	case NORMAL, EDIT, ADD:
		return mainView(&m)
	case NOTES:
		return notesView(&m)
	}
	return "View not found"
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
