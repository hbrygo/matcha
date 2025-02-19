package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"matcha/cookieGestion"
	"net/http"
	"strconv"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var userResp struct {
		User struct {
			UID int `json:"uid"`
		} `json:"user"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Username: %s\n", req.Username)
	fmt.Printf("Email: %s\n", req.Email)
	fmt.Printf("Password: %s\n", req.Password)
	// Process registration logic here...

	postBody, _ := json.Marshal(map[string]string{
		"username": req.Username,
		"email":    req.Email,
		"password": req.Password,
	})
	responseBody := bytes.NewBuffer(postBody)
	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8181/create_user", "application/json", responseBody)
	// Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	// Read the response body

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// fmt.Printf("Response: %s\n", responseData)
	if err := json.Unmarshal(responseData, &userResp); err != nil {
		log.Printf("Error parsing response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Maintenant vous pouvez accéder à l'UID
	uid := userResp.User.UID
	fmt.Printf("User UID: %d\n", uid)

	var response map[string]interface{}
	if resp.StatusCode == 200 {
		// Décommenter et utiliser l'UID récupéré
		cookieGestion.SetCookie(w, strconv.Itoa(uid), "uid")
		response = map[string]interface{}{
			"success": true,
			"message": "Registration successful",
			"uid":     uid,
		}
	} else {
		// ... reste du code ...
		response = map[string]interface{}{
			"success": false,
			"message": "Registration failed",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
