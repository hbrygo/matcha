package chatRoom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func GetMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET MESSAGE\n")
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "Vous n'êtes pas connecté", http.StatusUnauthorized)
		return
	}
	userID := uidCookie.Value

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Erreur lors du parsing du JSON", http.StatusBadRequest)
		return
	}

	// convert chatRoom en entier
	chatRoom, ok := requestBody["chatRoom"].(string)
	if !ok {
		http.Error(w, "Invalid chatRoom", http.StatusBadRequest)
		return
	}
	requestBody["chatRoom"], err = strconv.Atoi(chatRoom)
	if err != nil {
		http.Error(w, "Invalid chatRoom", http.StatusBadRequest)
		return
	}

	// Convertir userID en entier
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	requestBody["userID"] = userIDInt

	modifiedBody, err := json.Marshal(requestBody)
	if err != nil {
		http.Error(w, "Erreur lors de la création du JSON", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://localhost:8181/get_message", "application/json", bytes.NewBuffer(modifiedBody))
	if err != nil {
		http.Error(w, "Erreur lors de la requête", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture de la réponse", http.StatusInternalServerError)
		return
	}

	fmt.Printf("GET MESSAGE RESPONSE: %s\n", response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(response)
}
