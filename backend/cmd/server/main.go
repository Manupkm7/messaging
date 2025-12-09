package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"
	"backend/internal/websocket"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuración
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Conectar a la base de datos
	if err := database.Connect(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer database.Close()

	// Inicializar WebSocket Hub
	websocket.GlobalHub = websocket.NewHub()
	go websocket.GlobalHub.Run()

	// Iniciar servicio de limpieza
	cleanupService := services.NewCleanupService()
	cleanupService.StartCleanupJob()

	// Configurar rutas
	router := mux.NewRouter()

	// Handlers
	authHandler := handlers.NewAuthHandler()
	chatHandler := handlers.NewChatHandler()
	userHandler := handlers.NewUserHandler()
	adminGroupHandler := handlers.NewAdminGroupHandler()
	workspaceHandler := handlers.NewWorkspaceHandler()

	// Rutas públicas
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Rutas protegidas
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/auth/me", authHandler.GetMe).Methods("GET")

	// Rutas de usuarios (requieren admin o super_admin)
	protected.HandleFunc("/users/admin", userHandler.CreateAdmin).Methods("POST")
	protected.Use(middleware.RequireRole("admin", "super_admin"))

	// Rutas de chats (todos los usuarios autenticados)
	apiProtected := api.PathPrefix("").Subrouter()
	apiProtected.Use(middleware.AuthMiddleware)

	apiProtected.HandleFunc("/chats", chatHandler.GetChats).Methods("GET")
	apiProtected.HandleFunc("/chats", chatHandler.CreateChat).Methods("POST")
	apiProtected.HandleFunc("/chats/{id}/delegate", chatHandler.DelegateChat).Methods("POST")
	apiProtected.HandleFunc("/chats/{id}/messages", chatHandler.GetMessages).Methods("GET")
	apiProtected.HandleFunc("/chats/{id}/messages", chatHandler.SendMessage).Methods("POST")

	// Handler para archivos de mensajes
	messageHandler := handlers.NewMessageHandler()
	apiProtected.HandleFunc("/messages/{id}/file", messageHandler.GetMessageFile).Methods("GET")

	// Rutas de grupos de admins
	apiProtected.HandleFunc("/admin-groups", adminGroupHandler.CreateGroup).Methods("POST")
	apiProtected.HandleFunc("/admin-groups", adminGroupHandler.GetGroups).Methods("GET")
	apiProtected.HandleFunc("/admin-groups/{id}/replace", adminGroupHandler.ReplaceAdmin).Methods("POST")

	// Rutas de workspaces (requieren admin o super_admin)
	apiProtected.HandleFunc("/workspaces", workspaceHandler.CreateWorkspace).Methods("POST", "OPTIONS")
	apiProtected.HandleFunc("/workspaces", workspaceHandler.GetWorkspaces).Methods("GET", "OPTIONS")
	apiProtected.HandleFunc("/workspaces/{id}", workspaceHandler.GetWorkspace).Methods("GET", "OPTIONS")
	apiProtected.HandleFunc("/workspaces/{id}/admins", workspaceHandler.AssignAdmin).Methods("POST", "OPTIONS")
	apiProtected.HandleFunc("/workspaces/{id}/admins/{admin_id}", workspaceHandler.RemoveAdmin).Methods("DELETE", "OPTIONS")
	apiProtected.HandleFunc("/clients/{client_id}/workspace", workspaceHandler.UpdateClientWorkspace).Methods("PUT", "OPTIONS")

	// WebSocket
	router.HandleFunc("/ws", websocket.HandleWebSocket)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Aplicar CORS middleware como wrapper del router (ANTES de iniciar el servidor)
	corsHandler := middleware.CORSMiddleware(router)

	// Iniciar servidor
	port := config.AppConfig.ServerPort
	cfg := config.AppConfig

	if cfg.EnableTLS {
		log.Printf("Server starting with TLS on port %s", port)

		// Configurar TLS
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}

		server := &http.Server{
			Addr:         ":" + port,
			Handler:      corsHandler,
			TLSConfig:    tlsConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		}

		log.Fatal(server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile))
	} else {
		log.Printf("Server starting on port %s (HTTP only)", port)
		log.Fatal(http.ListenAndServe(":"+port, corsHandler))
	}
}
