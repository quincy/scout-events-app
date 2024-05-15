package com.troop77eagle.events

import com.troop77eagle.getEventsResource
import com.troop77eagle.plugins.configureRouting
import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.shouldBe
import io.ktor.client.call.body
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation as ClientContentNegotiation
import io.ktor.client.request.accept
import io.ktor.client.request.get
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.http.ContentType.Application
import io.ktor.http.HttpStatusCode
import io.ktor.http.HttpStatusCode.Companion.NotFound
import io.ktor.http.HttpStatusCode.Companion.OK
import io.ktor.http.contentType
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.install
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation as ServerContentNegotiation
import io.ktor.server.testing.testApplication
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.UtcOffset
import kotlinx.datetime.toInstant
import kotlinx.serialization.json.Json

class EventsAPITest : AnnotationSpec() {
  private val jdbi = CockroachDbContainer.jdbi

  @Test
  fun `event can be fetched after it is created`() = testApplication {
    application {
      configureRouting(jdbi, getEventsResource(jdbi))
      install(ServerContentNegotiation) { json(Json { prettyPrint = true }) }
    }
    val client = createClient { install(ClientContentNegotiation) { json() } }

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

    client
        .get("/api/v1/events/1") { accept(Application.Json) }
        .apply {
          status shouldBe OK
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

  @Test
  fun `attempt to create incomplete event returns 400 Bad Request`() = testApplication {
    application {
      configureRouting(jdbi, getEventsResource(jdbi))
      install(ServerContentNegotiation) { json(Json { prettyPrint = true }) }
    }
    val client = createClient { install(ClientContentNegotiation) { json() } }

    client
        .put("/api/v1/events") {
          contentType(Application.Json)
          accept(Application.Json)
          setBody("""{ "name": "most required fields are missing" }""")
        }
        .apply { status shouldBe HttpStatusCode.BadRequest }
  }

  @Test
  fun `fetching non-existent event returns 404 Not Found`() = testApplication {
    application {
      configureRouting(jdbi, getEventsResource(jdbi))
      install(ServerContentNegotiation) { json(Json { prettyPrint = true }) }
    }
    val client = createClient { install(ClientContentNegotiation) { json() } }

    client.get("/api/v1/events/404") { accept(Application.Json) }.apply { status shouldBe NotFound }
  }
}
