package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	db *sql.DB
}

func NewApi(file string) (*Api, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.WithField("dbfile", file).Error("Failed to open database.")
		return nil, err
	}

	return &Api{db: db}, nil
}

func (a *Api) Handle(r *http.Request) ([]byte, error) {
	l := log.WithField("method", r.Method).WithField("url", r.URL)
	comps := strings.Split(r.URL.Path, "/")
	l.Debug(comps[2])
	switch r.Method {
	case http.MethodGet:
		switch comps[2] {
		case "seasons":
			return a.GetSeasons()
		default:
			return nil, fmt.Errorf("not found")
		}
	case http.MethodPost:
		switch comps[2] {
		case "seasons":
			name := r.FormValue("season")
			return []byte("{}"), a.AddSeason(name)
		default:
			return nil, fmt.Errorf("not found")
		}
	}
	return nil, fmt.Errorf("not found")
}
