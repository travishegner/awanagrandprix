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
	Car        *Car
	Lane       *Lane
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
select r.id, c.id, c.number, c.name, c.weight, c.driver, cls.id, cls.name, l.id, l.color, r.heat, r.time
from runs r
inner join cars c on r.car_id=c.id
inner join classes cls on c.class_id=cls.id
inner join lanes l on r.lane_id=l.id
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
		var carnumber string
		var carname sql.NullString
		var carweight float64
		var driver string
		var classid int64
		var classname string
		var laneid int64
		var color string
		var heat sql.NullInt64
		var time sql.NullFloat64
		err = rows.Scan(&id, &carid, &carnumber, &carname, &carweight, &driver, &classid, &classname, &laneid, &color, &heat, &time)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		r := &Run{
			Id: id,
			Car: &Car{
				Id:     carid,
				Number: carnumber,
				Name:   carname.String,
				Weight: carweight,
				Driver: driver,
				Class: &Class{
					Id:   classid,
					Name: classname,
				},
			},
			Lane: &Lane{
				Id:    laneid,
				Color: color,
			},
			HeatNumber: heat,
			Time:       time,
		}
		runs = append(runs, r)
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
