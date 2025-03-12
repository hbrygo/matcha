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
func NewMessageHandler(w http.ResponseWriter, r *http.Request) {
	// check if method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get request data
	req, err := parseNewMessageRequest(r)
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

	// check if user exists
	if !handlers_utils.CheckUsersExist(db, []int{req.UserID}) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// check if chatroom exists
	if !handlers_utils.CheckChatroomExists(db, req.ChatRoom) {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}

	// check if user is in chatroom
	if !handlers_utils.IsUserInChatroom(db, req.ChatRoom, req.UserID) {
		http.Error(w, "User not in chatroom", http.StatusForbidden)
		return
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// add message to db
	if err = handlers_utils.AddMessageToChatroom(tx, req.ChatRoom, req.UserID, req.Message); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// send success response
	sendNewMessageResponse(w)
}

// parse request data
func parseNewMessageRequest(r *http.Request) (*models.NewMessageRequest, error) {
	var req models.NewMessageRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// validate request
	if req.Message == "" {
		return nil, errors.New("Message cannot be empty")
	}

	if req.ChatRoom <= 0 {
		return nil, errors.New("Invalid chatroom ID")
	}

	if req.UserID <= 0 {
		return nil, errors.New("Invalid user ID")
	}

	return &req, nil
}

// send success response
func sendNewMessageResponse(w http.ResponseWriter) {
	response := models.NewMessageResponse{
		Status:  "success",
		Message: "Message added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
