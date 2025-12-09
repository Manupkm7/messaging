package services

import (
	"database/sql"
	"errors"
	"time"

	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/google/uuid"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) CreateMessage(chatID, senderID uuid.UUID, content *string, audioData, imageData []byte, mimeType *string) (*models.Message, error) {
	// Verificar que el chat existe
	var chatAdminID uuid.UUID
	var chatClientID sql.NullString
	var chatOtherAdminID sql.NullString
	var isAdminChat bool
	err := database.DB.QueryRow(`
		SELECT admin_id, client_id, other_admin_id, is_admin_chat
		FROM chats WHERE id = $1
	`, chatID).Scan(&chatAdminID, &chatClientID, &chatOtherAdminID, &isAdminChat)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	// Verificar que el sender tiene permiso para enviar en este chat
	canSend := false
	if senderID == chatAdminID {
		canSend = true
	} else if chatClientID.Valid && senderID.String() == chatClientID.String {
		canSend = true
	} else if chatOtherAdminID.Valid && senderID.String() == chatOtherAdminID.String {
		canSend = true
	}

	if !canSend {
		return nil, errors.New("you don't have permission to send messages in this chat")
	}

	// Validar que hay contenido, audio o imagen
	isAudio := len(audioData) > 0
	isImage := len(imageData) > 0

	if content == nil && !isAudio && !isImage {
		return nil, errors.New("message must have content, audio, or image")
	}

	// Cifrar contenido de texto si existe
	var encryptedContent *string
	if content != nil && *content != "" {
		encrypted, err := utils.EncryptMessage(*content)
		if err != nil {
			return nil, errors.New("failed to encrypt message content")
		}
		encryptedContent = &encrypted
	}

	// Cifrar datos binarios si existen
	var encryptedAudioData []byte
	var encryptedImageData []byte
	if isAudio && len(audioData) > 0 {
		encrypted, err := utils.EncryptData(audioData)
		if err != nil {
			return nil, errors.New("failed to encrypt audio data")
		}
		encryptedAudioData = encrypted
	}
	if isImage && len(imageData) > 0 {
		encrypted, err := utils.EncryptData(imageData)
		if err != nil {
			return nil, errors.New("failed to encrypt image data")
		}
		encryptedImageData = encrypted
	}

	message := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   encryptedContent,
		AudioData: encryptedAudioData,
		ImageData: encryptedImageData,
		MimeType:  mimeType,
		IsAudio:   isAudio,
		IsImage:   isImage,
		CreatedAt: time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO messages (id, chat_id, sender_id, content, audio_data, image_data, mime_type, is_audio, is_image, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, message.ID, message.ChatID, message.SenderID, message.Content,
		message.AudioData, message.ImageData, message.MimeType,
		message.IsAudio, message.IsImage, message.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Actualizar last_message_at del chat
	_, err = database.DB.Exec(`
		UPDATE chats SET last_message_at = $1, updated_at = $1 WHERE id = $2
	`, time.Now(), chatID)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetMessagesByChatID(chatID uuid.UUID, limit, offset int) ([]models.Message, error) {
	rows, err := database.DB.Query(`
		SELECT id, chat_id, sender_id, content, audio_data, image_data, mime_type, is_audio, is_image, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content,
			&msg.AudioData, &msg.ImageData, &msg.MimeType,
			&msg.IsAudio, &msg.IsImage, &msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Descifrar contenido de texto si existe
		if msg.Content != nil && *msg.Content != "" {
			decrypted, err := utils.DecryptMessage(*msg.Content)
			if err != nil {
				// Si falla el descifrado, mantener el contenido cifrado
				// pero loguear el error
				continue
			}
			msg.Content = &decrypted
		}

		// Descifrar datos binarios si existen
		if msg.IsAudio && len(msg.AudioData) > 0 {
			decrypted, err := utils.DecryptData(msg.AudioData)
			if err == nil {
				msg.AudioData = decrypted
			}
		}
		if msg.IsImage && len(msg.ImageData) > 0 {
			decrypted, err := utils.DecryptData(msg.ImageData)
			if err == nil {
				msg.ImageData = decrypted
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *MessageService) DeleteOldMessages() error {
	// Eliminar mensajes de chats que ya fueron eliminados o mensajes antiguos
	// Esto se ejecuta automáticamente por CASCADE cuando se eliminan los chats
	_, err := database.DB.Exec(`
		DELETE FROM messages WHERE created_at < NOW() - INTERVAL '24 hours'
	`)
	return err
}
