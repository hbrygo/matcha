package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// GetMyChatroomHandler récupère toutes les chatrooms auxquelles un utilisateur participe
func GetMyChatroomHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	req, err := parseGetMyChatroomRequest(r)
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

	// Vérifier que l'utilisateur existe
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE uid = ?)", req.UserID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Récupérer toutes les chatrooms de l'utilisateur
	chatrooms, err := getUserChatrooms(db, req.UserID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Préparer la réponse
	response := models.GetMyChatroomResponse{
		User: struct {
			ChatRoomIDs []int `json:"chatRoomID"`
		}{
			ChatRoomIDs: chatrooms,
		},
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// parseGetMyChatroomRequest analyse la requête pour extraire l'ID utilisateur
func parseGetMyChatroomRequest(r *http.Request) (*models.GetMyChatroomRequest, error) {
	var req models.GetMyChatroomRequest

	// Décoder le corps de la requête
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// Valider les données
	if req.UserID <= 0 {
		return nil, errors.New("User ID is required and must be positive")
	}

	return &req, nil
}

// getUserChatrooms récupère toutes les chatrooms auxquelles l'utilisateur participe
func getUserChatrooms(db *sql.DB, userID int) ([]int, error) {
	// Exécuter la requête
	rows, err := db.Query(`
        SELECT chatroom_id 
        FROM chat_participants 
        WHERE user_uid = ?
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	var chatrooms []int
	for rows.Next() {
		var chatroomID int
		if err := rows.Scan(&chatroomID); err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, chatroomID)
	}

	// Vérifier les erreurs de parcours
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatrooms, nil
}
