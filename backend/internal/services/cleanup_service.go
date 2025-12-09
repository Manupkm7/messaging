package services

import (
	"log"
	"time"
)

type CleanupService struct {
	chatService    *ChatService
	messageService *MessageService
}

func NewCleanupService() *CleanupService {
	return &CleanupService{
		chatService:    NewChatService(),
		messageService: NewMessageService(),
	}
}

func (s *CleanupService) StartCleanupJob() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			log.Println("Starting cleanup job: deleting chats older than 24 hours...")
			
			if err := s.chatService.DeleteOldChats(); err != nil {
				log.Printf("Error deleting old chats: %v", err)
			} else {
				log.Println("Successfully deleted old chats")
			}

			if err := s.messageService.DeleteOldMessages(); err != nil {
				log.Printf("Error deleting old messages: %v", err)
			} else {
				log.Println("Successfully deleted old messages")
			}
		}
	}()
}

