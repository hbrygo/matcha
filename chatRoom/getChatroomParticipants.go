package chatRoom

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type GetChatroomParticipantsResponse struct {
	Participants []int `json:"participants"`
}

func GetChatroomParticipants(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getChatroomParticipants\n")
	// get my cookie
	_, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "You are not connected", 401)
		return
	}

	// get chatroom id dans le body du POST
	chatroomID, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	responseBody := bytes.NewBuffer(chatroomID)
	fmt.Printf("chatroomID: %v\n", responseBody)

	// Transmettre ce body à la nouvelle requête
	resp, err := http.Post("http://localhost:8181/get_chatroom_participants", "application/json", responseBody)
	if err != nil {
		http.Error(w, "Error making request to get_chatroom_participants", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// fmt.Printf("Response from get_chatroom_participants: %v\n", resp)

	// Lire la réponse de la nouvelle requête
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response from get_chatroom_participants: %s\n", response)

	// Renvoyer cette réponse telle quelle à ton front
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(response)
}
