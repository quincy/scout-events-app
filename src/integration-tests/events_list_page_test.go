package integration_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/playwright-community/playwright-go"
	"github.com/quincy/scout-events-app/database"
	"github.com/quincy/scout-events-app/events"
	"net/http"
	"strings"
	"testing"
	"time"
)

var navigationTimeout = 5_000.0

type pageKey struct{} // playwright.Page

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(InitPlaywright)
	ctx.Before(InitBrowser)

	ctx.After(ClosePlaywright)
	ctx.After(CloseBrowser)

	ctx.Given(`^the events app is running$`, appIsRunning)
	ctx.Given(`^I am logged in$`, userIsLoggedIn)

	ctx.Given(`^no events exist$`, clearEvents)
	ctx.When(`^I visit the events list page$`, visitEventsPage)

	ctx.When(`^I click on the "([^"]*)" event row$`, clickEventRow)
	ctx.Then(`^I should see "([^"]*)"$`, validatePageContainsText)
	ctx.Step(`^the following events have been created:$`, theFollowingEventsHaveBeenCreated)
}

func TestEvents_list_page(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests
		},
	}

	code := suite.Run()

	if code != 0 {
		t.Fatalf("failed to run feature tests.  exit code: %d", code)
	}
}

func appIsRunning(ctx context.Context) error {
	response, err := http.Get("http://localhost:8080/healthcheck")
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("app is not healthy. Got status code: %d", response.StatusCode)
	}

	return nil
}

func userIsLoggedIn() error {
	return nil // FIXME need to implement login
}

func clearEvents() error {
	conn, err := createLocalDbConn()
	defer conn.Close()
	if err != nil {
		return err
	}
	dao := events.NewEventDao(conn)
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

func tomorrowAt(hour int) time.Time {
	mt, err := time.LoadLocation("America/Boise")
	if err != nil {
		mt = time.UTC // fallback to UTC if timezone loading fails
	}

	now := time.Now().In(mt)
	date := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, mt)

	return date.AddDate(0, 0, 1)
}

var timestamps = map[string]time.Time{
	"tomorrow at 7pm": tomorrowAt(19),
	"tomorrow at 9pm": tomorrowAt(21),
}

func theFollowingEventsHaveBeenCreated(eventList *godog.DocString) error {
	var eventsToCreate []events.Event

	jsonStr := strings.TrimSpace(eventList.Content)
	for key, t := range timestamps {
		jsonStr = strings.ReplaceAll(jsonStr, key, t.Format(time.RFC3339))
	}

	if err := json.Unmarshal([]byte(jsonStr), &eventsToCreate); err != nil {
		return fmt.Errorf("could not unmarshal events list: %v", err)
	}

	conn, err := createLocalDbConn()
	defer conn.Close()
	if err != nil {
		return err
	}
	dao := events.NewEventDao(conn)
	err = dao.CreateEvents(context.Background(), eventsToCreate)
	if err != nil {
		return err
	}

	return nil
}

func createLocalDbConn() (database.PgxClient, error) {
	return database.CreateDbConnection(
		context.Background(),
		&database.DbConfig{
			Username: "admin",
			Password: "admin",
			Hostname: "localhost",
			Port:     26257,
			Dbname:   "scouting",
		})
}

func visitEventsPage(ctx context.Context) (context.Context, error) {
	browser, ok := ctx.Value(browserKey{}).(playwright.Browser)
	if !ok {
		return nil, errors.New("no browser found in context")
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create new page: %v", err)
	}
	response, err := page.Goto("http://localhost:8080/")
	if err != nil {
		return nil, fmt.Errorf("failed to visit page: %v", err)
	}

	if !response.Ok() {
		body, _ := response.Body()
		return nil, fmt.Errorf("failed to visit page. status='%v %v' body=%v", response.Status(), response.StatusText(), body)
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func validatePageContainsText(ctx context.Context, expected string) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return ctx, errors.New("no HTML found in context")
	}

	content, err := page.Content()
	if err != nil {
		return nil, err
	}

	if !strings.Contains(content, expected) {
		return nil, fmt.Errorf("expected to see '%s' but got: %s", expected, BodyAsString(page))
	}

	return ctx, nil
}

func clickEventRow(ctx context.Context, eventName string) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	err := page.GetByRole("row", playwright.PageGetByRoleOptions{Name: eventName}).Click()
	if err != nil {
		return nil, fmt.Errorf("failed to click on event row: %v", err)
	}

	err = page.WaitForURL(
		"**/events/*",
		playwright.PageWaitForURLOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   &navigationTimeout,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for URL: %v", err)
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}
