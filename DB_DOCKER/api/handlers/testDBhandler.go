package handlers

import (
	"encoding/json"
	"matcha/database"
	"net/http"
)

// Handler pour la route "/testDBavailability"
func TestDBHandler(w http.ResponseWriter, r *http.Request) {
	// Vérification de la méthode HTTP
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	db, err := database.InitDB()
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Erreur lors de l'initialisation de la base de données: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer db.Close()

	// Vérifier si la connexion à la DB est active
	err = db.Ping()
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "La base de données n'est pas accessible: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Si tout va bien
	response := Response{
		Status:  "success",
		Message: "Base de données opérationnelle et prête à l'emploi",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
