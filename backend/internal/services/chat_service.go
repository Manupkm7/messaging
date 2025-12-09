package services

import (
	"database/sql"
	"errors"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type ChatService struct{}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) CreateChat(adminID, clientID uuid.UUID) (*models.Chat, error) {
	// Verificar que el admin existe y es admin o super_admin
	var adminRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		adminID,
	).Scan(&adminRole)
	if err != nil {
		return nil, errors.New("admin not found")
	}
	if adminRole != string(models.RoleAdmin) && adminRole != string(models.RoleSuperAdmin) {
		return nil, errors.New("user is not an admin")
	}

	// Verificar que el cliente existe y es cliente, y obtener su workspace_id
	var clientRole string
	var clientWorkspaceID uuid.UUID
	err = database.DB.QueryRow(
		"SELECT role, workspace_id FROM users WHERE id = $1",
		clientID,
	).Scan(&clientRole, &clientWorkspaceID)
	if err != nil {
		return nil, errors.New("client not found")
	}
	if clientRole != string(models.RoleClient) {
		return nil, errors.New("user is not a client")
	}

	// Verificar que el admin está asignado al workspace del cliente (si no es super_admin)
	if adminRole != string(models.RoleSuperAdmin) {
		var isAssigned bool
		err = database.DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM workspace_admins 
				WHERE workspace_id = $1 AND admin_id = $2
			)
		`, clientWorkspaceID, adminID).Scan(&isAssigned)
		if err != nil || !isAssigned {
			return nil, errors.New("admin is not assigned to client's workspace")
		}
	}

	// Verificar si ya existe un chat activo
	var existingChatID uuid.UUID
	err = database.DB.QueryRow(`
		SELECT id FROM chats 
		WHERE admin_id = $1 AND client_id = $2 AND is_admin_chat = FALSE
		ORDER BY created_at DESC LIMIT 1
	`, adminID, clientID).Scan(&existingChatID)
	if err == nil {
		// Retornar el chat existente
		return s.GetChatByID(existingChatID)
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Crear nuevo chat con workspace_id
	chat := &models.Chat{
		ID:          uuid.New(),
		AdminID:     adminID,
		ClientID:    &clientID,
		IsAdminChat: false,
		WorkspaceID: &clientWorkspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO chats (id, admin_id, client_id, is_admin_chat, workspace_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, chat.ID, chat.AdminID, chat.ClientID, chat.IsAdminChat, chat.WorkspaceID, chat.CreatedAt, chat.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

// CreateChatForClient crea un chat para un cliente, asignándolo automáticamente a un admin del workspace
func (s *ChatService) CreateChatForClient(clientID uuid.UUID) (*models.Chat, error) {
	// Verificar que el cliente existe y obtener su workspace_id
	var clientWorkspaceID uuid.UUID
	err := database.DB.QueryRow(
		"SELECT workspace_id FROM users WHERE id = $1 AND role = 'client'",
		clientID,
	).Scan(&clientWorkspaceID)
	if err != nil {
		return nil, errors.New("client not found or not associated with a workspace")
	}

	// Verificar si ya existe un chat activo para este cliente
	var existingChatID uuid.UUID
	err = database.DB.QueryRow(`
		SELECT id FROM chats 
		WHERE client_id = $1 AND is_admin_chat = FALSE
		ORDER BY created_at DESC LIMIT 1
	`, clientID).Scan(&existingChatID)
	if err == nil {
		// Retornar el chat existente
		return s.GetChatByID(existingChatID)
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Buscar un admin disponible del workspace
	var adminID uuid.UUID
	err = database.DB.QueryRow(`
		SELECT wa.admin_id 
		FROM workspace_admins wa
		JOIN users u ON wa.admin_id = u.id
		WHERE wa.workspace_id = $1 AND u.role = 'admin'
		ORDER BY RANDOM()
		LIMIT 1
	`, clientWorkspaceID).Scan(&adminID)
	if err != nil {
		return nil, errors.New("no admin available in workspace")
	}

	// Crear nuevo chat con el admin asignado
	chat := &models.Chat{
		ID:           uuid.New(),
		AdminID:      adminID,
		ClientID:     &clientID,
		AssignedToID: &adminID,
		IsAdminChat: false,
		WorkspaceID:  &clientWorkspaceID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO chats (id, admin_id, client_id, assigned_to_id, is_admin_chat, workspace_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, chat.ID, chat.AdminID, chat.ClientID, chat.AssignedToID, chat.IsAdminChat, chat.WorkspaceID, chat.CreatedAt, chat.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) CreateAdminChat(adminID, otherAdminID uuid.UUID) (*models.Chat, error) {
	// Para chats entre admins, buscar un workspace común o usar el primero disponible
	var workspaceID uuid.UUID
	var workspaceIDPtr *uuid.UUID
	err := database.DB.QueryRow(`
		SELECT workspace_id FROM workspace_admins 
		WHERE admin_id = $1 
		AND workspace_id IN (SELECT workspace_id FROM workspace_admins WHERE admin_id = $2)
		LIMIT 1
	`, adminID, otherAdminID).Scan(&workspaceID)
	
	// Si no hay workspace común, usar el primero del admin que crea el chat
	if err != nil {
		err = database.DB.QueryRow(`
			SELECT workspace_id FROM workspace_admins WHERE admin_id = $1 LIMIT 1
		`, adminID).Scan(&workspaceID)
		if err == nil {
			workspaceIDPtr = &workspaceID
		}
		// Si aún no hay, crear sin workspace (NULL permitido para admin chats)
	} else {
		workspaceIDPtr = &workspaceID
	}

	chat := &models.Chat{
		ID:           uuid.New(),
		AdminID:      adminID,
		OtherAdminID: &otherAdminID,
		IsAdminChat:  true,
		WorkspaceID:  workspaceIDPtr,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO chats (id, admin_id, other_admin_id, is_admin_chat, workspace_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, chat.ID, chat.AdminID, chat.OtherAdminID, chat.IsAdminChat, workspaceIDPtr, chat.CreatedAt, chat.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) GetChatByID(chatID uuid.UUID) (*models.Chat, error) {
	chat := &models.Chat{}
	err := database.DB.QueryRow(`
		SELECT id, admin_id, client_id, assigned_to_id, is_admin_chat, other_admin_id, 
		       workspace_id, created_at, updated_at, last_message_at
		FROM chats WHERE id = $1
	`, chatID).Scan(
		&chat.ID, &chat.AdminID, &chat.ClientID, &chat.AssignedToID,
		&chat.IsAdminChat, &chat.OtherAdminID, &chat.WorkspaceID, &chat.CreatedAt, &chat.UpdatedAt, &chat.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) GetChatsByUserID(userID uuid.UUID, role string) ([]models.ChatWithParticipant, error) {
	var rows *sql.Rows
	var err error

	if role == string(models.RoleAdmin) || role == string(models.RoleSuperAdmin) {
		if role == string(models.RoleSuperAdmin) {
			// Super admin puede ver todos los chats
			rows, err = database.DB.Query(`
				SELECT c.id, c.admin_id, c.client_id, c.assigned_to_id, c.is_admin_chat, 
				       c.other_admin_id, c.workspace_id, c.created_at, c.updated_at, c.last_message_at,
				       CASE 
				         WHEN c.is_admin_chat THEN u2.username
				         ELSE u1.username
				       END as participant_name,
				       CASE 
				         WHEN c.is_admin_chat THEN COALESCE(c.other_admin_id, c.admin_id)
				         ELSE COALESCE(c.client_id, c.admin_id)
				       END as participant_id
				FROM chats c
				LEFT JOIN users u1 ON c.client_id = u1.id
				LEFT JOIN users u2 ON c.other_admin_id = u2.id
				WHERE c.admin_id = $1 OR c.assigned_to_id = $1 OR c.other_admin_id = $1
				ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
			`, userID)
		} else {
			// Admin normal: ver chats donde es admin/assigned Y chats de workspaces donde está asignado
			rows, err = database.DB.Query(`
				SELECT DISTINCT c.id, c.admin_id, c.client_id, c.assigned_to_id, c.is_admin_chat, 
				       c.other_admin_id, c.workspace_id, c.created_at, c.updated_at, c.last_message_at,
				       CASE 
				         WHEN c.is_admin_chat THEN u2.username
				         ELSE u1.username
				       END as participant_name,
				       CASE 
				         WHEN c.is_admin_chat THEN COALESCE(c.other_admin_id, c.admin_id)
				         ELSE COALESCE(c.client_id, c.admin_id)
				       END as participant_id
				FROM chats c
				LEFT JOIN users u1 ON c.client_id = u1.id
				LEFT JOIN users u2 ON c.other_admin_id = u2.id
				LEFT JOIN workspace_admins wa ON c.workspace_id = wa.workspace_id
				WHERE (
					-- Chats donde el admin es dueño o está asignado
					c.admin_id = $1 OR c.assigned_to_id = $1 OR c.other_admin_id = $1
					OR
					-- Chats de workspaces donde el admin está asignado (solo chats con clientes)
					(c.is_admin_chat = FALSE AND c.workspace_id IS NOT NULL AND wa.admin_id = $1)
				)
				ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
			`, userID)
		}
	} else {
		// Cliente: ver solo sus chats
		rows, err = database.DB.Query(`
			SELECT c.id, c.admin_id, c.client_id, c.assigned_to_id, c.is_admin_chat, 
			       c.other_admin_id, c.workspace_id, c.created_at, c.updated_at, c.last_message_at,
			       u.username as participant_name,
			       COALESCE(c.assigned_to_id, c.admin_id) as participant_id
			FROM chats c
			JOIN users u ON COALESCE(c.assigned_to_id, c.admin_id) = u.id
			WHERE c.client_id = $1
			ORDER BY COALESCE(c.last_message_at, c.created_at) DESC
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatWithParticipant
	for rows.Next() {
		var chat models.ChatWithParticipant
		err := rows.Scan(
			&chat.ID, &chat.AdminID, &chat.ClientID, &chat.AssignedToID,
			&chat.IsAdminChat, &chat.OtherAdminID, &chat.WorkspaceID, &chat.CreatedAt, &chat.UpdatedAt, &chat.LastMessageAt,
			&chat.ParticipantName, &chat.ParticipantID,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *ChatService) DelegateChat(chatID, fromAdminID, toAdminID uuid.UUID) error {
	// Verificar que el chat existe y pertenece al admin
	var currentAdminID uuid.UUID
	err := database.DB.QueryRow(`
		SELECT admin_id FROM chats WHERE id = $1
	`, chatID).Scan(&currentAdminID)
	if err != nil {
		return errors.New("chat not found")
	}

	// Verificar que el admin que delega es el dueño o está asignado
	var assignedToID sql.NullString
	err = database.DB.QueryRow(`
		SELECT assigned_to_id FROM chats WHERE id = $1
	`, chatID).Scan(&assignedToID)
	if err != nil {
		return err
	}

	if currentAdminID != fromAdminID && (!assignedToID.Valid || assignedToID.String != fromAdminID.String()) {
		return errors.New("you don't have permission to delegate this chat")
	}

	// Verificar que el admin destino existe y es admin o super_admin
	var toAdminRole string
	err = database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		toAdminID,
	).Scan(&toAdminRole)
	if err != nil {
		return errors.New("target admin not found")
	}
	if toAdminRole != string(models.RoleAdmin) && toAdminRole != string(models.RoleSuperAdmin) {
		return errors.New("target user is not an admin")
	}

	// Actualizar el chat
	_, err = database.DB.Exec(`
		UPDATE chats SET assigned_to_id = $1, updated_at = $2 WHERE id = $3
	`, toAdminID, time.Now(), chatID)

	return err
}

func (s *ChatService) DeleteOldChats() error {
	// Eliminar chats creados hace más de 24 horas
	_, err := database.DB.Exec(`
		DELETE FROM chats WHERE created_at < NOW() - INTERVAL '24 hours'
	`)
	return err
}
