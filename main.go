package main

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
	"time"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		AuthorNameHeader  string `flag:"author-name-header" default:"" description:"Header to use as Author name"`
		AuthorEmailHeader string `flag:"author-email-header" default:"" description:"Header to use as Author email"`
		DataDir           string `flag:"data-dir" default:"./data/" description:"Directory to store data to"`
		Listen            string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel          string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit    bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return fmt.Errorf("parsing CLI options: %w", err)
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log-level: %w", err)
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error

	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		fmt.Printf("wiki %s\n", version) //nolint:forbidigo
		os.Exit(0)
	}

	r := mux.NewRouter()

	r.HandleFunc("/_content/{page}", handlePageRead).Methods(http.MethodGet)
	r.HandleFunc("/_content/{page}", handlePageWrite).Methods(http.MethodPost)

	r.NotFoundHandler = http.HandlerFunc(handleIndexPage)

	var handler http.Handler = r
	handler = httpHelper.GzipHandler(handler)
	handler = httpHelper.NewHTTPLogHandler(handler)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithFields(logrus.Fields{
		"addr":    cfg.Listen,
		"version": version,
	}).Info("wiki starting")

	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("listening for HTTP traffic")
	}
}

func handleIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/app.js" && r.URL.Path != "/app.css" {
		r.URL.Path = "/index.html"
	}

	var (
		filename = path.Join("frontend", r.URL.Path)
		src      io.Reader
	)

	if _, err := os.Stat(filename); err == nil {
		f, err := os.Open(filename) //#nosec:G304 // Path is sanitized
		if err != nil {
			logrus.WithError(err).Error("Unable to open base asset")
		}
		defer func() {
			if err := f.Close(); err != nil {
				logrus.WithError(err).Error("closing frontend file (leaked fd)")
			}
		}()

		src = f
	} else if asset, err := assets.ReadFile(filename); err == nil {
		src = bytes.NewReader(asset)
	} else {
		logrus.WithField("asset", filename).Error("Asset not found in frontend dir or bundled assets")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(filename)))
	if _, err := io.Copy(w, src); err != nil {
		logrus.WithError(err).Debug("copying data to HTTP client")
	}
}

func handlePageRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	file, err := loadStoredFile(sanitizeFilename(vars["page"]))
	switch err {
	case nil:
		// All okay, render follows

	case errFileNotFound:
		initContent, err := assets.ReadFile(path.Join("default_files", sanitizeFilename(vars["page"])))
		if err != nil {
			http.Error(w, "Page not yet exists", http.StatusNotFound)
			return
		}

		// Deliver initial "Home" page
		file = &storedFile{Content: string(initContent)}

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	if err := json.NewEncoder(w).Encode(file); err != nil {
		logrus.WithError(err).Error("Unable to marshal file for JSON")
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

	if cfg.AuthorNameHeader != "" {
		file.AuthorName = r.Header.Get(cfg.AuthorNameHeader)
	}

	if cfg.AuthorEmailHeader != "" {
		file.AuthorEmail = r.Header.Get(cfg.AuthorEmailHeader)
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
