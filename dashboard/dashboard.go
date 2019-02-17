package dashboard

import (
	"database/sql"
	"html/template"
	"mime"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	dbfile = "agp.db"
)

//go:generate go-bindata -debug -prefix "pub/" -pkg dashboard -o assets.go pub/...

type Dashboard struct {
	db *sql.DB
}

func NewDashboard() (*Dashboard, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.WithField("dbfile", dbfile).Error("Failed to open database.")
		return nil, err
	}
	return &Dashboard{db: db}, nil
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

	//base, _ := Asset("tpl/base.html")
	switch url {
	case "":
		tpl, _ := template.New("").ParseFiles("dashboard/pub/tpl/seasons.html", "dashboard/pub/tpl/base.html")
		data, _ := dash.GetSeasons()
		tpl.ExecuteTemplate(w, "base", data)
		return
	case "season":
		return
	}

	a, err := Asset(url)
	if err != nil {
		http.NotFound(w, r)
		l.Error("failed to load asset")
		return
	}
	f, err := AssetInfo(url)
	if err != nil {
		http.Error(w, "failed to load asset info", 500)
		l.Error("failed to load asset info")
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
