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
	"github.com/cyrusn/lineup-system/ws"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

const (
	contextKeyName = "authClaim"
	jwtKeyName     = "jwt"
	roleKeyName    = "Role"
	privateKey     = "skill-vein-planet-neigh-envoi"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {
		paths := []string{dbPath, staticFolderLocation}
		checkPathExist(paths)

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		secret := auth_helper.New(contextKeyName, jwtKeyName, roleKeyName, []byte(privateKey))
		store := createStore(db, secret)

		r := mux.NewRouter()
		h := ws.NewHub()
		go h.Run()

		handleRoute(r, store, h)
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

func handleRoute(r *mux.Router, s *route.Store, h *ws.Hub) {
	for _, ro := range route.Routes(s, h) {
		handler := http.HandlerFunc(ro.Handler)

		if ro.Auth {
			handler = s.Secret.Authenticate(handler).(http.HandlerFunc)
		}

		r.
			PathPrefix("/api/").
			Methods(ro.Methods...).
			Path(ro.Path).
			HandlerFunc(handler)
	}
}

func createStore(db *sql.DB, secret auth_helper.Secret) *route.Store {
	return &route.Store{
		&auth.DB{db},
		&schedule.DB{db},
		secret,
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
