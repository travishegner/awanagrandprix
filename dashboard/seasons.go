package dashboard

import (
	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id   int
	Name string
}

func (dash *Dashboard) GetSeasons() ([]*Season, error) {
	rows, err := dash.db.Query("select id, name from seasons")
	if err != nil {
		log.Error("Failed to execute query.")
		return nil, err
	}

	seasons := []*Season{}
	for rows.Next() {
		s := &Season{}
		err = rows.Scan(&s.Id, &s.Name)
		if err != nil {
			log.Error("Failed to read row.")
			return nil, err
		}
		seasons = append(seasons, s)
	}

	return seasons, nil
}
