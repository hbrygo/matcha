package handlers

import (
	"database/sql"
	"encoding/json"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode request
	var req models.GetUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// val UID
	if req.UID <= 0 {
		http.Error(w, "Invalid UID", http.StatusBadRequest)
		return
	}

	// initi db connexion
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// prepare rseponse
	var response models.GetUserResponse

	err = db.QueryRow(`
    SELECT 
        COALESCE(nom, '') as nom,
        COALESCE(prenom, '') as prenom,
        COALESCE(dob, '') as dob,
        COALESCE(gender, '') as gender,
        COALESCE(preference, '') as preference,
        COALESCE(bio, '') as bio
    FROM users 
    WHERE uid = ?`, req.UID).Scan(
		&response.User.Nom,
		&response.User.Prenom,
		&response.User.DOB,
		&response.User.Gender,
		&response.User.Preference,
		&response.User.Bio,
	)

	// handle database query res
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// user interests
	rows, err := db.Query("SELECT interest FROM user_interests WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// process interests results
	for rows.Next() {
		var interest string
		if err := rows.Scan(&interest); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		response.User.Interests = append(response.User.Interests, interest)
	}

	// get user pictures
	rows, err = db.Query("SELECT picture_path FROM user_pictures WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// process pictures results
	for rows.Next() {
		var picture string
		if err := rows.Scan(&picture); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		response.User.Pictures = append(response.User.Pictures, picture)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
