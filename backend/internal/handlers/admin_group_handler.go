package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AdminGroupHandler struct {
	groupService *services.AdminGroupService
}

func NewAdminGroupHandler() *AdminGroupHandler {
	return &AdminGroupHandler{
		groupService: services.NewAdminGroupService(),
	}
}

type CreateGroupRequest struct {
	Name     string      `json:"name"`
	AdminIDs []uuid.UUID `json:"admin_ids"`
}

type ReplaceAdminRequest struct {
	OldAdminID uuid.UUID `json:"old_admin_id"`
	NewAdminID uuid.UUID `json:"new_admin_id"`
}

func (h *AdminGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group, err := h.groupService.CreateGroup(req.Name, req.AdminIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *AdminGroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupService.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *AdminGroupHandler) ReplaceAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var req ReplaceAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.groupService.ReplaceAdminInGroup(groupID, req.OldAdminID, req.NewAdminID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin replaced successfully"})
}

