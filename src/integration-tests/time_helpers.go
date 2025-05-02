package integration_tests

import "time"

func TomorrowAt(hour int) time.Time {
	timezone, err := time.LoadLocation("America/Boise")
	if err != nil {
		panic("Could not load timezone: " + err.Error())
	}

	now := time.Now().In(timezone)
	date := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, timezone)

	return date.AddDate(0, 0, 1)
}

var Timestamps = map[string]time.Time{
	"tomorrow at 7pm": TomorrowAt(19),
	"tomorrow at 9pm": TomorrowAt(21),
}
