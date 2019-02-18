package dashboard

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id   int64
	Name string
}

func (dash *Dashboard) NewSeason(name string) (int64, error) {
	stmt, err := dash.db.Prepare("insert into seasons (name) values (?)")
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

func (dash *Dashboard) GetSeasons() ([]*Season, error) {
	rows, err := dash.db.Query("select id, name from seasons")
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

func (dash *Dashboard) GetSeason(id int64) (*Season, error) {
	fmt.Println(id)
	stmt, err := dash.db.Prepare("select name from seasons where id=?")
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

func (dash *Dashboard) seasonsHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := Asset("tpl/seasons.html")
	tpl, _ := template.New("seasons").Parse(string(b))
	tpl.New("head").Parse(dash.head)
	tpl.New("foot").Parse(dash.foot)
	data, _ := dash.GetSeasons()
	tpl.Execute(w, data)
}

func (dash *Dashboard) seasonHandler(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("failed to convert %v to int64", sid)
		return
	}
	b, _ := Asset("tpl/season.html")
	tpl, _ := template.New("season").Parse(string(b))
	tpl.New("head").Parse(dash.head)
	tpl.New("foot").Parse(dash.foot)
	data, _ := dash.GetSeason(id)
	tpl.Execute(w, data)
}
