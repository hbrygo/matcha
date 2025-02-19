package router

import (
	"log"
	"matcha/api/handlers"
	"matcha/middleware"
	"net/http"
)

/**

here set the headers rules and the handlers for each route

**/

func Router(mux *http.ServeMux) http.Handler {
	return middleware.GeneralMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Router log %s %s %s", r.RemoteAddr, r.Method, r.URL)

		switch r.URL.Path {
		case "/":

			/*
				w.Header().Set("Cache-Control", "no-cache")
				if r.Method != http.MethodGet {
					http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
					return
				}*/
			handlers.RootHandler(w, r)

		case "/testDBavailability":
			// Headers spécifiques pour le test DB
			/*
				w.Header().Set("Cache-Control", "no-store")
				if r.Method != http.MethodGet {
					http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
					return
				}*/
			handlers.TestDBHandler(w, r)

		case "/update":
			handlers.UpdateHandler(w, r)
		case "/me":
			handlers.MeHandler(w, r)
		case "/login":
			handlers.LoginHandler(w, r)
		case "/get_user":
			handlers.GetUserHandler(w, r)

		case "/create_user":
			handlers.CreateUserHandler(w, r)

		default:
			http.NotFound(w, r)
		}
	}))
}
