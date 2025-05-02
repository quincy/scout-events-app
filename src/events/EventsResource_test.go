package events

import (
	"github.com/pashagolub/pgxmock/v4"
	"github.com/quincy/scout-events-app/src/config"
	"github.com/quincy/scout-events-app/src/templates"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var testConfig = map[string]string{
	"DatabaseUsername": "testuser",
	"DatabasePassword": "testpassword",
	"DatabaseHostname": "localhost",
	"DatabasePort":     "5432",
	"DatabaseName":     "testdb",
	"Timezone":         "America/Boise",
}
var TestConfig = config.NewTestConfig(testConfig)

func Test_CanRenderEventsList(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Could not create mock database: %s", err)
	}
	defer db.Close()

	db.ExpectQuery(`SELECT \* FROM events WHERE end_time > NOW\(\) - INTERVAL '1 week'`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "start_time", "end_time", "summary", "description", "event_location", "assembly_location", "pickup_location"}).
			AddRow("1", "Troop Meeting1", startTime(-6), endTime(-6), "summary", "description", "event_location", "assembly_location", "pickup_location").
			AddRow("2", "Troop Meeting2", startTime(4), endTime(4), "summary", "description", "event_location", "assembly_location", "pickup_location"))

	dao := NewEventDao(config.TestConfig, db)
	rootResource := NewEventsResource(TestConfig, dao, templates.New())

	response := httptest.NewRecorder()

	rootResource.EventsListPage(response, httptest.NewRequest("GET", "/events/1", nil))

	if response.Code != 200 {
		t.Errorf("Expected status code 200, got %d", response.Code)
	}

	if err := db.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	body := response.Body.String()
	if !strings.Contains(body, "Troop Meeting1") {
		t.Errorf("Expected body to contain Troop Meeting1, got %s", body)
	}
	if !strings.Contains(body, "Troop Meeting2") {
		t.Errorf("Expected body to contain Troop Meeting, got %s", body)
	}
}

func startTime(daysFromNow int) time.Time {
	timezone := TestConfig.Timezone()

	now := time.Now().In(timezone)
	date := time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, timezone)

	return date.AddDate(0, 0, daysFromNow)
}

func endTime(daysFromNow int) time.Time {
	timezone := TestConfig.Timezone()

	now := time.Now().In(timezone)
	date := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, timezone)

	return date.AddDate(0, 0, daysFromNow)
}
