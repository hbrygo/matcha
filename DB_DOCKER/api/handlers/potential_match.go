package handlers

import (
	"database/sql"
	"encoding/json"
	//"errors"
	"matcha/api/handlers_utils.go"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// PotentialMatchHandler gère l'ajout d'un "like" potentiel entre utilisateurs
func PotentialMatchHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire et valider les données de la requête
	var req models.PotentialMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Valider les UIDs
	if req.MyUID <= 0 || req.OtherUID <= 0 {
		http.Error(w, "Invalid user IDs", http.StatusBadRequest)
		return
	}

	if req.MyUID == req.OtherUID {
		http.Error(w, "Cannot match with yourself", http.StatusBadRequest)
		return
	}

	// Connexion à la base de données
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Vérifier que les deux utilisateurs existent
	if !handlers_utils.CheckUsersExist(db, []int{req.MyUID, req.OtherUID}) {
		http.Error(w, "One or both users not found", http.StatusNotFound)
		return
	}

	// Traiter le potential match
	result, err := handlePotentialMatch(db, req.MyUID, req.OtherUID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Envoyer la réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handlePotentialMatch ajoute un potential match et vérifie s'il y a un match réciproque
func handlePotentialMatch(db *sql.DB, likerUID, likedUID int) (*models.PotentialMatchResponse, error) {
	// Commencer une transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Vérifier si un match réciproque existe déjà
	var matchExists bool
	err = tx.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM potential_matches
            WHERE liker_uid = ? AND liked_uid = ?
        )
    `, likedUID, likerUID).Scan(&matchExists)
	if err != nil {
		return nil, err
	}

	// Vérifier si ce liker a déjà liké cette personne
	var alreadyLiked bool
	err = tx.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM potential_matches
            WHERE liker_uid = ? AND liked_uid = ?
        )
    `, likerUID, likedUID).Scan(&alreadyLiked)
	if err != nil {
		return nil, err
	}

	response := &models.PotentialMatchResponse{}

	if alreadyLiked {
		// Déjà liké, pas besoin de faire quoi que ce soit
		response.Status = "success"
		response.Message = "Already liked this user"
		response.IsMatch = matchExists
	} else {
		// Ajouter le nouveau "like"
		_, err = tx.Exec(`
            INSERT INTO potential_matches (liker_uid, liked_uid)
            VALUES (?, ?)
        `, likerUID, likedUID)
		if err != nil {
			return nil, err
		}

		if matchExists {
			// C'est un match!
			response.Status = "success"
			response.Message = "It's a match!"
			response.IsMatch = true

			// Créer une chatroom pour les deux utilisateurs s'ils n'en ont pas déjà une
			var chatroomExists bool
			err = tx.QueryRow(`
                SELECT EXISTS(
                    SELECT 1 FROM chat_participants cp1
                    JOIN chat_participants cp2 ON cp1.chatroom_id = cp2.chatroom_id
                    WHERE cp1.user_uid = ? AND cp2.user_uid = ? AND cp1.user_uid != cp2.user_uid
                )
            `, likerUID, likedUID).Scan(&chatroomExists)
			if err != nil {
				return nil, err
			}

			if !chatroomExists {
				// Créer une nouvelle chatroom pour les deux utilisateurs
				chatroomName := "Match Chat"
				chatroomID, err := handlers_utils.CreateNewChatroom(tx, chatroomName)
				if err != nil {
					return nil, err
				}

				// Ajouter les deux utilisateurs à la chatroom
				err = handlers_utils.AddParticipantsToRoom(tx, chatroomID, []int{likerUID, likedUID})
				if err != nil {
					return nil, err
				}
			}
		} else {
			// Juste un potential match pour l'instant
			response.Status = "success"
			response.Message = "Potential match added"
			response.IsMatch = false
		}
	}

	// Committer la transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return response, nil
}
