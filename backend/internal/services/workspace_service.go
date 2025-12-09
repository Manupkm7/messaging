package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type WorkspaceService struct{}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

// GenerateSlug genera un slug único a partir del nombre
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remover caracteres especiales, mantener solo alfanuméricos y guiones
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// CreateWorkspace crea un nuevo workspace
func (s *WorkspaceService) CreateWorkspace(name string, createdByID uuid.UUID) (*models.Workspace, error) {
	// Verificar que el creador es super_admin o admin
	var creatorRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		createdByID,
	).Scan(&creatorRole)
	if err != nil {
		return nil, errors.New("creator not found")
	}
	if creatorRole != string(models.RoleSuperAdmin) && creatorRole != string(models.RoleAdmin) {
		return nil, errors.New("only admins and super admins can create workspaces")
	}

	// Generar slug único
	baseSlug := generateSlug(name)
	slug := baseSlug
	counter := 1

	// Verificar que el slug sea único
	for {
		var existingID uuid.UUID
		err := database.DB.QueryRow(
			"SELECT id FROM workspaces WHERE slug = $1",
			slug,
		).Scan(&existingID)
		if err == sql.ErrNoRows {
			break // Slug único encontrado
		} else if err != nil {
			return nil, err
		}
		// Slug existe, agregar número
		slug = baseSlug + "-" + string(rune('0'+counter))
		counter++
	}

	workspace := &models.Workspace{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = database.DB.Exec(`
		INSERT INTO workspaces (id, name, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, workspace.ID, workspace.Name, workspace.Slug, workspace.CreatedAt, workspace.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return workspace, nil
}

// GetWorkspaceByID obtiene un workspace por ID
func (s *WorkspaceService) GetWorkspaceByID(id uuid.UUID) (*models.Workspace, error) {
	workspace := &models.Workspace{}
	err := database.DB.QueryRow(`
		SELECT id, name, slug, created_at, updated_at
		FROM workspaces WHERE id = $1
	`, id).Scan(
		&workspace.ID, &workspace.Name, &workspace.Slug,
		&workspace.CreatedAt, &workspace.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

// GetWorkspaceBySlug obtiene un workspace por slug
func (s *WorkspaceService) GetWorkspaceBySlug(slug string) (*models.Workspace, error) {
	workspace := &models.Workspace{}
	err := database.DB.QueryRow(`
		SELECT id, name, slug, created_at, updated_at
		FROM workspaces WHERE slug = $1
	`, slug).Scan(
		&workspace.ID, &workspace.Name, &workspace.Slug,
		&workspace.CreatedAt, &workspace.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return workspace, nil
}

// GetAllWorkspaces obtiene todos los workspaces
func (s *WorkspaceService) GetAllWorkspaces(userID uuid.UUID) ([]models.WorkspaceWithAdmins, error) {
	// Verificar que el usuario es super_admin o admin
	var userRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		userID,
	).Scan(&userRole)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var rows *sql.Rows
	if userRole == string(models.RoleSuperAdmin) {
		// Super admin puede ver todos los workspaces
		rows, err = database.DB.Query(`
			SELECT w.id, w.name, w.slug, w.created_at, w.updated_at,
			       COALESCE(
			         json_agg(DISTINCT wa.admin_id) FILTER (WHERE wa.admin_id IS NOT NULL),
			         '[]'::json
			       ) as admin_ids
			FROM workspaces w
			LEFT JOIN workspace_admins wa ON w.id = wa.workspace_id
			GROUP BY w.id, w.name, w.slug, w.created_at, w.updated_at
			ORDER BY w.created_at DESC
		`)
	} else {
		// Admin solo puede ver workspaces donde está asignado
		rows, err = database.DB.Query(`
			SELECT w.id, w.name, w.slug, w.created_at, w.updated_at,
			       COALESCE(
			         json_agg(DISTINCT wa.admin_id) FILTER (WHERE wa.admin_id IS NOT NULL),
			         '[]'::json
			       ) as admin_ids
			FROM workspaces w
			INNER JOIN workspace_admins wa ON w.id = wa.workspace_id
			WHERE wa.admin_id = $1
			GROUP BY w.id, w.name, w.slug, w.created_at, w.updated_at
			ORDER BY w.created_at DESC
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []models.WorkspaceWithAdmins
	for rows.Next() {
		var workspace models.WorkspaceWithAdmins
		var adminIDsJSON []byte
		err := rows.Scan(
			&workspace.ID, &workspace.Name, &workspace.Slug,
			&workspace.CreatedAt, &workspace.UpdatedAt, &adminIDsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parsear JSON de admin_ids
		if len(adminIDsJSON) > 0 && string(adminIDsJSON) != "[]" {
			var adminIDs []string
			if err := json.Unmarshal(adminIDsJSON, &adminIDs); err == nil {
				for _, idStr := range adminIDs {
					if id, err := uuid.Parse(idStr); err == nil {
						workspace.AdminIDs = append(workspace.AdminIDs, id)
					}
				}
			}
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

// AssignAdminToWorkspace asigna un admin a un workspace
func (s *WorkspaceService) AssignAdminToWorkspace(workspaceID, adminID, assignedByID uuid.UUID) error {
	// Verificar que el asignador tiene permisos
	var assignerRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		assignedByID,
	).Scan(&assignerRole)
	if err != nil {
		return errors.New("assigner not found")
	}
	if assignerRole != string(models.RoleSuperAdmin) && assignerRole != string(models.RoleAdmin) {
		return errors.New("only admins and super admins can assign admins to workspaces")
	}

	// Verificar que el workspace existe
	var workspaceExists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1)",
		workspaceID,
	).Scan(&workspaceExists)
	if err != nil || !workspaceExists {
		return errors.New("workspace not found")
	}

	// Verificar que el admin existe y es admin
	var adminRole string
	err = database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		adminID,
	).Scan(&adminRole)
	if err != nil {
		return errors.New("admin not found")
	}
	if adminRole != string(models.RoleAdmin) {
		return errors.New("user is not an admin")
	}

	// Verificar que no está ya asignado
	var exists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM workspace_admins WHERE workspace_id = $1 AND admin_id = $2)",
		workspaceID, adminID,
	).Scan(&exists)
	if err == nil && exists {
		return errors.New("admin is already assigned to this workspace")
	}

	// Asignar admin al workspace
	_, err = database.DB.Exec(`
		INSERT INTO workspace_admins (workspace_id, admin_id, created_at)
		VALUES ($1, $2, $3)
	`, workspaceID, adminID, time.Now())

	return err
}

