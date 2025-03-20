package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
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

// Gestion des connexions SSE
var (
	sseClients = make(map[int][]chan string) // Map des chatrooms vers leurs canaux SSE
	mu         sync.Mutex                    // Mutex pour protéger l'accès à sseClients
)

// Fonction pour gérer les connexions SSE
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'UID de l'utilisateur depuis les cookies
	_, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "User not connected", http.StatusUnauthorized)
		return
	}

	// Récupérer l'ID de la chatroom depuis les paramètres de la requête
	roomIDStr := r.URL.Query().Get("roomID")
	if roomIDStr == "" {
		http.Error(w, "Chatroom ID is required", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid chatroom ID", http.StatusBadRequest)
		return
	}

	// Configurer les en-têtes pour SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Créer un canal pour envoyer les messages à cet utilisateur
	messageChan := make(chan string)
	defer close(messageChan)

	// Ajouter le canal à la liste des clients pour cette chatroom
	mu.Lock()
	sseClients[roomID] = append(sseClients[roomID], messageChan)
	mu.Unlock()

	// Envoyer les messages au client
	for msg := range messageChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		w.(http.Flusher).Flush()
	}
}

func broadcastMessageToRoom(roomID int, message string) {
	mu.Lock()
	defer mu.Unlock()

	if clients, ok := sseClients[roomID]; ok {
		// Create a new slice to hold active clients
		activeClients := make([]chan string, 0, len(clients))

		for _, client := range clients {
			select {
			case client <- message:
				// Message sent successfully, keep this client
				activeClients = append(activeClients, client)
			default:
				// Channel is blocked or closed, discard this client
				close(client)
			}
		}

		// Replace with only active clients
		sseClients[roomID] = activeClients
	}
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
