package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type GetMyChatroomResponse struct {
	User struct {
		ChatRoomIDs []int `json:"chatRoomID"`
	} `json:"user"`
}

func GetMyChatRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetMyChatRoom\n")
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "You are not connected", 401)
		return
	}
	uid, _ := strconv.Atoi(UID.Value)
	fmt.Printf("Cookie: %v\n", uid)
	postBody, err := json.Marshal(map[string]int{
		"userID": uid,
	})
	if err != nil {
		http.Error(w, "Error: json", 400)
		return
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8181/get_my_chatroom", "application/json", responseBody)
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

	var message GetMyChatroomResponse
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func GetAllChatRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetMyChatRoom\n")
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "You are not connected", 401)
		return
	}
	uid, _ := strconv.Atoi(UID.Value)
	fmt.Printf("Cookie: %v\n", uid)
	postBody, err := json.Marshal(map[string]int{
		"userID": uid,
	})
	if err != nil {
		http.Error(w, "Error: json", 400)
		return
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8181/get_all_chatroom", "application/json", responseBody)
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

	var message GetMyChatroomResponse
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}
