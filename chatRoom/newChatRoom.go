package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ChatRoom struct {
	ChatRoomID int `json:"chatRoomID"`
}

func NewChatRoom(w http.ResponseWriter, r *http.Request) {
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "You are not connected", 401)
		return
	}
	uid, _ := strconv.Atoi(UID.Value)
	postBody, err := json.Marshal(map[string]int{
		"uid": uid,
	})
	if err != nil {
		http.Error(w, "Error: json", 400)
		return
	}

	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8181/create_chatroom", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	var chatRoom ChatRoom
	err = json.NewDecoder(resp.Body).Decode(&chatRoom)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	fmt.Printf("ChatRoomID: %v\n", chatRoom.ChatRoomID)
	http.Redirect(w, r, "/chatRoom.html", http.StatusSeeOther)
}
