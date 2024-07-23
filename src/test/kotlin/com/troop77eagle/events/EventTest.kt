package com.troop77eagle.events

import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import kotlinx.datetime.Instant
import org.junit.jupiter.api.Test

class EventTest {
  private val event =
      Event(
          id = 1L,
          name = "event1",
          startTime = Instant.parse("1111-11-11T01:11-01:00"),
          endTime = Instant.parse("2222-02-22T11:11-01:00"),
          summary = "summary1",
          description = "description1",
          eventLocation = "location1",
          assemblyLocation = "assemblyLocation1",
          pickupLocation = "pickupLocation1",
      )

  @Test
  fun `can get startDate as ISO8601 string`() {
    assertThat(event.startDate, equalTo("1111-11-11"))
  }

  @Test
  fun `can get endDate as ISO8601 string`() {
    assertThat(event.endDate, equalTo("2222-02-22"))
  }
}
