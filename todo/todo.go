package todo

import "fmt"

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func NewTodo(text string) Todo {
	return Todo{
		Text: text,
		Done: false,
	}
}

func (t Todo) String() string {
	doneLabel := " "
	if t.Done {
		doneLabel = "x"
	}
	return fmt.Sprintf("[%s] %s", doneLabel, t.Text)
}

func (t *Todo) Toggle() {
	t.Done = !t.Done
}
