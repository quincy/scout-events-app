package events

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quincy/scout-events-app/src/config"
	"github.com/quincy/scout-events-app/src/templates"
	"log"
	"net/http"
	"time"
)

type Resource interface {
	EventsListPage(w http.ResponseWriter, r *http.Request)
	EventDetailsPage(w http.ResponseWriter, r *http.Request)
	CreateEventPage(w http.ResponseWriter, r *http.Request)
	CreateEvent(w http.ResponseWriter, r *http.Request)
}

type resource struct {
	cfg       config.AppConfig
	eventDao  EventDao
	templates templates.Templates
}

func NewEventsResource(cfg config.AppConfig, db EventDao, templates templates.Templates) Resource {
	return &resource{
		cfg:       cfg,
		eventDao:  db,
		templates: templates,
	}
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
	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("No event ID provided in URL")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res.eventDetailsPage(w, r, id)
}

func (res *resource) eventDetailsPage(w http.ResponseWriter, r *http.Request, id string) {
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

func (res *resource) CreateEventPage(w http.ResponseWriter, r *http.Request) {
	t, err := res.templates.GetCreateEventTemplate()
	if err != nil {
		log.Printf("parsing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("executing template failed. err=%s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (res *resource) CreateEvent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("Form values: %v", r.PostForm)

	startTime, err := time.ParseInLocation("2006-01-02T15:04", r.FormValue("start_time"), res.cfg.Timezone())
	if err != nil {
		log.Printf("Error parsing start time: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02T15:04", r.FormValue("end_time"), res.cfg.Timezone())
	if err != nil {
		log.Printf("Error parsing end time: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Could not generate new UUID: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	event := Event{
		Id:               id.String(),
		Name:             r.FormValue("name"),
		StartTime:        startTime,
		EndTime:          endTime,
		Summary:          r.FormValue("summary"),
		Description:      r.FormValue("description"),
		EventLocation:    r.FormValue("event_location"),
		AssemblyLocation: r.FormValue("assembly_location"),
		PickupLocation:   r.FormValue("pickup_location"),
	}
	log.Printf("Creating event: %v", event)

	// Create the event in the database
	createdEvent, err := res.eventDao.CreateEvent(r.Context(), event)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Created event: %v", createdEvent)

	eventId := createdEvent.Id

	// Redirect to the event details page for the new event
	w.Header().Set("HX-Replace-Url", "/events/"+eventId)
	res.eventDetailsPage(w, r, eventId)
}
