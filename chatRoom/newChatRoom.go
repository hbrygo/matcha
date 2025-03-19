package chatRoom

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type ChatRoom struct {
	ChatRoomID int `json:"chatRoomID"`
}

func NewChatRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("NewChatRoom\n")

	// Lire le body de la requête entrante
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Body: %s\n", body)

	// Transmettre ce body à la nouvelle requête
	resp, err := http.Post("http://localhost:8181/create_chatroom", "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error making request to create_chatroom", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lire la réponse de la nouvelle requête
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response: %s\n", responseBody)

	// Renvoyer cette réponse telle quelle à ton front
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}
