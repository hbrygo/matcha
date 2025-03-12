package router

import (
	"log"
	"matcha/api/handlers"
	"matcha/middleware"
	"net/http"
)

/**

here set the headers rules and the handlers for each route

**/

func Router(mux *http.ServeMux) http.Handler {
	return middleware.GeneralMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Router log %s %s %s", r.RemoteAddr, r.Method, r.URL)

		switch r.URL.Path {
		case "/":
			handlers.RootHandler(w, r)
		case "/testDBavailability":
			handlers.TestDBHandler(w, r)
		case "/update":
			handlers.UpdateHandler(w, r)
		case "/me":
			handlers.MeHandler(w, r)
		case "/login":
			handlers.LoginHandler(w, r)
		case "/get_user":
			handlers.GetUserHandler(w, r)
		case "/create_user":
			handlers.CreateUserHandler(w, r)
		// in the Router function, add this case:
		case "/create_chatroom":
			handlers.CreateChatroomHandler(w, r)
			// Dans router.go, ajoutez ces cases dans la fonction switch
		case "/add_user_to_chatroom":
			handlers.AddUserToChatroomHandler(w, r)
		case "/remove_user_from_chatroom":
			handlers.RemoveUserFromChatroomHandler(w, r)
			// Dans router.go, ajoutez cette case
		case "/new_message":
			handlers.NewMessageHandler(w, r)
		case "/get_message":
			handlers.GetMessageHandler(w, r)
		case "/get_user_by_username":
			handlers.GetUserByUsernameHandler(w, r)
		case "/get_my_chatroom":
			handlers.GetMyChatroomHandler(w, r)
		case "/get_all_chatroom":
			handlers.GetAllChatroomHandler(w, r)
		case "/change_password":
			handlers.ChangePasswordHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
}
