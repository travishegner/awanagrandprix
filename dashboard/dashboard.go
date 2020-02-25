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
	err  string
}

func NewDashboard() (*Dashboard, error) {
	bHead, _ := Asset("tpl/head.html")
	bFoot, _ := Asset("tpl/foot.html")
	bError, _ := Asset("tpl/error.html")
	return &Dashboard{head: string(bHead), foot: string(bFoot), err: string(bError)}, nil
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
	_, err = tpl.New("head").Parse(dash.head)
	if err != nil {
		log.WithError(err).Error("error parsing head template")
		return nil, err
	}
	_, err = tpl.New("foot").Parse(dash.foot)
	if err != nil {
		log.WithError(err).Error("error parsing foot template")
		return nil, err
	}
	_, err = tpl.New("error").Parse(dash.err)
	if err != nil {
		log.WithError(err).Error("error parsing error template")
		return nil, err
	}

	if name == "heats" {
		hl, err := Asset("tpl/heatlist.html")
		if err != nil {
			log.WithField("name", name).WithError(err).Error("failed to load heatlist template asset")
			return nil, err
		}
		he, err := Asset("tpl/heat.html")
		if err != nil {
			log.WithField("name", name).WithError(err).Error("failed to load heat template asset")
			return nil, err
		}
		_, err = tpl.New("heatlist").Parse(string(hl))
		if err != nil {
			log.WithError(err).Error("error parsing heatlist template")
			return nil, err
		}
		_, err = tpl.New("heat").Parse(string(he))
		if err != nil {
			log.WithError(err).Error("error parsing heat template")
			return nil, err
		}
	}

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
	errs := []string{}
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

	s, err := FetchSeason(id)
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
			err = AddClass(s.Id, cName)
			if err != nil {
				errs = append(errs, err.Error())
			}
		case "addcar":
			sClassId := r.FormValue("classid")
			carNumber := r.FormValue("carnumber")
			if len(carNumber) > 3 {
				carNumber = carNumber[:3]
			}
			carName := r.FormValue("carname")
			sCarWeight := r.FormValue("carweight")
			driver := r.FormValue("driver")

			if sClassId == "" || sCarWeight == "" || driver == "" {
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
			err = AddCar(s.Id, classId, carNumber, carName, carWeight, driver)
			if err != nil {
				errs = append(errs, err.Error())
			}
		case "generateheats":
			err = GenerateHeats(s)
			if err != nil {
				errs = append(errs, err.Error())
			}
		case "heattimes":
			he := r.FormValue("heatedit")
			hn, err := strconv.ParseInt(he, 10, 64)
			if err != nil {
				log.WithError(err).Warning("failed to parse heat number")
			}

			h, err := FetchHeat(s.Id, hn)
			if err != nil {
				log.WithField("heatNumber", hn).WithError(err).Error("failed to fetch heat")
				http.Error(w, "failed to fetch heat", 500)
				return
			}

			rt := r.FormValue("redtime")
			if rt != "" {
				frt, err := strconv.ParseFloat(rt, 64)
				if err != nil {
					log.WithError(err).Error("failed to parse red time")
					http.Error(w, "failed to parse red time", 500)
					return
				}
				h.Red.SetTime(frt)
			}

			bt := r.FormValue("bluetime")
			if bt != "" {
				fbt, err := strconv.ParseFloat(bt, 64)
				if err != nil {
					log.WithError(err).Error("failed to parse blue time")
					http.Error(w, "failed to parse blue time", 500)
					return
				}
				h.Blue.SetTime(fbt)
			}

			gt := r.FormValue("greentime")
			if gt != "" {
				fgt, err := strconv.ParseFloat(gt, 64)
				if err != nil {
					log.WithError(err).Error("failed to parse green time")
					http.Error(w, "failed to parse green time", 500)
					return
				}
				h.Green.SetTime(fgt)
			}

			yt := r.FormValue("yellowtime")
			if yt != "" {
				fyt, err := strconv.ParseFloat(yt, 64)
				if err != nil {
					log.WithError(err).Error("failed to parse blue time")
					http.Error(w, "failed to parse blue time", 500)
					return
				}
				h.Yellow.SetTime(fyt)
			}

			http.Redirect(w, r, fmt.Sprintf("season?id=%v&tab=heats&heatedit=%v", s.Id, hn+1), 301)
			return
		}
	}

	tname := r.FormValue("tab")

	if tname == "" {
		tname = "cars"
	}

	tpl, err := dash.getTemplate(tname)
	if err != nil {
		log.WithError(err).Error("failed to get template")
		http.Error(w, "failed to get template", 500)
		return
	}

	sp, err := NewSeasonPage(errs, s, tname)
	if err != nil {
		log.WithError(err).Error("failed to generate season page")
		http.Error(w, "failed to generate season page", 500)
		return
	}

	he := r.FormValue("heatedit")
	if he != "" {
		hid, err := strconv.ParseInt(he, 10, 64)
		sp.HeatEdit = hid
		if err != nil {
			log.WithError(err).Warningf("failed to convert %v to int64", he)
			sp.HeatEdit = 0
		}
	}

	err = tpl.Execute(w, sp)
	if err != nil {
		log.WithError(err).Error("failed to execute template")
		return
	}
}
