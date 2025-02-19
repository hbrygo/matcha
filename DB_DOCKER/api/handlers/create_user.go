package handlers

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"matcha/database"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest defines the expected request structure
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserResponse defines the API response structure
type CreateUserResponse struct {
	User struct {
		UID int `json:"uid"`
	} `json:"user"`
}

// validatePassword checks for password strength requirements
func validatePassword(password string) bool {
	// Minimum 8 characters
	if len(password) < 8 {
		fmt.Printf("Password too short")
		return false
	}

	// Check for at least one number
	num := regexp.MustCompile(`[0-9]`)
	if !num.MatchString(password) {
		fmt.Printf("Password no number")
		return false
	}

	// Check for at least one uppercase letter
	upper := regexp.MustCompile(`[A-Z]`)
	if !upper.MatchString(password) {
		fmt.Printf("Password no upper")
		return false
	}

	// Check for at least one special character
	special := regexp.MustCompile(`[!@#$%^&*]`)
	if !special.MatchString(password) {
		fmt.Printf("Password no special")
		return false
	}

	return true
}

// validateEmail checks if email format is valid
func validateEmail(email string) bool {
	/*
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		return emailRegex.MatchString(email)*/
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// CreateUserHandler handles user creation requests
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Check HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !validateEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if !validatePassword(req.Password) {
		http.Error(w, "Password does not meet security requirements", http.StatusBadRequest)
		return
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Check for existing email and username
	var emailExists, usernameExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", req.Email).Scan(&emailExists)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", req.Username).Scan(&usernameExists)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Handle existing email/username cases
	if emailExists && usernameExists {
		http.Error(w, "Email and username already exist", 411)
		return
	} else if emailExists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	} else if usernameExists {
		http.Error(w, "Username already exists", 410)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Insert new user
	result, err := db.Exec(`
        INSERT INTO users (username, email, password, created_at) 
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		req.Username, req.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Get the newly created user ID
	uid, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Prepare and send response
	response := CreateUserResponse{}
	response.User.UID = int(uid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
