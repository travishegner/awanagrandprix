package api

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (a *Api) GetSeasons() ([]byte, error) {
	rows, err := a.db.Query("select id, name from seasons")
	if err != nil {
		log.Error("Failed to execute query.")
		return nil, err
	}

	var name string
	var id int
	Seasons := []*Season{}
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Error("Failed to read row.")
			return nil, err
		}
		Seasons = append(Seasons, &Season{Id: id, Name: name})
	}

	js, err := json.Marshal(Seasons)
	if err != nil {
		log.Error("Failed to marshal json.")
		return nil, err
	}

	return js, nil
}

func (a *Api) AddSeason(name string) error {
	log.WithField("season", name).Debug("Adding season")

	stmt, err := a.db.Prepare("insert into seasons (name) values (?)")
	if err != nil {
		log.WithField("season", name).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(name)
	if err != nil {
		log.WithField("season", name).Error("Failed to exec statement.")
		return err
	}

	return nil
}

func (a *Api) DeleteSeason(s *Season) error {
	log.WithField("season", s).Debug("Deleting season")

	stmt, err := a.db.Prepare("delete from seasons where id=?")
	if err != nil {
		log.WithField("season", s).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(s.Id)
	if err != nil {
		log.WithField("season", s).Error("Failed to exec statement.")
		return err
	}

	return nil
}

func (a *Api) RenameSeason(s *Season) error {
	log.WithField("season", s).Debug("Renaming season")

	stmt, err := a.db.Prepare("update seasons set name=? where id=?")
	if err != nil {
		log.WithField("season", s).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(s.Name, s.Id)
	if err != nil {
		log.WithField("season", s).Error("Failed to exec statement.")
		return err
	}

	return nil
}
