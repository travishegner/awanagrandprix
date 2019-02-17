package dashboard

import (
	"html/template"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Season struct {
	Id   int
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
	stmt, err := dash.db.Prepare("select id, name from seasons where id=?")
	if err != nil {
		log.WithError(err).Error("failed to prepare")
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		log.WithError(err).Error("failed to get season")
		return nil, err
	}

	season := &Season{}
	for rows.Next() {
		err = rows.Scan(&season.Id, &season.Name)
		if err != nil {
			log.WithError(err).Error("failed to get session info")
			return nil, err
		}
	}

	return season, nil
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
	id, _ := strconv.ParseInt(r.FormValue("id"), 64, 0)
	b, _ := Asset("tpl/season.html")
	tpl, _ := template.New("season").Parse(string(b))
	tpl.New("head").Parse(dash.head)
	tpl.New("foot").Parse(dash.foot)
	data, _ := dash.GetSeason(id)
	tpl.Execute(w, data)
}
