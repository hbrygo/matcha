package cookieGestion

import "net/http"

func SetCookie(w http.ResponseWriter, value string, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true, // Empêche l'accès via JS
		Secure:   true, // Seulement HTTPS
		SameSite: http.SameSiteStrictMode,
	})
}
