package events

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/quincy/scout-events-app/src/database"
	"log"
)

type EventDao interface {
	GetEventById(ctx context.Context, id string) (*Event, error)
	CreateEvents(ctx context.Context, entities []Event) error
	GetUpcomingEvents(ctx context.Context) ([]Event, error)
	Truncate(ctx context.Context) error
	Count(ctx context.Context) (int, error)
}

type dao struct {
	conn database.PgxClient
}

func NewEventDao(conn database.PgxClient) EventDao {
	return &dao{conn: conn}
}

// GetEventById retrieves an event by its ID if it exists.
func (d *dao) GetEventById(ctx context.Context, id string) (*Event, error) {
	rows, err := d.conn.Query(ctx, "SELECT * FROM events WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var event Event
	if rows.Next() {
		err = rows.Scan(
			&event.Id,
			&event.Name,
			&event.StartTime,
			&event.EndTime,
			&event.Summary,
			&event.Description,
			&event.EventLocation,
			&event.AssemblyLocation,
			&event.PickupLocation,
		)
		if err != nil {
			return nil, err
		}
	}

	return &event, nil
}

// CreateEvents inserts a batch of events into the database.
func (d *dao) CreateEvents(ctx context.Context, entities []Event) error {
	batch := &pgx.Batch{}

	for _, event := range entities {
		batch.Queue(`
	            INSERT INTO events (
	                id,
	                name,
	                start_time,
	                end_time,
	                summary,
	                description,
	                event_location,
	                assembly_location,
	                pickup_location
	            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			event.Id,
			event.Name,
			event.StartTime,
			event.EndTime,
			event.Summary,
			event.Description,
			event.EventLocation,
			event.AssemblyLocation,
			event.PickupLocation,
		)
	}

	results := d.conn.SendBatch(ctx, batch)
	var execErr error
	defer func() {
		if closeErr := results.Close(); closeErr != nil {
			if execErr != nil {
				execErr = fmt.Errorf("exec error: %v, close error: %v", execErr, closeErr)
			} else {
				execErr = closeErr
			}
		}
	}()

	_, execErr = results.Exec()
	return execErr
}

// GetUpcomingEvents retrieves all upcoming events.
// An event is upcoming if its end time is 1 week prior to today or any point in the future.
func (d *dao) GetUpcomingEvents(ctx context.Context) ([]Event, error) {
	rows, err := d.conn.Query(ctx, `SELECT * FROM events WHERE end_time > NOW() - INTERVAL '1 week'`)
	if err != nil {
		return nil, fmt.Errorf("failed to query upcoming events: %v", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(
			&event.Id,
			&event.Name,
			&event.StartTime,
			&event.EndTime,
			&event.Summary,
			&event.Description,
			&event.EventLocation,
			&event.AssemblyLocation,
			&event.PickupLocation,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event row: %v", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %v", err)
	}

	return events, nil
}

// Truncate clears the events table.  Only meant to be used in tests.
func (d *dao) Truncate(ctx context.Context) error {
	_, err := d.conn.Exec(ctx, `TRUNCATE TABLE events`)
	if err != nil {
		log.Printf("failed to truncate events table: %v", err)
		return fmt.Errorf("failed to truncate events table: %v", err)
	}

	log.Println("truncated events table")
	return nil
}

// Count returns the number of events in the database.
func (d *dao) Count(ctx context.Context) (int, error) {
	var count int
	err := d.conn.QueryRow(ctx, `SELECT COUNT(*) FROM events`).Scan(&count)
	if err != nil {
		log.Printf("failed to count events: %v", err)
		return 0, fmt.Errorf("failed to count events: %v", err)
	}

	log.Printf("counted %d events", count)
	return count, nil
}
