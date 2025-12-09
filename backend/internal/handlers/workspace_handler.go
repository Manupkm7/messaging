package handlers

import (
	"encoding/json"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WorkspaceHandler struct {
	workspaceService *services.WorkspaceService
}

func NewWorkspaceHandler() *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: services.NewWorkspaceService(),
	}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

type AssignAdminRequest struct {
	AdminID uuid.UUID `json:"admin_id"`
}

type UpdateClientWorkspaceRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id"`
}

func (h *WorkspaceHandler) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	var req CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	if req.Name == "" {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Name is required"))
		return
	}

	workspace, err := h.workspaceService.CreateWorkspace(req.Name, userID)
	if err != nil {
		if err.Error() == "only admins and super admins can create workspaces" {
			utils.HandleError(w, r, utils.NewAppError(http.StatusForbidden, err.Error()))
		} else {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, err.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workspace)
}

func (h *WorkspaceHandler) GetWorkspaces(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	workspaces, err := h.workspaceService.GetAllWorkspaces(userID)
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusInternalServerError, "Error fetching workspaces", err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}

func (h *WorkspaceHandler) GetWorkspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workspaceID, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace ID"))
		return
	}

	workspace, err := h.workspaceService.GetWorkspaceByID(workspaceID)
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusNotFound, "Workspace not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspace)
}

func (h *WorkspaceHandler) AssignAdmin(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	workspaceID, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace ID"))
		return
	}

	var req AssignAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err = h.workspaceService.AssignAdminToWorkspace(workspaceID, req.AdminID, userID)
	if err != nil {
		if err.Error() == "only admins and super admins can assign admins to workspaces" {
			utils.HandleError(w, r, utils.NewAppError(http.StatusForbidden, err.Error()))
		} else {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin assigned successfully"})
}

func (h *WorkspaceHandler) RemoveAdmin(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	workspaceID, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace ID"))
		return
	}

	adminID, err := uuid.Parse(vars["admin_id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid admin ID"))
		return
	}

	err = h.workspaceService.RemoveAdminFromWorkspace(workspaceID, adminID, userID)
	if err != nil {
		if err.Error() == "only admins and super admins can remove admins from workspaces" {
			utils.HandleError(w, r, utils.NewAppError(http.StatusForbidden, err.Error()))
		} else {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin removed successfully"})
}

func (h *WorkspaceHandler) UpdateClientWorkspace(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.HandleError(w, r, utils.NewAppError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	clientID, err := uuid.Parse(vars["client_id"])
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid client ID"))
		return
	}

	var req UpdateClientWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	err = h.workspaceService.UpdateClientWorkspace(clientID, req.WorkspaceID, userID)
	if err != nil {
		if err.Error() == "only admins and super admins can update client workspace" {
			utils.HandleError(w, r, utils.NewAppError(http.StatusForbidden, err.Error()))
		} else {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client workspace updated successfully"})
}

