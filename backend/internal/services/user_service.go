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

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(username string, email *string, password, role string, createdByID *uuid.UUID, workspaceID *uuid.UUID) (*models.User, error) {
	// Validar que solo admins o super_admins pueden crear admins
	if role == string(models.RoleAdmin) && createdByID != nil {
		var creatorRole string
		err := database.DB.QueryRow(
			"SELECT role FROM users WHERE id = $1",
			createdByID,
		).Scan(&creatorRole)
		if err != nil {
			return nil, errors.New("error verifying creator")
		}
		if creatorRole != string(models.RoleAdmin) && creatorRole != string(models.RoleSuperAdmin) {
			return nil, errors.New("only admins and super admins can create admin users")
		}
	}

	// Validar workspace_id para clientes
	if role == string(models.RoleClient) && workspaceID == nil {
		return nil, errors.New("workspace_id is required for clients")
	}
	if role != string(models.RoleClient) && workspaceID != nil {
		return nil, errors.New("workspace_id can only be set for clients")
	}

	// Verificar si el usuario ya existe (solo por username)
	var existingID uuid.UUID
	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE username = $1",
		username,
	).Scan(&existingID)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Generar session token para clientes
	var sessionToken *string
	if role == string(models.RoleClient) {
		token, err := utils.GenerateSessionToken()
		if err != nil {
			return nil, err
		}
		sessionToken = &token
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.UserRole(role),
		WorkspaceID:  workspaceID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if sessionToken != nil {
		user.SessionToken = *sessionToken
	}

	_, err = database.DB.Exec(`
		INSERT INTO users (id, username, email, password_hash, role, session_token, workspace_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, sessionToken, workspaceID, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, session_token, workspace_id, created_at, updated_at
		FROM users WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.SessionToken, &user.WorkspaceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, session_token, workspace_id, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.SessionToken, &user.WorkspaceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserBySessionToken(token string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, session_token, workspace_id, created_at, updated_at
		FROM users WHERE session_token = $1
	`, token).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.SessionToken, &user.WorkspaceID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateSessionToken(userID uuid.UUID, token string) error {
	_, err := database.DB.Exec(
		"UPDATE users SET session_token = $1, updated_at = $2 WHERE id = $3",
		token, time.Now(), userID,
	)
	return err
}
