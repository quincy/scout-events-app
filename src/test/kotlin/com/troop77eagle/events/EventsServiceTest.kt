package com.troop77eagle.events

import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.equals.shouldBeEqual
import io.kotest.matchers.shouldBe
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import kotlinx.datetime.Instant

class EventsServiceTest : AnnotationSpec() {
  @Test
  fun `fetch existing event`() {
    val eventsDAO =
        mockk<EventsDAO> {
          every { getById(any()) } returns
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

    val event = EventsService.create(eventsDAO).getById(1L)!!

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

    verify { eventsDAO.getById(1L) }
  }

  @Test
  fun `create event happy path`() {
    val eventsDAO = mockk<EventsDAO> { every { create(any()) } returns 1L }

    val request =
        CreateEventRequest(
            name = "eventName",
            startTime = "2024-04-30T16:30-08:00",
            endTime = "2024-05-01T14:00-08:00",
            summary = "eventSummary",
            description = "eventDescription",
            eventLocation = "eventLocation",
            assemblyLocation = "assemblyLocation",
            pickupLocation = "pickupLocation",
        )

    val event = EventsService.create(eventsDAO).create(request)

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

    verify { eventsDAO.create(request) }
  }

  @Test
  fun `fetch non-existent event returns null`() {
    val eventsDAO = mockk<EventsDAO> { every { getById(any()) } returns null }

    val event = EventsService.create(eventsDAO).getById(404L)

    event shouldBe null
  }
}
