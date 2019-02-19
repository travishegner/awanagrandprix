package dashboard

import (
	"database/sql"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	dbfile = "agp.db"
	db     *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
}

//go:generate go-bindata -debug -pkg dashboard -o assets.go tpl/...

type Dashboard struct {
	head string
	foot string
}

func NewDashboard() (*Dashboard, error) {
	bHead, _ := Asset("tpl/head.html")
	bFoot, _ := Asset("tpl/foot.html")
	return &Dashboard{head: string(bHead), foot: string(bFoot)}, nil
}

func (dash *Dashboard) Start() error {
	http.HandleFunc("/", dash.handlePage)
	http.ListenAndServe(":8080", nil)

	return nil
}

func (dash *Dashboard) handlePage(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("url", r.URL.Path)
	url := r.URL.Path[1:]

	w.Header().Set("Content-Type", "text/html")

	switch url {
	case "":
		dash.handleSeasons(w, r)
		return
	case "season":
		dash.handleSeason(w, r)
		return
	}

	a, err := Asset(url)
	if err != nil {
		http.NotFound(w, r)
		l.WithError(err).Error("failed to load asset")
		return
	}
	f, err := AssetInfo(url)
	if err != nil {
		l.WithError(err).Error("failed to load asset info")
		http.Error(w, "failed to load asset info", 500)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(f.Name())))
	l.WithField("ct", w.Header().Get("Content-Type")).Debug("content-type")
	var b int
	b, err = w.Write(a)
	if err != nil {
		l.Error("error writing content")
	}
	l.WithField("bytes", b).Debug("bytes written")
}

func (dash *Dashboard) getTemplate(name string) (*template.Template, error) {
	b, err := Asset("tpl/" + name + ".html")
	if err != nil {
		log.WithField("name", name).WithError(err).Error("Error getting template asset.")
		return nil, err
	}
	tpl, err := template.New(name).Parse(string(b))
	if err != nil {
		log.WithField("name", name).WithError(err).Error("Error parsing template bytes.")
		return nil, err
	}
	tpl.New("head").Parse(dash.head)
	tpl.New("foot").Parse(dash.foot)

	return tpl, nil
}

func (dash *Dashboard) handleSeasons(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("seasonname")
		id, err := AddSeason(name)
		if err != nil {
			log.WithError(err).Error("Failed to create new season.")
			http.Error(w, "failed to create new season", 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("season?id=%v", id), 301)
		return
	}

	tpl, err := dash.getTemplate("seasons")
	if err != nil {
		log.WithError(err).Error("failed to get template")
		return
	}

	data, _ := GetSeasons()
	if err != nil {
		log.WithError(err).Error("failed to get seasons")
		return
	}

	tpl.Execute(w, data)
}

func (dash *Dashboard) handleSeason(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("id")
	if sid == "" {
		http.Error(w, "missing id", 400)
		return
	}
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("failed to convert %v to int64", sid)
		return
	}

	s, err := GetSeason(id)
	if err != nil {
		log.WithError(err).Error("failed to get season")
		http.Error(w, "failed to get season", 500)
		return
	}

	if r.Method == "POST" {
		switch r.FormValue("action") {
		case "addclass":
			cName := r.FormValue("classname")
			if cName == "" {
				break
			}
			s.AddClass(cName)
		case "addcar":
			sClassId := r.FormValue("classid")
			carName := r.FormValue("carname")
			sCarWeight := r.FormValue("carweight")
			driver := r.FormValue("driver")

			if sClassId == "" || carName == "" || sCarWeight == "" || driver == "" {
				break
			}

			classId, err := strconv.ParseInt(sClassId, 10, 64)
			if err != nil {
				log.WithError(err).Error("failed to parse class id")
				break
			}

			carWeight, err := strconv.ParseFloat(sCarWeight, 64)
			if err != nil {
				log.WithError(err).Error("failed to parse car weight")
				break
			}
			s.AddCar(classId, carName, carWeight, driver)
		}
	}

	tname := r.FormValue("tab")

	if tname == "" {
		tname = "season"
		s.LoadClasses()
		s.LoadCars()
	}

	tpl, err := dash.getTemplate(tname)
	if err != nil {
		log.WithError(err).Error("failed to get template")
		http.Error(w, "failed to get template", 500)
		return
	}

	tpl.Execute(w, s)
}
