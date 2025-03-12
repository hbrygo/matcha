package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// GetUserByUsernameHandler récupère les informations d'un utilisateur en fonction de son nom d'utilisateur
func GetUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	req, err := parseGetUserByUsernameRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Récupérer les informations de l'utilisateur
	response, err := getUserByUsername(db, req.Username)

	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Analyse la requête pour extraire le nom d'utilisateur
func parseGetUserByUsernameRequest(r *http.Request) (*models.GetUserByUsernameRequest, error) {
	var req models.GetUserByUsernameRequest

	// Décoder le corps de la requête
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// Valider les données
	if req.Username == "" {
		return nil, errors.New("Username is required")
	}

	return &req, nil
}

// Récupère les informations de l'utilisateur depuis la base de données
func getUserByUsername(db *sql.DB, username string) (*models.GetUserResponse, error) {
	// Préparation de la réponse
	var response models.GetUserResponse

	// Récupérer les données de base de l'utilisateur
	var uid int
	err := db.QueryRow(`
        SELECT uid, nom, prenom, dob, gender, preference, bio
        FROM users
        WHERE username = ?
    `, username).Scan(&uid, &response.User.Nom, &response.User.Prenom, &response.User.DOB, &response.User.Gender, &response.User.Preference, &response.User.Bio)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Récupérer les intérêts
	rows, err := db.Query(`
        SELECT interest 
        FROM user_interests
        WHERE user_uid = ?
    `, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []string
	for rows.Next() {
		var interest string
		if err := rows.Scan(&interest); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	response.User.Interests = interests

	// Récupérer les photos
	rows, err = db.Query(`
        SELECT picture_path 
        FROM user_pictures
        WHERE user_uid = ?
    `, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pictures []string
	for rows.Next() {
		var picture string
		if err := rows.Scan(&picture); err != nil {
			return nil, err
		}
		pictures = append(pictures, picture)
	}
	response.User.Pictures = pictures

	return &response, nil
}
