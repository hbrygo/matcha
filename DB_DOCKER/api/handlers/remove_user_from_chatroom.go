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
func RemoveUserFromChatroomHandler(w http.ResponseWriter, r *http.Request) {
	// check if method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get request data
	req, err := parseRemoveUserRequest(r)
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

	// check if chatroom exists
	if !handlers_utils.CheckChatroomExists(db, req.ChatroomID) {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}

	// check if user is in the chatroom
	if !handlers_utils.IsUserInChatroom(db, req.ChatroomID, req.UserID) {
		http.Error(w, "User not in chatroom", http.StatusBadRequest)
		return
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// remove user from chatroom
	if err = handlers_utils.RemoveUserFromChatroom(tx, req.ChatroomID, req.UserID); err != nil {
		http.Error(w, "Failed to remove user from chatroom", http.StatusInternalServerError)
		return
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// send success response
	sendRemoveUserResponse(w)
}

// parse request data
func parseRemoveUserRequest(r *http.Request) (*models.RemoveUserFromChatroomRequest, error) {
	var req models.RemoveUserFromChatroomRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// validate request
	if req.ChatroomID <= 0 {
		return nil, errors.New("Invalid chatroom ID")
	}

	if req.UserID <= 0 {
		return nil, errors.New("Invalid user ID")
	}

	return &req, nil
}

// send success response
func sendRemoveUserResponse(w http.ResponseWriter) {
	response := models.RemoveUserFromChatroomResponse{
		Status:  "success",
		Message: "User removed from chatroom successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
