package model

import "os"

var TicketUrl string
var TicketPrefix string

func init() {
	TicketUrl = os.Getenv("TODO_TICKET_URL")
	TicketPrefix = os.Getenv("TODO_TICKET_PREFIX")
}
