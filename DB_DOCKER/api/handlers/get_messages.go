package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"matcha/api/handlers_utils.go"
	"matcha/api/models"
	"matcha/database"
	"net/http"
)

// main handler function
func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	// check if method is post
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lire et afficher le corps brut de la requête
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Corps brut de la requête : %s\n", string(body))

	// Réinitialiser le lecteur pour le décodage JSON
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Décoder le JSON
	req, err := parseGetMessageRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("JSON parsé : %+v\n", req)

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

	// get messages from db
	messages, err := handlers_utils.GetChatroomMessages(db, req.ChatRoom)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	// send success response
	sendGetMessageResponse(w, messages)
}

// parse request data
func parseGetMessageRequest(r *http.Request) (*models.GetMessageRequest, error) {
	var req models.GetMessageRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	// validate request
	if req.ChatRoom <= 0 {
		return nil, errors.New("Invalid chatroom ID")
	}

	if req.UserID <= 0 {
		return nil, errors.New("Invalid user ID")
	}

	return &req, nil
}

// send success response
func sendGetMessageResponse(w http.ResponseWriter, messages []models.Message) {
	response := models.GetMessageResponse{
		Status:   "success",
		Messages: messages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
