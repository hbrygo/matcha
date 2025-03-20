package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Structures existantes
type IncomingMessage struct {
	UID    string `json:"userID"`
	Msg    string `json:"message"`
	RoomID string `json:"chatRoom"`
}

type Message struct {
	UID    int    `json:"userID"`
	Msg    string `json:"message"`
	RoomID int    `json:"chatRoom"`
}

// Fonction existante NewMessage (ajout de la diffusion SSE)
func NewMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("NEW MESSAGE\n")

	// Vérifier si l'utilisateur est connecté
	_, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "User not connected", http.StatusUnauthorized)
		return
	}

	// Lire le corps de la requête entrante
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Désérialiser le JSON reçu
	var incomingMessage IncomingMessage
	if err := json.Unmarshal(body, &incomingMessage); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Convertir les champs en entiers
	roomID, err := strconv.Atoi(incomingMessage.RoomID)
	if err != nil {
		http.Error(w, "Invalid chatRoom ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(incomingMessage.UID)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Créer un message avec les types corrects
	outgoingMessage := Message{
		UID:    userID,
		Msg:    incomingMessage.Msg,
		RoomID: roomID,
	}

	// Convertir le message en JSON
	messageJSON, err := json.Marshal(outgoingMessage)
	if err != nil {
		http.Error(w, "Failed to encode message", http.StatusInternalServerError)
		return
	}

	// Transmettre la requête à l'API externe
	resp, err := http.Post("http://localhost:8181/new_message", "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		http.Error(w, "Failed to contact external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lire la réponse de l'API externe
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external API", http.StatusInternalServerError)
		return
	}

	// Diffuser le message via SSE à tous les utilisateurs de la chatroom
	broadcastMessageToRoom(roomID, string(messageJSON))

	// Renvoyer la réponse au client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
