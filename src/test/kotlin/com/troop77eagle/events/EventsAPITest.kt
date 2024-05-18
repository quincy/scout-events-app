package com.troop77eagle.events

import com.troop77eagle.getEventsResource
import com.troop77eagle.plugins.jsonConfig
import com.troop77eagle.plugins.statusPagesConfig
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
import io.ktor.server.config.MapApplicationConfig
import io.ktor.server.config.mergeWith
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation as ServerContentNegotiation
import io.ktor.server.plugins.statuspages.StatusPages
import io.ktor.server.routing.route
import io.ktor.server.testing.testApplication
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.UtcOffset
import kotlinx.datetime.toInstant
import kotlinx.serialization.ExperimentalSerializationApi

class EventsAPITest : AnnotationSpec() {
  private val jdbi = CockroachDbContainer.jdbi

  private fun testConfig() =
      MapApplicationConfig(
          "db.host" to CockroachDbContainer.dbContainer.host,
          "db.port" to CockroachDbContainer.dbContainer.firstMappedPort.toString(),
          "db.username" to CockroachDbContainer.dbContainer.username,
          "db.password" to CockroachDbContainer.dbContainer.password,
          "db.liquibase.username" to CockroachDbContainer.dbContainer.password,
          "db.liquibase.password" to CockroachDbContainer.dbContainer.password,
      )

  @Test
  fun `event can be fetched after it is created`() = testApplication {
    environment { config = config.mergeWith(testConfig()) }
    install(ServerContentNegotiation) { json(jsonConfig()) }
    routing { route("/api/v1") { eventsRoute(getEventsResource(jdbi)) } }
    val client = createClient { install(ClientContentNegotiation) { json(jsonConfig()) } }

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

  @ExperimentalSerializationApi
  @Test
  fun `attempt to create incomplete event returns 400 Bad Request`() = testApplication {
    environment { config = config.mergeWith(testConfig()) }
    install(ServerContentNegotiation) { json(jsonConfig()) }
    install(StatusPages) { statusPagesConfig() }
    routing { route("/api/v1") { eventsRoute(getEventsResource(jdbi)) } }
    val client = createClient { install(ClientContentNegotiation) { json(jsonConfig()) } }

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
    environment { config = config.mergeWith(testConfig()) }
    install(ServerContentNegotiation) { json(jsonConfig()) }
    routing { route("/api/v1") { eventsRoute(getEventsResource(jdbi)) } }
    val client = createClient { install(ClientContentNegotiation) { json(jsonConfig()) } }

    client.get("/api/v1/events/404") { accept(Application.Json) }.apply { status shouldBe NotFound }
  }
}
