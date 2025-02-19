package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"matcha/cookieGestion"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Check login start\n")
	fmt.Printf("____________________________________________________\n")
	var req struct {
		Email    string `json:"email,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//parse le json et prepare postBody pour la requête si il y'a un email
	var postBody []byte
	var err error
	if req.Username != "" {
		postBody, err = json.Marshal(map[string]string{
			"username": req.Username,
			"password": req.Password,
		})
	} else {
		postBody, err = json.Marshal(map[string]string{
			"email":    req.Email,
			"password": req.Password,
		})
	}
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// postBody, err := json.Marshal(map[string]string{
	// 	"username": req.Username,
	// 	"password": req.Password,
	// })
	// if err != nil {
	// 	log.Fatalf("An Error Occured %v", err)
	// }

	responseBody := bytes.NewBuffer(postBody)
	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8181/login", "application/json", responseBody)
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

	fmt.Printf("Response: %s\n", responseData)

	// var userResp UserResponse
	var userResp struct {
		User struct {
			UID       int  `json:"uid"`
			FirstStep bool `json:"first_step"`
		} `json:"user"`
	}
	if err := json.Unmarshal(responseData, &userResp); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	uid := userResp.User.UID
	fmt.Printf("Extracted UID: %d\n", uid)
	firstStep := userResp.User.FirstStep
	fmt.Printf("First step: %t\n", firstStep)
	if firstStep {
		fmt.Printf("Setting firstStep cookie\n")
		cookieGestion.SetCookie(w, "1", "firstStep")
	}

	// Créer un cookie sécurisé
	cookieGestion.SetCookie(w, fmt.Sprintf("%d", uid), "uid")

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(responseData, &apiResponse); err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	var response map[string]interface{}
	if resp.StatusCode == 200 {
		response = map[string]interface{}{
			"success": true,
			"message": "Login successful",
		}
	} else {
		response = map[string]interface{}{
			"success": false,
			"message": apiResponse["message"],
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("____________________________________________________\n")
	fmt.Printf("Check login end\n")
}
