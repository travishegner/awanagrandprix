package dashboard

import log "github.com/sirupsen/logrus"

type Run struct {
	Id     int64
	CarId  int64
	LaneId int64
	Heat   int64
	Time   float64
}

func FetchRuns(seasonId int64) ([]*Run, error) {
	stmt, err := db.Prepare("select id, car_id, lane_id, heat, time from runs where season_id=?")
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
		var heat int64
		var time float64
		err = rows.Scan(&id, &carid, &laneid, &heat, &time)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		runs = append(runs, &Run{Id: id, CarId: carid, LaneId: laneid, Heat: heat, Time: time})
	}

	return runs, nil
}
