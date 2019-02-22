package dashboard

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id   int64
	Name string
}

type SeasonPage struct {
	Errors []string
	Tabs   map[string]*Tab
	Season *Season
}

type Tab struct {
	Name   string
	Active bool
}

func NewSeasonPage(errs []string, season *Season, active string) (*SeasonPage, error) {
	tabs := map[string]*Tab{
		"cars":        &Tab{"Cars", false},
		"heats":       &Tab{"Heats", false},
		"leaderboard": &Tab{"Leaderboard", false},
	}

	if _, ok := tabs[active]; !ok {
		return nil, fmt.Errorf("unknown tab id")
	}

	tabs[active].Active = true

	sp := &SeasonPage{
		Errors: errs,
		Tabs:   tabs,
		Season: season,
	}

	return sp, nil
}

func AddSeason(name string) (int64, error) {
	stmt, err := db.Prepare("insert into seasons (name) values (:name)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return -1, err
	}

	res, err := stmt.Exec(sql.Named("name", name))
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return -1, err
	}

	li, err := res.LastInsertId()
	if err != nil {
		log.WithError(err).Error("Failed to get last insert id.")
		return -1, err
	}

	return li, err
}

func FetchSeason(id int64) (*Season, error) {
	stmt, err := db.Prepare("select name from seasons where id=:sid")
	if err != nil {
		log.WithError(err).Error("failed to prepare")
		return nil, err
	}

	var name string
	err = stmt.QueryRow(sql.Named("sid", id)).Scan(&name)
	if err != nil {
		log.WithError(err).Error("failed to get season info")
		return nil, err
	}

	return &Season{Id: id, Name: name}, nil
}

func (s *Season) Cars() ([]*Car, error) {
	return FetchCars(s.Id)
}

func (s *Season) Classes() ([]*Class, error) {
	return FetchClasses(s.Id)
}

func (s *Season) Runs() ([]*Run, error) {
	return FetchRuns(s.Id)
}

func (s *Season) Heats() ([]*Heat, error) {
	return FetchHeats(s.Id)
}
