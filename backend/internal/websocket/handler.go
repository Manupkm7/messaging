package websocket

import (
	"log"
	"net/http"

	"backend/internal/services"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// En producción, verificar el origen
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Obtener token del header o query parameter
	tokenString := ""
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}
	} else if token := r.URL.Query().Get("token"); token != "" {
		tokenString = token
	}

	if tokenString == "" {
		http.Error(w, "Unauthorized: token required", http.StatusUnauthorized)
		return
	}

	// Validar token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:      uuid.New(),
		UserID:  userID,
		Conn:    conn,
		Hub:     GlobalHub,
		Send:    make(chan []byte, 256),
		ChatIDs: make(map[uuid.UUID]bool),
	}

	// Obtener chats del usuario y suscribirse automáticamente
	role := claims.Role
	chatService := services.NewChatService()
	chats, err := chatService.GetChatsByUserID(userID, role)
	if err == nil {
		for _, chat := range chats {
			client.ChatIDs[chat.ID] = true
		}
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
