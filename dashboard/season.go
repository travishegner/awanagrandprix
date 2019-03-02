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
	Errors   []string
	Tabs     map[string]*Tab
	HeatEdit int64
	Season   *Season
}

type Tab struct {
	Name   string
	Active bool
}

type ResultSet struct {
	Class   *Class
	Results []*Result
}

func NewSeasonPage(errs []string, season *Season, active string) (*SeasonPage, error) {
	tabs := map[string]*Tab{
		"cars":        &Tab{"Cars", false},
		"heats":       &Tab{"Heats", false},
		"leaderboard": &Tab{"Leaderboard", false},
		"results":     &Tab{"Results", false},
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

func (s *Season) CurrentHeat() (*Heat, error) {
	heats, err := s.Heats()
	if err != nil {
		log.WithError(err).Error("failed to get heats")
		return nil, err
	}

	for _, h := range heats {
		if h.Complete() {
			continue
		}
		return h, nil
	}

	return nil, nil
}

func (s *Season) PreviousHeat() (*Heat, error) {
	heats, err := s.Heats()
	if err != nil {
		log.WithError(err).Error("failed to get heats")
		return nil, err
	}

	for i, h := range heats {
		if h.Complete() {
			continue
		}
		if i == 0 {
			return nil, nil
		}
		return heats[i-1], nil
	}

	if len(heats) == 0 {
		return nil, nil
	}

	return heats[len(heats)-1], nil
}

func (s *Season) NextHeat() (*Heat, error) {
	heats, err := s.Heats()
	if err != nil {
		log.WithError(err).Error("failed to get heats")
		return nil, err
	}

	for i, h := range heats {
		if h.Complete() {
			continue
		}
		if i == len(heats)-1 {
			return nil, nil
		}
		return heats[i+1], nil
	}

	return nil, nil
}

func (s *Season) ResultSets() ([]*ResultSet, error) {
	classes, err := s.Classes()
	if err != nil {
		log.WithError(err).Error("failed to get classes")
		return nil, err
	}

	sets := make([]*ResultSet, len(classes))
	for i, c := range classes {
		rs, err := c.Results()
		if err != nil {
			log.WithField("class", c.Name).WithError(err).Error("failed to get results from class")
			return nil, err
		}
		sets[i] = &ResultSet{Class: c, Results: rs}
	}

	return sets, nil
}
