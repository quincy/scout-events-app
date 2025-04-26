package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/quincy/scout-events-app/src/database"
	"github.com/quincy/scout-events-app/src/events"
	"github.com/quincy/scout-events-app/src/healthcheck"
	"github.com/quincy/scout-events-app/src/templates"
	"log"
	"net/http"
	"strings"
)

func main() {
	ctx := context.Background()

	dbConfig, err := database.MakeDatabaseConfig()
	if err != nil {
		log.Fatalf("Could not create database config: %v", err)
	}
	conn, err := database.CreateDbConnection(ctx, dbConfig)
	defer conn.Close()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(routingMiddleware)
	registerRoutes(r, events.NewEventDao(conn), templates.New())

	log.Println("Starting server at http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}
}

func registerRoutes(r *mux.Router, db events.EventDao, templates templates.Templates) {
	r.HandleFunc("/healthcheck", healthcheck.Handler).Methods("GET")
	r.HandleFunc("/deepcheck", healthcheck.DeepcheckHandler).Methods("GET")
	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	eventsResource := events.NewEventsResource(db, templates)
	r.HandleFunc("/events", eventsResource.EventsListPage).Methods("GET")
	r.HandleFunc("/events/{id}", eventsResource.EventDetailsPage).Methods("GET")

	rootResource := newRootResource(db, eventsResource)
	r.HandleFunc("/", rootResource.Home).Methods("GET") // will be replaced with a homepage
}

func routingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := *r.URL
		url.Path = strings.TrimSuffix(r.URL.Path, "/")
		r.URL = &url

		h.ServeHTTP(w, r)
	})
}

type RootResource struct {
	eventDao       events.EventDao
	eventsResource events.Resource
}

func newRootResource(db events.EventDao, eventsResource events.Resource) *RootResource {
	return &RootResource{eventDao: db, eventsResource: eventsResource}
}

func (root RootResource) Home(w http.ResponseWriter, r *http.Request) {
	root.eventsResource.EventsListPage(w, r)
}
