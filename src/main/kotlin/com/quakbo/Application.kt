package com.quakbo

import com.quakbo.plugins.configureHTTP
import com.quakbo.plugins.configureMonitoring
import com.quakbo.plugins.configureRouting
import com.quakbo.plugins.configureSecurity
import com.quakbo.plugins.configureSerialization
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
