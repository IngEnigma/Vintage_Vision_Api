package domain

import "time"

type PartyMessage struct {
	ID       uint
	PartyID  uint
	SenderID uint
	Message  string
	SentAt   time.Time

	Sender Profile
}
