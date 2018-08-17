package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	helper "github.com/cyrusn/goHTTPHelper"

	auth_helper "github.com/cyrusn/goJWTAuthHelper"

	"github.com/cyrusn/lineup-server/model/auth"
	"github.com/cyrusn/lineup-server/model/schedule"
	"github.com/cyrusn/lineup-server/route"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Line-Up System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {

		auth.UpdateLifeTime(lifeTime)

		checkPathExist(staticFolderLocation)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		store := createStore(db, &secret)
		r := mux.NewRouter()

		handleRoute(r, store)
		serveStaticFolder(r, staticFolderLocation)

		location := "localhost" + port
		printNotice(dsn, staticFolderLocation, location)
		http.ListenAndServe(location, helper.Logger(r))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(
		&port,
		"port",
		"p",
		DEFAULT_PORT,
		"port value",
	)
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))

	serveCmd.PersistentFlags().StringVarP(
		&staticFolderLocation,
		"static",
		"s",
		DEFAULT_PUBLIC_FOLDER,
		"location of static folder for serving",
	)
	viper.BindPFlag("static", serveCmd.PersistentFlags().Lookup("static"))

	serveCmd.PersistentFlags().Int64VarP(
		&lifeTime,
		"time",
		"t",
		DEFAULT_JWT_EXPIRE_TIME,
		"update the life time (minutes) of jwt",
	)
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

func printNotice(dsn, staticFolderLocation, location string) {
	fmt.Printf(
		"Static folder: \"%s\"\nServer start on http://%s\n",
		staticFolderLocation,
		location,
	)
}
