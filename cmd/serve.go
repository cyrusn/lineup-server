package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	helper "github.com/cyrusn/goHTTPHelper"
	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/cyrusn/lineup-system/auth"
	"github.com/cyrusn/lineup-system/route"
	"github.com/cyrusn/lineup-system/schedule"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if database file is already exist
		locations := []string{dbPath, staticFolderLocation}
		for _, location := range locations {
			if _, err := os.Stat(location); os.IsNotExist(err) {
				log.Fatalf(`"%s" doesn't exist`, location)
			}
		}

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}

		var secret = auth_helper.New(
			"authClaim",
			"jwt",
			"Role",
			[]byte("skill-vein-planet-neigh-envoi"),
		)

		r := mux.NewRouter()

		store := route.Store{
			&auth.DB{db},
			&schedule.DB{db},
			secret,
		}

		for _, ro := range route.Routes(&store) {
			handler := http.HandlerFunc(ro.Handler)

			if ro.Auth {
				handler = secret.Authenticate(handler).(http.HandlerFunc)
			}

			r.
				PathPrefix("/api/").
				Methods(ro.Methods...).
				Path(ro.Path).
				HandlerFunc(handler)
		}

		staticFolder := http.Dir(staticFolderLocation)

		// serve static file
		r.PathPrefix("/").Handler(
			http.FileServer(staticFolder),
		)

		location := "localhost" + port
		fmt.Printf(
			"Database: \"%s\"\nStatic folder: \"%s\"\nServer start on http://%s\n",
			dbPath,
			staticFolderLocation,
			location,
		)

		http.ListenAndServe(location, helper.Logger(r))
	},
}
