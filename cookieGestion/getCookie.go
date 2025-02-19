package cookieGestion

import (
	"encoding/json"
	"net/http"
)

func GetCookie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(r.Cookies())
}
