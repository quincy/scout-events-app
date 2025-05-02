package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/quincy/scout-events-app/src/config"
	"github.com/quincy/scout-events-app/src/database"
	"github.com/quincy/scout-events-app/src/events"
	"strings"
	"time"
)

func ClearEvents() error {
	conn, err := CreateLocalDbConn()
	defer conn.Close()
	if err != nil {
		return err
	}
	dao := events.NewEventDao(config.TestConfig, conn)
	err = dao.Truncate(context.Background())
	if err != nil {
		return err
	}

	count, err := dao.Count(context.Background())
	if err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("expected no events, but found %d", count)
	}

	return nil
}

func TheFollowingEventsHaveBeenCreated(eventList *godog.DocString) error {
	var eventsToCreate []events.Event

	jsonStr := strings.TrimSpace(eventList.Content)
	for key, t := range Timestamps {
		jsonStr = strings.ReplaceAll(jsonStr, key, t.Format(time.RFC3339))
	}

	if err := json.Unmarshal([]byte(jsonStr), &eventsToCreate); err != nil {
		return fmt.Errorf("could not unmarshal events list: %v", err)
	}

	conn, err := CreateLocalDbConn()
	defer conn.Close()
	if err != nil {
		return err
	}
	dao := events.NewEventDao(config.TestConfig, conn)
	err = dao.CreateEvents(context.Background(), eventsToCreate)
	if err != nil {
		return err
	}

	return nil
}

func CreateLocalDbConn() (database.PgxClient, error) {
	return database.CreateDbConnection(context.Background(), config.TestConfig)
}
