package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// AppError representa un error de la aplicación
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError crea un nuevo error de aplicación
func NewAppError(code int, message string, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// ErrorResponse estructura de respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Code    int    `json:"code"`
}

// HandleError maneja errores y envía respuesta HTTP apropiada
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int
	var errorMsg string
	var details string

	switch e := err.(type) {
	case *AppError:
		statusCode = e.Code
		errorMsg = e.Message
		details = e.Details
	default:
		// Error genérico
		statusCode = http.StatusInternalServerError
		errorMsg = "Internal server error"
		details = err.Error()
		log.Printf("Error: %v - Path: %s - Method: %s", err, r.URL.Path, r.Method)
	}

	response := ErrorResponse{
		Error:   errorMsg,
		Details: details,
		Code:    statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Errores comunes
var (
	ErrUnauthorized       = NewAppError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden          = NewAppError(http.StatusForbidden, "Forbidden")
	ErrNotFound           = NewAppError(http.StatusNotFound, "Resource not found")
	ErrBadRequest         = NewAppError(http.StatusBadRequest, "Bad request")
	ErrInternalError      = NewAppError(http.StatusInternalServerError, "Internal server error")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid credentials")
	ErrUserExists         = NewAppError(http.StatusConflict, "User already exists")
	ErrInvalidToken       = NewAppError(http.StatusUnauthorized, "Invalid or expired token")
)
