package com.troop77eagle.plugins

import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.resources.Resource
import io.ktor.server.application.Application
import io.ktor.server.application.call
import io.ktor.server.application.install
import io.ktor.server.http.content.staticResources
import io.ktor.server.plugins.autohead.AutoHeadResponse
import io.ktor.server.plugins.statuspages.StatusPages
import io.ktor.server.resources.Resources
import io.ktor.server.resources.get
import io.ktor.server.response.respond
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.routing
import io.ktor.server.webjars.Webjars
import kotlinx.serialization.Serializable

fun Application.configureRouting() {
  install(Webjars) {
    path = "/webjars" // defaults to /webjars
  }
  install(StatusPages) {
    exception<Throwable> { call, cause ->
      call.respondText(text = "500: $cause", status = HttpStatusCode.InternalServerError)
    }
  }
  install(Resources)
  install(AutoHeadResponse)
  routing {
    get("/") { call.respondText("Hello World!") }
    get("/webjars") {
      call.respondText("<script src='/webjars/jquery/jquery.js'></script>", ContentType.Text.Html)
    }
    // Static plugin. Try to access `/static/index.html`
    staticResources(remotePath = "/static", basePackage = "static")
    get<Articles> { article ->
      // Get all articles ...
      call.respond("List of articles sorted starting from ${article.sort}")
    }
  }
}

@Serializable @Resource("/articles") class Articles(val sort: String? = "new")
