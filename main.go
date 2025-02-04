package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func sendIndex(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("template/index.html")
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	http.ServeContent(w, r, file.Name(), fileInfo.ModTime(), file)
}

func sendRegister(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("template/register.html")
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	http.ServeContent(w, r, file.Name(), fileInfo.ModTime(), file)
}

func sendLogin(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("template/login.html")
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Fichier introuvable", 404)
		return
	}
	http.ServeContent(w, r, file.Name(), fileInfo.ModTime(), file)
}

func checkRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Username: %s\n", req.Username)
	fmt.Printf("Email: %s\n", req.Email)
	fmt.Printf("Password: %s\n", req.Password)
	// Process registration logic here...

	// if (/*all good*/) {
	response := map[string]interface{}{
		"success": true,
		"message": "Registration successful",
		"user":    req, // Example user data
	}
	// } else {
	// 	errorMessage := selectErrorMessage(req)
	// 	response := map[string]interface{}{
	// 		"success": false,
	// 		"message": errorMessage,
	// 	}
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Email: %s\n", req.Email)
	fmt.Printf("Password: %s\n", req.Password)
	// Check dans la db si le login est correct

	// if (/*all good*/) {
	response := map[string]interface{}{
		"success": true,
		"message": "Login successful",
	}
	// } else {
	// 	let errorMessage = selectErrorMessage(req)
	// 	response := map[string]interface{}{
	// 		"success": false,
	// 		"message": errorMessage,
	// 	}
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", sendIndex)
	http.HandleFunc("/register", sendRegister)
	http.HandleFunc("/login", sendLogin)
	http.HandleFunc("POST /checkRegister", checkRegister)
	http.HandleFunc("POST /checkLogin", checkLogin)
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur")
		return
	}
}
