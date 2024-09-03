package com.troop77eagle.events

import com.ninja_squad.dbsetup.DbSetup
import com.ninja_squad.dbsetup.Operations
import com.ninja_squad.dbsetup.destination.DataSourceDestination
import com.troop77eagle.events.EventMatchers.haveAssemblyLocation
import com.troop77eagle.events.EventMatchers.haveDescription
import com.troop77eagle.events.EventMatchers.haveEndTime
import com.troop77eagle.events.EventMatchers.haveEventLocation
import com.troop77eagle.events.EventMatchers.haveName
import com.troop77eagle.events.EventMatchers.havePickupLocation
import com.troop77eagle.events.EventMatchers.haveStartTime
import com.troop77eagle.events.EventMatchers.haveSummary
import com.troop77eagle.getEventsResource
import com.troop77eagle.plugins.jsonConfig
import com.troop77eagle.plugins.statusPagesConfig
import io.kotest.core.spec.style.AnnotationSpec
import io.kotest.matchers.Matcher
import io.kotest.matchers.MatcherResult
import io.kotest.matchers.should
import io.kotest.matchers.shouldBe
import io.ktor.client.call.body
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation as ClientContentNegotiation
import io.ktor.client.request.accept
import io.ktor.client.request.get
import io.ktor.client.request.put
import io.ktor.client.request.setBody
import io.ktor.client.statement.bodyAsText
import io.ktor.http.ContentType.Application
import io.ktor.http.ContentType.Text
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
import kotlinx.datetime.Instant
import kotlinx.datetime.LocalDateTime
import kotlinx.datetime.UtcOffset
import kotlinx.datetime.toInstant
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.json.Json

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

  @BeforeEach
  fun prepare() {
    val operation = Operations.sequenceOf(CockroachDbContainer.DELETE_ALL)
    val dbSetup = DbSetup(DataSourceDestination(CockroachDbContainer.datasource), operation)
    dbSetup.launch()
  }

  @AfterClass
  fun teardown() {
    CockroachDbContainer.dbContainer.stop()
  }

  @Test
  fun `event can be fetched after it is created`() = testApplication {
    environment { config = config.mergeWith(testConfig()) }
    install(ServerContentNegotiation) { json(jsonConfig()) }
    routing { route("/api/v1") { eventsRoute(getEventsResource(jdbi)) } }
    val client = createClient { install(ClientContentNegotiation) { json(jsonConfig()) } }
    var id: Long

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
              ))
        }
        .apply {
          id = Json.decodeFromString<Event>(bodyAsText()).id

          status shouldBe HttpStatusCode.Created
          body<Event>() should haveName("Wagon Wheel Camporee")
          body<Event>() should
              haveStartTime(LocalDateTime(2024, 4, 12, 16, 30).toInstant(UtcOffset(-8)))
          body<Event>() should
              haveEndTime(LocalDateTime(2024, 4, 14, 15, 0).toInstant(UtcOffset(-8)))
          body<Event>() should haveSummary("First annual Wagon Wheel Camporee")
          body<Event>() should haveDescription("Such camporee")
          body<Event>() should haveEventLocation("Vale, OR")
          body<Event>() should haveAssemblyLocation("Eagle Hills Church")
          body<Event>() should havePickupLocation("Eagle Hills Church")
        }

    client
        .get("/api/v1/events/$id") { accept(Application.Json) }
        .apply {
          status shouldBe OK
          body<Event>() should haveName("Wagon Wheel Camporee")
          body<Event>() should
              haveStartTime(LocalDateTime(2024, 4, 12, 16, 30).toInstant(UtcOffset(-8)))
          body<Event>() should
              haveEndTime(LocalDateTime(2024, 4, 14, 15, 0).toInstant(UtcOffset(-8)))
          body<Event>() should haveSummary("First annual Wagon Wheel Camporee")
          body<Event>() should haveDescription("Such camporee")
          body<Event>() should haveEventLocation("Vale, OR")
          body<Event>() should haveAssemblyLocation("Eagle Hills Church")
          body<Event>() should havePickupLocation("Eagle Hills Church")
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

  @Test
  fun `can fetch all events`() = testApplication {
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
              ))
        }
        .apply { status shouldBe HttpStatusCode.Created }
    client
        .put("/api/v1/events") {
          contentType(Application.Json)
          accept(Application.Json)
          setBody<CreateEventRequest>(
              CreateEventRequest(
                  name = "City of Rocks Campout",
                  startTime = "2024-06-10T16:30-08:00",
                  endTime = "2024-06-12T15:00-08:00",
                  summary = "Rock climbing campout",
                  description = "We'll do some climbing in City of Rocks Nat'l Monument",
                  eventLocation = "City of Rocks National Monument",
                  assemblyLocation = "Eagle Hills Church",
                  pickupLocation = "Eagle Hills Church",
              ))
        }
        .apply { status shouldBe HttpStatusCode.Created }

    client
        .get("/api/v1/events") { accept(Application.Json) }
        .apply {
          status shouldBe OK
          val events = body<List<Event>>()
          events[0].let {
            it should haveName("Wagon Wheel Camporee")
            it should haveStartTime(LocalDateTime(2024, 4, 12, 16, 30).toInstant(UtcOffset(-8)))
            it should haveEndTime(LocalDateTime(2024, 4, 14, 15, 0).toInstant(UtcOffset(-8)))
            it should haveSummary("First annual Wagon Wheel Camporee")
            it should haveDescription("Such camporee")
            it should haveEventLocation("Vale, OR")
            it should haveAssemblyLocation("Eagle Hills Church")
            it should havePickupLocation("Eagle Hills Church")
          }
          events[1].let {
            it should haveName("City of Rocks Campout")
            it should haveStartTime(LocalDateTime(2024, 6, 10, 16, 30).toInstant(UtcOffset(-8)))
            it should haveEndTime(LocalDateTime(2024, 6, 12, 15, 0).toInstant(UtcOffset(-8)))
            it should haveSummary("Rock climbing campout")
            it should haveDescription("We'll do some climbing in City of Rocks Nat'l Monument")
            it should haveEventLocation("City of Rocks National Monument")
            it should haveAssemblyLocation("Eagle Hills Church")
            it should havePickupLocation("Eagle Hills Church")
          }
        }
  }

  @Test
  fun `can fetch upcoming events list as html`() = testApplication {
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
              ))
        }
        .apply { status shouldBe HttpStatusCode.Created }

    client
        .put("/api/v1/events") {
          contentType(Application.Json)
          accept(Application.Json)
          setBody<CreateEventRequest>(
              CreateEventRequest(
                  name = "City of Rocks Campout",
                  startTime = "2024-06-10T16:30-08:00",
                  endTime = "2024-06-12T15:00-08:00",
                  summary = "Rock climbing campout",
                  description = "We'll do some climbing in City of Rocks Nat'l Monument",
                  eventLocation = "City of Rocks National Monument",
                  assemblyLocation = "Eagle Hills Church",
                  pickupLocation = "Eagle Hills Church",
              ))
        }
        .apply { status shouldBe HttpStatusCode.Created }

    client
        .get("/api/v1/events") { accept(Text.Html) }
        .apply {
          status shouldBe OK
          bodyAsText().trim() shouldBe
              """
              <thead>Upcoming Events</thead>
              <tr>
                <td>Date</td>
                <td>Event</td>
                <td>Summary</td>
                <td>Location</td>
              </tr>
              <tr>
                <td>2024-04-12</td>
                <td>Wagon Wheel Camporee</td>
                <td>First annual Wagon Wheel Camporee</td>
                <td>Vale, OR</td>
              </tr>
              <tr>
                <td>2024-06-10</td>
                <td>City of Rocks Campout</td>
                <td>Rock climbing campout</td>
                <td>City of Rocks National Monument</td>
              </tr>
              """
                  .trimIndent()
        }
  }
}

