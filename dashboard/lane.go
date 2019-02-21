package dashboard

import log "github.com/sirupsen/logrus"

type Lane struct {
	Id    int64
	Color string
}

func FetchLanes() ([]*Lane, error) {
	stmt, err := db.Prepare("select id, color from lanes")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return nil, err
	}

	lanes := []*Lane{}
	for rows.Next() {
		var id int64
		var color string
		err = rows.Scan(&id, &color)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		l := &Lane{
			Id:    id,
			Color: color,
		}
		lanes = append(lanes, l)
	}

	return lanes, nil
}
