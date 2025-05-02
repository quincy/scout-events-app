package integration_tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/playwright-community/playwright-go"
	"os"
)

type playwrightKey struct{} // *playwright.Playwright
type browserKey struct{}    // playwright.Browser

// InitBrowser initializes a Playwright browser instance and stores it in the context.
func InitBrowser(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	var pw *playwright.Playwright
	var err error
	var browser playwright.Browser

	// default to headless mode unless overridden by environment
	headless, ok := os.LookupEnv("PLAYWRIGHT_HEADLESS")
	if !ok {
		headless = "true"
	}

	browser, ok = ctx.Value(browserKey{}).(playwright.Browser)
	if !ok {
		pw, err = getPlaywrightInstance(ctx)
		if err != nil {
			return nil, err
		}

		browser, err = pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: playwright.Bool(headless == "true")})
		if err != nil {
			return nil, fmt.Errorf("failed to launch browser: %v", err)
		}
	}

	return context.WithValue(ctx, browserKey{}, browser), nil
}

// CloseBrowser closes the Playwright browser instance stored in the context.
func CloseBrowser(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	browser, ok := ctx.Value(browserKey{}).(playwright.Browser)
	if ok {
		err := browser.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close browser: %v", err)
		}
	}

	return ctx, nil
}

// InitPlaywright initializes a Playwright instance and stores it in the context.
func InitPlaywright(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run playwright: %v", err)
	}

	return context.WithValue(ctx, playwrightKey{}, pw), nil
}

// ClosePlaywright closes the Playwright instance stored in the context.
func ClosePlaywright(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	pw, ok := ctx.Value(playwrightKey{}).(*playwright.Playwright)
	if ok {
		err := pw.Stop()
		if err != nil {
			return nil, fmt.Errorf("failed to stop playwright: %v", err)
		}
	}

	return ctx, nil
}

func getPlaywrightInstance(ctx context.Context) (*playwright.Playwright, error) {
	pw, ok := ctx.Value(playwrightKey{}).(*playwright.Playwright)
	if !ok {
		return nil, errors.New("no playwright instance exists")
	}

	return pw, nil
}

// BodyAsString retrieves the body content of the page as a string,
// or returns the error as the content if retrieving the content failed.
func BodyAsString(page playwright.Page) string {
	content, err := page.Content()
	if err != nil {
		return err.Error()
	}

	return content
}
