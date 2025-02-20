package cookieGestion

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCookie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("Cookies: %v\n", r.Cookies())
	json.NewEncoder(w).Encode(r.Cookies())
}

func GetCookieForBackend(r *http.Request) []*http.Cookie {
	return r.Cookies()
}
