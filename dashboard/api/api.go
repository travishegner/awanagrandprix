package api

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("method", r.Method).WithField("url", r.URL)
	comps := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case http.MethodGet:
		switch comps[1] {
		case "seasons":

		default:
			http.NotFound(w, r)
			l.Error("Not Found")
			return
		}
	case http.MethodPost:
	default:
		http.Error(w, "Unsupported Method", 405)
		l.Error("Unsupported Method")
		return
	}
}
