package api

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	dbfile = "agp.db"
)

func HandleApi(r *http.Request) ([]byte, error) {
	l := log.WithField("method", r.Method).WithField("url", r.URL)
	comps := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case http.MethodGet:
		l.Debug(comps[2])
		switch comps[2] {
		case "seasons":
			return GetSeasons()
		default:
			return nil, fmt.Errorf("not found")
		}
	case http.MethodPost:
	}
	return nil, fmt.Errorf("not found")
}
