package com.troop77eagle.events

import com.ninja_squad.dbsetup.Operations.deleteAllFrom
import com.ninja_squad.dbsetup.operation.Operation
import com.troop77eagle.getJdbi
import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import liquibase.command.CommandScope
import liquibase.command.core.UpdateCommandStep.CHANGELOG_FILE_ARG
import liquibase.command.core.UpdateCommandStep.COMMAND_NAME
import liquibase.command.core.helpers.DbUrlConnectionArgumentsCommandStep.DATABASE_ARG
import liquibase.database.DatabaseFactory
import liquibase.database.jvm.JdbcConnection
import org.testcontainers.containers.CockroachContainer

object CockroachDbContainer {
  val dbContainer =
      CockroachContainer("cockroachdb/cockroach:v23.1.17")
          .apply {
            startupAttempts = 1
            withUsername("app")
            withPassword("password")
            withDatabaseName("scouting")
          }
          .also { it.start() }

  val datasource =
      HikariConfig()
          .apply {
            jdbcUrl = dbContainer.jdbcUrl
            driverClassName = dbContainer.driverClassName
            username = dbContainer.username
            password = dbContainer.password
          }
          .let { HikariDataSource(it) }
          .also {
            val db =
                DatabaseFactory.getInstance()
                    .findCorrectDatabaseImplementation(JdbcConnection(it.connection))
            CommandScope(*COMMAND_NAME)
                .addArgumentValue(DATABASE_ARG, db)
                //                .addArgumentValue(USERNAME_ARG, "liquibase")
                .addArgumentValue(CHANGELOG_FILE_ARG, "/db/changelog/changelog-root.yaml")
                .execute()
          }

  val jdbi = getJdbi(datasource)

  val DELETE_ALL: Operation = deleteAllFrom("events")
}
