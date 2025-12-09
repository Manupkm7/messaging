package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/google/uuid"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userService: services.NewUserService(),
	}
}

type RegisterRequest struct {
	Username      string  `json:"username"`
	Email         *string `json:"email,omitempty"`
	Password      string  `json:"password"`
	Role          string  `json:"role"`
	WorkspaceSlug string  `json:"workspace_slug,omitempty"` // Para clientes que se registran con slug
}

type LoginRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	SessionToken  string `json:"session_token,omitempty"`  // Para clientes
	WorkspaceSlug string `json:"workspace_slug,omitempty"` // Para clientes que se registran/login con slug
}

type AuthResponse struct {
	Token        string       `json:"token"`
	User         *models.User `json:"user"`
	SessionToken string       `json:"session_token,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Validar rol
	if req.Role != string(models.RoleAdmin) && req.Role != string(models.RoleClient) && req.Role != string(models.RoleSuperAdmin) {

		utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, `Invalid role value. Must be "super_admin", "admin", or "client". Got: `+req.Role))
		return
	}

	// Obtener el ID del usuario que crea (si es admin creando admin)
	var createdByID *uuid.UUID
	var workspaceID *uuid.UUID
	if req.Role == string(models.RoleAdmin) || req.Role == string(models.RoleSuperAdmin) {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			utils.HandleError(w, r, utils.NewAppError(http.StatusForbidden, "Only admins can create admin users"))
			return
		}
		createdByID = &userID
	}

	// Si es cliente y tiene workspace_slug, obtener el workspace_id
	if req.Role == string(models.RoleClient) && req.WorkspaceSlug != "" {
		workspaceService := services.NewWorkspaceService()
		workspace, err := workspaceService.GetWorkspaceBySlug(req.WorkspaceSlug)
		if err != nil {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace slug"))
			return
		}
		workspaceID = &workspace.ID
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password, req.Role, createdByID, workspaceID)
	if err != nil {
		if err.Error() == "username already exists" {
			utils.HandleError(w, r, utils.ErrUserExists)
		} else {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, err.Error()))
		}
		return
	}

	// Generar JWT
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		utils.HandleError(w, r, utils.NewAppError(http.StatusInternalServerError, "Error generating token", err.Error()))
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	// Para clientes, incluir session token y establecer cookie
	if user.Role == models.RoleClient && user.SessionToken != "" {
		response.SessionToken = user.SessionToken
		secure := config.AppConfig.EnableTLS
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    user.SessionToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure, // true cuando TLS está habilitado
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 año
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Para clientes, verificar session token primero
	if req.SessionToken != "" {
		user, err := h.userService.GetUserBySessionToken(req.SessionToken)
		if err == nil && user.Role == models.RoleClient {
			// Si hay workspace_slug, verificar que coincida
			if req.WorkspaceSlug != "" {
				workspaceService := services.NewWorkspaceService()
				workspace, err := workspaceService.GetWorkspaceBySlug(req.WorkspaceSlug)
				if err != nil || workspace.ID != *user.WorkspaceID {
					utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace slug"))
					return
				}
			}

			token, err := utils.GenerateToken(user.ID, string(user.Role))
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			response := AuthResponse{
				Token:        token,
				User:         user,
				SessionToken: user.SessionToken,
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    user.SessionToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(365 * 24 * time.Hour),
			})

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Login normal con username y password
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		//console 
		log.Println("Login error:", err)

		utils.HandleError(w, r, utils.ErrInvalidCredentials)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		log.Println("Login error: invalid password for user", req.Username)
		
		utils.HandleError(w, r, utils.ErrInvalidCredentials)
		return
	}

	// Para clientes, verificar workspace_slug si se proporciona
	if user.Role == models.RoleClient && req.WorkspaceSlug != "" {
		if user.WorkspaceID == nil {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "User is not associated with a workspace"))
			return
		}
		workspaceService := services.NewWorkspaceService()
		workspace, err := workspaceService.GetWorkspaceBySlug(req.WorkspaceSlug)
		if err != nil || workspace.ID != *user.WorkspaceID {
			utils.HandleError(w, r, utils.NewAppError(http.StatusBadRequest, "Invalid workspace slug"))
			return
		}
	}

	// Para clientes, verificar o actualizar session token
	if user.Role == models.RoleClient {
		if user.SessionToken == "" {
			sessionToken, err := utils.GenerateSessionToken()
			if err != nil {
				http.Error(w, "Error generating session token", http.StatusInternalServerError)
				return
			}
			user.SessionToken = sessionToken
			h.userService.UpdateSessionToken(user.ID, sessionToken)
		}
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	if user.Role == models.RoleClient && user.SessionToken != "" {
		response.SessionToken = user.SessionToken
		secure := config.AppConfig.EnableTLS
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    user.SessionToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(365 * 24 * time.Hour),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Eliminar cookie de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// No enviar password hash
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
