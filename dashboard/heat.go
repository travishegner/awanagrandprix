package dashboard

import (
	"database/sql"
	"sort"

	log "github.com/sirupsen/logrus"
)

type Heat struct {
	Number int64
	Red    *Run
	Green  *Run
	Blue   *Run
	Yellow *Run
}

func FetchHeats(seasonId int64) ([]*Heat, error) {
	runs, err := FetchRuns(seasonId)
	if err != nil {
		log.WithError(err).Error("failed to get runs")
		return nil, err
	}

	heatMap := map[int64]*Heat{}
	for i, r := range runs {
		n := r.HeatNumber.Int64
		if !r.HeatNumber.Valid {
			log.WithError(err).Error("failed to get heat number from run")
			return nil, err
		}
		if _, ok := heatMap[n]; !ok {
			heatMap[n] = &Heat{Number: n}
		}

		switch r.LaneId {
		case 1: //red
			heatMap[n].Red = runs[i]
		case 2: //green
			heatMap[n].Green = runs[i]
		case 3: //blue
			heatMap[n].Blue = runs[i]
		case 4: //yellow
			heatMap[n].Yellow = runs[i]
		}
	}

	index := 0
	heats := make([]*Heat, len(heatMap))
	for _, h := range heatMap {
		heats[index] = h
		index++
	}

	sort.Slice(heats, func(i, j int) bool {
		return heats[i].Number < heats[j].Number
	})

	return heats, nil
}

func GenerateHeats(s *Season) error {
	err := GenerateRuns(s)
	if err != nil {
		log.WithError(err).Error("failed to generate runs")
		return err
	}

	classes, err := s.Classes()
	if err != nil {
		log.WithError(err).Error("failed to get classes for season")
		return err
	}

	q := `
with red as (
	select r.id, r.car_id
	from runs r
	inner join cars c on r.car_id=c.id
	where r.lane_id=1
	and r.heat is null
	and c.class_id=:clsid
limit 1
),
blue as (
	select r.id, r.car_id
	from runs r
	inner join cars c on r.car_id=c.id
	where r.lane_id=2
	and r.heat is null
	and c.class_id=:clsid
	and r.car_id not in (select car_id from red)
	limit 1
),
green as (
	select r.id, r.car_id
	from runs r
	inner join cars c on r.car_id=c.id
	where r.lane_id=3
	and r.heat is null
	and c.class_id=:clsid
	and r.car_id not in (select car_id from red union select car_id from blue)
	limit 1
),
yellow as (
	select r.id, r.car_id
	from runs r inner join cars c on r.car_id=c.id
	where r.lane_id=4
	and r.heat is null
	and c.class_id=:clsid
	and r.car_id not in (select car_id from red union select car_id from blue union select car_id from green)
	limit 1
)

update runs set heat=:ht where id in ((select id from red), (select id from blue), (select id from green), (select id from yellow))
`
	stmt, err := db.Prepare(q)
	if err != nil {
		log.WithError(err).Error("failed to prepare statement")
		return err
	}

	dones := map[int]bool{}
	for i, _ := range classes {
		dones[i] = false
	}
	attempt := 0
	heat := 1

	for {
		cls := attempt % len(classes)

		res, err := stmt.Exec(sql.Named("clsid", classes[cls].Id), sql.Named("ht", heat))
		if err != nil {
			log.WithError(err).Error("failed to execute statement")
			return err
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.WithError(err).Error("failed to count affected rows")
			return err
		}

		attempt += 1
		if ra == 0 {
			dones[cls] = true
			done := true
			for _, b := range dones {
				done = done && b
				if !done {
					break
				}
			}
			if done {
				break
			}
			continue
		}

		heat += 1
	}

	return nil
}
