package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/ilhammramadhan/gabble/internal/database"
	"github.com/ilhammramadhan/gabble/internal/models"
)

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	DB         *database.DB
	mu         sync.RWMutex
}

func NewHub(db *database.DB) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		DB:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)

				if client.RoomID != "" {
					h.removeFromRoom(client, client.RoomID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) HandleMessage(client *Client, msg *WSMessage) {
	switch msg.Type {
	case EventJoinRoom:
		h.handleJoinRoom(client, msg)
	case EventLeaveRoom:
		h.handleLeaveRoom(client, msg)
	case EventSendMessage:
		h.handleSendMessage(client, msg)
	case EventTyping:
		h.handleTyping(client, msg)
	}
}

func (h *Hub) handleJoinRoom(client *Client, msg *WSMessage) {
	payloadBytes, _ := json.Marshal(msg.Payload)
	var payload JoinRoomPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	h.mu.Lock()
	if client.RoomID != "" {
		h.removeFromRoom(client, client.RoomID)
	}

	if h.Rooms[payload.RoomID] == nil {
		h.Rooms[payload.RoomID] = make(map[*Client]bool)
	}
	h.Rooms[payload.RoomID][client] = true
	client.RoomID = payload.RoomID
	h.mu.Unlock()

	h.broadcastToRoom(payload.RoomID, &WSMessage{
		Type: EventUserJoined,
		Payload: UserEventPayload{
			RoomID: payload.RoomID,
			User:   client.User,
		},
	}, client)

	h.sendOnlineUsers(payload.RoomID)
}

func (h *Hub) handleLeaveRoom(client *Client, msg *WSMessage) {
	payloadBytes, _ := json.Marshal(msg.Payload)
	var payload LeaveRoomPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return
	}

	h.mu.Lock()
	h.removeFromRoom(client, payload.RoomID)
	client.RoomID = ""
	h.mu.Unlock()

	h.broadcastToRoom(payload.RoomID, &WSMessage{
		Type: EventUserLeft,
		Payload: UserEventPayload{
			RoomID: payload.RoomID,
			User:   client.User,
		},
	}, nil)

	h.sendOnlineUsers(payload.RoomID)
}

func (h *Hub) handleSendMessage(client *Client, msg *WSMessage) {
	payloadBytes, _ := json.Marshal(msg.Payload)
	var payload SendMessagePayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	if payload.Content == "" || payload.RoomID == "" {
		h.sendError(client, "Message content and room ID are required")
		return
	}

	dbMsg, err := h.DB.CreateMessage(context.Background(), payload.RoomID, client.User.ID, payload.Content)
	if err != nil {
		log.Printf("error creating message: %v", err)
		h.sendError(client, "Failed to save message")
		return
	}

	h.broadcastToRoom(payload.RoomID, &WSMessage{
		Type: EventMessage,
		Payload: MessagePayload{
			ID:        dbMsg.ID,
			RoomID:    dbMsg.RoomID,
			Content:   dbMsg.Content,
			User:      client.User,
			CreatedAt: dbMsg.CreatedAt,
		},
	}, nil)
}

func (h *Hub) handleTyping(client *Client, msg *WSMessage) {
	payloadBytes, _ := json.Marshal(msg.Payload)
	var payload TypingPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return
	}

	h.broadcastToRoom(payload.RoomID, &WSMessage{
		Type: EventTyping,
		Payload: TypingEventPayload{
			RoomID:   payload.RoomID,
			User:     client.User,
			IsTyping: payload.IsTyping,
		},
	}, client)
}

func (h *Hub) removeFromRoom(client *Client, roomID string) {
	if room, ok := h.Rooms[roomID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, msg *WSMessage, exclude *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for client := range room {
		if client != exclude {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(room, client)
			}
		}
	}
}

func (h *Hub) sendOnlineUsers(roomID string) {
	h.mu.RLock()
	room, ok := h.Rooms[roomID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	users := make([]*models.User, 0, len(room))
	for client := range room {
		users = append(users, client.User)
	}
	h.mu.RUnlock()

	msg := &WSMessage{
		Type: EventOnlineUsers,
		Payload: OnlineUsersPayload{
			RoomID: roomID,
			Users:  users,
		},
	}

	h.broadcastToRoom(roomID, msg, nil)
}

func (h *Hub) sendError(client *Client, message string) {
	data, _ := json.Marshal(&WSMessage{
		Type:    EventError,
		Payload: ErrorPayload{Message: message},
	})
	client.Send <- data
}
