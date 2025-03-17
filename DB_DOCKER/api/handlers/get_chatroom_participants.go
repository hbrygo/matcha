package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"matcha/api/handlers_utils.go"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// GetChatroomParticipantsHandler récupère tous les participants d'une chatroom spécifique
func GetChatroomParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	req, err := parseGetChatroomParticipantsRequest(r)
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

	//  chatroom existe
	if !handlers_utils.CheckChatroomExists(db, req.ChatroomID) {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}

	// Récupérer tous les participants de la chatroom
	participants, err := getChatroomParticipants(db, req.ChatroomID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Préparer la réponse
	response := models.GetChatroomParticipantsResponse{
		Participants: participants,
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// parseGetChatroomParticipantsRequest analyse la requête pour extraire l'ID de la chatroom
func parseGetChatroomParticipantsRequest(r *http.Request) (*models.GetChatroomParticipantsRequest, error) {
	var req models.GetChatroomParticipantsRequest

	// Décoder le corps de la requête
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// Valider les données
	if req.ChatroomID <= 0 {
		return nil, errors.New("Chatroom ID is required and must be positive")
	}

	return &req, nil
}

// getChatroomParticipants récupère les IDs de tous les participants d'une chatroom
func getChatroomParticipants(db *sql.DB, chatroomID int) ([]int, error) {
	rows, err := db.Query(`
        SELECT user_uid 
        FROM chat_participants 
        WHERE chatroom_id = ?
    `, chatroomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	var participants []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}
