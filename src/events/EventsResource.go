package events

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Resource interface {
	EventDetailsPage(w http.ResponseWriter, r *http.Request)
	UpcomingEvents(w http.ResponseWriter, r *http.Request)
}

type resource struct {
	eventDao EventDao
}

func NewEventsResource(db EventDao) Resource {
	return &resource{db}
}

func (res *resource) EventDetailsPage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	event, err := res.eventDao.GetEventById(r.Context(), id)
	if err != nil {
		log.Printf("Error loading event details: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t, err := template.New("event-details.gohtml").ParseFiles("./templates/event-details.gohtml")
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

func (res *resource) UpcomingEvents(w http.ResponseWriter, r *http.Request) {
	_, err := loadEventListTemplate()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("parsing template failed. err=%s", err)
	}

}

func loadEventListTemplate() (*template.Template, error) {
	return template.New("eventList").Parse(`
	    {{if .}}
	      <table id="events-table" data-testid="events-table">
	        <tr>
	          <td>Date</td>
	          <td>Event</td>
	          <td>Summary</td>
	          <td>Location</td>
	        </tr>
	        {{range .}}
	        {{- /*gotype: github.com/quincy/scout-events-app/src/events.Event*/ -}}
	        <tr>
	          <td>{{.StartTime.Format "2006-01-02 15:04" }}</td>
	          <td>{{.Name}}</td>
	          <td>{{.Summary}}</td>
	          <td>{{.EventLocation}}</td>
	        </tr>
	        {{end}}
	      </table>
	    {{else}}
	      <h3>No upcoming events</h3>
	    {{end}}`)
}
