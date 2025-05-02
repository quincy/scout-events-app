package integration_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/playwright-community/playwright-go"
	"net/http"
	"strings"
	"testing"
)

var navigationTimeoutMillis = 5_000.0

type pageKey struct{} // playwright.Page

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(InitPlaywright)
	ctx.Before(InitBrowser)

	ctx.After(ClosePlaywright)
	ctx.After(CloseBrowser)

	ctx.Given(`^the events app is running$`, AppIsRunning)
	ctx.Given(`^I am logged in$`, UserIsLoggedIn)
	ctx.Given(`^I am logged in with the "([^"]*)" role$`, UserIsLoggedInWithRole)
	ctx.Given(`^no events exist$`, ClearEvents)

	ctx.When(`^I visit the events list page$`, VisitEventsPage)
	ctx.When(`^I click the "([^"]*)" event row$`, ClickEventRow)
	ctx.When(`^I click the "([^"]*)" button$`, ClickButton)
	ctx.Step(`^I fill in the form with the following values:$`, FillInForm)

	ctx.Then(`^I should see "([^"]*)"$`, ValidatePageContainsText)
	ctx.Step(`^the following events have been created:$`, TheFollowingEventsHaveBeenCreated)
	ctx.Then(`^I should see the "([^"]*)" button which GETs "([^"]*)"$`, ValidatePageContainsButton)
	ctx.Step(`^I should see a form POSTing to "([^"]*)" with the following fields:$`, ValidateForm)
	ctx.Step(`^I should see the event details page with these values:$`, ValidateEventDetailsPage)
	ctx.Step(`^I submit the form$`, SubmitForm)
}

func Test_events_list_page(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             createTestOptions(t, []string{"features/events_list_page.feature"}),
	}

	code := suite.Run()

	if code != 0 {
		t.Fatalf("failed to run feature tests.  exit code: %d", code)
	}
}

func Test_add_event(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             createTestOptions(t, []string{"features/add_event_page.feature"}),
	}

	code := suite.Run()

	if code != 0 {
		t.Fatalf("failed to run feature tests.  exit code: %d", code)
	}
}

func createTestOptions(t *testing.T, tests []string) *godog.Options {
	return &godog.Options{
		Strict:    true,
		Format:    "pretty",
		Paths:     tests,
		TestingT:  t, // Testing instance that will run subtests
		Randomize: -1,
	}
}

var Buttons = map[string]string{
	"Add Event": "button-add-event",
}

func AppIsRunning(ctx context.Context) error {
	response, err := http.Get("http://localhost:8080/healthcheck")
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("app is not healthy. Got status code: %d", response.StatusCode)
	}

	return nil
}

func VisitEventsPage(ctx context.Context) (context.Context, error) {
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

func ValidatePageContainsText(ctx context.Context, expected string) (context.Context, error) {
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

func ValidatePageContainsButton(ctx context.Context, buttonText string, buttonLink string) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return ctx, errors.New("no HTML found in context")
	}

	buttonId, ok := Buttons[buttonText]
	if !ok {
		return nil, fmt.Errorf("buttonId '%s' not found in button map", buttonText)
	}

	buttonLocator := page.Locator("css=#" + buttonId)
	if buttonLocator == nil {
		return nil, fmt.Errorf("no button with id '%s' found", buttonId)
	}

	actualText, err := buttonLocator.InnerText()
	if err != nil {
		return nil, fmt.Errorf("failed to get button text: %v", err)
	}
	if actualText != buttonText {
		return nil, fmt.Errorf("expected button text '%s' but got '%s'", buttonText, actualText)
	}

	actualLink, err := buttonLocator.GetAttribute("hx-get")
	if err != nil {
		return nil, fmt.Errorf("failed to get button link: %v", err)
	}
	if actualLink != buttonLink {
		return nil, fmt.Errorf("expected button link '%s' but got '%s'", buttonLink, actualLink)
	}

	return ctx, nil
}

