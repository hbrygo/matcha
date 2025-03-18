package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserResponse struct {
	User struct {
		LastName   string   `json:"nom"`
		FirstName  string   `json:"prenom"`
		dob        string   `json:"dob"`
		gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interest   []string `json:"interest"`
		Photos     []string `json:"photos"`
		Bio        string   `json:"bio"`
	} `json:"user"`
}

func GetUserByName(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetUserByName\n")
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
	resp, err := http.Post("http://localhost:8181/get_user_by_name", "application/json", responseBody)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 400)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error: "+resp.Status, resp.StatusCode)
		return
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
