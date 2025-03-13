package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"matcha/chatRoom"
	"matcha/connection"
	"matcha/cookieGestion"
	"net/http"
	"os"
	"strconv"
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

func getMe(w http.ResponseWriter, r *http.Request) {
	UID, err := r.Cookie("uid")
	if err != nil {
		http.Error(w, "Vous n'êtes pas connecté", 401)
		return
	}
	// fmt.Printf("Cookie: %v\n", UID.Value)
	uid, _ := strconv.Atoi(UID.Value)
	postBody, err := json.Marshal(map[string]int{
		"uid": uid,
	})

	// postBody, err := json.Marshal(map[string]string{
	// 	"uid": UID.Value,
	// })
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	// fmt.Printf("postBody (string): %s\n", string(postBody))
	responseBody := bytes.NewBuffer(postBody)
	// fmt.Printf("postBody: %v\n", responseBody)
	resp, err := http.Post("http://localhost:8181/me", "application/json", responseBody)
	if err != nil {
		// Gérer l'erreur de la requête
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("Response: %s\n", body)
	// Définir le type de contenu de la réponse
	w.Header().Set("Content-Type", "application/json")

	// Écrire le status code
	w.WriteHeader(resp.StatusCode)

	// Écrire le corps de la réponse au client
	w.Write(body)
}

func main() {
	http.HandleFunc("/", sendPage)
	http.HandleFunc("/logout", connection.Logout)
	http.HandleFunc("GET /getCookie", cookieGestion.GetCookie)
	http.HandleFunc("POST /me", getMe)
	http.HandleFunc("POST /register", connection.Register)
	http.HandleFunc("POST /login", connection.Login)
	http.HandleFunc("POST /setData", connection.SetData)
	http.HandleFunc("POST /getMyChatRoom", chatRoom.GetMyChatRoom)
	http.HandleFunc("POST /getChatRoom", chatRoom.GetMessage)
	http.HandleFunc("POST /sendMessage", chatRoom.NewMessage)
	http.HandleFunc("POST /createChatRoom", chatRoom.NewChatRoom)
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
