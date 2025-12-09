package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatService    *services.ChatService
	messageService *services.MessageService
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatService:    services.NewChatService(),
		messageService: services.NewMessageService(),
	}
}

type CreateChatRequest struct {
	ClientID     *uuid.UUID `json:"client_id,omitempty"`
	OtherAdminID *uuid.UUID `json:"other_admin_id,omitempty"`
}

func (h *ChatHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := middleware.GetUserRole(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := h.chatService.GetChatsByUserID(userID, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := middleware.GetUserRole(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var chat *models.Chat
	var err error

	if role == string(models.RoleAdmin) || role == string(models.RoleSuperAdmin) {
		if req.OtherAdminID != nil {
			// Crear chat entre admins
			chat, err = h.chatService.CreateAdminChat(userID, *req.OtherAdminID)
		} else if req.ClientID != nil {
			// Crear chat con cliente
			chat, err = h.chatService.CreateChat(userID, *req.ClientID)
		} else {
			http.Error(w, "client_id or other_admin_id required", http.StatusBadRequest)
			return
		}
	} else if role == string(models.RoleClient) {
		// Cliente creando chat - asignar automáticamente a un admin del workspace
		chat, err = h.chatService.CreateChatForClient(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) DelegateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ToAdminID uuid.UUID `json:"to_admin_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.chatService.DelegateChat(chatID, userID, req.ToAdminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat delegated successfully"})
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Verificar que el usuario tiene acceso al chat
	chat, err := h.chatService.GetChatByID(chatID)
	if err != nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	// Verificar permisos
	hasAccess := false
	userRole, _ := middleware.GetUserRole(r)
	
	if chat.AdminID == userID {
		hasAccess = true
	} else if chat.ClientID != nil && *chat.ClientID == userID {
		hasAccess = true
	} else if chat.OtherAdminID != nil && *chat.OtherAdminID == userID {
		hasAccess = true
	} else if chat.AssignedToID != nil && *chat.AssignedToID == userID {
		hasAccess = true
	} else if (userRole == string(models.RoleAdmin) || userRole == string(models.RoleSuperAdmin)) && chat.WorkspaceID != nil {
		// Verificar si el admin está asignado al workspace del chat
		var isAssigned bool
		err := database.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM workspace_admins 
				WHERE workspace_id = $1 AND admin_id = $2
			)
		`, *chat.WorkspaceID, userID).Scan(&isAssigned)
		if err == nil && (isAssigned || userRole == string(models.RoleSuperAdmin)) {
			hasAccess = true
		}
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Obtener parámetros de paginación
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	messages, err := h.messageService.GetMessagesByChatID(chatID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No enviar los datos binarios completos en la lista, solo metadata
	type MessageResponse struct {
		ID        uuid.UUID `json:"id"`
		ChatID    uuid.UUID `json:"chat_id"`
		SenderID  uuid.UUID `json:"sender_id"`
		Content   *string   `json:"content,omitempty"`
		HasAudio  bool      `json:"has_audio"`
		HasImage  bool      `json:"has_image"`
		MimeType  *string   `json:"mime_type,omitempty"`
		CreatedAt string    `json:"created_at"`
	}

	response := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = MessageResponse{
			ID:        msg.ID,
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			HasAudio:  msg.IsAudio,
			HasImage:  msg.IsImage,
			MimeType:  msg.MimeType,
			CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Parsear multipart form para archivos
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var content *string
	if c := r.FormValue("content"); c != "" {
		content = &c
	}

	var audioData []byte
	var imageData []byte
	var mimeType *string

	// Procesar audio
	if audioFile, header, err := r.FormFile("audio"); err == nil {
		defer audioFile.Close()
		audioData = make([]byte, header.Size)
		audioFile.Read(audioData)
		mime := header.Header.Get("Content-Type")
		mimeType = &mime
	}

	// Procesar imagen
	if imageFile, header, err := r.FormFile("image"); err == nil {
		defer imageFile.Close()
		imageData = make([]byte, header.Size)
		imageFile.Read(imageData)
		mime := header.Header.Get("Content-Type")
		mimeType = &mime
	}

	message, err := h.messageService.CreateMessage(chatID, userID, content, audioData, imageData, mimeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast del mensaje vía WebSocket
	websocket.BroadcastMessageToChat(chatID, message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
