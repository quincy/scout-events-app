package events

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Summary          string    `json:"summary"`
	Description      string    `json:"description"`
	EventLocation    string    `json:"event_location"`
	AssemblyLocation string    `json:"assembly_location"`
	PickupLocation   string    `json:"pickup_location"`
}

func NewEvent(
	name string,
	startTime time.Time,
	endTime time.Time,
	summary string,
	description string,
	eventLocation string,
	assemblyLocation string,
	pickupLocation string,
) *Event {
	id := uuid.New()
	return &Event{
		Id:               id.String(),
		Name:             name,
		StartTime:        startTime,
		EndTime:          endTime,
		Summary:          summary,
		Description:      description,
		EventLocation:    eventLocation,
		AssemblyLocation: assemblyLocation,
		PickupLocation:   pickupLocation,
	}
}
