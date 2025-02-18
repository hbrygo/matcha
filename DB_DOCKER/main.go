package main

import (
	"log"
	//"matcha/api/handlers"
	"matcha/database"
	//"matcha/middleware"
	"matcha/router"
	"net/http"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	// server mux
	mux := http.NewServeMux()

	// middleware
	//finalHandler := middleware.GlobalMiddleware(mux)
	finalHandler := router.Router(mux)

	// routes
	//mux.Handle("/", http.HandlerFunc(handlers.RootHandler))
	//mux.Handle("/testDBavailability", http.HandlerFunc(handlers.TestDBHandler))

	// server configuration
	server := &http.Server{
		Addr:    ":8181",
		Handler: finalHandler,
	}

	// Logs
	log.Printf("Serveur démarré sur http://localhost:8181")
	log.Printf("Test DB disponible sur http://localhost:8181/testDBavailability")

	// start
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
