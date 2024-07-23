package com.troop77eagle.events

import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.collections.shouldHaveSize
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

  @Test
  fun `fetch all events`() {
    val eventsDAO =
        mockk<EventsDAO> {
          every { getAll() } returns
              listOf(
                  Event(
                      id = 1L,
                      name = "event1",
                      startTime = Instant.parse("1111-11-11T01:11-01:00"),
                      endTime = Instant.parse("1111-11-11T11:11-01:00"),
                      summary = "summary1",
                      description = "description1",
                      eventLocation = "location1",
                      assemblyLocation = "assemblyLocation1",
                      pickupLocation = "pickupLocation1",
                  ),
                  Event(
                      id = 2L,
                      name = "event2",
                      startTime = Instant.parse("2222-02-22T02:22-02:00"),
                      endTime = Instant.parse("2222-02-22T22:22-02:00"),
                      summary = "summary2",
                      description = "description2",
                      eventLocation = "location2",
                      assemblyLocation = "assemblyLocation2",
                      pickupLocation = "pickupLocation2",
                  ),
              )
        }

    val events = EventsService.create(eventsDAO).getAll()

    events shouldHaveSize 2
    verify { eventsDAO.getAll() }
  }
}
