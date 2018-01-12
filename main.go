package main

import (
	"github.com/cyrusn/lineup-system/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}

/*
import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/client"
	"github.com/cyrusn/lineup-system/hub"
	"github.com/cyrusn/lineup-system/route"
	"github.com/gorilla/mux"
	// https://github.com/spf13/cobra
)

// TODO: Use cobra to provide commands as listed below
// 	- serve (use for serve the backend server)
// 	- init (create database and import students list)

const version = "2.0.0"

var (
	port                 string
	staticFolderLocation string
	versionFlag          bool
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "Prints current version")

	flag.StringVar(&port, "port", ":5000", "Port value")
	flag.StringVar(&staticFolderLocation, "static", "../static/dist", "Location of static folder to be served")
	flag.Usage = func() {
		fmt.Printf("Server of line-up system:\n")
		fmt.Printf("Usage:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	if versionFlag {
		fmt.Printf("Version %s\n", version)
		os.Exit(0)
	}

	r := mux.NewRouter()
	h := hub.New()

	go h.Run()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client.ServeWS(h, w, r)
	})

	for _, ro := range route.Routes(h) {
		r.
			PathPrefix("/api/").
			Methods(ro.Methods...).
			Path(ro.Path).
			HandlerFunc(ro.Handler)
	}

	staticFolder := http.Dir(staticFolderLocation)
	r.PathPrefix("/").Handler(
		http.FileServer(staticFolder),
	)

	location := "localhost" + port
	fmt.Printf("Server start on http://%s\n", location)
	http.ListenAndServe(location, helper.Logger(r))
}
*/
