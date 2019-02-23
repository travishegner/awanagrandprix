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

type Car struct {
	Id     int64
	Number string
	Name   string
	Weight float64
	Driver string
	Class  *Class
}

func FetchCars(seasonId int64) ([]*Car, error) {
	stmt, err := db.Prepare(`
select c.id, c.number, c.name, c.weight, c.driver, cls.id, cls.Name
from cars c
inner join classes cls on c.class_id=cls.id
where cls.season_id=:sid
order by c.Id
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

	cars := []*Car{}
	for rows.Next() {
		var id int64
		var number string
		var name string
		var weight float64
		var driver string
		var classId int64
		var className string
		err = rows.Scan(&id, &number, &name, &weight, &driver, &classId, &className)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		c := &Car{
			Id:     id,
			Number: number,
			Name:   name,
			Weight: weight,
			Driver: driver,
			Class:  &Class{Id: classId, Name: className},
		}
		cars = append(cars, c)
	}

	return cars, nil
}

func AddCar(seasonId, classId int64, number string, name string, weight float64, driver string) error {
	cars, err := FetchCars(seasonId)
	if err != nil {
		return err
	}

	//populate a map with car numbers already in use
	usedNums := map[string]struct{}{}
	for _, c := range cars {
		usedNums[c.Number] = struct{}{}
	}

	if number == "" {
		//generate a new random car number
		//this can crash (infinite loop) if 1-99 are already used
		for {
			number = fmt.Sprintf("%v", rand.Intn(98)+1)
			if _, ok := usedNums[number]; ok {
				number = ""
			}
			if number != "" {
				break
			}
		}
	}

	stmt, err := db.Prepare("insert into cars (class_id,number,name,weight,driver) values (:cid,:num,:name,:wt,:drv)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(
		sql.Named("cid", classId),
		sql.Named("num", number),
		sql.Named("name", name),
		sql.Named("wt", weight),
		sql.Named("drv", driver),
	)
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return err
	}

	return nil
}
