package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"matcha/database"
	"net/http"
)

type MeRequest struct {
	UID int `json:"uid"`
}

type MeResponse struct {
	User struct {
		Username   string   `json:"username"`
		Email      string   `json:"email"`
		Nom        string   `json:"nom"`
		Prenom     string   `json:"prenom"`
		DOB        string   `json:"dob"`
		Gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interests  []string `json:"interests"`
		Pictures   []string `json:"pictures"`
		Bio        string   `json:"bio"`
	} `json:"user"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// verify http method
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// decode request
	var req MeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("je suis ici\n")
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	// UID
	if req.UID <= 0 {
		http.Error(w, "UID invalide", http.StatusBadRequest)
		return
	}

	// connexion to db
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// create reponse
	var response MeResponse

	// get db info
	err = db.QueryRow(`
    SELECT 
        COALESCE(username, '') as username,
        COALESCE(email, '') as email,
        COALESCE(nom, '') as nom,
        COALESCE(prenom, '') as prenom,
        COALESCE(dob, '') as dob,
        COALESCE(gender, '') as gender,
        COALESCE(preference, '') as preference,
        COALESCE(bio, '') as bio
    FROM users 
    WHERE uid = ?`, req.UID).Scan(
		&response.User.Username,
		&response.User.Email,
		&response.User.Nom,
		&response.User.Prenom,
		&response.User.DOB,
		&response.User.Gender,
		&response.User.Preference,
		&response.User.Bio,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// get interests
	rows, err := db.Query("SELECT interest FROM user_interests WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var interests []string
	for rows.Next() {
		var interest string
		if err := rows.Scan(&interest); err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		interests = append(interests, interest)
	}
	response.User.Interests = interests

	// get pictures
	rows, err = db.Query("SELECT picture_path FROM user_pictures WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pictures []string
	for rows.Next() {
		var picture string
		if err := rows.Scan(&picture); err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		pictures = append(pictures, picture)
	}
	response.User.Pictures = pictures

	// send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
