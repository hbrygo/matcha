package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"matcha/cookieGestion"
	"net/http"
	"strconv"
)

func SetData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Set data start\n")
	fmt.Printf("____________________________________________________\n")
	var req struct {
		Uid        string   `json:"uid"`
		FirstName  string   `json:"firstName"`
		LastName   string   `json:"lastName"`
		Dob        string   `json:"dob"`
		Gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interest   []string `json:"interest"`
		Photos     []string `json:"photos"` // Will contain base64 strings
		Bio        string   `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("An Error Occured: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uid, err := strconv.Atoi(req.Uid)
	if err != nil {
		fmt.Printf("Invalid UID format: %v\n", err)
		http.Error(w, "Invalid UID format", http.StatusBadRequest)
		return
	}

	dbReq := struct {
		Uid        int      `json:"uid"`
		FirstName  string   `json:"firstName"`
		LastName   string   `json:"lastName"`
		Dob        string   `json:"dob"`
		Gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interest   []string `json:"interest"`
		Photos     []string `json:"photos"`
		Bio        string   `json:"bio"`
	}{
		Uid:        uid,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Dob:        req.Dob,
		Gender:     req.Gender,
		Preference: req.Preference,
		Interest:   req.Interest,
		Photos:     req.Photos,
		Bio:        req.Bio,
	}

	// fmt.Printf("req: %v\n", dbReq)
	postBody, _ := json.Marshal(dbReq)
	responseBody := bytes.NewBuffer(postBody)

	// send data to database
	resp, err := http.Post("http://localhost:8181/update", "application/json", responseBody)
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Read the response body
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	fmt.Printf("Response: %s\n", responseData)
	var response map[string]interface{}
	if resp.StatusCode == 200 {
		cookieGestion.SetCookie(w, "1", "firstStep")
		response = map[string]interface{}{
			"success": true,
			"message": "Registration successful",
		}
	} else {
		response = map[string]interface{}{
			"success": false,
			"message": "Registration failed",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("____________________________________________________\n")
	fmt.Printf("Set data end\n")
}
