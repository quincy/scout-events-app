package com.troop77eagle.events

import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.equals.shouldBeEqual
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import kotlinx.datetime.Instant

class EventsResourceTest : AnnotationSpec() {
  @Test
  fun `create event happy path`() {
    val eventsService =
        mockk<EventsService> {
          every { create(any()) } returns
              Event(
                  id = 1L,
                  name = "eventName",
                  startTime = Instant.parse("2024-04-30T16:30-08:00"),
                  endTime = Instant.parse("2024-05-01T14:00-08:00"),
                  summary = "eventSummary",
                  description = "eventDescription",
                  eventLocation = "eventLocation",
                  assemblyLocation = "assemblyLocation",
                  pickupLocation = "pickupLocation",
              )
        }

    val request =
        CreateEventRequest(
            name = "eventName",
            startTime = "2024-04-30T16:30-08:00",
            endTime = "2024-05-01T14:30-08:00",
            summary = "eventSummary",
            description = "eventDescription",
            eventLocation = "eventLocation",
            assemblyLocation = "assemblyLocation",
            pickupLocation = "pickupLocation",
        )

    val event = EventsResource(eventsService).create(request)

    event shouldBeEqual
        Event(
            id = 1L,
            name = "eventName",
            startTime = Instant.parse("2024-04-30T16:30-08:00"),
            endTime = Instant.parse("2024-05-01T14:00-08:00"),
            summary = "eventSummary",
            description = "eventDescription",
            eventLocation = "eventLocation",
            assemblyLocation = "assemblyLocation",
            pickupLocation = "pickupLocation",
        )
    verify { eventsService.create(request) }
  }
}
