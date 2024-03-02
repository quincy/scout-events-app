package com.troop77eagle

import com.troop77eagle.plugins.configureHTTP
import com.troop77eagle.plugins.configureMonitoring
import com.troop77eagle.plugins.configureRouting
import com.troop77eagle.plugins.configureSecurity
import com.troop77eagle.plugins.configureSerialization
import io.ktor.server.application.Application
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty

fun main() {
  embeddedServer(Netty, port = 8080, host = "0.0.0.0", module = Application::module)
      .start(wait = true)
}

fun Application.module() {
  configureSerialization()
  configureMonitoring()
  configureHTTP()
  configureSecurity()
  configureRouting()
}
