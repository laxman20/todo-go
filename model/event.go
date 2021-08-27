package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
