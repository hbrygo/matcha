package handlers

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req models.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	if req.UID <= 0 {
		http.Error(w, "UID invalide", http.StatusBadRequest)
		return
	}

	// verif input
	if req.Nom == "" || req.Prenom == "" || req.DOB == "" || req.Gender == "" ||
		req.Preference == "" || req.Bio == "" {
		fmt.Printf("req.nom: %v\n", req.Nom)
		fmt.Printf("req.prenom: %v\n", req.Prenom)
		fmt.Printf("req.DOB: %v\n", req.DOB)
		fmt.Printf("req.Gender: %v\n", req.Gender)
		fmt.Printf("req.Preference: %v\n", req.Preference)
		fmt.Printf("req.bio: %v\n", req.Bio)
		http.Error(w, "Champs obligatoires manquants", http.StatusBadRequest)
		return
	}

	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var exists bool // look for user before update
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE uid = ?)", req.UID).Scan(&exists)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // roolback sort of retrive

	// update base
	result, err := tx.Exec(`
        UPDATE users 
        SET nom = ?, prenom = ?, dob = ?, gender = ?, preference = ?, bio = ?
        WHERE uid = ?`,
		req.Nom, req.Prenom, req.DOB, req.Gender, req.Preference, req.Bio, req.UID)

	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// verify update of a line
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Mise à jour impossible", http.StatusInternalServerError)
		return
	}

	// remove old interests
	_, err = tx.Exec("DELETE FROM user_interests WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// insert new interests
	for _, interest := range req.Interests {
		_, err = tx.Exec("INSERT INTO user_interests (user_uid, interest) VALUES (?, ?)",
			req.UID, interest)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
	}

	// remove old pictures
	_, err = tx.Exec("DELETE FROM user_pictures WHERE user_uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// insert new pics
	for _, picture := range req.Pictures {
		_, err = tx.Exec("INSERT INTO user_pictures (user_uid, picture_path) VALUES (?, ?)",
			req.UID, picture)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec("UPDATE users SET first_step = TRUE WHERE uid = ?", req.UID)
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// validate transactions
	if err = tx.Commit(); err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// prepare rep
	response := models.UpdateResponse{
		Status:  "success",
		Message: "Informations mises à jour avec succès",
		User:    &req,
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
