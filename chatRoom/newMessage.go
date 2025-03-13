package chatRoom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Message struct {
	UID    int    `json:"uid"`
	Msg    string `json:"msg"`
	RoomID int    `json:"room_id"`
}

func NewMessage(w http.ResponseWriter, r *http.Request) {
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "User not connected", 401)
		return
	}
	uid, _ := strconv.Atoi(UID.Value)

	var msg Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Error: decode", 400)
		return
	}

	if msg.UID != uid {
		http.Error(w, "U are not this person, stop this pls", 401)
		return
	}

	if msg.RoomID == 0 {
		http.Error(w, "No room specified", 400)
		return
	}

	resp, err := http.Post("http://localhost:8181/new_message", "application/json", r.Body)
	if err != nil {
		http.Error(w, "Error: message not send", 500)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "Message send")
}
