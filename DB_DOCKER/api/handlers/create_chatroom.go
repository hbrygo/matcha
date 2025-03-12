package handlers

import (
	"encoding/json"
	"errors"
	"matcha/api/handlers_utils.go"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// main handler function
func CreateChatroomHandler(w http.ResponseWriter, r *http.Request) {
	// check if method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get request data
	req, err := parseCreateChatroomRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// connect to database
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// check if users exist using utils
	if !handlers_utils.CheckUsersExist(db, req.UserIDs) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// create new chatroom using utils (MODIFIÉ)
	chatroomID, err := handlers_utils.CreateNewChatroom(tx, req.Name)
	if err != nil {
		http.Error(w, "Failed to create chatroom", http.StatusInternalServerError)
		return
	}

	// add participants using utils
	if err = handlers_utils.AddParticipantsToRoom(tx, chatroomID, req.UserIDs); err != nil {
		http.Error(w, "Failed to add participants", http.StatusInternalServerError)
		return
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// send response
	sendChatroomResponse(w, chatroomID)
}

// get data from request (MODIFIÉ)
func parseCreateChatroomRequest(r *http.Request) (*models.CreateChatroomRequest, error) {
	var req models.CreateChatroomRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// validate request
	if len(req.UserIDs) == 0 {
		return nil, errors.New("User IDs are required")
	}

	// check for private chat
	if len(req.UserIDs) < 2 {
		return nil, errors.New("Chatroom requires at least 2 users")
	}

	return &req, nil
}

// send success response
func sendChatroomResponse(w http.ResponseWriter, chatroomID int) {
	response := models.CreateChatroomResponse{
		Status:     "success",
		ChatroomID: chatroomID,
		Message:    "Chatroom created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
