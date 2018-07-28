package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	helper "github.com/cyrusn/goHTTPHelper"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/cyrusn/lineup-system/model/auth"
	"github.com/cyrusn/lineup-system/model/schedule"
	"github.com/cyrusn/lineup-system/route"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {

		auth.UpdateLifeTime(lifeTime)

		paths := []string{dbPath, staticFolderLocation}
		checkPathExist(paths)

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}

		store := createStore(db, &secret)
		r := mux.NewRouter()

		handleRoute(r, store)
		serveStaticFolder(r, staticFolderLocation)

		location := "localhost" + port
		printNotice(dbPath, staticFolderLocation, location)
		http.ListenAndServe(location, helper.Logger(r))
	},
}

func serveStaticFolder(r *mux.Router, staticFolderLocation string) {
	staticFolder := http.Dir(staticFolderLocation)

	// serve static file
	r.PathPrefix("/").Handler(
		http.FileServer(staticFolder),
	)
}

func handleRoute(r *mux.Router, s *route.Store) {
	for _, ro := range route.Routes(s) {
		handler := http.HandlerFunc(ro.Handler)

		if len(ro.Scopes) != 0 {
			handler = secret.Access(ro.Scopes, handler).(http.HandlerFunc)
		}

		if ro.Auth {
			handler = secret.Authenticate(handler).(http.HandlerFunc)
		}

		r.
			PathPrefix("/api/").
			Methods(ro.Methods...).
			Path(ro.Path).
			HandlerFunc(handler)
	}
}

func createStore(db *sql.DB, secret *auth_helper.Secret) *route.Store {
	return &route.Store{
		&auth.DB{db, secret},
		&schedule.DB{db},
	}
}

func printNotice(dbPath, staticFolderLocation, location string) {
	fmt.Printf(
		"Database: \"%s\"\nStatic folder: \"%s\"\nServer start on http://%s\n",
		dbPath,
		staticFolderLocation,
		location,
	)
}
