package main

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -modtime 1 -md5checksum ./frontend/...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		DataDir        string `flag:"data-dir" default:"./data/" description:"Directory to store data to"`
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("wiki %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/_content/{page}", handlePageRead).Methods(http.MethodGet)
	r.HandleFunc("/_content/{page}", handlePageWrite).Methods(http.MethodPost)

	r.NotFoundHandler = http.HandlerFunc(handleIndexPage)

	var handler http.Handler = r
	handler = httpHelper.GzipHandler(handler)
	handler = httpHelper.NewHTTPLogHandler(handler)

	http.ListenAndServe(cfg.Listen, handler)
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/app.js" {
		r.URL.Path = "/index.html"
	}

	var (
		filename = path.Join("frontend", r.URL.Path)
		src      io.Reader
	)

	if _, err := os.Stat(filename); err == nil {
		f, err := os.Open(filename)
		if err != nil {
			log.WithError(err).Error("Unable to open base asset")
		}
		defer f.Close()
		src = f
	} else if asset, err := Asset(filename); err == nil {
		src = bytes.NewReader(asset)
	} else {
		log.WithField("asset", filename).Error("Asset not found in frontend dir or bundled assets")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(filename)))
	io.Copy(w, src)
}

func handlePageRead(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)

	file, err := loadStoredFile(sanitizeFilename(vars["page"]))
	switch err {

	case nil:
		// All okay, render follows

	case errFileNotFound:
		http.Error(w, "Page not yet exists", http.StatusNotFound)
		return

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(file); err != nil {
		log.WithError(err).Error("Unable to marshal file for JSON")
	}
}

func handlePageWrite(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		file = &storedFile{}
	)

	if err := json.NewDecoder(r.Body).Decode(file); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := file.Save(sanitizeFilename(vars["page"])); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func sanitizeFilename(page string) string {
	return strings.Join([]string{slug.Make(page), "md"}, ".")
}
