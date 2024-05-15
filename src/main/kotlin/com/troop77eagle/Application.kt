package com.troop77eagle

import com.troop77eagle.events.EventsDAO
import com.troop77eagle.events.EventsResource
import com.troop77eagle.events.EventsService
import com.troop77eagle.events.InstantMapper
import com.troop77eagle.plugins.configureHTTP
import com.troop77eagle.plugins.configureMonitoring
import com.troop77eagle.plugins.configureRouting
import com.troop77eagle.plugins.configureSecurity
import com.troop77eagle.plugins.configureSerialization
import io.ktor.server.application.Application
import io.ktor.server.config.ApplicationConfig
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import javax.sql.DataSource
import kotlinx.datetime.Instant
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.kotlin.KotlinPlugin
import org.jdbi.v3.sqlobject.SqlObjectPlugin
import org.jdbi.v3.sqlobject.kotlin.KotlinSqlObjectPlugin
import org.jdbi.v3.sqlobject.kotlin.onDemand
import org.postgresql.ds.PGSimpleDataSource

fun main() {
  embeddedServer(Netty, port = 8080, host = "0.0.0.0", module = Application::module)
      .start(wait = true)
}

fun Application.module() {
  val datasource = createDatasource(environment.config.config("db"))
  val jdbi = getJdbi(datasource)

  configureSerialization()
  configureMonitoring()
  configureHTTP()
  configureSecurity()
  configureRouting(jdbi, getEventsResource(jdbi))
}

fun Application.createDatasource(dbConfig: ApplicationConfig): DataSource =
    PGSimpleDataSource().apply {
      applicationName = dbConfig.property("app").getString()
      serverNames = arrayOf(dbConfig.property("host").getString())
      portNumbers = IntArray(dbConfig.property("port").getString().toInt())
      databaseName = dbConfig.property("database").getString()
      user = dbConfig.property("username").getString()
      password = dbConfig.property("password").getString()
      sslmode = "verify-full"
    }

fun Application.getEventsResource(jdbi: Jdbi): EventsResource =
    EventsResource(EventsService.create(jdbi.onDemand<EventsDAO>()))

fun getJdbi(datasource: DataSource): Jdbi =
    Jdbi.create(datasource).apply {
      installPlugin(KotlinPlugin())
      installPlugin(SqlObjectPlugin())
      installPlugin(KotlinSqlObjectPlugin())
      registerColumnMapper(Instant::class.java, InstantMapper())
    }
