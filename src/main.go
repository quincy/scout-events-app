package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/quincy/scout-events-app/database"
	"github.com/quincy/scout-events-app/events"
	"github.com/quincy/scout-events-app/healthcheck"
	"html/template"
	"log"
	"net/http"
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
	registerRoutes(r, events.NewEventDao(conn))

	log.Println("Starting server at http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Unhandled error: %s", err)
	}
}

func registerRoutes(r *mux.Router, db events.EventDao) {
	r.HandleFunc("/healthcheck", healthcheck.Handler).Methods("GET")
	r.HandleFunc("/deepcheck", healthcheck.DeepcheckHandler).Methods("GET")
	fileServer := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	rootResource := newRootResource(db)
	r.HandleFunc("/", rootResource.Home).Methods("GET") // will be replaced with a homepage

	eventsResource := events.NewEventsResource(db)
	r.HandleFunc("/events", eventsResource.UpcomingEvents).Methods("GET")
	r.HandleFunc("/events/{id}", eventsResource.EventDetailsPage).Methods("GET")
}

type RootResource struct {
	eventDao events.EventDao
}

func newRootResource(db events.EventDao) *RootResource {
	return &RootResource{db}
}

func (root RootResource) Home(w http.ResponseWriter, r *http.Request) {
	upcomingEvents, err := root.eventDao.GetUpcomingEvents(r.Context())
	if err != nil {
		log.Printf("Error getting upcoming events: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t, err := template.New("events-list.gohtml").ParseFiles("./templates/events-list.gohtml")
	if err != nil {
		log.Printf("parsing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = t.Execute(w, upcomingEvents)
	if err != nil {
		log.Printf("executing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
