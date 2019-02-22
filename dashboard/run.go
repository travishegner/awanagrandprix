package dashboard

import (
	"database/sql"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

type Run struct {
	Id         int64
	CarId      int64
	LaneId     int64
	HeatNumber sql.NullInt64
	Time       sql.NullFloat64
}

func SetRunTime(r int64, t float64) error {
	stmt, err := db.Prepare("update runs set time=:t where id=:rid")
	if err != nil {
		log.WithError(err).Error("failed to prepare statement")
		return err
	}

	_, err = stmt.Exec(sql.Named("t", t), sql.Named("rid", r))
	if err != nil {
		log.WithError(err).Error("failed to update time")
		return err
	}

	return nil
}

func (r *Run) SetTime(t float64) error {
	return SetRunTime(r.Id, t)
}

func FetchRuns(seasonId int64) ([]*Run, error) {
	stmt, err := db.Prepare(`
select r.id, r.car_id, r.lane_id, r.heat, r.time
from runs r
inner join cars c on r.car_id=c.id
inner join classes cls on c.class_id=cls.id
where cls.season_id=:sid
`)
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return nil, err
	}

	rows, err := stmt.Query(sql.Named("sid", seasonId))
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return nil, err
	}

	runs := []*Run{}
	for rows.Next() {
		var id int64
		var carid int64
		var laneid int64
		var heat sql.NullInt64
		var time sql.NullFloat64
		err = rows.Scan(&id, &carid, &laneid, &heat, &time)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		runs = append(runs, &Run{Id: id, CarId: carid, LaneId: laneid, HeatNumber: heat, Time: time})
	}

	return runs, nil
}

func GenerateRuns(s *Season) error {
	stmt, err := db.Prepare("insert into runs (car_id,lane_id) values (:cid,:lid)")
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
			_, err = stmt.Exec(sql.Named("cid", c.Id), sql.Named("lid", l.Id))
			if err != nil {
				log.WithError(err).Error("failed insert run")
				continue
			}
		}
	}

	return nil
}
