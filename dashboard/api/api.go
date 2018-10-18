package api

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	dbfile = "agp.db"
)

func HandleApi(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("method", r.Method).WithField("url", r.URL)
	comps := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case http.MethodGet:
		l.Debug(comps[2])
		switch comps[2] {
		case "seasons":
			s, err := GetSeasons()
			if err != nil {
				http.Error(w, "Error getting seasons", 500)
				l.Error("Error getting seasons")
				return
			}
			b, err := w.Write(s)
			if err != nil {
				http.Error(w, "Error writing output", 500)
				l.WithError(err).Error("Error writing output")
				return
			}
			l.WithField("bytes", b).Debug("Wrote data")
			return
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
