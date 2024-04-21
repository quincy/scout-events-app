package com.troop77eagle.events

import com.troop77eagle.getEventsResource
import com.troop77eagle.plugins.configureRouting
import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.shouldBe
import io.ktor.client.call.body
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.accept
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.http.ContentType.Application
import io.ktor.http.HttpStatusCode
import io.ktor.http.contentType
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.testing.testApplication
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.UtcOffset
import kotlinx.datetime.toInstant

class EventsAPITest : AnnotationSpec() {
  private val jdbi = CockroachDbContainer.jdbi

  @Test
  fun `event can be fetched after it is created`() = testApplication {
    application { configureRouting(jdbi, getEventsResource(jdbi)) }
    val client = createClient { install(ContentNegotiation) { json() } }

    client
        .put("/api/v1/events") {
          contentType(Application.Json)
          accept(Application.Json)
          setBody<CreateEventRequest>(
              CreateEventRequest(
                  name = "Wagon Wheel Camporee",
                  startTime = "2024-04-12T16:30-08:00",
                  endTime = "2024-04-14T15:00-08:00",
                  summary = "First annual Wagon Wheel Camporee",
                  description = "Such camporee",
                  eventLocation = "Vale, OR",
                  assemblyLocation = "Eagle Hills Church",
                  pickupLocation = "Eagle Hills Church",
              ),
          )
        }
        .apply {
          status shouldBe HttpStatusCode.Created
          body<Event>() shouldBe
              Event(
                  id = 1L,
                  name = "Wagon Wheel Camporee",
                  startTime = LocalDateTime(2024, 4, 12, 16, 30).toInstant(UtcOffset(-8)),
                  endTime = LocalDateTime(2024, 4, 14, 15, 0).toInstant(UtcOffset(-8)),
                  summary = "First annual Wagon Wheel Camporee",
                  description = "Such camporee",
                  eventLocation = "Vale, OR",
                  assemblyLocation = "Eagle Hills Church",
                  pickupLocation = "Eagle Hills Church",
              )
        }
  }
}
