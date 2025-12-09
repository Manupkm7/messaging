package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	AdminID       uuid.UUID  `json:"admin_id" db:"admin_id"`
	ClientID      *uuid.UUID `json:"client_id,omitempty" db:"client_id"`
	AssignedToID  *uuid.UUID `json:"assigned_to_id,omitempty" db:"assigned_to_id"`
	IsAdminChat   bool       `json:"is_admin_chat" db:"is_admin_chat"`
	OtherAdminID  *uuid.UUID `json:"other_admin_id,omitempty" db:"other_admin_id"`
	WorkspaceID   *uuid.UUID `json:"workspace_id,omitempty" db:"workspace_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty" db:"last_message_at"`
}

type Message struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChatID    uuid.UUID `json:"chat_id" db:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id" db:"sender_id"`
	Content   *string   `json:"content,omitempty" db:"content"`
	AudioData []byte    `json:"-" db:"audio_data"`
	ImageData []byte    `json:"-" db:"image_data"`
	MimeType  *string   `json:"mime_type,omitempty" db:"mime_type"`
	IsAudio   bool      `json:"is_audio" db:"is_audio"`
	IsImage   bool      `json:"is_image" db:"is_image"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatWithParticipant struct {
	Chat
	ParticipantName string    `json:"participant_name"`
	ParticipantID   uuid.UUID `json:"participant_id"`
}
