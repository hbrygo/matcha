package connection

import (
	"matcha/cookieGestion"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookieGestion.DeleteCookie(w, "uid")
	cookieGestion.DeleteCookie(w, "firstStep")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