func ClickButton(ctx context.Context, buttonText string) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	buttonId, ok := Buttons[buttonText]
	if !ok {
		return nil, fmt.Errorf("buttonId '%s' not found in button map", buttonText)
	}

	buttonLocator := page.Locator("css=#" + buttonId)
	if buttonLocator == nil {
		return nil, fmt.Errorf("no button with id '%s' found", buttonId)
	}

	err := buttonLocator.Click()
	if err != nil {
		return nil, fmt.Errorf("failed to click on button: %v", err)
	}

	err = page.WaitForURL(
		"**/events:create",
		playwright.PageWaitForURLOptions{
			WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			Timeout:   &navigationTimeoutMillis,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for URL: %v", err)
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func ClickEventRow(ctx context.Context, eventName string) (context.Context, error) {
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
			Timeout:   &navigationTimeoutMillis,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for URL: %v", err)
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func ValidateForm(ctx context.Context, expectedUrl string, formFields *godog.Table) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	form := page.Locator("css=form")
	actualUrl, err := form.GetAttribute("hx-post")
	if err != nil {
		return nil, fmt.Errorf("failed to get form action: %v", err)
	}
	if actualUrl != expectedUrl {
		return nil, fmt.Errorf("expected form action '%s' but got '%s'", expectedUrl, actualUrl)
	}

	// Skip header row (index 0) and iterate through data rows
	for _, row := range formFields.Rows[1:] {
		fieldName := row.Cells[0].Value
		expectedType := row.Cells[1].Value

		// Convert field name to ID format (lowercase, hyphens)
		fieldId := strings.ToLower(strings.ReplaceAll(fieldName, " ", "_"))

		input := form.Locator("css=#" + fieldId)
		actualType, err := input.GetAttribute("type")
		if err != nil && expectedType != "textarea" {
			return nil, fmt.Errorf("failed to get input type for field '%s': %v", fieldName, err)
		}

		if expectedType == "textarea" {
			exists, err := page.Locator("css=textarea#" + fieldId).Count()
			if err != nil || exists == 0 {
				return nil, fmt.Errorf("textarea '%s' not found", fieldName)
			}
		} else if actualType != expectedType {
			return nil, fmt.Errorf("field '%s' expected type '%s' but got '%s'", fieldName, expectedType, actualType)
		}
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func FillInForm(ctx context.Context, formData *godog.Table) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	// Skip header row (index 0) and iterate through data rows
	for _, row := range formData.Rows[1:] {
		fieldName := row.Cells[0].Value
		value := row.Cells[1].Value

		// Convert field name to ID format (lowercase, underscore)
		fieldId := strings.ToLower(strings.ReplaceAll(fieldName, " ", "_"))

		// Handle different input types
		element := page.Locator("css=#" + fieldId)
		elementType, err := element.GetAttribute("type")
		if err != nil {
			return nil, fmt.Errorf("failed to get element type for %s: %v", fieldName, err)
		}

		// Fill the field
		if elementType == "datetime-local" {
			// Parse relative date expressions
			parsedTime, ok := Timestamps[value]
			if !ok {
				return nil, fmt.Errorf("failed to parse time for %s: %v", fieldName, err)
			}
			value = parsedTime.Format("2006-01-02T15:04")
		}

		err = element.Fill(value)
		if err != nil {
			return nil, fmt.Errorf("failed to fill field %s: %v", fieldName, err)
		}
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func ValidateEventDetailsPage(ctx context.Context, expectedJson *godog.DocString) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	var expected map[string]string
	if err := json.Unmarshal([]byte(expectedJson.Content), &expected); err != nil {
		return nil, fmt.Errorf("failed to parse expected JSON: %v", err)
	}

	// Validate each field
	for field, expectedValue := range expected {
		// Convert JSON field name to HTML ID format

		// Get the value from the page
		element := page.Locator(fmt.Sprintf("css=[data-field='%s']", field))
		actualValue, err := element.InnerText()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for field '%s': %v", field, err)
		}

		// For datetime fields, handle timestamp replacements
		if field == "start-and-end-times" {
			times := strings.Split(expectedValue, " - ")
			startTime, ok := Timestamps[times[0]]
			if !ok {
				return nil, fmt.Errorf("timestamp '%s' not found in timestamps map", expectedValue)
			}
			endTime, ok := Timestamps[times[1]]
			if !ok {
				return nil, fmt.Errorf("timestamp '%s' not found in timestamps map", expectedValue)
			}
			expectedValue = fmt.Sprintf("%s - %s", startTime.Format("2006-01-02 3:04 PM"), endTime.Format("2006-01-02 3:04 PM"))
		}

		if strings.TrimSpace(actualValue) != strings.TrimSpace(expectedValue) {
			return nil, fmt.Errorf("field '%s': expected '%s' but got '%s'", field, expectedValue, actualValue)
		}
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}

func SubmitForm(ctx context.Context) (context.Context, error) {
	page, ok := ctx.Value(pageKey{}).(playwright.Page)
	if !ok {
		return nil, errors.New("no page found in context")
	}

	// Click the submit button
	err := page.Locator("css=button[type=submit]").Click()
	if err != nil {
		return nil, fmt.Errorf("failed to click submit button: %v", err)
	}

	// Wait for navigation to complete (event details page)
	//err = page.WaitForURL(
	//	"/events/*",
	//	playwright.PageWaitForURLOptions{
	//		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	//		Timeout:   &navigationTimeoutMillis,
	//	},
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to wait for navigation after form submission: %v", err)
	//}

	// Wait for HTMX to update the URL
	_, err = page.WaitForFunction(`() => window.location.pathname.startsWith('/events/')`, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for URL update: %v", err)
	}

	return context.WithValue(ctx, pageKey{}, page), nil
}
