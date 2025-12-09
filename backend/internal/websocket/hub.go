package websocket

import (
	"encoding/json"
	"log"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Conn    *websocket.Conn
	Hub     *Hub
	Send    chan []byte
	ChatIDs map[uuid.UUID]bool
}

type Hub struct {
	Clients    map[uuid.UUID]*Client
	Broadcast  chan BroadcastMessage
	Register   chan *Client
	Unregister chan *Client
}

type BroadcastMessage struct {
	ChatID  uuid.UUID
	Message *models.Message
	Data    []byte
}

var GlobalHub *Hub

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan BroadcastMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
			log.Printf("Client connected: %s (User: %s)", client.ID, client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				log.Printf("Client disconnected: %s", client.ID)
			}

		case broadcast := <-h.Broadcast:
			messageData, err := json.Marshal(broadcast.Message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			// Enviar a todos los clientes que están en este chat
			for _, client := range h.Clients {
				if client.ChatIDs[broadcast.ChatID] {
					select {
					case client.Send <- messageData:
					default:
						close(client.Send)
						delete(h.Clients, client.ID)
					}
				}
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Procesar mensajes del cliente si es necesario
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if action, ok := msg["action"].(string); ok {
				switch action {
				case "subscribe":
					if chatIDStr, ok := msg["chat_id"].(string); ok {
						if chatID, err := uuid.Parse(chatIDStr); err == nil {
							c.ChatIDs[chatID] = true
						}
					}
				case "unsubscribe":
					if chatIDStr, ok := msg["chat_id"].(string); ok {
						if chatID, err := uuid.Parse(chatIDStr); err == nil {
							delete(c.ChatIDs, chatID)
						}
					}
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

func BroadcastMessageToChat(chatID uuid.UUID, message *models.Message) {
	if GlobalHub != nil {
		GlobalHub.Broadcast <- BroadcastMessage{
			ChatID:  chatID,
			Message: message,
		}
	}
}
