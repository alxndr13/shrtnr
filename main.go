package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/cli/v2"

	"github.com/charmbracelet/log"

	bolt "go.etcd.io/bbolt"
)

var (
	//go:embed pages/* base.html
	html embed.FS
)

type App struct {
	Db      *bolt.DB
	DbPath  string
	RootUrl string
	Port    string
}

func main() {
	var a App

	app := &cli.App{
		Name:  "shrtnr",
		Usage: "starts the shrtnr application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "root-url",
				Destination: &a.RootUrl,
				EnvVars:     []string{"SH_ROOT_URL"},
				DefaultText: "http://localhost:8000/",
				Value:       "http://localhost:8000/",
			},
			&cli.StringFlag{
				Name:        "db-path",
				Destination: &a.DbPath,
				EnvVars:     []string{"SH_DB_PATH"},
				Value:       "./shrtnr.db",
				DefaultText: "./shrtnr.db",
			},
			&cli.StringFlag{
				Name:        "port",
				Destination: &a.Port,
				EnvVars:     []string{"SH_PORT"},
				Value:       "8000",
				DefaultText: "8000",
			},
		},
		Action: func(*cli.Context) error {
			return a.startApp()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func (a *App) startApp() error {
	if !strings.HasSuffix(a.RootUrl, "/") {
		a.RootUrl += "/"
	}

	err := a.createDatabaseIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting Server", "port", a.Port)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// define handlers
	r.Get("/", defaultHandler)
	r.Get("/r/{id}", a.redirectHandler)
	r.Post("/shorten", a.shortenHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", a.Port), r))
	return nil
}
