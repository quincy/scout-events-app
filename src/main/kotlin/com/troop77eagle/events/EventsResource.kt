package com.troop77eagle.events

import io.ktor.http.ContentType.Application
import io.ktor.http.HttpStatusCode.Companion.Created
import io.ktor.server.application.call
import io.ktor.server.request.receiveText
import io.ktor.server.response.respondText
import io.ktor.server.routing.Route
import io.ktor.server.routing.get
import io.ktor.server.routing.put
import io.ktor.server.routing.route
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json

fun Route.eventsRoute(eventsResource: EventsResource) {
  route("/events") {
    put {
      val request = Json.decodeFromString<CreateEventRequest>(call.receiveText())
      val event = eventsResource.create(request)
      call.respondText(contentType = Application.Json, status = Created) {
        Json.encodeToString(event)
      }
    }

    route("/{id}") { get {} }
  }
}

class EventsResource(private val eventsService: EventsService) {
  fun create(request: CreateEventRequest): Event {
    return eventsService.create(request)
  }
}
