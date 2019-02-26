package dashboard

import (
	"database/sql"
	"sort"

	log "github.com/sirupsen/logrus"
)

type Class struct {
	Id       int64
	SeasonId int64
	Name     string
}

type Result struct {
	Place       int64
	Car         *Car
	AverageTime float64
}

func AddClass(seasonId int64, name string) error {
	stmt, err := db.Prepare("insert into classes (season_id,name) values (:sid,:name)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(sql.Named("sid", seasonId), sql.Named("name", name))
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return err
	}

	return nil
}

func FetchClasses(seasonId int64) ([]*Class, error) {
	stmt, err := db.Prepare("select id, name from classes where season_id=:sid")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return nil, err
	}

	rows, err := stmt.Query(sql.Named("sid", seasonId))
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return nil, err
	}

	classes := []*Class{}
	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		classes = append(classes, &Class{Id: id, SeasonId: seasonId, Name: name})
	}

	return classes, nil
}

func (c *Class) Results() ([]*Result, error) {
	cars, err := FetchCars(c.SeasonId)
	if err != nil {
		log.WithError(err).Error("failed to fetch cars")
		return nil, err
	}

	set := make([]*Result, 0)
	for _, car := range cars {
		if car.Class.Id != c.Id {
			continue
		}
		//f, err := car.FetchBestThreeResult()
		f, err := car.FetchResult()
		if err != nil {
			log.WithError(err).Error("failed to get result")
			return nil, err
		}
		set = append(set, &Result{Car: car, AverageTime: f})
	}

	sort.Slice(set, func(i, j int) bool {
		return set[i].AverageTime < set[j].AverageTime
	})

	for i, r := range set {
		r.Place = int64(i + 1)
	}

	return set, nil
}
