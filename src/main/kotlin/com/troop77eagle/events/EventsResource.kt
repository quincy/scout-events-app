package com.troop77eagle.events

import io.ktor.http.HttpStatusCode.Companion.Created
import io.ktor.server.application.call
import io.ktor.server.plugins.BadRequestException
import io.ktor.server.plugins.NotFoundException
import io.ktor.server.request.receiveText
import io.ktor.server.response.respond
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

    route("/{id}") {
      get {
        val id = call.parameters["id"]?.toLong() ?: throw BadRequestException("id cannot be null")
        call.respond(eventsResource.fetch(id))
      }
    }
  }
}

class EventsResource(private val eventsService: EventsService) {
  fun create(request: CreateEventRequest): Event {
    return eventsService.create(request)
  }

  fun fetch(id: Long): Event {
    return eventsService.getById(id) ?: throw NotFoundException("No such event with id=$id")
  }
}
