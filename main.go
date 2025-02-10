package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func sendPage(w http.ResponseWriter, r *http.Request) {
	pageName := "template/" + r.URL.Path[1:] + ".html"

	if r.URL.Path == "/" {
		pageName = "template/index.html"
	}

	file, err := os.Open(pageName)
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

	postBody, _ := json.Marshal(map[string]string{
		"nom":      req.Username, // Le nom de famille de l'utilisateur
		"prenom":   req.Username, // Le prénom de l'utilisateur
		"email":    req.Email,    // L'email unique de l'utilisateur
		"password": req.Password, // Le mot de passe de l'utilisateur
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

	fmt.Printf("Response: %s\n", responseData)

	// if (/*all good*/) {
	response := map[string]interface{}{}

	if resp.StatusCode == 200 {
		response = map[string]interface{}{
			"success": true,
			"message": "Registration successful",
			"user":    req, // Example user data
		}
	} else {
		response = map[string]interface{}{
			"success": false,
			"message": "Registration failed",
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Check login start\n")
	fmt.Printf("____________________________________________________")
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

	postBody, _ := json.Marshal(map[string]string{
		"email":    req.Email,
		"password": req.Password,
	})
	responseBody := bytes.NewBuffer(postBody)
	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:8181/get_user", "application/json", responseBody)
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
	fmt.Printf("____________________________________________________")
	fmt.Printf("Check login end\n")
}

func updateData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Update data start\n")
	var req struct {
		Birth_date          string   `json:"birthdate"`
		Gender              string   `json:"gender"`
		Preferred_gender    []string `json:"preferred_gender"`
		Sexual_orientation  string   `json:"sexual_orientation"`
		Relationship_status string   `json:"relationship_status"`
		Physical_appearance string   `json:"physical_appearance"`
		Size                string   `json:"size"`
		Weight              string   `json:"weight"`
		Hobbies             string   `json:"hobbies"`
		Location            string   `json:"location"`
		Search_distance     string   `json:"search_distance"`
		Age_range           string   `json:"age_range"`
		Bio                 string   `json:"bio"`
		Notifications       string   `json:"notifications"`
		Terms               string   `json:"terms"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	log.Printf("Raw request body: %s", body)

	// Reset the body so it can be read again
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %v\n", req)

	// Process registration logic here...

	// postBody, _ := json.Marshal(map[string]string{
	// 	"nom":      req.Username, // Le nom de famille de l'utilisateur
	// 	"prenom":   req.Username, // Le prénom de l'utilisateur
	// 	"email":    req.Email,    // L'email unique de l'utilisateur
	// 	"password": req.Password, // Le mot de passe de l'utilisateur
	// })
	// responseBody := bytes.NewBuffer(postBody)
	// // Leverage Go's HTTP Post function to make request
	// resp, err := http.Post("http://localhost:8181/create_user", "application/json", responseBody)
	// // Handle Error
	// if err != nil {
	// 	log.Fatalf("An Error Occured %v", err)
	// }
	// defer resp.Body.Close()
	// // Read the response body

	// responseData, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("An Error Occured %v", err)
	// }

	// fmt.Printf("Response: %s\n", responseData)

	// // if (/*all good*/) {
	// response := map[string]interface{}{}

	// if resp.StatusCode == 200 {
	response := map[string]interface{}{
		"success": true,
		"message": "Registration successful",
		"user":    req, // Example user data
	}
	// } else {
	// 	response = map[string]interface{}{
	// 		"success": false,
	// 		"message": "Registration failed",
	// 	}
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", sendPage)
	http.HandleFunc("POST /checkRegister", checkRegister)
	http.HandleFunc("POST /checkLogin", checkLogin)
	http.HandleFunc("POST /updateData", updateData)
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur")
		return
	}
}
