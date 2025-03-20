package chatRoom

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

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
