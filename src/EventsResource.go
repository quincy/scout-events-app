package main

import (
	"github.com/quincy/scout-events-app/src/events"
	"html/template"
	"log"
	"net/http"
)

type RootResource struct {
	eventDao events.EventDao
}

func NewRootResource(db events.EventDao) *RootResource {
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
