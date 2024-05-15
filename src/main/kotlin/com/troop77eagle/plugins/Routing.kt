package com.troop77eagle.plugins

import com.troop77eagle.checks.BasicHealthcheck
import com.troop77eagle.checks.DeepHealthcheck
import com.troop77eagle.events.EventsResource
import com.troop77eagle.events.eventsRoute
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.HttpStatusCode.Companion.BadRequest
import io.ktor.http.HttpStatusCode.Companion.InternalServerError
import io.ktor.http.HttpStatusCode.Companion.NotFound
import io.ktor.server.application.Application
import io.ktor.server.application.call
import io.ktor.server.application.install
import io.ktor.server.http.content.staticResources
import io.ktor.server.plugins.NotFoundException
import io.ktor.server.plugins.autohead.AutoHeadResponse
import io.ktor.server.plugins.statuspages.StatusPages
import io.ktor.server.resources.Resources
import io.ktor.server.response.respond
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.route
import io.ktor.server.routing.routing
import io.ktor.server.webjars.Webjars
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.MissingFieldException
import org.jdbi.v3.core.Jdbi

@OptIn(ExperimentalSerializationApi::class)
fun Application.configureRouting(jdbi: Jdbi, eventsResource: EventsResource) {
  install(Webjars) {
    path = "/webjars" // defaults to /webjars
  }
  install(StatusPages) {
    exception<Throwable> { call, cause ->
      when (cause) {
        is NotFoundException -> call.respondText(text = "$cause", status = NotFound)
        is MissingFieldException -> call.respondText(text = "$cause", status = BadRequest)
        else -> call.respondText(text = "500: $cause", status = InternalServerError)
      }
    }
  }
  install(Resources)
  install(AutoHeadResponse)

  val basicHealthcheck = BasicHealthcheck()
  val deepHealthcheck = DeepHealthcheck(jdbi)

  routing {
    get("/") { call.respondText("Hello World!") }
    get("/healthcheck") {
      basicHealthcheck.check().apply {
        call.respond(HttpStatusCode.fromValue(this.status), this.body)
      }
    }
    get("/deepcheck") {
      deepHealthcheck.check().apply { call.respond(HttpStatusCode.fromValue(status), body) }
    }

    get("/webjars") {
      call.respondText("<script src='/webjars/jquery/jquery.js'></script>", ContentType.Text.Html)
    }
    // Static plugin. Try to access `/static/index.html`
    staticResources(remotePath = "/static", basePackage = "static")

    // REST API
    route("/api/v1") { eventsRoute(eventsResource) }
  }
}
