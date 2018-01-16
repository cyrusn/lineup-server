package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/cyrusn/goHTTPHelper"
	"github.com/cyrusn/lineup-system/client"
	"github.com/cyrusn/lineup-system/hub"
	"github.com/cyrusn/lineup-system/route"
	"github.com/gorilla/mux"
)

var (
	port                 string
	staticFolderLocation string
)

func init() {
	flag.StringVar(&port, "port", ":5000", "Port value")
	flag.StringVar(&staticFolderLocation, "static", "../static/dist", "location of static folder to be served")
	flag.Parse()
}

func main() {
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
