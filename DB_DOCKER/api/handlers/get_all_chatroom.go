package handlers

import (
	"database/sql"
	"encoding/json"
	"matcha/database"
	"net/http"
)

// GetAllChatroomHandler récupère toutes les chatrooms existantes
func GetAllChatroomHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Connexion à la base de données
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Récupérer toutes les chatrooms
	chatroomIDs, err := getAllChatrooms(db)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Préparer la réponse
	response := struct {
		ChatRoom struct {
			ChatRoomID []int `json:"chatRoomID"`
		} `json:"chatRoom"`
	}{
		ChatRoom: struct {
			ChatRoomID []int `json:"chatRoomID"`
		}{
			ChatRoomID: chatroomIDs,
		},
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// getAllChatrooms récupère les IDs de toutes les chatrooms dans la base de données
func getAllChatrooms(db *sql.DB) ([]int, error) {
	// Exécuter la requête
	rows, err := db.Query(`
        SELECT chatroom_id 
        FROM chat_rooms
        ORDER BY chatroom_id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les résultats
	var chatroomIDs []int
	for rows.Next() {
		var chatroomID int
		if err := rows.Scan(&chatroomID); err != nil {
			return nil, err
		}
		chatroomIDs = append(chatroomIDs, chatroomID)
	}

	// Vérifier les erreurs de parcours
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatroomIDs, nil
}
