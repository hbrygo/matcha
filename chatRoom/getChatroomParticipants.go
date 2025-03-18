package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type GetChatroomParticipantsResponse struct {
	Participants []int `json:"participants"`
}

func GetChatroomParticipants(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetChatroomParticipants\n")
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "You are not connected", 401)
		return
	}
	uid, _ := strconv.Atoi(UID.Value)
	fmt.Printf("Cookie: %v\n", uid)
	chatroomID, err := strconv.Atoi(r.URL.Query().Get("chatroomID"))
	if err != nil {
		http.Error(w, "Error: chatroomID", 400)
		return
	}
	postBody, err := json.Marshal(map[string]int{
		"chatroomID": chatroomID,
	})
	if err != nil {
		http.Error(w, "Error: json", 400)
		return
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8181/get_chatroom_participants", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error: "+resp.Status, resp.StatusCode)
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Message{})
		return
	}

	var message GetChatroomParticipantsResponse
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message.Participants)
}
