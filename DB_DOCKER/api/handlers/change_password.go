package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"matcha/api/models"
	"matcha/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// ChangePasswordHandler gère la modification du mot de passe d'un utilisateur
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	req, err := parseChangePasswordRequest(r)
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

	// Changer le mot de passe
	if err := changeUserPassword(db, req); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else if err.Error() == "incorrect password" {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
		} else if err.Error() == "invalid new password" {
			http.Error(w, "New password does not meet requirements", http.StatusBadRequest)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Envoyer la réponse
	response := models.Response{
		Status:  "success",
		Message: "Password changed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Structure de requête pour changer le mot de passe
type changePasswordRequest struct {
	UserID      int    `json:"userID"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Analyse la requête pour extraire les données de changement de mot de passe
func parseChangePasswordRequest(r *http.Request) (*changePasswordRequest, error) {
	var req changePasswordRequest

	// Décoder le corps de la requête
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// Valider les données
	if req.UserID <= 0 {
		return nil, errors.New("User ID is required and must be positive")
	}

	if req.OldPassword == "" {
		return nil, errors.New("Old password is required")
	}

	if req.NewPassword == "" {
		return nil, errors.New("New password is required")
	}

	// Vérifier que le nouveau mot de passe est suffisamment complexe
	if len(req.NewPassword) < 8 {
		return nil, errors.New("New password must be at least 8 characters long")
	}

	return &req, nil
}

// Change le mot de passe d'un utilisateur après vérification de l'ancien
func changeUserPassword(db *sql.DB, req *changePasswordRequest) error {
	// Récupérer le mot de passe actuel de la base de données
	var currentHashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE uid = ?", req.UserID).Scan(&currentHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Vérifier que l'ancien mot de passe correspond
	if !checkPasswordHash(req.OldPassword, currentHashedPassword) {
		return errors.New("incorrect password")
	}

	// Hacher le nouveau mot de passe
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Mettre à jour le mot de passe dans la base de données
	_, err = db.Exec("UPDATE users SET password = ? WHERE uid = ?", hashedPassword, req.UserID)
	if err != nil {
		return err
	}

	return nil
}

// Hacher un mot de passe
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Vérifier si un mot de passe correspond à son hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
