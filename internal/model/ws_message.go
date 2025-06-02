package model

type WebSocketMessage struct {
	Event      string `json:"event"`
	Timestamp  string `json:"timestamp,omitempty"`
	SenderID   int    `json:"sender_id,omitempty"`
	SenderName string `json:"sender_name,omitempty"`
	Message    string `json:"message,omitempty"`
}
