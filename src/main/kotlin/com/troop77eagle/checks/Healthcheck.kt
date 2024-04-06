package com.troop77eagle.checks

import io.github.oshai.kotlinlogging.KotlinLogging
import io.ktor.http.HttpStatusCode
import org.jdbi.v3.core.ConnectionException
import org.jdbi.v3.core.Jdbi

private val logger = KotlinLogging.logger {}

interface Healthcheck {
  fun check(): Response

  data class Response(val status: Int, val body: Map<String, String>) {
    constructor(status: HttpStatusCode, body: Map<String, String>) : this(status.value, body)
  }
}

class BasicHealthcheck : Healthcheck {
  override fun check(): Healthcheck.Response =
      Healthcheck.Response(HttpStatusCode.OK, mapOf("service" to "ok"))
}

class DeepHealthcheck(private val jdbi: Jdbi) : Healthcheck {
  override fun check(): Healthcheck.Response {
    return try {
      jdbi.open().use {
        Healthcheck.Response(HttpStatusCode.OK, mapOf("service" to "ok", "db" to "ok"))
      }
      // good
    } catch (e: ConnectionException) {
      logger.error(e) { "Database connection failed" }
      Healthcheck.Response(
          HttpStatusCode.ServiceUnavailable, mapOf("service" to "unhealthy", "db" to "unhealthy"))
    }
  }
}
