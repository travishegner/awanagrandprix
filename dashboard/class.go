package dashboard

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Class struct {
	Id   int64
	Name string
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
		classes = append(classes, &Class{Id: id, Name: name})
	}

	return classes, nil
}
