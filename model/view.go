package model

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

var encodedEnter = "¬"

func textInputView(m *Model) string {
	decoded := strings.Replace(m.textInput.View(), encodedEnter, "\n", -1)
	return wordwrap.String(decoded+"\n", 80)
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
