package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	notificationClients = make(map[int]chan string) // Map des utilisateurs vers leurs canaux SSE de notification
	notificationMu      sync.Mutex                  // Mutex pour protéger l'accès à notificationClients
)

// Fonction pour gérer les connexions SSE de notification
func NotificationSSEHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'UID de l'utilisateur depuis les cookies
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "User not connected", http.StatusUnauthorized)
		return
	}

	uid, err := strconv.Atoi(uidCookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Configurer les en-têtes pour SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Créer un canal pour envoyer les notifications à cet utilisateur
	notificationChan := make(chan string)
	defer close(notificationChan)

	// Ajouter le canal à la map des notifications
	notificationMu.Lock()
	// Fermer l'ancien canal si l'utilisateur se reconnecte
	if oldChan, exists := notificationClients[uid]; exists {
		close(oldChan)
	}
	notificationClients[uid] = notificationChan
	notificationMu.Unlock()

	// Envoyer les notifications au client
	for notification := range notificationChan {
		fmt.Fprintf(w, "data: %s\n\n", notification)
		w.(http.Flusher).Flush()
	}
}

// Fonction pour envoyer une notification à un utilisateur spécifique
func sendNotificationToUser(userID int, notification string) {
	notificationMu.Lock()
	defer notificationMu.Unlock()

	if client, ok := notificationClients[userID]; ok {
		select {
		case client <- notification:
			// Notification envoyée avec succès
		default:
			// Canal bloqué ou fermé, on le ferme
			close(client)
			delete(notificationClients, userID)
		}
	}
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

	// Diffuser le message aux clients connectés à la chatroom
	if clients, ok := sseClients[roomID]; ok {
		activeClients := make([]chan string, 0, len(clients))
		for _, client := range clients {
			select {
			case client <- message:
				activeClients = append(activeClients, client)
			default:
				close(client)
			}
		}
		sseClients[roomID] = activeClients
	}
	mu.Unlock()

	// Obtenir les participants de la chatroom pour les notifications
	// Note: on doit faire cette partie en dehors du mutex pour éviter un deadlock
	participants, err := getChatroomParticipantsForNotification(roomID)
	if err != nil {
		fmt.Printf("Erreur lors de la récupération des participants: %v\n", err)
		return
	}

	// Convertir le message en objet pour extraire l'expéditeur
	var msgObj Message
	if err := json.Unmarshal([]byte(message), &msgObj); err != nil {
		fmt.Printf("Erreur lors du parsing du message: %v\n", err)
		return
	}

	// Créer un objet de notification
	notification := map[string]interface{}{
		"type":      "new_message",
		"roomID":    roomID,
		"senderID":  msgObj.UID,
		"senderMsg": msgObj.Msg,
		"timestamp": time.Now().Unix(),
	}

	notificationJSON, _ := json.Marshal(notification)

	// Envoyer des notifications à tous les participants (sauf l'expéditeur)
	for _, participantID := range participants {
		if participantID != msgObj.UID {
			sendNotificationToUser(participantID, string(notificationJSON))
		}
	}
}

// Fonction pour obtenir les participants d'une chatroom pour les notifications
func getChatroomParticipantsForNotification(roomID int) ([]int, error) {
	// Créer une requête pour obtenir les participants de la chatroom
	postBody, err := json.Marshal(map[string]int{
		"chatRoomID": roomID,
	})
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la création du corps de la requête: %v", err)
	}

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:8181/get_chatroom_participants", "application/json", responseBody)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la requête pour obtenir les participants: %v", err)
	}
	defer resp.Body.Close()

	// Lire la réponse de la requête
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la lecture de la réponse: %v", err)
	}

	// Parser la réponse pour obtenir les participants
	var participantsResponse GetChatroomParticipantsResponse
	if err := json.Unmarshal(response, &participantsResponse); err != nil {
		return nil, fmt.Errorf("Erreur lors du parsing de la réponse: %v", err)
	}

	return participantsResponse.Participants, nil
}
