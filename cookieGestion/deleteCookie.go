package cookieGestion

import "net/http"

func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:  name,
		Value: "",
	})
}
