package events

import (
	"context"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"testing"
	"time"
)

func Test_canFetchSingleEventById(t *testing.T) {
	mt, err := time.LoadLocation("America/Boise")
	if err != nil {
		mt = time.UTC // fallback to UTC if timezone loading fails
	}

	db, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Could not create mock database: %s", err)
	}
	defer db.Close()

	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Could not create mock uuid: %s", err)
	}
	db.ExpectQuery(`SELECT \* FROM events WHERE id=\$1`).
		WithArgs(id.String()).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"name",
				"start_time",
				"end_time",
				"summary",
				"description",
				"event_location",
				"assembly_location",
				"pickup_location",
			}).AddRow(
				id.String(),
				"Troop Meeting",
				time.Date(2025, 1, 1, 19, 0, 0, 0, mt),
				time.Date(2025, 1, 1, 21, 0, 0, 0, mt),
				"summary",
				"description",
				"event_location",
				"assembly_location",
				"pickup_location",
			))

	dao := NewEventDao(db)
	event, err := dao.GetEventById(context.Background(), id.String())
	if err != nil {
		t.Fatalf("Error getting event by id: %s", err)
	}

	if event.Name != "Troop Meeting" {
		t.Errorf("Expecting Troop Meeting, got %s", event.Name)
	}
}

func Test_canTruncateTable(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Could not create mock database: %s", err)
	}
	defer db.Close()

	db.ExpectExec(`TRUNCATE TABLE events`).
		WillReturnResult(pgxmock.NewResult("TRUNCATE TABLE", 0))

	dao := NewEventDao(db)
	err = dao.Truncate(context.Background())
	if err != nil {
		t.Fatalf("Error truncating table: %s", err)
	}

	if err := db.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_canCountEvents(t *testing.T) {
	db, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Could not create mock database: %s", err)
	}
	defer db.Close()

	db.ExpectQuery(`SELECT COUNT\(\*\) FROM events`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))

	dao := NewEventDao(db)
	count, err := dao.Count(context.Background())
	if err != nil {
		t.Fatalf("Error counting events: %s", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}
