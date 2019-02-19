package dashboard

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id      int64
	Name    string
	Classes []*Class
	Cars    []*Car
}

func AddSeason(name string) (int64, error) {
	stmt, err := db.Prepare("insert into seasons (name) values (?)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return -1, err
	}

	res, err := stmt.Exec(name)
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return -1, err
	}

	li, err := res.LastInsertId()
	if err != nil {
		log.WithError(err).Error("Failed to get last insert id.")
		return -1, err
	}

	return li, err
}

func GetSeasons() ([]*Season, error) {
	rows, err := db.Query("select id, name from seasons")
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return nil, err
	}

	seasons := []*Season{}
	for rows.Next() {
		s := &Season{}
		err = rows.Scan(&s.Id, &s.Name)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			return nil, err
		}
		seasons = append(seasons, s)
	}

	return seasons, nil
}

func GetSeason(id int64) (*Season, error) {
	stmt, err := db.Prepare("select name from seasons where id=?")
	if err != nil {
		log.WithError(err).Error("failed to prepare")
		return nil, err
	}

	var name string
	err = stmt.QueryRow(id).Scan(&name)
	if err != nil {
		log.WithError(err).Error("failed to get season info")
		return nil, err
	}

	return &Season{Id: id, Name: name}, nil
}

func (s *Season) AddClass(name string) error {
	stmt, err := db.Prepare("insert into classes (season_id,name) values (?,?)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(s.Id, name)
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return err
	}

	return nil
}

func (s *Season) LoadClasses() error {
	stmt, err := db.Prepare("select id, name from classes where season_id=?")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	rows, err := stmt.Query(s.Id)
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return err
	}

	for rows.Next() {
		var id int64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.WithError(err).Error("Failed to read row.")
			continue
		}
		s.Classes = append(s.Classes, &Class{Id: id, Name: name})
	}

	return nil
}

func (s *Season) AddCar(classId int64, name string, weight float64, driver string) error {
	//make sure we've got all the cars from the db
	if len(s.Cars) == 0 {
		s.LoadCars()
	}

	if len(s.Cars) >= 99 {
		err := fmt.Errorf("Reached the maximum number of cars. Rebuild with a higher random max.")
		log.WithError(err).Error("Couldn't add another car.")
		return err
	}

	//populate a map with car numbers already in use
	usedNums := map[int64]struct{}{}
	for _, c := range s.Cars {
		usedNums[c.Number] = struct{}{}
	}

	//generate a new random car number
	number := int64(0)
	for {
		number = int64(rand.Intn(98) + 1)
		if _, ok := usedNums[number]; ok {
			number = 0
		}
		if number > 0 {
			break
		}
	}

	stmt, err := db.Prepare("insert into cars (season_id,class_id,number,name,weight,driver) values (?,?,?,?,?,?)")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	_, err = stmt.Exec(s.Id, classId, number, name, weight, driver)
	if err != nil {
		log.WithError(err).Error("Failed to execute insert statement.")
		return err
	}

	return nil
}

func (s *Season) LoadCars() error {
	s.Cars = []*Car{}
	stmt, err := db.Prepare("select c.id, c.number, c.name, c.weight, c.driver, cls.Name from cars c inner join classes cls on c.class_id=cls.id where c.season_id=?")
	if err != nil {
		log.WithError(err).Error("Failed to prepare statement.")
		return err
	}

	rows, err := stmt.Query(s.Id)
	if err != nil {
		log.WithError(err).Error("Failed to execute query.")
		return err
	}

	for rows.Next() {
		var id int64
		var number int64
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
		s.Cars = append(s.Cars, c)
	}

	return nil
}
