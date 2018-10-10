package dashboard

import (
	"mime"
	"net/http"
	"net/http/httputil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

//go:generate go-bindata -debug -prefix "pub/" -pkg dashboard -o assets.go pub/...

type Dashboard struct{}

func NewDashboard() (*Dashboard, error) {
	return &Dashboard{}, nil
}

func (db *Dashboard) Start() error {
	http.HandleFunc("/", db.handlePage)
	http.ListenAndServe(":8080", nil)

	return nil
}

func (db *Dashboard) handlePage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	if url[len(url)-1:] == "/" {
		url = url[1:] + "index.html"
	}

	w.Write([]byte(url))
	l := log.WithField("url", url)
	l.Debug("dashboard request")

	a, err := Asset(url)
	if err != nil {
		http.NotFound(w, r)
		l.Error("failed to load asset")
		return
	}
	f, err := AssetInfo(url)
	if err != nil {
		http.NotFound(w, r)
		l.Error("failed to load asset info")
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(f.Name())))
	var b int
	b, err = w.Write(a)
	if err != nil {
		l.Error("error writing content")
	}
	l.WithField("bytes", b).Debug("bytes written")
}
