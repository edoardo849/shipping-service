package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func withBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}
		//dXNlcm5hbWU6cGFzc3dvcmQ=
		if pair[0] != "username" || pair[1] != "password" {
			respondWithError(w, http.StatusUnauthorized, "Not authorized")
			return
		}

		h(w, r)
	}
}

func handle404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
