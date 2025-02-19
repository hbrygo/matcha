package middleware

import (
	"net/http"
)

const frontendURL = "http://localhost:8080"

func GeneralMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// basic security
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Configuration CORS
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// option headers
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Passe au handler suivant
		next.ServeHTTP(w, r)
	})
}

/*
// SecurityConfig contient les configurations de sécurité
type SecurityConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
}
*/
/*
// GlobalMiddleware gère la sécurité générale et les conditions spécifiques par route
func GlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// base securit config headers for all routes
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// CORS for front end
		w.Header().Set("Access-Control-Allow-Origin", frontendurl)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// treatement based on route path
		switch r.URL.Path {
		case "/testDBavailability":
			// test db endpoint
			if r.Method != http.MethodGet {
				http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
				return
			}

		case "/create_user":
			// user creation
			if r.Method != http.MethodPost {
				http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
				return
			}
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Content-Type doit être application/json", http.StatusUnsupportedMediaType)
				return
			}

		case "/me":
			// authentification
			if r.Header.Get("Authorization") == "" {
				http.Error(w, "Token d'authentification requis", http.StatusUnauthorized)
				return
			}
		}

		// go to handler if everything is ok
		next.ServeHTTP(w, r)
	})
}
*/
