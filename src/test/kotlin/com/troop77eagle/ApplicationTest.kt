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
import io.mockk.Runs
import io.mockk.every
import io.mockk.just
import io.mockk.mockk
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlinx.serialization.json.Json
import org.jdbi.v3.core.ConnectionException
import org.jdbi.v3.core.Handle
import org.jdbi.v3.core.Jdbi

class ApplicationTest {
  @Test
  fun testRoot() = testApplication {
    application { configureRouting(mockk(), mockk()) }
    client.get("/").apply {
      assertEquals(HttpStatusCode.OK, status)
      assertEquals("Hello World!", bodyAsText())
    }
  }

  @Test
  fun `healthcheck is healthy`() = testApplication {
    install(ContentNegotiation) { json(Json { prettyPrint = true }) }
    application { configureRouting(mockk(), mockk()) }
    client.get("/healthcheck").apply {
      assertEquals(HttpStatusCode.OK, status)
      assertEquals(mapOf("service" to "ok"), Json.decodeFromString(bodyAsText()))
    }
  }

  @Test
  fun `deepcheck is healthy`() = testApplication {
    install(ContentNegotiation) { json(Json { prettyPrint = true }) }
    application {
      configureRouting(
          mockk<Jdbi> {
            every<Handle> { open() } returns mockk<Handle> { every { close() } just Runs }
          },
          mockk())
    }
    client
        .get("/deepcheck") { this.header(HttpHeaders.Accept, ContentType.Application.Json) }
        .apply {
          assertEquals(HttpStatusCode.OK, status)
          assertEquals(mapOf("service" to "ok", "db" to "ok"), Json.decodeFromString(bodyAsText()))
        }
  }

  @Test
  fun `deepcheck is because database is unhealthy`() = testApplication {
    install(ContentNegotiation) { json(Json { prettyPrint = true }) }
    application {
      configureRouting(
          mockk<Jdbi> {
            every<Handle?> { open() } throws ConnectionException(IllegalStateException())
          },
          mockk())
    }
    client
        .get("/deepcheck") { this.header(HttpHeaders.Accept, ContentType.Application.Json) }
        .apply {
          assertEquals(HttpStatusCode.ServiceUnavailable, status)
          assertEquals(
              mapOf("service" to "unhealthy", "db" to "unhealthy"),
              Json.decodeFromString(bodyAsText()))
        }
  }
}
