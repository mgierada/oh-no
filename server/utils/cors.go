package utils

import (
	"net/http"
	"os"
)

var uiRootUrl = os.Getenv("UI_ROOT_URL")

var allowedOrigins = []string{
	"http://localhost:3000",
	uiRootUrl,
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

func EnableCors(w *http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if isOriginAllowed(origin) {
		(*w).Header().Set("Access-Control-Allow-Origin", origin)
	}
}