object EventMatchers {
  fun haveName(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.name == expected,
        { "event had name '${value.name}' but we expected name '$expected'" },
        { "event should not have name '$expected'" },
    )
  }

  fun haveStartTime(expected: Instant): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.startTime == expected,
        { "event had startTime '${value.startTime}' but we expected startTime '$expected'" },
        { "event should not have startTime '$expected'" },
    )
  }

  fun haveEndTime(expected: Instant): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.endTime == expected,
        { "event had endTime '${value.endTime}' but we expected endTime '$expected'" },
        { "event should not have endTime '$expected'" },
    )
  }

  fun haveSummary(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.summary == expected,
        { "event had summary '${value.summary}' but we expected summary '$expected'" },
        { "event should not have summary '$expected'" },
    )
  }

  fun haveDescription(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.description == expected,
        { "event had description '${value.description}' but we expected description '$expected'" },
        { "event should not have description '$expected'" },
    )
  }

  fun haveEventLocation(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.eventLocation == expected,
        {
          "event had eventLocation '${value.eventLocation}' but we expected eventLocation '$expected'"
        },
        { "event should not have eventLocation '$expected'" },
    )
  }

  fun haveAssemblyLocation(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.assemblyLocation == expected,
        {
          "event had assemblyLocation '${value.assemblyLocation}' but we expected assemblyLocation '$expected'"
        },
        { "event should not have assemblyLocation '$expected'" },
    )
  }

  fun havePickupLocation(expected: String): Matcher<Event> = Matcher { value ->
    MatcherResult(
        value.pickupLocation == expected,
        {
          "event had pickupLocation '${value.pickupLocation}' but we expected pickupLocation '$expected'"
        },
        { "event should not have pickupLocation '$expected'" },
    )
  }
}
