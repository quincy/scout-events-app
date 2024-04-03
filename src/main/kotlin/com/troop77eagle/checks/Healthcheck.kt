package com.troop77eagle.checks

import io.ktor.http.HttpStatusCode

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

class DeepHealthcheck : Healthcheck {
  override fun check(): Healthcheck.Response =
      Healthcheck.Response(HttpStatusCode.OK, mapOf("service" to "ok", "db" to "ok"))
}
