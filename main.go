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

type UserResponse struct {
	User struct {
		Email  string `json:"email"`
		Nom    string `json:"nom"`
		Prenom string `json:"prenom"`
		UID    int    `json:"uid"`
	} `json:"user"`
}

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
	// fmt.Printf("Response: %v\n", resp)

	// if (/*all good*/) {
	response := map[string]interface{}{}

	if resp.StatusCode == 200 {
		// setCookie(w, "1", "firstStep")
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
	fmt.Printf("____________________________________________________\n")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Printf("Email: %s\n", req)
	// fmt.Printf("Password: %s\n", req.Password)
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

	var userResp UserResponse
	if err := json.Unmarshal(responseData, &userResp); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	uid := userResp.User.UID
	fmt.Printf("Extracted UID: %d\n", uid)

	// Créer un cookie sécurisé
	setCookie(w, fmt.Sprintf("%d", uid), "uid")

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

func setData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Set data start\n")
	fmt.Printf("____________________________________________________\n")
	var req struct {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %v\n", req.Interest)
	// Process base64 images if needed
	// for i, photo := range req.Photos {
	//     // Remove base64 header if needed
	//     // Store images...
	// }

	setCookie(w, "1", "firstStep")
	response := map[string]interface{}{
		"success": true,
		"message": "Registration successful",
		"user":    req,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Printf("____________________________________________________\n")
	fmt.Printf("Set data end\n")
}

func setCookie(w http.ResponseWriter, value string, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true, // Empêche l'accès via JS
		Secure:   true, // Seulement HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r.Cookies())
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: "",
	})
}

func logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w, "uid")
	deleteCookie(w, "firstStep")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", sendPage)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("GET /getCookie", getCookie)
	http.HandleFunc("POST /checkRegister", checkRegister)
	http.HandleFunc("POST /checkLogin", checkLogin)
	http.HandleFunc("POST /setData", setData)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	fmt.Println("Serveur démarré sur : http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur")
		return
	}
}
