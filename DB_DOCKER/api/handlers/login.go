package handlers

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"matcha/database"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"Email,omitempty"`
	Username string `json:"Username,omitempty"`
	Password string `json:"Mot de passe"`
}

type LoginResponse struct {
	User struct {
		UID       int  `json:"uid"`
		FirstStep bool `json:"first_step"`
	} `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// verifiy method
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// decode request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	// verification
	if req.Password == "" || (req.Email == "" && req.Username == "") {
		http.Error(w, "Champs obligatoires manquants", http.StatusBadRequest)
		return
	}

	// Connexion to db
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var (
		uid       int
		hashedPwd string
		firstStep bool
	)

	// look by email or username
	if req.Email != "" {
		err = db.QueryRow("SELECT uid, password, first_step FROM users WHERE email = ?", req.Email).
			Scan(&uid, &hashedPwd, &firstStep)
	} else {
		err = db.QueryRow("SELECT uid, password, first_step FROM users WHERE username = ?", req.Username).
			Scan(&uid, &hashedPwd, &firstStep)
	}

	if err == sql.ErrNoRows {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(req.Password)); err != nil {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// prepare resp
	response := LoginResponse{}
	response.User.UID = uid
	response.User.FirstStep = firstStep

	// send respo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
