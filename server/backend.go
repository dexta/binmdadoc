package serv

import (
	"context"
  "flag"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var ExecPath = "./"

type spaHandler struct {
	staticPath string
	indexPath  string
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func StartRouter() {
	// TODO make it confi
	ExecPath = "./"
	// ex, err := os.Executable()
	// check(err)
	// ExecPath = filepath.Dir(ex)


	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/echo/markdown", ProxyMarkdown).Methods("POST")
	
	router.HandleFunc("/api/echo/asciidoc", ProxyAsciidoc).Methods("POST")
	
	router.HandleFunc("/api/add/doc", NewFileWriter).Methods("POST")
	
	router.HandleFunc("/api/update/doc", UpdateFileWirter).Methods("POST")
	
	router.HandleFunc("/api/doc/index", GetDocumentIndex).Methods("GET")
	
	router.HandleFunc("/api/doc/by/id/{docId}", GetDocumentById).Methods("GET")
	
	router.HandleFunc("/view/{docId}", HTMLById).Methods("GET")

	spa := spaHandler{staticPath: "frontend", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
        Handler: router,
	}

	go func() {
	    if err := srv.ListenAndServe(); err != nil {
	        log.Println(err)
	    }
    }()

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    // Block until we receive our signal.
    <-c

    ctx, cancel := context.WithTimeout(context.Background(), wait)
    defer cancel()
    srv.Shutdown(ctx)
    log.Println("shutting down")
    os.Exit(0)
}
