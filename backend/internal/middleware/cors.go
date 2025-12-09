package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	// Construir lista de orígenes permitidos
	allowed := make(map[string]struct{})
	
	// Agregar localhost por defecto (siempre permitido en desarrollo)
	allowed["http://localhost:5173"] = struct{}{}
	allowed["http://127.0.0.1:5173"] = struct{}{}
	allowed["https://localhost:5173"] = struct{}{}
	allowed["https://127.0.0.1:5173"] = struct{}{}
	
	// Agregar BASE_URL_FRONT si está configurado
	baseUrlFront := strings.TrimRight(strings.TrimSpace(os.Getenv("BASE_URL_FRONT")), "/")
	if baseUrlFront != "" {
		allowed[baseUrlFront] = struct{}{}
		log.Printf("CORS: Added BASE_URL_FRONT to allowed origins: %s", baseUrlFront)
	}
	fmt.Println("allowed", allowed)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimRight(strings.TrimSpace(r.Header.Get("Origin")), "/")

		fmt.Println("origin", origin)
		fmt.Println("allowed", allowed)
		// Manejar preflight OPTIONS PRIMERO, antes de verificar origen
		fmt.Println("r.Method", r.Method)
		if r.Method == http.MethodOptions {
			// Si hay Origin, verificar si está permitido
			if origin != "" {
				if _, ok := allowed[origin]; !ok {
					log.Printf("CORS: Origin not allowed for preflight: %s", origin)
					w.Header().Set("Access-Control-Allow-Origin", "null")
					w.Header().Set("Vary", "Origin")
					http.Error(w, "CORS origin not allowed", http.StatusForbidden)
					return
				}
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// Si no hay Origin en preflight, usar localhost por defecto (puede ser same-origin)
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			}
			
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			
			// Headers permitidos: usar los solicitados o los por defecto
			reqHeaders := r.Header.Get("Access-Control-Request-Headers")
			if reqHeaders == "" {
				reqHeaders = "Content-Type, Authorization, X-Requested-With, Accept, Origin"
			}
			w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Para peticiones no-OPTIONS, verificar origen
		if origin == "" {
			// Si no hay Origin header, permitir la petición (puede ser same-origin)
			next.ServeHTTP(w, r)
			return
		}

		// Verificar si el origen está permitido
		if _, ok := allowed[origin]; !ok {
			log.Printf("CORS: Origin not allowed: %s", origin)
			w.Header().Set("Access-Control-Allow-Origin", "null")
			w.Header().Set("Vary", "Origin")
			http.Error(w, "CORS origin not allowed", http.StatusForbidden)
			return
		}

		// Establecer headers CORS para peticiones normales
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}
