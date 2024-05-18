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
import io.ktor.server.netty.EngineMain
import javax.sql.DataSource
import kotlinx.datetime.Instant
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.kotlin.KotlinPlugin
import org.jdbi.v3.sqlobject.SqlObjectPlugin
import org.jdbi.v3.sqlobject.kotlin.KotlinSqlObjectPlugin
import org.jdbi.v3.sqlobject.kotlin.onDemand
import org.postgresql.ds.PGSimpleDataSource

fun main(args: Array<String>): Unit = EngineMain.main(args)

fun Application.module() {
  val datasource = createDatasource(environment.config.config("db"))
  val jdbi = getJdbi(datasource)

  configureSerialization()
  configureMonitoring()
  configureHTTP()
  configureSecurity()
  configureRouting(jdbi, getEventsResource(jdbi))
}

fun createDatasource(dbConfig: ApplicationConfig): DataSource =
    PGSimpleDataSource().apply {
      applicationName = dbConfig.property("app").getString()
      serverNames = arrayOf(dbConfig.property("host").getString())
      portNumbers = intArrayOf(dbConfig.property("port").getString().toInt())
      sslRootCert = dbConfig.property("sslRootCert").getString()
      databaseName = dbConfig.property("database").getString()
      user = dbConfig.property("username").getString()
      password = dbConfig.property("password").getString()
      sslmode = "verify-full"
    }

fun getEventsResource(jdbi: Jdbi): EventsResource =
    EventsResource(EventsService.create(jdbi.onDemand<EventsDAO>()))

fun getJdbi(datasource: DataSource): Jdbi =
    Jdbi.create(datasource).apply {
      installPlugin(KotlinPlugin())
      installPlugin(SqlObjectPlugin())
      installPlugin(KotlinSqlObjectPlugin())
      registerColumnMapper(Instant::class.java, InstantMapper())
    }
