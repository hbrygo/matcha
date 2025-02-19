package handlers

import (
	"encoding/json"
	"net/http"
)

// structure of response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// handler for root
func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := Response{
		Status:  "success",
		Message: "API Matcha en ligne. Utilisez /testDBavailability pour tester la base de données.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
