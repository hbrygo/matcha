package handlers

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// GetRandomUserHandler récupère un utilisateur aléatoire selon les préférences
func GetRandomUserHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	var req models.GetRandomUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Valider l'UID
	if req.UID <= 0 {
		http.Error(w, "Invalid UID", http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Vérifier que l'utilisateur existe
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE uid = ?)", req.UID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Récupérer un utilisateur aléatoire
	randomUser, err := getRandomUser(db, req.UID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No compatible users found", http.StatusNotFound)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(randomUser)
}

// getRandomUser récupère un utilisateur aléatoire compatible avec les préférences
func getRandomUser(db *sql.DB, userID int) (*models.GetRandomUserResponse, error) {
	// Récupérer les préférences de l'utilisateur pour le filtrage
	var gender string
	err := db.QueryRow("SELECT gender FROM users WHERE uid = ?", userID).Scan(&gender)
	if err != nil {
		return nil, err
	}

	// Récupérer les préférences de l'utilisateur (maintenant dans une table séparée)
	rows, err := db.Query("SELECT preference FROM user_preferences WHERE user_uid = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []string
	for rows.Next() {
		var pref string
		if err := rows.Scan(&pref); err != nil {
			return nil, err
		}
		preferences = append(preferences, pref)
	}

	// Construire la partie de requête pour les préférences
	var prefCondition string
	if len(preferences) > 0 {
		prefCondition = "AND ("
		for i, pref := range preferences {
			if i > 0 {
				prefCondition += " OR "
			}
			prefCondition += "u.gender = '" + pref + "'"
		}
		prefCondition += ")"
	}

	// Récupérer un utilisateur aléatoire qui:
	// 1. N'est pas l'utilisateur actuel
	// 2. N'est pas déjà un match potentiel
	// 3. Correspond aux préférences de genre
	query := `
        SELECT u.uid, COALESCE(u.nom, '') as nom, COALESCE(u.prenom, '') as prenom, 
               COALESCE(u.dob, '') as dob, COALESCE(u.bio, '') as bio,
               COALESCE(u.latitude, '') as latitude, COALESCE(u.longitude, '') as longitude
        FROM users u
        WHERE u.uid != ? 
        AND NOT EXISTS (
            SELECT 1 FROM potential_matches 
            WHERE (liker_uid = ? AND liked_uid = u.uid)
            OR (liker_uid = u.uid AND liked_uid = ?)
        )
        ` + prefCondition + `
        ORDER BY RANDOM()
        LIMIT 1
    `

	var response models.GetRandomUserResponse
	var uid int
	var nom, prenom, dob, bio, latitude, longitude string

	err = db.QueryRow(query, userID, userID, userID).Scan(
		&uid, &nom, &prenom, &dob, &bio, &latitude, &longitude,
	)
	if err != nil {
		return nil, err
	}

	response.User.UID = uid
	response.User.Nom = nom
	response.User.Prenom = prenom
	response.User.DOB = dob
	response.User.Bio = bio
	response.User.Latitude = latitude
	response.User.Longitude = longitude

	// Récupérer les intérêts
	interestsRows, err := db.Query("SELECT interest FROM user_interests WHERE user_uid = ?", uid)
	if err != nil {
		return nil, err
	}
	defer interestsRows.Close()

	for interestsRows.Next() {
		var interest string
		if err := interestsRows.Scan(&interest); err != nil {
			return nil, err
		}
		response.User.Interests = append(response.User.Interests, interest)
	}

	// Récupérer les photos
	picturesRows, err := db.Query("SELECT picture_path FROM user_pictures WHERE user_uid = ?", uid)
	if err != nil {
		return nil, err
	}
	defer picturesRows.Close()

	for picturesRows.Next() {
		var picture string
		if err := picturesRows.Scan(&picture); err != nil {
			return nil, err
		}
		response.User.Pictures = append(response.User.Pictures, picture)
	}

	return &response, nil
}
