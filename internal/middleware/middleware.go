package middleware

import (
	"net/http"

	"github.com/stephenafamo/isbot"
)


func CheckBot (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isBot := isbot.Check(r.UserAgent())
		if isBot {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			next.ServeHTTP(w,r)
		}
	})
}