package com.troop77eagle.events

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import kotlinx.datetime.Instant
import org.junit.jupiter.api.Test

class TemplatorTest {
  @Test
  fun `templating empty event list into HTML table returns only thead`() {
    assertThat(
        Templator.create().toHtmlTable(listOf()).trim(),
        equalTo(
            """
            <thead>Upcoming Events</thead>
            <tr>
              <td>Date</td>
              <td>Event</td>
              <td>Summary</td>
              <td>Location</td>
            </tr>
            """
                .trimIndent()))
  }

  @Test
  fun `templating event list into HTML table happy path`() {
    val events =
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

    assertThat(
        Templator.create().toHtmlTable(events).trim(),
        equalTo(
            """
            <thead>Upcoming Events</thead>
            <tr>
              <td>Date</td>
              <td>Event</td>
              <td>Summary</td>
              <td>Location</td>
            </tr>
            <tr>
              <td>1111-11-11</td>
              <td>event1</td>
              <td>summary1</td>
              <td>location1</td>
            </tr>
            <tr>
              <td>2222-02-22</td>
              <td>event2</td>
              <td>summary2</td>
              <td>location2</td>
            </tr>
            """
                .trimIndent()))
  }
}
