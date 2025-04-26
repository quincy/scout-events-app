package events

import (
	"github.com/gorilla/mux"
	"github.com/quincy/scout-events-app/src/templates"
	"log"
	"net/http"
)

type Resource interface {
	EventsListPage(w http.ResponseWriter, r *http.Request)
	EventDetailsPage(w http.ResponseWriter, r *http.Request)
}

type resource struct {
	eventDao  EventDao
	templates templates.Templates
}

func NewEventsResource(db EventDao, templates templates.Templates) Resource {
	return &resource{eventDao: db, templates: templates}
}

func (res *resource) EventsListPage(w http.ResponseWriter, r *http.Request) {
	upcomingEvents, err := res.eventDao.GetUpcomingEvents(r.Context())
	if err != nil {
		log.Printf("Error getting upcoming events: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t, err := res.templates.GetEventListTemplate()
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

func (res *resource) EventDetailsPage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	event, err := res.eventDao.GetEventById(r.Context(), id)
	if err != nil {
		log.Printf("Error loading event details: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t, err := res.templates.GetEventDetailsTemplate()
	if err != nil {
		log.Printf("parsing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = t.Execute(w, event)
	if err != nil {
		log.Printf("executing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
