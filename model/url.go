package model

import (
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/laxman20/todo-go/todo"
	"github.com/pkg/browser"
)

type urlOpened struct{}

var TicketUrl string
var TicketPrefix string

func init() {
	TicketUrl = os.Getenv("TODO_TICKET_URL")
	TicketPrefix = os.Getenv("TODO_TICKET_PREFIX")
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
