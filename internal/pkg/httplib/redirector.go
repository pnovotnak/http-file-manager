package httplib

import "net/http"

func Redirector(url string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, code)
	}
}
