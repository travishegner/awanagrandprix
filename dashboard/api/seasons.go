package api

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Season struct {
	Name string `json:name`
}

func GetSeasons() ([]byte, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.WithField("dbfile", dbfile).Error("Failed to open database.")
		return nil, err
	}

	rows, err := db.Query("select name from seasons")
	if err != nil {
		log.Error("Failed to execute query.")
		return nil, err
	}

	var name string
	Seasons := []*Season{}
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			log.Error("Failed to read row.")
			return nil, err
		}
		Seasons = append(Seasons, &Season{Name: name})
	}

	js, err := json.Marshal(Seasons)
	if err != nil {
		log.Error("Failed to marshal json.")
		return nil, err
	}

	return js, nil
}
