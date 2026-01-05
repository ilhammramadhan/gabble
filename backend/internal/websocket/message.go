package websocket

import (
	"time"

	"github.com/ilhammramadhan/gabble/internal/models"
)

type EventType string

const (
	EventJoinRoom    EventType = "join_room"
	EventLeaveRoom   EventType = "leave_room"
	EventSendMessage EventType = "send_message"
	EventTyping      EventType = "typing"
	EventMessage     EventType = "message"
	EventUserJoined  EventType = "user_joined"
	EventUserLeft    EventType = "user_left"
	EventOnlineUsers EventType = "online_users"
	EventError       EventType = "error"
)

type WSMessage struct {
	Type    EventType   `json:"type"`
	Payload interface{} `json:"payload"`
}

type JoinRoomPayload struct {
	RoomID string `json:"room_id"`
}

type LeaveRoomPayload struct {
	RoomID string `json:"room_id"`
}

type SendMessagePayload struct {
	RoomID  string `json:"room_id"`
	Content string `json:"content"`
}

type TypingPayload struct {
	RoomID   string `json:"room_id"`
	IsTyping bool   `json:"is_typing"`
}

type MessagePayload struct {
	ID        string       `json:"id"`
	RoomID    string       `json:"room_id"`
	Content   string       `json:"content"`
	User      *models.User `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

type UserEventPayload struct {
	RoomID string       `json:"room_id"`
	User   *models.User `json:"user"`
}

type TypingEventPayload struct {
	RoomID   string       `json:"room_id"`
	User     *models.User `json:"user"`
	IsTyping bool         `json:"is_typing"`
}

type OnlineUsersPayload struct {
	RoomID string         `json:"room_id"`
	Users  []*models.User `json:"users"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
