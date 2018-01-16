package chatroom

import "time"

// Message store information of message send
type Message struct {
	Form    string    `json:"form"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
	IssueAt time.Time `json:"issueAt"`
}
