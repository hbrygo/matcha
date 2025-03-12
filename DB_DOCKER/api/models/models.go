package models

// Message struct shared between handlers and utils
type Message struct {
	UserID    int    `json:"userID"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

// CreateChatroomRequest représente la requête pour créer une chatroom
type CreateChatroomRequest struct {
	UserIDs []int  `json:"user_ids"`
	Name    string `json:"name"`
}

// CreateChatroomResponse représente la réponse à la création d'une chatroom
type CreateChatroomResponse struct {
	Status     string `json:"status"`
	ChatroomID int    `json:"chatroom_id"`
	Message    string `json:"message"`
}

// AddUserToChatroomRequest représente la requête pour ajouter un utilisateur à une chatroom
type AddUserToChatroomRequest struct {
	ChatroomID int `json:"chatroom_id"`
	UserID     int `json:"user_id"`
}

// AddUserToChatroomResponse représente la réponse à l'ajout d'un utilisateur à une chatroom
type AddUserToChatroomResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RemoveUserFromChatroomRequest représente la requête pour retirer un utilisateur d'une chatroom
type RemoveUserFromChatroomRequest struct {
	ChatroomID int `json:"chatroom_id"`
	UserID     int `json:"user_id"`
}

// RemoveUserFromChatroomResponse représente la réponse au retrait d'un utilisateur d'une chatroom
type RemoveUserFromChatroomResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewMessageRequest représente la requête pour ajouter un message
type NewMessageRequest struct {
	Message  string `json:"message"`
	ChatRoom int    `json:"chatRoom"`
	UserID   int    `json:"userID"`
}

// response struct
type NewMessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// request struct
type GetMessageRequest struct {
	ChatRoom int `json:"chatRoom"`
	UserID   int `json:"userID"`
}

// response struct
type GetMessageResponse struct {
	Status   string    `json:"status"`
	Messages []Message `json:"messages"`
}

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

// Request structure to get user
type GetUserRequest struct {
	UID int `json:"uid"`
}

// Response structure
type GetUserResponse struct {
	User struct {
		Nom        string   `json:"nom"`
		Prenom     string   `json:"prenom"`
		DOB        string   `json:"dob"`
		Gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interests  []string `json:"interests"`
		Pictures   []string `json:"pictures"`
		Bio        string   `json:"bio"`
	} `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User struct {
		UID       int  `json:"uid"`
		FirstStep bool `json:"first_step"`
	} `json:"user"`
}

type MeRequest struct {
	UID int `json:"uid"`
}

type MeResponse struct {
	User struct {
		Username   string   `json:"username"`
		Email      string   `json:"email"`
		Nom        string   `json:"nom"`
		Prenom     string   `json:"prenom"`
		DOB        string   `json:"dob"`
		Gender     string   `json:"gender"`
		Preference string   `json:"preference"`
		Interests  []string `json:"interests"`
		Pictures   []string `json:"pictures"`
		Bio        string   `json:"bio"`
	} `json:"user"`
}

// structure of response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UpdateRequest struct {
	UID        int      `json:"uid"`
	Nom        string   `json:"lastName"`
	Prenom     string   `json:"firstName"`
	DOB        string   `json:"dob"`
	Gender     string   `json:"gender"`
	Preference string   `json:"preference"`
	Interests  []string `json:"interest"`
	Pictures   []string `json:"photos"`
	Bio        string   `json:"bio"`
}

type UpdateResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	User    *UpdateRequest `json:"user,omitempty"`
}

// Structure de requête pour obtenir un utilisateur par son nom d'utilisateur
type GetUserByUsernameRequest struct {
	Username string `json:"username"`
}

// GetMyChatroomRequest représente la requête pour obtenir les chatrooms d'un utilisateur
type GetMyChatroomRequest struct {
	UserID int `json:"userID"`
}

// GetMyChatroomResponse représente la réponse contenant les chatrooms d'un utilisateur
type GetMyChatroomResponse struct {
	User struct {
		ChatRoomIDs []int `json:"chatRoomID"`
	} `json:"user"`
}
