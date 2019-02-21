package dashboard

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

type Run struct {
	Id     int64
	CarId  int64
	LaneId int64
}

func (r *Run) FetchHeat() (int64, error) {
	stmt, err := db.Prepare("select heat from runs where id=?")
	if err != nil {
		log.WithError(err).Error("failed to prepare statement")
		return -1, err
	}

	var heat sql.NullInt64
	err = stmt.QueryRow(r.Id).Scan(&heat)
	if err != nil {
		log.WithError(err).Error("failed to execute query")
		return -1, err
	}

	if !heat.Valid {
		return -1, fmt.Errorf("heat undefined for this run")
	}

	return heat.Int64, nil
}

func (r *Run) FetchTime() (float64, error) {
	stmt, err := db.Prepare("select time from runs where id=?")
	if err != nil {
		log.WithError(err).Error("failed to prepare statement")
		return -1, err
	}

	var time sql.NullFloat64
	err = stmt.QueryRow(r.Id).Scan(&time)
	if err != nil {
		log.WithError(err).Error("failed to execute query")
		return -1, err
	}

	if !time.Valid {
		return -1, fmt.Errorf("heat undefined for this run")
	}

	return time.Float64, nil
}

func FetchRuns(seasonId int64) ([]*Run, error) {
	stmt, err := db.Prepare(`
select r.id, r.car_id, r.lane_id
from runs r
inner join cars c on r.car_id=c.id
inner join classes cls on c.class_id=cls.id
where cls.season_id=?
`)
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return nil, err
	}

	rows, err := stmt.Query(seasonId)
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return nil, err
	}

	runs := []*Run{}
	for rows.Next() {
		var id int64
		var carid int64
		var laneid int64
		err = rows.Scan(&id, &carid, &laneid)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		runs = append(runs, &Run{Id: id, CarId: carid, LaneId: laneid})
	}

	return runs, nil
}

func GenerateHeats(s *Season) error {
	err := GenerateRuns(s)
	if err != nil {
		log.WithError(err).Error("failed to generate runs")
		return err
	}

	//TODO: assign runs to heats
	//1. each heat should only have one class
	//2. heats should alternate classes

	return nil
}

func GenerateRuns(s *Season) error {
	stmt, err := db.Prepare("insert into runs (car_id,lane_id) values (?,?)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	cars, err := s.Cars()
	if err != nil {
		log.WithError(err).Error("failed to fetch cars")
		return err
	}
	lanes, err := FetchLanes()
	if err != nil {
		log.WithError(err).Error("failed to fetch lanes")
		return err
	}

	for _, l := range lanes {
		for _, index := range rand.Perm(len(cars)) {
			c := cars[index]
			_, err = stmt.Exec(c.Id, l.Id)
			if err != nil {
				log.WithError(err).Error("failed insert run")
				continue
			}
		}
	}

	return nil
}
