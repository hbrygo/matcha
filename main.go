package main

import (
	"fmt"
	"matcha/connection"
	"matcha/cookieGestion"
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

func main() {
	http.HandleFunc("/", sendPage)
	http.HandleFunc("/logout", connection.Logout)
	http.HandleFunc("GET /getCookie", cookieGestion.GetCookie)
	http.HandleFunc("POST /register", connection.Register)
	http.HandleFunc("POST /login", connection.Login)
	http.HandleFunc("POST /setData", connection.SetData)
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
