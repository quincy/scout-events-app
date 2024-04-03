package com.troop77eagle

import com.troop77eagle.plugins.configureRouting
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.statement.bodyAsText
import io.ktor.http.ContentType
import io.ktor.http.HttpHeaders
import io.ktor.http.HttpStatusCode
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.testing.testApplication
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlinx.serialization.json.Json

class ApplicationTest {
  @Test
  fun testRoot() = testApplication {
    application { configureRouting() }
    client.get("/").apply {
      assertEquals(HttpStatusCode.OK, status)
      assertEquals("Hello World!", bodyAsText())
    }
  }

  @Test
  fun `healthcheck is healthy`() = testApplication {
    install(ContentNegotiation) { json(Json { prettyPrint = true }) }
    application { configureRouting() }
    client.get("/healthcheck").apply {
      assertEquals(HttpStatusCode.OK, status)
      assertEquals(mapOf("service" to "ok"), Json.decodeFromString(bodyAsText()))
    }
  }

  @Test
  fun `deepcheck is healthy`() = testApplication {
    install(ContentNegotiation) { json(Json { prettyPrint = true }) }
    application { configureRouting() }
    client
        .get("/deepcheck") { this.header(HttpHeaders.Accept, ContentType.Application.Json) }
        .apply {
          assertEquals(HttpStatusCode.OK, status)
          assertEquals(mapOf("service" to "ok", "db" to "ok"), Json.decodeFromString(bodyAsText()))
        }
  }
}