// RemoveAdminFromWorkspace remueve un admin de un workspace
func (s *WorkspaceService) RemoveAdminFromWorkspace(workspaceID, adminID, removedByID uuid.UUID) error {
	// Verificar que el remover tiene permisos
	var removerRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		removedByID,
	).Scan(&removerRole)
	if err != nil {
		return errors.New("remover not found")
	}
	if removerRole != string(models.RoleSuperAdmin) && removerRole != string(models.RoleAdmin) {
		return errors.New("only admins and super admins can remove admins from workspaces")
	}

	// Remover admin del workspace
	_, err = database.DB.Exec(`
		DELETE FROM workspace_admins WHERE workspace_id = $1 AND admin_id = $2
	`, workspaceID, adminID)

	return err
}

// GetWorkspaceAdmins obtiene los admins de un workspace
func (s *WorkspaceService) GetWorkspaceAdmins(workspaceID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := database.DB.Query(`
		SELECT admin_id FROM workspace_admins WHERE workspace_id = $1
	`, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adminIDs []uuid.UUID
	for rows.Next() {
		var adminID uuid.UUID
		if err := rows.Scan(&adminID); err != nil {
			return nil, err
		}
		adminIDs = append(adminIDs, adminID)
	}

	return adminIDs, nil
}

// UpdateClientWorkspace actualiza el workspace de un cliente (delegación)
func (s *WorkspaceService) UpdateClientWorkspace(clientID, newWorkspaceID, updatedByID uuid.UUID) error {
	// Verificar que el actualizador tiene permisos
	var updaterRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		updatedByID,
	).Scan(&updaterRole)
	if err != nil {
		return errors.New("updater not found")
	}
	if updaterRole != string(models.RoleSuperAdmin) && updaterRole != string(models.RoleAdmin) {
		return errors.New("only admins and super admins can update client workspace")
	}

	// Verificar que el cliente existe y es cliente
	var clientRole string
	err = database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		clientID,
	).Scan(&clientRole)
	if err != nil {
		return errors.New("client not found")
	}
	if clientRole != string(models.RoleClient) {
		return errors.New("user is not a client")
	}

	// Verificar que el nuevo workspace existe
	var workspaceExists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1)",
		newWorkspaceID,
	).Scan(&workspaceExists)
	if err != nil || !workspaceExists {
		return errors.New("workspace not found")
	}

	// Actualizar workspace del cliente
	_, err = database.DB.Exec(`
		UPDATE users SET workspace_id = $1, updated_at = $2 WHERE id = $3
	`, newWorkspaceID, time.Now(), clientID)

	return err
}

