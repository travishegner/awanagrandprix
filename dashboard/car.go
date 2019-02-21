package dashboard

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

type Car struct {
	Id     int64
	Number string
	Name   string
	Weight float64
	Driver string
	Class  string
}

func FetchCars(seasonId int64) ([]*Car, error) {
	stmt, err := db.Prepare("select c.id, c.number, c.name, c.weight, c.driver, cls.Name from cars c inner join classes cls on c.class_id=cls.id where c.season_id=?")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return nil, err
	}

	rows, err := stmt.Query(seasonId)
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
		var class string
		err = rows.Scan(&id, &number, &name, &weight, &driver, &class)
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
			Class:  class,
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

	stmt, err := db.Prepare("insert into cars (season_id,class_id,number,name,weight,driver) values (?,?,?,?,?,?)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(seasonId, classId, number, name, weight, driver)
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return err
	}

	return nil
}
