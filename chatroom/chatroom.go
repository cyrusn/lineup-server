package chatroom

import "time"

// Message store information of message send
type Message struct {
	Form    string
	Name    string
	Content string
	IssueAt time.Time
}
