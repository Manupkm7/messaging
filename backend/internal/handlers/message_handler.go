package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/database"
	"backend/internal/middleware"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type MessageHandler struct {
	messageService *services.MessageService
	chatService    *services.ChatService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: services.NewMessageService(),
		chatService:    services.NewChatService(),
	}
}

func (h *MessageHandler) GetMessageFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.ErrUnauthorized)
		return
	}

	vars := mux.Vars(r)
	messageID, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid message ID", err.Error()))
		return
	}

	// Obtener el mensaje
	var chatID uuid.UUID
	var senderID uuid.UUID
	var audioData []byte
	var imageData []byte
	var mimeType string
	var isAudio bool
	var isImage bool

	err = database.DB.QueryRow(`
		SELECT chat_id, sender_id, audio_data, image_data, mime_type, is_audio, is_image
		FROM messages WHERE id = $1
	`, messageID).Scan(&chatID, &senderID, &audioData, &imageData, &mimeType, &isAudio, &isImage)

	if err != nil {
		utils.HandleError(w, r, utils.ErrNotFound)
		return
	}

	// Verificar que el usuario tiene acceso al chat
	chat, err := h.chatService.GetChatByID(chatID)
	if err != nil {
		utils.HandleError(w, r, utils.ErrNotFound)
		return
	}

	hasAccess := false
	if chat.AdminID == userID {
		hasAccess = true
	} else if chat.ClientID != nil && *chat.ClientID == userID {
		hasAccess = true
	} else if chat.OtherAdminID != nil && *chat.OtherAdminID == userID {
		hasAccess = true
	} else if chat.AssignedToID != nil && *chat.AssignedToID == userID {
		hasAccess = true
	}

	if !hasAccess {
		utils.HandleError(w, r, utils.ErrForbidden)
		return
	}

	// Determinar qué archivo devolver
	var fileData []byte
	var contentType string

	if isAudio && len(audioData) > 0 {
		// Descifrar audio
		decrypted, err := utils.DecryptData(audioData)
		if err != nil {
			utils.HandleError(w, r, utils.NewAppError(http.StatusInternalServerError, "Failed to decrypt audio file"))
			return
		}
		fileData = decrypted
		contentType = mimeType
		if contentType == "" {
			contentType = "audio/mpeg"
		}
	} else if isImage && len(imageData) > 0 {
		// Descifrar imagen
		decrypted, err := utils.DecryptData(imageData)
		if err != nil {
			utils.HandleError(w, r, utils.NewAppError(http.StatusInternalServerError, "Failed to decrypt image file"))
			return
		}
		fileData = decrypted
		contentType = mimeType
		if contentType == "" {
			contentType = "image/jpeg"
		}
	} else {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Message has no file"))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))
	w.Write(fileData)
}
