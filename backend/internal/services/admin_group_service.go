package services

import (
	"encoding/json"
	"errors"
	"time"

	"backend/internal/database"
	"backend/internal/models"

	"github.com/google/uuid"
)

type AdminGroupService struct{}

func NewAdminGroupService() *AdminGroupService {
	return &AdminGroupService{}
}

func (s *AdminGroupService) CreateGroup(name string, adminIDs []uuid.UUID) (*models.AdminGroup, error) {
	// Validar que todos los IDs son admins
	for _, adminID := range adminIDs {
		var role string
		err := database.DB.QueryRow(
			"SELECT role FROM users WHERE id = $1",
			adminID,
		).Scan(&role)
		if err != nil {
			return nil, errors.New("admin not found: " + adminID.String())
		}
		if role != string(models.RoleAdmin) {
			return nil, errors.New("user is not an admin: " + adminID.String())
		}
	}

	group := &models.AdminGroup{
		ID:        uuid.New(),
		Name:      name,
		AdminIDs:  adminIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	adminIDsJSON, err := json.Marshal(adminIDs)
	if err != nil {
		return nil, err
	}

	_, err = database.DB.Exec(`
		INSERT INTO admin_groups (id, name, admin_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, group.ID, group.Name, adminIDsJSON, group.CreatedAt, group.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *AdminGroupService) GetGroupByID(groupID uuid.UUID) (*models.AdminGroup, error) {
	group := &models.AdminGroup{}
	var adminIDsJSON []byte
	err := database.DB.QueryRow(`
		SELECT id, name, admin_ids, created_at, updated_at
		FROM admin_groups WHERE id = $1
	`, groupID).Scan(
		&group.ID, &group.Name, &adminIDsJSON, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(adminIDsJSON, &group.AdminIDs)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *AdminGroupService) ReplaceAdminInGroup(groupID, oldAdminID, newAdminID uuid.UUID) error {
	// Verificar que el nuevo admin es admin
	var newAdminRole string
	err := database.DB.QueryRow(
		"SELECT role FROM users WHERE id = $1",
		newAdminID,
	).Scan(&newAdminRole)
	if err != nil {
		return errors.New("new admin not found")
	}
	if newAdminRole != string(models.RoleAdmin) {
		return errors.New("new user is not an admin")
	}

	// Obtener el grupo
	group, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	// Reemplazar el admin en la lista
	found := false
	for i, adminID := range group.AdminIDs {
		if adminID == oldAdminID {
			group.AdminIDs[i] = newAdminID
			found = true
			break
		}
	}

	if !found {
		return errors.New("old admin not found in group")
	}

	// Actualizar en la base de datos
	adminIDsJSON, err := json.Marshal(group.AdminIDs)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(`
		UPDATE admin_groups SET admin_ids = $1, updated_at = $2 WHERE id = $3
	`, adminIDsJSON, time.Now(), groupID)

	return err
}

func (s *AdminGroupService) GetAllGroups() ([]models.AdminGroup, error) {
	rows, err := database.DB.Query(`
		SELECT id, name, admin_ids, created_at, updated_at
		FROM admin_groups
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.AdminGroup
	for rows.Next() {
		var group models.AdminGroup
		var adminIDsJSON []byte
		err := rows.Scan(
			&group.ID, &group.Name, &adminIDsJSON, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(adminIDsJSON, &group.AdminIDs)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}
