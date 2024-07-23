package com.troop77eagle.events

import io.ktor.http.ContentType
import io.ktor.http.ContentType.Application
import io.ktor.http.ContentType.Text
import io.ktor.http.HttpStatusCode.Companion.Created
import io.ktor.server.application.call
import io.ktor.server.plugins.BadRequestException
import io.ktor.server.plugins.NotFoundException
import io.ktor.server.request.ApplicationRequest
import io.ktor.server.request.accept
import io.ktor.server.request.receiveText
import io.ktor.server.response.respond
import io.ktor.server.response.respondText
import io.ktor.server.routing.Route
import io.ktor.server.routing.get
import io.ktor.server.routing.put
import io.ktor.server.routing.route
import kotlinx.serialization.json.Json

fun Route.eventsRoute(eventsResource: EventsResource) {
  route("/events") {
    put {
      val request = Json.decodeFromString<CreateEventRequest>(call.receiveText())
      val event = eventsResource.create(request)
      call.respond(Created, event)
    }
    get {
      val contentType = requestedContentType(call.request)
      val events = eventsResource.fetchAll()

      if (contentType == Text.Html) {
        call.respondText(contentType = contentType) { Templator.create().toHtmlTable(events) }
      } else {
        call.respond(eventsResource.fetchAll())
      }
    }

    route("/{id}") {
      get {
        val id = call.parameters["id"]?.toLong() ?: throw BadRequestException("id cannot be null")
        call.respond(eventsResource.fetch(id))
      }
    }
  }
}

fun requestedContentType(request: ApplicationRequest): ContentType {
  val isHxRequest = request.headers["HX-Request"]?.let { v -> v.toBoolean() } ?: false
  return if (isHxRequest) {
    Text.Html
  } else {
    val acceptHeader = request.accept() ?: Application.Json.toString()
    acceptHeader
      .split(Regex(","))
      .first { it in allowedTypesForFetchAll }
      .let { ContentType.parse(it) }
  }
}

private val allowedTypesForFetchAll = setOf(Application.Json, Text.Html).map { it.toString() }

class EventsResource(private val eventsService: EventsService) {
  fun create(request: CreateEventRequest): Event {
    return eventsService.create(request)
  }

  fun fetch(id: Long): Event {
    return eventsService.getById(id) ?: throw NotFoundException("No such event with id=$id")
  }

  fun fetchAll(): List<Event> {
    return eventsService.getAll()
  }
}
