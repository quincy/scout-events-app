package com.troop77eagle

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
import org.jdbi.v3.core.Jdbi
import org.postgresql.ds.PGSimpleDataSource

fun main() {
  embeddedServer(Netty, port = 8080, host = "0.0.0.0", module = Application::module)
      .start(wait = true)
}

fun Application.module() {
  val datasource = createDatasource(environment.config.config("db"))

  configureSerialization()
  configureMonitoring()
  configureHTTP()
  configureSecurity()
  configureRouting(Jdbi.create(datasource))
}

fun createDatasource(dbConfig: ApplicationConfig): DataSource {
  return PGSimpleDataSource().apply {
    applicationName = dbConfig.property("app").getString()
    serverNames = arrayOf(dbConfig.property("host").getString())
    portNumbers = IntArray(dbConfig.property("port").getString().toInt())
    databaseName = dbConfig.property("database").getString()
    user = dbConfig.property("username").getString()
    password = dbConfig.property("password").getString()
    sslmode = "verify-full"
  }
}
