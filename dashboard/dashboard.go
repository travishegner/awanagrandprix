package dashboard

import (
	"database/sql"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	dbfile = "agp.db"
)

//go:generate go-bindata -debug -pkg dashboard -o assets.go tpl/...

type Dashboard struct {
	db   *sql.DB
	head string
	foot string
}

func NewDashboard() (*Dashboard, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.WithField("dbfile", dbfile).Error("Failed to open database.")
		return nil, err
	}
	bHead, _ := Asset("tpl/head.html")
	bFoot, _ := Asset("tpl/foot.html")
	return &Dashboard{db: db, head: string(bHead), foot: string(bFoot)}, nil
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
		if r.Method == "POST" {
			name := r.FormValue("seasonname")
			id, err := dash.NewSeason(name)
			if err != nil {
				l.WithError(err).Error("Failed to create new season.")
				http.Error(w, "failed to create new season", 500)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("season?id=%v", id), 301)
			return
		}
		dash.seasonsHandler(w, r)
		return
	case "season":
		dash.seasonHandler(w, r)
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
